package yeelight

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	timeout  = time.Duration(2 * time.Second)
	discover = "M-SEARCH * HTTP/1.1\r\nHOST:239.255.255.250:1982\r\nMAN:\"ssdp:discover\"\r\nST:wifi_bulb\r\n"
	// error messages
	errResolveTCP   = errors.New("Cannot resolve TCP address")
	errConnectLight = errors.New("Cannot connect to light")
	errInvalidParam = errors.New("Invalid parameter value")
)

// TODO: if duration > 0 -> effect = smooth otherwise sudden

//Yeelight are the light properties
type Yeelight struct {
	Location   string   `json:"location"`
	ID         string   `json:"id,omitempty"`
	Model      string   `json:"model,omitempty"`
	FWVersion  int      `json:"fw_ver,omitempty"`
	Support    []string `json:"support,omitempty"`
	Power      string   `json:"power,omitempty"`
	Bright     int      `json:"bright,omitempty"`
	ColorMode  int      `json:"color_mode,omitempty"`
	ColorTemp  int      `json:"ct,omitempty"`
	RGB        int      `json:"rgb,omitempty"`
	Hue        int      `json:"hue,omitempty"`
	Saturation int      `json:"sat,omitempty"`
	Name       string   `json:"name"`
}

//Command to send to the light
type Command struct {
	ID     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

//Response sent by the light
type Response struct {
	ID     int         `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  Error       `json:"error,omitempty"`
}

//Error struct is used on the ResponseError payload
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//Discover uses SSDP to find and return the IP address of the lights
//credit: https://github.com/edgard/yeelight/blob/master/yeelight.go
func Discover(timeout time.Duration) ([]Yeelight, error) {
	laddr, err := net.ResolveUDPAddr("udp4", ":0")
	if err != nil {
		return nil, err
	}
	maddr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1982")
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	go func() {
		conn.WriteToUDP([]byte(discover), maddr)
	}()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	answers := make(map[string]string)
	for {
		answer := make([]byte, 1024)
		n, src, err := conn.ReadFromUDP(answer)
		if err != nil {
			break
		}
		answers[src.String()] = string(answer[:n])
	}

	var lights []Yeelight
	for _, answer := range answers {
		tp := textproto.NewReader(bufio.NewReader(strings.NewReader(answer)))
		tp.ReadLine()
		header, _ := tp.ReadMIMEHeader()

		var light Yeelight
		location, _ := url.Parse(header.Get("location"))
		light.Location = location.Host
		light.ID = header.Get("id")
		light.Model = header.Get("model")
		light.FWVersion, _ = strconv.Atoi(header.Get("fw_ver"))
		light.Support = strings.Split(header.Get("support"), " ")
		light.Power = header.Get("power")
		light.Bright, _ = strconv.Atoi(header.Get("bright"))
		light.ColorMode, _ = strconv.Atoi(header.Get("color_mode"))
		light.ColorTemp, _ = strconv.Atoi(header.Get("ct"))
		light.RGB, _ = strconv.Atoi(header.Get("rgb"))
		light.Hue, _ = strconv.Atoi(header.Get("hue"))
		light.Saturation, _ = strconv.Atoi(header.Get("sat"))
		light.Name = header.Get("name")

		lights = append(lights, light)
	}

	return lights, err
}

// Handles the request
func (y *Yeelight) request(cmd Command) (Response, error) {
	conn, err := net.DialTimeout("tcp", y.Location, timeout)
	defer conn.Close()
	if err != nil {
		return Response{}, errResolveTCP
	}

	cmdJSON, err := json.Marshal(cmd)
	if err != nil {
		return Response{}, errInvalidParam
	}

	if _, err := fmt.Fprintf(conn, "%s\r\n", cmdJSON); err != nil {
		return Response{}, errConnectLight
	}

	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return Response{}, errConnectLight
	}

	// parse response
	resp := Response{}
	err = json.Unmarshal([]byte(data), &resp)
	if err != nil {
		return Response{}, err
	}

	return resp, nil
}

//GetProp method is used to retrieve current property a light
func (y *Yeelight) GetProp() error {
	cmd := Command{
		ID:     1,
		Method: "get_prop",
		Params: []interface{}{"power", "bright", "ct", "rgb", "hue", "sat", "color_mode", "name"},
	}

	resp, err := y.request(cmd)
	if err != nil {
		return err
	}

	y.Power = resp.Result.([]interface{})[0].(string)
	y.Bright, _ = strconv.Atoi(resp.Result.([]interface{})[1].(string))
	y.ColorTemp, _ = strconv.Atoi(resp.Result.([]interface{})[2].(string))
	y.RGB, _ = strconv.Atoi(resp.Result.([]interface{})[3].(string))
	y.Hue, _ = strconv.Atoi(resp.Result.([]interface{})[4].(string))
	y.Saturation, _ = strconv.Atoi(resp.Result.([]interface{})[5].(string))
	y.ColorMode, _ = strconv.Atoi(resp.Result.([]interface{})[6].(string))
	y.Name, _ = resp.Result.([]interface{})[7].(string)

	return nil
}

//SetCtAbx method is used to change the color temperature of a smart LED.
func (y *Yeelight) SetCtAbx(value, duration int) (Response, error) {
	var effect string

	if duration > 0 {
		effect = "smooth"
	} else {
		effect = "sudden"
		duration = 0
	}

	cmd := Command{
		ID:     2,
		Method: "set_ct_abx",
		Params: []interface{}{value, effect, duration},
	}

	return y.request(cmd)
}

//SetRGB method is used to change the color RGB of a smart LED.
func (y *Yeelight) SetRGB(value, duration int) (Response, error) {
	var effect string

	if duration > 0 {
		effect = "smooth"
	} else {
		effect = "sudden"
		duration = 0
	}

	cmd := Command{
		ID:     3,
		Method: "set_rgb",
		Params: []interface{}{value, effect, duration},
	}

	return y.request(cmd)
}

//SetHSV method is used to change the color of a smart LED.
func (y *Yeelight) SetHSV(hue, sat, duration int) (Response, error) {
	var effect string

	if duration > 0 {
		effect = "smooth"
	} else {
		effect = "sudden"
		duration = 0
	}

	cmd := Command{
		ID:     4,
		Method: "set_hsv",
		Params: []interface{}{hue, sat, effect, duration},
	}

	return y.request(cmd)
}

//SetBright method is used to change the brightness of a smart LED.
func (y *Yeelight) SetBright(brightness, duration int) (Response, error) {
	var effect string

	if duration > 0 {
		effect = "smooth"
	} else {
		effect = "sudden"
		duration = 0
	}

	cmd := Command{
		ID:     5,
		Method: "set_bright",
		Params: []interface{}{brightness, effect, duration},
	}

	return y.request(cmd)
}

//SetPower method is used to switch on or off the smart LED (software managed on/off).
func (y *Yeelight) SetPower(power string, duration int) (Response, error) {
	var effect string

	if duration > 0 {
		effect = "smooth"
	} else {
		effect = "sudden"
		duration = 0
	}

	cmd := Command{
		ID:     6,
		Method: "set_power",
		Params: []interface{}{power, effect, duration},
	}

	return y.request(cmd)
}

//Toggle method is used to toggle the smart LED.
func (y *Yeelight) Toggle() (Response, error) {
	cmd := Command{
		ID:     7,
		Method: "toggle",
	}

	return y.request(cmd)
}

//SetDefault method is used to save current state of smart LED in persistent
//memory. So if user powers off and then powers on the smart LED again (hard power reset),
//the smart LED will show last saved state.
func (y *Yeelight) SetDefault() (Response, error) {
	cmd := Command{
		ID:     8,
		Method: "set_default",
	}

	return y.request(cmd)
}

//StartCf method is used to start a color flow. Color flow is a series of smart
//LED visible state changing. It can be brightness changing, color changing or color
//temperature changing.This is the most powerful command. All our recommended scenes,
//e.g. Sunrise/Sunset effect is implemented using this method. With the flow expression, user
//can actually “program” the light effect.
func (y *Yeelight) StartCf(count, action int, flowExpression string) (Response, error) {
	cmd := Command{
		ID:     9,
		Method: "start_cf",
		Params: []interface{}{count, action, flowExpression},
	}

	return y.request(cmd)
}

//StopCf method is used to stop a running color flow.
func (y *Yeelight) StopCf() (Response, error) {
	cmd := Command{
		ID:     10,
		Method: "stop_cf",
	}

	return y.request(cmd)
}

//SetScene method is used to set the smart LED directly to specified state.
//If the smart LED is off, then it will turn on the smart LED firstly and then
//apply the specified command.
func (y *Yeelight) SetScene(class, values string) (Response, error) {
	cmd := Command{
		ID:     11,
		Method: "set_scene",
		Params: []interface{}{class, values},
	}

	return y.request(cmd)
}

//CronAdd method is used to start a timer job on the smart LED.
func (y *Yeelight) CronAdd(t, value int) (Response, error) {
	cmd := Command{
		ID:     12,
		Method: "cron_add",
		Params: []interface{}{t, value},
	}

	return y.request(cmd)
}

//CronGet method is used to retrieve the setting of the current cron job of the specified type.
func (y *Yeelight) CronGet(t int) (Response, error) {
	cmd := Command{
		ID:     13,
		Method: "cron_get",
		Params: []interface{}{t},
	}

	return y.request(cmd)
}

//CronDel method is used to stop the specified cron job.
func (y *Yeelight) CronDel(t int) (Response, error) {
	cmd := Command{
		ID:     14,
		Method: "cron_del",
		Params: []interface{}{t},
	}

	return y.request(cmd)
}

//SetAdjust method is used to change brightness, CT or color of a smart LED
//without knowing the current value, it's main used by controllers.
func (y *Yeelight) SetAdjust(action, prop string) (Response, error) {
	cmd := Command{
		ID:     15,
		Method: "set_adjust",
		Params: []interface{}{action, prop},
	}

	return y.request(cmd)
}

//SetName method is used to name the device. The name will be stored on the
//device and reported in discovering response. User can also read the name
//through “get_prop” method
func (y *Yeelight) SetName(name string) (Response, error) {
	cmd := Command{
		ID:     16,
		Method: "set_name",
		Params: []interface{}{name},
	}

	return y.request(cmd)
}

//On method is used to switch on the smart LED
func (y *Yeelight) On() (Response, error) {
	return y.SetPower("on", 1000)
}

//Off method is used to switch off the smart LED
func (y *Yeelight) Off() (Response, error) {
	return y.SetPower("off", 1000)
}

package yeelight

import "strings"

// MatchingName checks if a Yeelight match a name
func MatchingName(light Yeelight, arg string) bool {
	return light.Name == strings.ToLower(arg)
}

// MatchingIP checks if a Yeelight match an ip
func MatchingIP(light Yeelight, arg string) bool {
	return strings.Split(light.Location, ":")[0] == arg
}

// Matching checks if a Yeelight match a name or an ip
func Matching(light Yeelight, arg string) bool {
	return MatchingName(light, arg) || MatchingIP(light, arg)
}

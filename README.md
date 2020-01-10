# Yeego

[![GoDoc](https://godoc.org/github.com/julienrbrt/yeego?status.svg)](https://godoc.org/github.com/julienrbrt/yeego/light/yeelight)

Control your Yeelight bulbs over LAN with this simple go tool.

Yeego includes a CLI for simple controlling and demonstrate what the `yeelight` library can do and `yeelight`, a library implementing the Yeelight API.

## Installation

``` bash
go get github.com/julienrbrt/yeego
```

## Usage

The "Developer Mode" need to be enabled to discover and operate the device.

### Yeego

**Discover lights in your network**
```
yeego discover
```

**Turn on a light**
```
yeego on bedroom
yeego on 192.168.2.1
```

**Togge a light**
```
yeego toggle plant
yeego toggle 192.168.2.5
```

**Exhaustive list of supported commands**
```
yeego help
```

### package yeelight

**Example usage of Yeelight Package**

``` go
package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/julienrbrt/yeego/light/yeelight"
)

func main() {
	// discover yeelight on network
	lights, err := yeelight.Discover(time.Duration(time.Second))
	if err != nil {a
		fmt.Println(err)
		os.Exit(1)
	}

	if len(lights) == 0 {
		fmt.Println(errors.New("No Yeelight found."))
		os.Exit(1)
	}

	for i := range lights {
		// turn a a light
		lights[i].Toggle()
	}
}
```

The list of supported commands is present on [![GoDoc](https://godoc.org/github.com/julienrbrt/yeego?status.svg)](https://godoc.org/github.com/julienrbrt/yeego/light/yeelight) 

## Feature and bugs

Please file feature requests and bugs at the [issue tracker](https://github.com/julienrbrt/yeego/issues/).

## More Info

More info about Yeelight API:
* [Yeelight Developer](https://www.yeelight.com/en_US/developer)
* [Yeelight Inter-Operation Specification](doc/Yeelight_Inter-Operation_Spec.pdf)

Yeelight API heavily based on:
* https://github.com/nunows/goyeelight
* https://github.com/edgard/yeelight
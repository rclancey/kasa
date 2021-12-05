package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/rclancey/kasa"
)

func main() {
	devices, err := kasa.Discover(5 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	ch := make(chan error, len(devices))
	for _, dev := range devices {
		if len(os.Args) > 1 && !strings.HasPrefix(strings.ToLower(dev.Alias()), strings.ToLower(os.Args[1])) {
			continue
		}
		log.Println("trying to turn off", dev.Alias())
		plug, isa := dev.(kasa.Switch)
		if isa && plug.IsOn() {
			count += 1
			go func(p kasa.Switch) {
				err := p.TurnOff()
				if err != nil {
					log.Println(err)
				}
				ch <- err
			}(plug)
		}
	}
	for count > 0 {
		err = <-ch
		count -= 1
	}
}

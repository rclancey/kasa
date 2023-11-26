package main

import (
	"bufio"
	"fmt"
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
	for {
		for i, dev := range devices {
			fmt.Printf("%d: %s\n", i + 1, dev.Alias())
		}
		fmt.Print("Which device? ")
		var devId int
		_, err := fmt.Scanln(&devId)
		if err != nil {
			log.Println(err)
			break
		}
		devId -= 1
		if devId < 0 || devId >= len(devices) {
			break
		}
		dev := devices[devId]
		strip, ok := dev.(*kasa.SmartStrip)
		if ok {
			fmt.Printf("%d: %s\n", 0, "Strip")
			for i, child := range strip.Children() {
				fmt.Printf("%d: %s\n", i+1, child.Alias())
			}
			fmt.Print("Which child? ")
			var childId int
			_, err := fmt.Scanln(&childId)
			if err != nil {
				log.Println(err)
				break
			}
			if childId > 0 && childId <= len(strip.Children()) {
				dev = strip.Children()[childId-1]
			}
		}
		buf := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			query, err := buf.ReadString('\n')
			if err != nil {
				log.Println(err)
				break
			}
			if query == "q" || query == "quit" {
				break
			}
			if strings.HasPrefix(query, "{") {
				res, err := dev.Repl(query)
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println(res)
				}
			} else {
				parts := strings.Fields(query)
				resp, err := kasa.ExecDeviceCommand(dev, parts[0], parts[1:]...)
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println(resp...)
				}
			}
		}
	}
}

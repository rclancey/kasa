package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
			if dev.DeviceType() == kasa.DeviceTypeStrip {
				fmt.Printf("%d: %s\n", i + 1, dev.Alias())
			}
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
		dev, ok := devices[devId].(*kasa.SmartStrip)
		if !ok {
			log.Println("not a SmartStrip")
			continue
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(dev.GetSysInfo())
		children := dev.Children()
		for i, child := range children {
			fmt.Printf("%d: %s\n", i+1, child.Alias())
		}
		fmt.Print("Which socket? ")
		var childIndex int
		_, err = fmt.Scanln(&childIndex)
		if err != nil {
			log.Println(err)
			break
		}
		childIndex -= 1
		if childIndex < 0 || childIndex >= len(children) {
			break
		}
		child := children[childIndex]
		if child.IsOn() {
			child.TurnOff()
		} else {
			child.TurnOn()
		}
	}
}

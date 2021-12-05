package main

import (
	"encoding/json"
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
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(devices)
	os.Stdout.Write([]byte("\n"))
}

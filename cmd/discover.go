package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rclancey/kasa"
)

func main() {
	retry := 10 * time.Second
	quitch := make(chan bool, 2)
	ch, err := kasa.DiscoverStream(retry, quitch)
	if err != nil {
		log.Fatal(err)
	}
	seen := map[string]bool{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	enc := json.NewEncoder(os.Stdout)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-quitch:
				return
			case dev, ok := <-ch:
				if !ok {
					return
				}
				log.Println(dev.IP(), dev.Alias())
				if !seen[dev.IP()] {
					enc.Encode(dev)
					seen[dev.IP()] = true
				}
			}
		}
	}()
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
	log.Println("shutting down")
	quitch <- true
	quitch <- true
	wg.Wait()
}

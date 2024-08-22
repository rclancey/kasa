package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rclancey/kasa"
)

func main() {
	kasa.Debug = true
	retry := 10 * time.Second
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	ch, err := kasa.DiscoverStream(ctx, retry)
	if err != nil {
		log.Fatal(err)
	}
	seen := map[string]bool{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
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
	wg.Wait()
	log.Println("shutting down")
}

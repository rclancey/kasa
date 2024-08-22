package kasa

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

var Debug = false

func DiscoverStream(ctx context.Context, retry time.Duration) (chan SmartDevice, error) {
	query := map[string]interface{}{
		"system": map[string]interface{}{
			"get_sysinfo": nil,
		},
	}
	plain, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	sendaddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9999")
	if err != nil {
		return nil, err
	}
	respaddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return nil, err
	}
	l, err := net.ListenUDP("udp", respaddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("listening on", l.LocalAddr().String())
	maxSize := 8192
	l.SetReadBuffer(maxSize)
	ch := make(chan SmartDevice, 10)
	quit := &atomic.Bool{}
	go func() {
		defer close(ch)
		for {
			if quit.Load() {
				return
			}
			b := make([]byte, maxSize)
			l.SetReadDeadline(time.Now().Add(time.Second))
			n, src, err := l.ReadFromUDP(b)
			if err != nil {
				if !strings.Contains(err.Error(), "i/o timeout") {
					log.Println(err)
				}
				continue
			}
			if quit.Load() {
				return
			}
			plain := decrypt(b[:n])
			if Debug {
				log.Println(string(plain))
			}
			info := &Query{}
			err = json.Unmarshal(plain, &info)
			if err != nil {
				log.Println(err)
				continue
			}
			dev := &BaseDevice{Addr: src.IP.String(), Info: info}
			dev.setUpdateTime()
			ch <- dev.AsConcrete()
		}
	}()

	go func() {
		defer quit.Store(true)
		requery := func() error {
			_, err = l.WriteTo(encrypt(plain), sendaddr)
			if err != nil {
				return err
			}
			return nil
		}
		err := requery()
		if err != nil {
			log.Println("error querying:", err)
			return
		}
		ticker := time.NewTicker(retry)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err = requery()
				if err != nil {
					log.Println("error querying:", err)
					return
				}
			}
		}
	}()
	return ch, nil
}

func Discover(timeout time.Duration) ([]SmartDevice, error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	quitch := make(chan bool, 2)
	ch, err := DiscoverStream(ctx, timeout)
	if err != nil {
		return nil, err
	}
	devices := []SmartDevice{}
	timer := time.NewTimer(timeout)
	for {
		select {
		case dev, ok := <-ch:
			if !ok {
				quitch <- true
				return devices, nil
			}
			devices = append(devices, dev)
		case <-timer.C:
			quitch <- true
			return devices, nil
		}
	}
	return devices, nil
}

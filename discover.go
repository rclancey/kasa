package kasa

import (
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"
)

func DiscoverStream(retry time.Duration, quitch chan bool) (chan SmartDevice, error) {
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
	quit := false
	go func() {
		defer close(ch)
		for {
			if quit {
				return
			}
			b := make([]byte, maxSize)
			l.SetReadDeadline(time.Now().Add(retry))
			n, src, err := l.ReadFromUDP(b)
			if err != nil {
				if !strings.Contains(err.Error(), "i/o timeout") {
					log.Println(err)
				}
				continue
			}
			plain := decrypt(b[:n])
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
		defer func() {
			quit = true
		}()
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
			case <-quitch:
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
	quitch := make(chan bool, 2)
	ch, err := DiscoverStream(timeout, quitch)
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

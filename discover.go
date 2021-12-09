package kasa

import (
	"encoding/json"
	"log"
	"net"
	//"strings"
	"time"
)

func Discover(timeout time.Duration) ([]SmartDevice, error) {
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
	maxSize := 8192
	l.SetReadBuffer(maxSize)
	ch := make(chan SmartDevice, 10)
	go func() {
		defer close(ch)
		for {
			b := make([]byte, maxSize)
			l.SetReadDeadline(time.Now().Add(timeout))
			n, src, err := l.ReadFromUDP(b)
			if err != nil {
					log.Println(err)
				/*
				if !strings.Contains(err.Error(), "i/o timeout") {
				}
				*/
				return
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

	_, err = l.WriteTo(encrypt(plain), sendaddr)
	if err != nil {
		return nil, err
	}
	devices := []SmartDevice{}
	for {
		dev, ok := <-ch
		if !ok {
			break
		}
		devices = append(devices, dev)
	}
	return devices, nil
}

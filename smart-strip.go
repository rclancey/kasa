package kasa

import (
	"log"
)

type SmartStrip struct {
	*BaseDevice
}

func (strip *SmartStrip) Children() []*SmartStripSocket {
	sysinfo := strip.GetSysInfo()
	if sysinfo == nil {
		log.Println("strip has no sysinfo")
		return nil
	}
	if sysinfo.Children == nil {
		log.Println("strip has no children")
		return nil
	}
	children := make([]*SmartStripSocket, len(sysinfo.Children))
	for i, childinfo := range sysinfo.Children {
		children[i] = &SmartStripSocket{strip, childinfo.ID}
	}
	//log.Printf("strip has %d children", len(children))
	return children
}

func (strip *SmartStrip) IsOff() bool {
	return !strip.IsOn()
}

func (strip *SmartStrip) IsOn() bool {
	for _, child := range strip.Children() {
		if child.IsOn() {
			return true
		}
	}
	return false
}

func (strip *SmartStrip) TurnOn() error {
	for _, child := range strip.Children() {
		err := child.TurnOn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (strip *SmartStrip) TurnOff() error {
	for _, child := range strip.Children() {
		err := child.TurnOff()
		if err != nil {
			return err
		}
	}
	return nil
}

func (strip *SmartStrip) SetLED(state bool) error {
	return nil
}

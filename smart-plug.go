package kasa

import (
	"log"
)

type SmartPlug struct {
	*BaseDevice
}

func (plug *SmartPlug) IsOff() bool {
	return !plug.IsOn()
}

func (plug *SmartPlug) IsOn() bool {
	sysinfo := plug.GetSysInfo()
	if sysinfo == nil {
		return false
	}
	return sysinfo.RelayState > 0
}

func (plug *SmartPlug) TurnOn() error {
	var res interface{}
	err := plug.Query(&res, "system", "set_relay_state", map[string]interface{}{"state": 1})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (plug *SmartPlug) TurnOff() error {
	var res interface{}
	err := plug.Query(&res, "system", "set_relay_state", map[string]interface{}{"state": 0})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (plug *SmartPlug) SetLED(state bool) error {
	return nil
}

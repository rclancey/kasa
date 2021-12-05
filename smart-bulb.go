package kasa

import (
	"log"
)

type SmartBulb struct {
	*BaseDevice
}

func (plug *SmartBulb) GetLightState() *LightState {
	sysinfo := plug.GetSysInfo()
	if sysinfo == nil {
		return nil
	}
	return sysinfo.LightState
}

func (plug *SmartBulb) IsOff() bool {
	return !plug.IsOn()
}

func (plug *SmartBulb) IsOn() bool {
	light := plug.GetLightState()
	if light == nil {
		return false
	}
	return light.OnOff > 0
}

func (plug *SmartBulb) TurnOn() error {
	var res interface{}
	err := plug.Query(&res, "smartlife.iot.smartbulb.lightingservice", "transition_light_state", map[string]interface{}{"on_off": 1, "ignore_default": 1})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (plug *SmartBulb) TurnOff() error {
	var res interface{}
	err := plug.Query(&res, "smartlife.iot.smartbulb.lightingservice", "transition_light_state", map[string]interface{}{"on_off": 0, "ignore_default": 1})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (plug *SmartBulb) SetLED(state bool) error {
	return nil
}

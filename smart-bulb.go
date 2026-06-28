package kasa

import (
	"errors"
	"log"
)

type SmartBulb struct {
	*BaseDevice
}

func (bulb *SmartBulb) GetLightState() *LightState {
	sysinfo := bulb.GetSysInfo()
	if sysinfo == nil {
		return nil
	}
	return sysinfo.LightState
}

func (bulb *SmartBulb) IsOff() bool {
	return !bulb.IsOn()
}

func (bulb *SmartBulb) IsOn() bool {
	light := bulb.GetLightState()
	if light == nil {
		return false
	}
	return light.OnOff > 0
}

func (bulb *SmartBulb) TurnOn() error {
	var res interface{}
	err := bulb.Query(&res, "smartlife.iot.smartbulb.lightingservice", "transition_light_state", map[string]interface{}{"on_off": 1, "ignore_default": 1})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (bulb *SmartBulb) TurnOff() error {
	var res interface{}
	err := bulb.Query(&res, "smartlife.iot.smartbulb.lightingservice", "transition_light_state", map[string]interface{}{"on_off": 0, "ignore_default": 1})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (bulb *SmartBulb) IsDimmable() bool {
	sysinfo := bulb.GetSysInfo()
	if sysinfo == nil {
		return false
	}
	return sysinfo.IsDimmable > 0
}

func (bulb *SmartBulb) SetBrightness(b int) error {
	if !bulb.IsDimmable() {
		return errors.New("device is not dimmable")
	}
	if b <= 0 {
		return bulb.TurnOff()
	}
	if b > 100 {
		b = 100
	}
	var res interface{}
	err := bulb.Query(&res, "smartlife.iot.smartbulb.lightingservice", "transition_light_state", map[string]interface{}{"on_off": 1, "ignore_default": 1, "brightness": b})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (bulb *SmartBulb) IsVariableColorTemp() bool {
	sysinfo := bulb.GetSysInfo()
	if sysinfo == nil {
		return false
	}
	return sysinfo.IsVariableColorTemp > 0
}

func (bulb *SmartBulb) SetColorTemp(temp int) error {
	if !bulb.IsVariableColorTemp() {
		return errors.New("color not controllable")
	}
	params := map[string]any{
		"color_temp": max(2500, min(6500, temp)),
		"hue": 0,
		"saturation": 0,
	}
	var res any
	err := bulb.Query(&res, "smartlife.iot.smartbulb.lightingservice", "transition_light_state", params)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (bulb *SmartBulb) IsColor() bool {
	sysinfo := bulb.GetSysInfo()
	if sysinfo == nil {
		return false
	}
	return sysinfo.IsColor > 0
}

func (bulb *SmartBulb) SetColor(hue, saturation int) error {
	if !bulb.IsColor() {
		return errors.New("color not controllable")
	}
	params := map[string]any{
		"color_temp": 0,
		"hue": max(0, min(360, hue)),
		"saturation": max(0, min(100, saturation)),
	}
	var res any
	err := bulb.Query(&res, "smartlife.iot.smartbulb.lightingservice", "transition_light_state", params)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (bulb *SmartBulb) SetLED(state bool) error {
	return nil
}

func (bulb *SmartBulb) GetDetails() (interface{}, error) {
	var res interface{}
	err := bulb.Query(&res, "smartlife.iot.smartbulb.lightingservice", "get_light_details", nil)
	return res, err
}

func (bulb *SmartBulb) QueryLightState() (interface{}, error) {
	var res interface{}
	err := bulb.Query(&res, "smartlife.iot.smartbulb.lightingservice", "get_light_state", nil)
	return res, err
}

func (bulb *SmartBulb) QueryTurnOnBehavior() (interface{}, error) {
	var res interface{}
	err := bulb.Query(&res, "smartlife.iot.smartbulb.lightingservice", "get_default_behavior", nil)
	return res, err
}

func (bulb *SmartBulb) GetLightService() string {
	return "smartlife.iot.smartbulb.lightingservice"
}

func (bulb *SmartBulb) GetTimeService() string {
	return "smartlife.iot.common.timesetting"
}

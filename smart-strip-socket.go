package kasa

import (
	"encoding/json"
	"log"
	"time"
)

type SmartStripSocket struct {
	*SmartStrip
	id string
}

func (plug *SmartStripSocket) SetAlias(alias string) error {
	var res interface{}
	args := &SetAliasRequest{Alias: alias}
	err := plug.Query(&res, "system", "set_dev_alias", args, plug.id)
	if err != nil {
		log.Println("error in SetAlias():", err)
		return err
	}
	data, _ := json.Marshal(res)
	log.Println("SetAlias() =>", string(data))
	plug.Update()
	return nil
}

func (plug *SmartStripSocket) Alias() string {
	sysinfo := plug.GetSysInfo()
	if sysinfo == nil {
		return ""
	}
	return sysinfo.Alias
}

func (plug *SmartStripSocket) DeviceID() string {
	return plug.SmartStrip.DeviceID() + plug.id
}

func (plug *SmartStripSocket) DeviceType() DeviceType {
	return DeviceTypeStripSocket
}

func (plug *SmartStripSocket) GetSysInfo() *SysInfo {
	sysinfo := plug.SmartStrip.GetSysInfo()
	if sysinfo == nil {
		return nil
	}
	if sysinfo.Children == nil {
		return nil
	}
	for _, child := range sysinfo.Children {
		if child.ID == plug.id {
			return child
		}
	}
	return nil
}

func (plug *SmartStripSocket) IsOff() bool {
	return !plug.IsOn()
}

func (plug *SmartStripSocket) IsOn() bool {
	sysinfo := plug.GetSysInfo()
	if sysinfo == nil {
		return false
	}
	return sysinfo.State > 0
}

func (plug *SmartStripSocket) IsStripSocket() bool {
	return true
}

func (plug *SmartStripSocket) OnSince() *time.Time {
	sysinfo := plug.GetSysInfo()
	if sysinfo == nil {
		return nil
	}
	if sysinfo.OnTime == 0 {
		return nil
	}
	t := time.Now().Add(-1 * time.Second * time.Duration(sysinfo.OnTime))
	return &t
}

func (plug *SmartStripSocket) TurnOn() error {
	var res interface{}
	err := plug.Query(&res, "system", "set_relay_state", map[string]interface{}{"state": 1}, plug.DeviceID())
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (plug *SmartStripSocket) TurnOff() error {
	var res interface{}
	err := plug.Query(&res, "system", "set_relay_state", map[string]interface{}{"state": 0}, plug.DeviceID())
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(res)
	return nil
}

func (plug *SmartStripSocket) SetLED(state bool) error {
	return nil
}

package kasa

import (
	"strings"
	"time"
)

type DeviceType string

const (
	DeviceTypeDimmer = DeviceType("SmartDimmer")
	DeviceTypeStrip = DeviceType("SmartStrip")
	DeviceTypeStripSocket = DeviceType("SmartStripSocket")
	DeviceTypePlug = DeviceType("SmartPlug")
	DeviceTypeLightStrip = DeviceType("SmartLightStrip")
	DeviceTypeBulb = DeviceType("SmartBulb")
	DeviceTypeUnknown = DeviceType("Unknown")
)

type Action struct {
	Type int `json:"type"`
}

type LightState struct {
	OnOff int `json:"on_off"`
	ColorTemp int `json:"color_temp"`
	Hue int `json:"hue"`
	Saturation int `json:"saturation"`
	Brightness int `json:"brightness"`
	Mode string `json:"mode"`
}

type SysInfo struct {
	ActiveMode string `json:"active_mode"`
	Alias string `json:"alias"`
	ChildNum int `json:"child_num"`
	DeviceName string `json:"dev_name"`
	DeviceID string `json:"deviceId"`
	ErrorCode int `json:"err_code"`
	Feature string `json:"feature"`
	HardwareID string `json:"hwId"`
	HardwareVersion string `json:"hw_ver"`
	ID string `json:"id"`
	IconHash string `json:"icon_hash"`
	IsColor int `json:"is_color"`
	IsDimmable int `json:"is_dimmable"`
	IsFactory bool `json:"is_factory"`
	IsVariableColorTemp int `json:"is_variable_color_temp"`
	LEDOff int `json:"led_off"`
	Latitude int `json:"latitude_i"`
	Length int `json:"length"`
	LightState *LightState `json:"light_state"`
	Longitude int `json:"longitude_i"`
	MACAddr string `json:"mac"`
	MicType string `json:"mic_type"`
	Model string `json:"model"`
	NextAction *Action `json:"next_action"`
	ObdSrc string `json:"obd_src"`
	OEMID string `json:"oemId"`
	OnTime int `json:"on_time"`
	PreferredState []*LightState `json:"preferred_state"`
	RelayState int `json:"relay_state"`
	RSSI int `json:"rssi"`
	State int `json:"state"`
	Status string `json:"status"`
	SoftwareVersion string `json:"sw_ver"`
	Updating int `json:"updating"`
	Children []*SysInfo `json:"children"`
}

type System struct {
	SysInfo *SysInfo `json:"get_sysinfo"`
}

type Query struct {
	System *System `json:"system"`
}

type LatLon struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type SmartDevice interface {
	Update() error
	GetSysInfo() *SysInfo
	GetCurrentConsumption() (float64, error)
	GetTime() (time.Time, error)
	GetTimezone() (*time.Location, error)
	Reboot() error
	SetAlias(string) error
	SetMAC(string) error
	Alias() string
	DeviceID() string
	DeviceName() string
	DeviceType() DeviceType
	Features() []string
	IsBulb() bool
	IsColor() bool
	IsDimmable() bool
	IsDimmer() bool
	IsLightStrip() bool
	IsOff() bool
	IsOn() bool
	IsPlug() bool
	IsStrip() bool
	IsStripSocket() bool
	IsVariableColorTemp() bool
	Location() *LatLon
	MAC() string
	Model() string
	OnSince() *time.Time
	RSSI() int
}

type Switch interface {
	SmartDevice
	TurnOn() error
	TurnOff() error
}

type BaseDevice struct {
	Addr string `json:"addr"`
	Info *Query `json:"info"`
	self SmartDevice
	Responses []interface{}
}

func NewDevice(addr string) (SmartDevice, error) {
	dev := &BaseDevice{Addr: addr}
	dev.self = dev
	err := dev.Update()
	if err != nil {
		return nil, err
	}
	return dev.AsConcrete(), nil
}

func (dev *BaseDevice) AsConcrete() SmartDevice {
	switch dev.DeviceType() {
	case DeviceTypeDimmer:
		xdev := &SmartDimmer{BaseDevice: dev}
		dev.self = xdev
		return xdev
	case DeviceTypeStrip:
		xdev := &SmartStrip{BaseDevice: dev}
		dev.self = xdev
		return xdev
	case DeviceTypeStripSocket:
		xdev := &SmartStripSocket{BaseDevice: dev}
		dev.self = xdev
		return xdev
	case DeviceTypePlug:
		xdev := &SmartPlug{BaseDevice: dev}
		dev.self = xdev
		return xdev
	case DeviceTypeLightStrip:
		xdev := &SmartLightStrip{BaseDevice: dev}
		dev.self = xdev
		return xdev
	case DeviceTypeBulb:
		xdev := &SmartBulb{BaseDevice: dev}
		dev.self = xdev
		return xdev
	}
	return dev
}

func (dev *BaseDevice) makeQuery(target, cmd string, arg interface{}, childIds ...interface{}) map[string]interface{} {
	if len(childIds) > 0 {
		return map[string]interface{}{
			"context": map[string]interface{}{
				"child_ids": childIds,
				target: map[string]interface{}{
					cmd: arg,
				},
			},
		}
	}
	return map[string]interface{}{
		target: map[string]interface{}{
			cmd: arg,
		},
	}
}

func (dev *BaseDevice) Query(res interface{}, target, cmd string, arg interface{}, childIds ...interface{}) error {
	req := dev.makeQuery(target, cmd, arg, childIds...)
	dev.Responses = append(dev.Responses, res)
	return query(dev.Addr, req, res)
}

func (dev *BaseDevice) Update() error {
	req := &Query{
		System: &System{ SysInfo: nil },
	}
	err := query(dev.Addr, req, &req)
	if err != nil {
		return err
	}
	dev.Info = req
	return nil
}

func (dev *BaseDevice) GetSysInfo() *SysInfo {
	if dev.Info == nil || dev.Info.System == nil {
		return nil
	}
	return dev.Info.System.SysInfo
}


func (dev *BaseDevice) GetCurrentConsumption() (float64, error) {
	return 0, nil
}

func (dev *BaseDevice) GetTime() (time.Time, error) {
	return time.Now(), nil
}

func (dev *BaseDevice) GetTimezone() (*time.Location, error) {
	return time.UTC, nil
}

func (dev *BaseDevice) Reboot() error {
	return nil
}

func (dev *BaseDevice) SetAlias(alias string) error {
	return nil
}

func (dev *BaseDevice) SetMAC(mac string) error {
	return nil
}

func (dev *BaseDevice) Alias() string {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return ""
	}
	return sysinfo.Alias
}

func (dev *BaseDevice) DeviceID() string {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return ""
	}
	return sysinfo.DeviceID
}

func (dev *BaseDevice) DeviceName() string {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return ""
	}
	return sysinfo.DeviceName
}

func (dev *BaseDevice) DeviceType() DeviceType {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return DeviceTypeUnknown
	}
	name := sysinfo.DeviceName
	if strings.Contains(name, "Dimmer") {
		return DeviceTypeDimmer
	}
	micType := sysinfo.MicType
	if strings.Contains(strings.ToLower(micType), "smartplug") {
		if sysinfo.Children != nil && len(sysinfo.Children) > 0 {
			return DeviceTypeStrip
		}
		return DeviceTypePlug
	}
	if strings.Contains(strings.ToLower(micType), "smartbulb") {
		if sysinfo.Length > 0 {
			return DeviceTypeLightStrip
		}
		return DeviceTypeBulb
	}
	return DeviceTypeUnknown
}

func (dev *BaseDevice) Features() []string {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return []string{}
	}
	if sysinfo.Feature == "" {
		return []string{}
	}
	return strings.Split(sysinfo.Feature, ":")
}

func (dev *BaseDevice) IsBulb() bool {
	return dev.DeviceType() == DeviceTypeBulb
}

func (dev *BaseDevice) IsColor() bool {
	return false
}

func (dev *BaseDevice) IsDimmable() bool {
	return false
}

func (dev *BaseDevice) IsDimmer() bool {
	return dev.DeviceType() == DeviceTypeDimmer
}

func (dev *BaseDevice) IsLightStrip() bool {
	return dev.DeviceType() == DeviceTypeLightStrip
}

func (dev *BaseDevice) IsOff() bool {
	return false
}

func (dev *BaseDevice) IsOn() bool {
	return false
}

func (dev *BaseDevice) IsPlug() bool {
	return dev.DeviceType() == DeviceTypePlug
}

func (dev *BaseDevice) IsStrip() bool {
	return dev.DeviceType() == DeviceTypeStrip
}

func (dev *BaseDevice) IsStripSocket() bool {
	return dev.DeviceType() == DeviceTypeStripSocket
}

func (dev *BaseDevice) IsVariableColorTemp() bool {
	return false
}

func (dev *BaseDevice) Location() *LatLon {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return nil
	}
	return &LatLon{
		Lat: float64(sysinfo.Latitude) / 10000,
		Lon: float64(sysinfo.Longitude) / 10000,
	}
}

func (dev *BaseDevice) MAC() string {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return ""
	}
	return sysinfo.MACAddr
}

func (dev *BaseDevice) Model() string {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return ""
	}
	return sysinfo.Model
}

func (dev *BaseDevice) OnSince() *time.Time {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return nil
	}
	if sysinfo.OnTime == 0 {
		return nil
	}
	t := time.Now().Add(-1 * time.Second * time.Duration(sysinfo.OnTime))
	return &t
}

func (dev *BaseDevice) RSSI() int {
	sysinfo := dev.GetSysInfo()
	if sysinfo == nil {
		return 0
	}
	return sysinfo.RSSI
}

package kasa

import (
	"encoding/json"
	"log"
	"math"
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
	ActiveMode string `json:"active_mode,omitempty"`
	Alias string `json:"alias,omitempty"`
	ChildNum int `json:"child_num,omitempty"`
	DeviceName string `json:"dev_name,omitempty"`
	DeviceID string `json:"deviceId,omitempty"`
	ErrorCode int `json:"err_code,omitempty"`
	Feature string `json:"feature,omitempty"`
	HardwareID string `json:"hwId,omitempty"`
	HardwareVersion string `json:"hw_ver,omitempty"`
	ID string `json:"id,omitempty"`
	IconHash string `json:"icon_hash,omitempty"`
	IsColor int `json:"is_color,omitempty"`
	IsDimmable int `json:"is_dimmable,omitempty"`
	IsFactory bool `json:"is_factory,omitempty"`
	IsVariableColorTemp int `json:"is_variable_color_temp,omitempty"`
	LEDOff int `json:"led_off,omitempty"`
	Latitude int `json:"latitude_i,omitempty"`
	Length int `json:"length,omitempty"`
	LightState *LightState `json:"light_state,omitempty"`
	Longitude int `json:"longitude_i,omitempty"`
	MACAddr string `json:"mac,omitempty"`
	MicType string `json:"mic_type,omitempty"`
	Model string `json:"model,omitempty"`
	NextAction *Action `json:"next_action,omitempty"`
	ObdSrc string `json:"obd_src,omitempty"`
	OEMID string `json:"oemId,omitempty"`
	OnTime int `json:"on_time,omitempty"`
	PreferredState []*LightState `json:"preferred_state,omitempty"`
	RelayState int `json:"relay_state,omitempty"`
	RSSI int `json:"rssi,omitempty"`
	State int `json:"state,omitempty"`
	Status string `json:"status,omitempty"`
	SoftwareVersion string `json:"sw_ver,omitempty"`
	Updating int `json:"updating,omitempty"`
	Children []*SysInfo `json:"children,omitempty"`
	LastUpdate time.Time `json:"last_update,omitempty"`
}

type TimeInfo struct {
	Year int `json:"year"`
	Month time.Month `json:"month"`
	Day int `json:"mday"`
	Hour int `json:"hour"`
	Minute int `json:"min"`
	Second int `json:"sec"`
	ErrorCode int `json:"err_code"`
}

type TimeZoneInfo struct {
	Timezone Timezone `json:"index"`
	ErrorCode int `json:"err_code"`
}

type TimeInfoResponse struct {
	TimeInfo *TimeInfo `json:"get_time,omitempty"`
	TimeZone *TimeZoneInfo `json:"get_timezone,omitempty"`
}

func (res *TimeInfoResponse) Date() time.Time {
	if res.TimeInfo == nil || res.TimeZone == nil {
		return time.Now()
	}
	loc := res.TimeZone.Timezone.Location()
	if loc == nil {
		return time.Now()
	}
	return time.Date(res.TimeInfo.Year, res.TimeInfo.Month, res.TimeInfo.Day, res.TimeInfo.Hour, res.TimeInfo.Minute, res.TimeInfo.Second, 0, loc)
}

type SysInfoResponse struct {
	SysInfo *SysInfo `json:"get_sysinfo,omitempty"`
}

type WifiAP struct {
	SSID string `json:"ssid"`
	KeyType int `json:"key_type"`
	Password string `json:"password,omitempty"`
}

type WifiScanInfo struct {
	Refresh int `json:"refresh,omitempty"`
	AccessPoints []*WifiAP `json:"ap_list,omitempty"`
	WPA3 int `json:"wpa3_support,omitempty"`
	ErrorCode int `json:"err_code,omitempty"`
}

type WifiScanResponse struct {
	WifiScan *WifiScanInfo `json:"get_scaninfo"`
}

type Query struct {
	System *SysInfoResponse `json:"system,omitempty"`
	TimeInfoStd *TimeInfoResponse `json:"time,omitempty"`
	TimeInfoCommon *TimeInfoResponse `json:"smartlife.iot.common.timesetting,omitempty"`
	NetStd *WifiScanResponse `json:"netif,omitempty"`
	NetCommon *WifiScanResponse `json:"smartlife.iot.common.softaponboarding,omitempty"`
}

type AliasRequest struct {
	Alias string `json:"alias"`
}

type LatLon struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type SmartDevice interface {
	IP() string
	Update() error
	GetSysInfo() *SysInfo
	GetCurrentConsumption() (float64, error)
	GetTime() (time.Time, error)
	//GetTimezone() (*time.Location, error)
	Reboot(time.Duration) error
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

	GetLightService() string
	GetTimeService() string
	Repl(string) (string, error)
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

func (dev *BaseDevice) IP() string {
	return dev.Addr
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

func (dev *BaseDevice) Repl(request string) (string, error) {
	var req interface{}
	var res interface{}
	err := json.Unmarshal([]byte(request), &req)
	if err != nil {
		return "", err
	}
	err = query(dev.Addr, req, &res)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (dev *BaseDevice) Query(res interface{}, target, cmd string, arg interface{}, childIds ...interface{}) error {
	req := dev.makeQuery(target, cmd, arg, childIds...)
	dev.Responses = append(dev.Responses, res)
	var err error
	// retry up to 3 times
	for i := 0; i < 3; i += 1 {
		if i > 0 {
			time.Sleep(time.Second)
		}
		err = query(dev.Addr, req, res)
		if err == nil {
			return nil
		}
		_, isa := err.(*netError)
		if !isa {
			return err
		}
	}
	return err
}

func (dev *BaseDevice) Update() error {
	res := &Query{}
	err := dev.Query(res, "system", "get_sysinfo", nil)
	if err != nil {
		return err
	}
	dev.Info = res
	dev.setUpdateTime()
	return nil
}

func (dev *BaseDevice) setUpdateTime() {
	sysinfo := dev.GetSysInfo()
	sysinfo.LastUpdate = time.Now().In(time.UTC)
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
	res := &Query{}
	err := dev.Query(&res, dev.GetTimeService(), "get_timezone", nil)
	if err != nil {
		return time.Now(), err
	}
	err = dev.Query(&res, dev.GetTimeService(), "get_time", nil)
	if err != nil {
		return time.Now(), err
	}
	data, _ := json.Marshal(res)
	log.Println("GetTime() =>", string(data))
	if res.TimeInfoStd != nil {
		return res.TimeInfoStd.Date(), nil
	}
	if res.TimeInfoCommon != nil {
		return res.TimeInfoCommon.Date(), nil
	}
	return time.Now(), nil
}

/*
func (dev *BaseDevice) GetTimezone() (*time.Location, error) {
	var res interface{}
	err := dev.Query(&res, "time", "get_timezone", nil)
	if err != nil {
		log.Println("error in GetTimezone():", err)
		return time.UTC, err
	}
	data, _ := json.Marshal(res)
	log.Println("GetTimezone() =>", string(data))
	return time.UTC, nil
}
*/

type SetAliasRequest struct {
	Alias string `json:"alias"`
}

func (dev *BaseDevice) SetAlias(alias string) error {
	var res interface{}
	args := &SetAliasRequest{Alias: alias}
	err := dev.Query(&res, "system", "set_dev_alias", args)
	if err != nil {
		log.Println("error in SetAlias():", err)
		return err
	}
	data, _ := json.Marshal(res)
	log.Println("SetAlias() =>", string(data))
	dev.Update()
	return nil
}

type SetMACRequest struct {
	MAC string `json:"mac"`
}

func (dev *BaseDevice) SetMAC(mac string) error {
	var res interface{}
	args := &SetMACRequest{MAC: mac}
	err := dev.Query(&res, "system", "set_mac_addr", args)
	if err != nil {
		log.Println("error in SetMAC():", err)
		return err
	}
	data, _ := json.Marshal(res)
	log.Println("SetMAC() =>", string(data))
	dev.Update()
	return nil
}

type RebootRequest struct {
	Delay int `json:"delay"`
}

func (dev *BaseDevice) Reboot(delay time.Duration) error {
	var res interface{}
	args := &RebootRequest{Delay: int(math.Ceil(delay.Seconds()))}
	err := dev.Query(&res, "system", "reboot", args)
	if err != nil {
		log.Println("error in Reboot():", err)
		return err
	}
	data, _ := json.Marshal(res)
	log.Println("Reboot() =>", string(data))
	return nil
}

func (dev *BaseDevice) WifiScan() (*WifiScanInfo, error) {
	scan := func(target string) (*Query, error) {
		res := &Query{}
		args := &WifiScanInfo{Refresh: 1}
		err := dev.Query(&res, target, "get_scaninfo", args)
		return res, err
	}
	info, err := scan("netif")
	if err != nil {
		log.Println("can't scan with netif, trying softaponboarding:", err)
		info, err = scan("smartlive.iot.common.softaponboarding")
	}
	if err != nil {
		log.Println("error in WifiScan():", err)
		return nil, err
	}
	data, _ := json.Marshal(info)
	log.Println("WifiScan() =>", string(data))
	if info.NetStd != nil {
		return info.NetStd.WifiScan, nil
	}
	if info.NetCommon != nil {
		return info.NetCommon.WifiScan, nil
	}
	return nil, nil
}

func (dev *BaseDevice) WifiJoin(ssid, password string, keytype ...int) error {
	join := func(target string, payload *WifiAP) (interface{}, error) {
		var res interface{}
		err := dev.Query(&res, target, "set_stainfo", payload)
		return res, err
	}
	payload := &WifiAP{
		SSID: ssid,
		Password: password,
	}
	if len(keytype) == 0 {
		payload.KeyType = 3
	} else {
		payload.KeyType = keytype[0]
	}
	res, err := join("netif", payload)
	if err != nil {
		log.Println("Can't join with netif, trying with softaponboarding:", err)
		res, err = join("smartlife.iot.common.softaponboarding", payload)
	}
	if err != nil {
		log.Println("error in WifiJoin():", err)
		return err
	}
	data, _ := json.Marshal(res)
	log.Println("WifiJoin() =>", string(data))
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

func (dev *BaseDevice) GetLightService() string {
	return "light"
}

func (dev *BaseDevice) GetTimeService() string {
	return "time"
}

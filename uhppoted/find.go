package uhppoted

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/uhppoted/uhppote-core/types"
)

// TODO rename Address to IpAddress and use Address for IP:Port
type DeviceSummary struct {
	DeviceType string `json:"device-type"`
	Address    net.IP `json:"ip-address"`
	Port       int    `json:"port"`
}

type GetDevicesRequest struct {
}

type GetDevicesResponse struct {
	Devices map[uint32]DeviceSummary `json:"devices"`
}

func (u *UHPPOTED) GetDevices(request GetDevicesRequest) (*GetDevicesResponse, error) {
	u.debug("get-devices", fmt.Sprintf("request  %+v", request))

	wg := sync.WaitGroup{}
	list := sync.Map{}
	devices := u.UHPPOTE.DeviceList()

	for id, _ := range devices {
		deviceID := id
		wg.Add(1)
		go func() {
			defer wg.Done()
			if device, err := u.UHPPOTE.GetDevice(deviceID); err != nil {
				u.warn("find", fmt.Errorf("get-devices: %v %v", deviceID, err))
			} else if device != nil {
				list.Store(uint32(device.SerialNumber), DeviceSummary{
					DeviceType: identify(device.SerialNumber),
					Address:    device.IpAddress,
					Port:       device.Address.Port,
				})
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if devices, err := u.UHPPOTE.GetDevices(); err != nil {
			u.warn("find", fmt.Errorf("get-devices: %v", err))
		} else {
			for _, d := range devices {
				list.Store(uint32(d.SerialNumber), DeviceSummary{
					DeviceType: identify(d.SerialNumber),
					Address:    d.IpAddress,
					Port:       d.Address.Port,
				})
			}
		}
	}()

	wg.Wait()

	response := GetDevicesResponse{
		Devices: map[uint32]DeviceSummary{},
	}

	list.Range(func(key, value interface{}) bool {
		response.Devices[key.(uint32)] = value.(DeviceSummary)
		return true
	})

	u.debug("get-devices", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type GetDeviceRequest struct {
	DeviceID DeviceID
}

type GetDeviceResponse struct {
	DeviceType string           `json:"device-type"`
	DeviceID   DeviceID         `json:"device-id"`
	IpAddress  net.IP           `json:"ip-address"`
	SubnetMask net.IP           `json:"subnet-mask"`
	Gateway    net.IP           `json:"gateway-address"`
	MacAddress types.MacAddress `json:"mac-address"`
	Version    types.Version    `json:"version"`
	Date       types.Date       `json:"date"`
	Address    net.UDPAddr      `json:"address"`
	TimeZone   *time.Location   `json:"timezone,omitempty"`
}

func (u *UHPPOTED) GetDevice(request GetDeviceRequest) (*GetDeviceResponse, error) {
	u.debug("get-device", fmt.Sprintf("request  %+v", request))

	device, err := u.UHPPOTE.GetDevice(uint32(request.DeviceID))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting device info for %v (%w)", device, err))
	}

	if device == nil {
		return nil, fmt.Errorf("%w: %v", NotFound, fmt.Errorf("No device found for device ID %v", device))
	}

	response := GetDeviceResponse{
		DeviceID:   DeviceID(device.SerialNumber),
		DeviceType: identify(device.SerialNumber),
		IpAddress:  device.IpAddress,
		SubnetMask: device.SubnetMask,
		Gateway:    device.Gateway,
		MacAddress: device.MacAddress,
		Version:    device.Version,
		Date:       device.Date,
		Address:    device.Address,
		TimeZone:   device.TimeZone,
	}

	u.debug("get-device", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func identify(deviceID types.SerialNumber) string {
	id := strconv.FormatUint(uint64(deviceID), 10)

	if strings.HasPrefix(id, "4") {
		return "UTO311-L04"
	}

	if strings.HasPrefix(id, "3") {
		return "UTO311-L03"
	}

	if strings.HasPrefix(id, "2") {
		return "UTO311-L02"
	}

	if strings.HasPrefix(id, "1") {
		return "UTO311-L01"
	}

	return "UTO311-L0x"
}

package uhppoted

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/uhppoted/uhppote-core/types"
)

// TODO rename Address to IpAddress and use Address for IP:Port
type DeviceSummary struct {
	DeviceType string `json:"device-type"`
	Address    net.IP `json:"ip-address"`
	Port       uint16 `json:"port"`
}

func (u *UHPPOTED) GetDevices(request GetDevicesRequest) (*GetDevicesResponse, error) {
	u.debug("get-devices", fmt.Sprintf("request  %+v", request))

	wg := sync.WaitGroup{}
	list := sync.Map{}
	devices := u.UHPPOTE.DeviceList()

	for id := range devices {
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
					Port:       device.Address.Port(),
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
					Port:       d.Address.Port(),
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

func (u *UHPPOTED) GetDevice(request GetDeviceRequest) (*GetDeviceResponse, error) {
	u.debug("get-device", fmt.Sprintf("request  %+v", request))

	device, err := u.UHPPOTE.GetDevice(uint32(request.DeviceID))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error getting device info for %v (%w)", device, err))
	}

	if device == nil {
		return nil, fmt.Errorf("%w: %v", ErrNotFound, fmt.Errorf("no device found for device ID %v", device))
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
		Address:    device.Address.Addr(),
		TimeZone:   device.TimeZone,
	}

	u.debug("get-device", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) SetEventListener(deviceID uint32, address types.ListenAddr) (bool, error) {
	u.debug("set-event-listener", fmt.Sprintf("%v %v", deviceID, address))

	if addr := net.UDPAddrFromAddrPort(address.AddrPort); addr == nil {
		return false, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("set-event-listener: %v %w", deviceID, fmt.Errorf("invalid address (%v)", address)))
	} else if result, err := u.UHPPOTE.SetListener(deviceID, *addr); err != nil {
		return false, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("set-event-listener: %v %w", deviceID, err))
	} else if result == nil {
		return false, fmt.Errorf("%w: %v", ErrNotFound, fmt.Errorf("set-event-listener: %v  no response", deviceID))
	} else if !result.Succeeded {
		return false, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("set-event-listener: %v  failed", deviceID))
	}

	return true, nil
}

// Unwraps the request and dispatches the corresponding controller command to restore the
// manufacturer default configuration.
func (u *UHPPOTED) RestoreDefaultParameters(controller uint32) error {
	u.debug("restore-default-parameters", fmt.Sprintf("%v", controller))

	reset, err := u.UHPPOTE.RestoreDefaultParameters(controller)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("%v  error resetting controller to manufacturer default configuration (%w)", controller, err))
	} else if !reset {
		return fmt.Errorf("%w: %v", ErrFailed, fmt.Errorf("%v  failed to reset controller to manufacturer default configuration", controller))
	}

	u.debug("restore-default-parameters", fmt.Sprintf("reset %v", reset))

	return nil
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

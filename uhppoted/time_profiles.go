package uhppoted

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/uhppoted/uhppote-core/types"
)

func (u *UHPPOTED) GetTimeProfiles(request GetTimeProfilesRequest) (*GetTimeProfilesResponse, error) {
	u.debug("get-time-profiles", fmt.Sprintf("request  %+v", request))

	deviceID := request.DeviceID
	from := 2
	to := 254

	if request.From >= 2 && request.From <= 254 {
		from = request.From
	}

	if request.To >= 2 && request.To <= 254 {
		to = request.To
	}

	profiles := []types.TimeProfile{}

	for i := from; i <= to; i++ {
		profile, err := u.UHPPOTE.GetTimeProfile(deviceID, uint8(i))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving time profile %v from %v (%w)", i, deviceID, err))
		}

		if profile != nil {
			profiles = append(profiles, *profile)
		}
	}

	response := GetTimeProfilesResponse{
		DeviceID: DeviceID(deviceID),
		Profiles: profiles,
	}

	u.debug("get-time-profiles", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) PutTimeProfiles(request PutTimeProfilesRequest) (*PutTimeProfilesResponse, int, error) {
	u.debug("put-time-profiles", fmt.Sprintf("request  %+v", request))

	deviceID := request.DeviceID
	profiles := request.Profiles

	// check for duplicate profiles
	prewarn := []error{}

	set := map[uint8]int{}
	for i, profile := range profiles {
		if index, ok := set[profile.ID]; ok {
			if !reflect.DeepEqual(profile, profiles[index-1]) {
				return nil, http.StatusBadRequest, fmt.Errorf("Profile %v has more than one definition (records %v and %v)", profile.ID, index, i+1)
			}

			prewarn = append(prewarn, fmt.Errorf("Profile %-3v is defined twice (records %v and %v)", profile.ID, index, i+1))
		}

		set[profile.ID] = i + 1
	}

	// loop until all profiles are either set or could not be set
	warnings := prewarn[:]
	remaining := map[uint8]struct{}{}
	for _, p := range profiles {
		remaining[p.ID] = struct{}{}
	}

	for len(remaining) > 0 {
		warnings = prewarn[:]
		count := 0

		for _, profile := range profiles {
			// already loaded?
			if _, ok := remaining[profile.ID]; !ok {
				continue
			}

			// profile ok?
			if err := validateTimeProfile(profile); err != nil {
				warnings = append(warnings, fmt.Errorf("profile %-3v: %v", profile.ID, err))
				continue
			}

			// verify linked profile exists
			if linked := profile.LinkedProfileID; linked != 0 {
				if p, err := u.UHPPOTE.GetTimeProfile(deviceID, linked); err != nil {
					return nil, http.StatusInternalServerError, err
				} else if p == nil {
					warnings = append(warnings, fmt.Errorf("profile %-3v: linked time profile %v is not defined", profile.ID, linked))
					continue
				}
			}

			// check for circular references
			if err := circularReference(u, deviceID, profile); err != nil {
				warnings = append(warnings, fmt.Errorf("profile %-3v: %v", profile.ID, err))
				continue
			}

			// good to go!
			if ok, err := u.UHPPOTE.SetTimeProfile(deviceID, profile); err != nil {
				return nil, http.StatusInternalServerError, err
			} else if !ok {
				warnings = append(warnings, fmt.Errorf("%v: could not create time profile %v", deviceID, profile.ID))
			} else {
				u.debug("set-time-profiles", fmt.Sprintf("created/update time profile %v\n", profile.ID))

				delete(remaining, profile.ID)
				count++
			}
		}

		if count == 0 {
			break
		}
	}

	// ... format response
	response := PutTimeProfilesResponse{
		DeviceID: DeviceID(deviceID),
		Warnings: warnings,
	}

	u.debug("put-time-profiles", fmt.Sprintf("response %+v", response))

	return &response, http.StatusOK, nil
}

func (u *UHPPOTED) GetTimeProfile(request GetTimeProfileRequest) (*GetTimeProfileResponse, error) {
	u.debug("get-time-profile", fmt.Sprintf("request  %+v", request))

	deviceID := request.DeviceID
	profileID := request.ProfileID

	profile, err := u.UHPPOTE.GetTimeProfile(deviceID, profileID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving time profile %v from %v (%w)", profileID, deviceID, err))
	}

	if profile == nil {
		return nil, fmt.Errorf("%w: %v", NotFound, fmt.Errorf("Error retrieving time profile %v from %v", profileID, deviceID))
	}

	response := GetTimeProfileResponse{
		DeviceID:    DeviceID(deviceID),
		TimeProfile: *profile,
	}

	u.debug("get-time-profile", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) PutTimeProfile(request PutTimeProfileRequest) (*PutTimeProfileResponse, error) {
	u.debug("put-time-profile", fmt.Sprintf("request  %+v", request))

	deviceID := request.DeviceID
	profile := request.TimeProfile
	linked := profile.LinkedProfileID

	if profile.ID < 2 || profile.ID > 254 {
		return nil, fmt.Errorf("Invalid time profile ID (%v) - valid range is [1..254]", profile.ID)
	}

	if linked != 0 {
		if linked == profile.ID {
			return nil, fmt.Errorf("Link to self creates circular reference")
		}

		if p, err := u.UHPPOTE.GetTimeProfile(deviceID, linked); err != nil {
			return nil, err
		} else if p == nil {
			return nil, fmt.Errorf("Linked time profile %v is not defined", linked)
		}

		profiles := map[uint8]bool{profile.ID: true}
		links := []uint8{profile.ID}
		for l := linked; l != 0; {
			if p, err := u.UHPPOTE.GetTimeProfile(deviceID, l); err != nil {
				return nil, err
			} else if p == nil {
				return nil, fmt.Errorf("Linked time profile %v is not defined", l)
			} else {
				links = append(links, p.ID)
				if profiles[p.ID] {
					return nil, fmt.Errorf("Linking to time profile %v creates a circular reference (%v)", linked, links)
				}

				profiles[p.ID] = true
				l = p.LinkedProfileID
			}
		}
	}

	ok, err := u.UHPPOTE.SetTimeProfile(deviceID, profile)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error writing time profile %v to %v (%w)", profile.ID, deviceID, err))
	}

	if !ok {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Failed to write time profile %v to %v", profile.ID, deviceID))
	}

	response := PutTimeProfileResponse{
		DeviceID:    DeviceID(deviceID),
		TimeProfile: profile,
	}

	u.debug("put-time-profile", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) ClearTimeProfiles(request ClearTimeProfilesRequest) (*ClearTimeProfilesResponse, error) {
	u.debug("clear-time-profiles", fmt.Sprintf("request  %+v", request))

	deviceID := request.DeviceID

	cleared, err := u.UHPPOTE.ClearTimeProfiles(deviceID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error clearing time profiles from %v (%w)", deviceID, err))
	}

	response := ClearTimeProfilesResponse{
		DeviceID: DeviceID(deviceID),
		Cleared:  cleared,
	}

	u.debug("clear-time-profiles", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func validateTimeProfile(profile types.TimeProfile) error {
	if profile.From == nil {
		return fmt.Errorf("invalid 'From' date (%v)", profile.From)
	}

	if profile.To == nil {
		return fmt.Errorf("invalid 'To' date (%v)", profile.To)
	}

	if profile.To.Before(*profile.From) {
		return fmt.Errorf("'To' date (%v) is before 'From' date (%v)", profile.To, profile.From)
	}

	for _, i := range []uint8{1, 2, 3} {
		segment := profile.Segments[i]

		if segment.End.Before(segment.Start) {
			return fmt.Errorf("segment %v 'End' (%v) is before 'Start' (%v)", i, segment.End, segment.Start)
		}
	}

	return nil
}

func circularReference(u *UHPPOTED, deviceID uint32, profile types.TimeProfile) error {
	if linked := profile.LinkedProfileID; linked != 0 {
		profiles := map[uint8]bool{profile.ID: true}
		chain := []uint8{profile.ID}

		for l := linked; l != 0; {
			if p, err := u.UHPPOTE.GetTimeProfile(deviceID, l); err != nil {
				return err
			} else if p == nil {
				return fmt.Errorf("linked time profile %v is not defined", l)
			} else {
				chain = append(chain, p.ID)
				if profiles[p.ID] {
					return fmt.Errorf("linking to time profile %v creates a circular reference %v", profile.LinkedProfileID, chain)
				}

				profiles[p.ID] = true
				l = p.LinkedProfileID
			}
		}
	}

	return nil
}

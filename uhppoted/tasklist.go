package uhppoted

import (
	"fmt"
	"net/http"
)

func (u *UHPPOTED) PutTaskList(request PutTaskListRequest) (*PutTaskListResponse, int, error) {
	u.debug("put-task-list", fmt.Sprintf("request  %+v", request))

	deviceID := request.DeviceID
	tasks := request.Tasks
	warnings := []error{}

	if ok, err := u.UHPPOTE.ClearTaskList(deviceID); err != nil {
		return nil, http.StatusInternalServerError, err
	} else if !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("%v: could not clear  task list", deviceID)
	}

	for i, task := range tasks {
		if ok, err := u.UHPPOTE.AddTask(deviceID, task); err != nil {
			warnings = append(warnings, fmt.Errorf("%v: could not add task %d to controller (%v)", deviceID, i+1, err))
		} else if !ok {
			warnings = append(warnings, fmt.Errorf("%v: could not add task %d to controller", deviceID, i+1))
		}
	}

	if ok, err := u.UHPPOTE.RefreshTaskList(deviceID); err != nil {
		return nil, http.StatusInternalServerError, err
	} else if !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("%v: could not refresh  task list on controller", deviceID)
	}

	// ... format response
	response := PutTaskListResponse{
		DeviceID: DeviceID(deviceID),
		Warnings: warnings,
	}

	u.debug("put-task-list", fmt.Sprintf("response %+v", response))

	return &response, http.StatusOK, nil
}

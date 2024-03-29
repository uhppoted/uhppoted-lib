package uhppoted

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
)

func (u *UHPPOTED) GetCardRecords(request GetCardRecordsRequest) (*GetCardRecordsResponse, error) {
	u.debug("get-card-records", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)

	N, err := u.UHPPOTE.GetCards(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error retrieving number of cards from %v (%w)", device, err))
	}

	response := GetCardRecordsResponse{
		DeviceID: DeviceID(device),
		Cards:    N,
	}

	u.debug("get-card-records", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) GetCards(request GetCardsRequest) (*GetCardsResponse, error) {
	u.debug("get-cards", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)

	N, err := u.UHPPOTE.GetCards(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error retrieving cards from %v (%w)", device, err))
	}

	cards := make([]uint32, 0)

	var index uint32 = 1
	for count := uint32(0); count < N; {
		record, err := u.UHPPOTE.GetCardByIndex(device, index)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error retrieving cards from %v (%w)", device, err))
		}

		if record != nil {
			cards = append(cards, record.CardNumber)
			count++
		}

		index++
	}

	response := GetCardsResponse{
		DeviceID: DeviceID(device),
		Cards:    cards,
	}

	u.debug("get-cards", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) DeleteCards(request DeleteCardsRequest) (*DeleteCardsResponse, error) {
	u.debug("delete-cards", fmt.Sprintf("request  %+v", request))

	deviceID := uint32(request.DeviceID)

	deleted, err := u.UHPPOTE.DeleteCards(deviceID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error deleting cards from %v (%w)", deviceID, err))
	}

	response := DeleteCardsResponse{
		DeviceID: DeviceID(deviceID),
		Deleted:  deleted,
	}

	u.debug("delete-cards", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) GetCard(request GetCardRequest) (*GetCardResponse, error) {
	u.debug("get-card", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	cardID := request.CardNumber

	card, err := u.UHPPOTE.GetCardByID(device, cardID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error retrieving card %v from %v (%w)", cardID, device, err))
	}

	if card == nil {
		return nil, fmt.Errorf("%w: %v", ErrNotFound, fmt.Errorf("error retrieving card %v from %v", request.CardNumber, device))
	}

	response := GetCardResponse{
		DeviceID: DeviceID(device),
		Card:     *card,
	}

	u.debug("get-card", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func (u *UHPPOTED) PutCard(deviceID uint32, card types.Card) (bool, error) {
	u.debug("put-card", fmt.Sprintf("%v card:%v", deviceID, card))

	if ok, err := u.UHPPOTE.PutCard(deviceID, card); err != nil {
		return false, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error writing card %v to %v (%w)", card.CardNumber, deviceID, err))
	} else if !ok {
		return false, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("failed to write card %v to %v (%w)", card.CardNumber, deviceID, err))
	} else {
		u.debug("put-card", fmt.Sprintf("response %+v", ok))

		return ok, nil
	}
}

func (u *UHPPOTED) DeleteCard(request DeleteCardRequest) (*DeleteCardResponse, error) {
	u.debug("delete-card", fmt.Sprintf("request  %+v", request))

	deviceID := uint32(request.DeviceID)
	cardNo := request.CardNumber

	deleted, err := u.UHPPOTE.DeleteCard(deviceID, cardNo)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServerError, fmt.Errorf("error deleting card %v from %v (%w)", cardNo, deviceID, err))
	}

	response := DeleteCardResponse{
		DeviceID:   DeviceID(deviceID),
		CardNumber: cardNo,
		Deleted:    deleted,
	}

	u.debug("delete-card", fmt.Sprintf("response %+v", response))

	return &response, nil
}

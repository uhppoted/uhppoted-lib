package uhppoted

import (
	"fmt"

	"github.com/uhppoted/uhppote-core/types"
)

type GetCardRecordsRequest struct {
	DeviceID DeviceID
}

type GetCardRecordsResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Cards    uint32   `json:"cards"`
}

func (u *UHPPOTED) GetCardRecords(request GetCardRecordsRequest) (*GetCardRecordsResponse, error) {
	u.debug("get-card-records", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)

	N, err := u.UHPPOTE.GetCards(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving number of cards from %v (%w)", device, err))
	}

	response := GetCardRecordsResponse{
		DeviceID: DeviceID(device),
		Cards:    N,
	}

	u.debug("get-card-records", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type GetCardsRequest struct {
	DeviceID DeviceID
}

type GetCardsResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Cards    []uint32 `json:"cards"`
}

func (u *UHPPOTED) GetCards(request GetCardsRequest) (*GetCardsResponse, error) {
	u.debug("get-cards", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)

	N, err := u.UHPPOTE.GetCards(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving cards from %v (%w)", device, err))
	}

	cards := make([]uint32, 0)

	var index uint32 = 1
	for count := uint32(0); count < N; {
		record, err := u.UHPPOTE.GetCardByIndex(device, index)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving cards from %v (%w)", device, err))
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

type DeleteCardsRequest struct {
	DeviceID DeviceID
}

type DeleteCardsResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Deleted  bool     `json:"deleted"`
}

func (u *UHPPOTED) DeleteCards(request DeleteCardsRequest) (*DeleteCardsResponse, error) {
	u.debug("delete-cards", fmt.Sprintf("request  %+v", request))

	deviceID := uint32(request.DeviceID)

	deleted, err := u.UHPPOTE.DeleteCards(deviceID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error deleting cards from %v (%w)", deviceID, err))
	}

	response := DeleteCardsResponse{
		DeviceID: DeviceID(deviceID),
		Deleted:  deleted,
	}

	u.debug("delete-cards", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type GetCardRequest struct {
	DeviceID   DeviceID
	CardNumber uint32
}

type GetCardResponse struct {
	DeviceID DeviceID   `json:"device-id"`
	Card     types.Card `json:"card"`
}

func (u *UHPPOTED) GetCard(request GetCardRequest) (*GetCardResponse, error) {
	u.debug("get-card", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	cardID := request.CardNumber

	card, err := u.UHPPOTE.GetCardByID(device, cardID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving card %v from %v (%w)", card.CardNumber, device, err))
	}

	if card == nil {
		return nil, fmt.Errorf("%w: %v", NotFound, fmt.Errorf("Error retrieving card %v from %v", request.CardNumber, device))
	}

	response := GetCardResponse{
		DeviceID: DeviceID(device),
		Card:     *card,
	}

	u.debug("get-card", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type PutCardRequest struct {
	DeviceID DeviceID
	Card     types.Card
}

type PutCardResponse struct {
	DeviceID DeviceID   `json:"device-id"`
	Card     types.Card `json:"card"`
}

func (u *UHPPOTED) PutCard(request PutCardRequest) (*PutCardResponse, error) {
	u.debug("put-card", fmt.Sprintf("request  %+v", request))

	deviceID := uint32(request.DeviceID)
	card := request.Card

	authorised, err := u.UHPPOTE.PutCard(deviceID, card)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error writing card %v to %v (%w)", card.CardNumber, deviceID, err))
	}

	if !authorised {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Failed to write card %v to %v (%w)", card.CardNumber, deviceID, err))
	}

	response := PutCardResponse{
		DeviceID: DeviceID(deviceID),
		Card:     card,
	}

	u.debug("put-card", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type DeleteCardRequest struct {
	DeviceID   DeviceID
	CardNumber uint32
}

type DeleteCardResponse struct {
	DeviceID   DeviceID `json:"device-id"`
	CardNumber uint32   `json:"card-number"`
	Deleted    bool     `json:"deleted"`
}

func (u *UHPPOTED) DeleteCard(request DeleteCardRequest) (*DeleteCardResponse, error) {
	u.debug("delete-card", fmt.Sprintf("request  %+v", request))

	deviceID := uint32(request.DeviceID)
	cardNo := request.CardNumber

	deleted, err := u.UHPPOTE.DeleteCard(deviceID, cardNo)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error deleting card %v from %v (%w)", cardNo, deviceID, err))
	}

	response := DeleteCardResponse{
		DeviceID:   DeviceID(deviceID),
		CardNumber: cardNo,
		Deleted:    deleted,
	}

	u.debug("delete-card", fmt.Sprintf("response %+v", response))

	return &response, nil
}

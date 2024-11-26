package network

import (
	"encoding/json"
	"errors"
)

func dataToMessage(bytes []byte) (InboundMessage, error) {
	mtype, err := readMessageType(bytes)
	if err != nil {
		return nil, err
	}

	var msg InboundMessage
	switch mtype {
	case "bid":
		msg = &BidMessage{}
	case "pass-bid":
		msg = &PassBidMessage{}
	case "choose-game":
		msg = &ChooseGameTypeMessage{}
	case "choose-discard":
		msg = &ChooseDiscardCardsMessage{}
	case "play-card":
		msg = &PlayCardMessage{}
	case "ready":
		msg = &ReadyMessage{}
	default:
		return nil, errors.New("type unrecognized")
	}

	err = json.Unmarshal(bytes, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func readMessageType(bytes []byte) (string, error) {
	var obj map[string]interface{}
	err := json.Unmarshal(bytes, &obj)
	if err != nil {
		return "", err
	}

	mtype, ok := obj["type"].(string)
	if !ok {
		return "", errors.New("Received message has no type attribute.")
	}

	return mtype, nil
}

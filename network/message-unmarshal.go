package network

import (
	"encoding/json"
	"errors"
	"ognjen/go-pref/game"
)

func dataToMessage(bytes []byte, client *Client) (InboundMessage, error) {
	mtype, err := readMessageType(bytes)
	if err != nil {
		return InboundMessage{}, err
	}

	var msg InboundMessage

	/* TODO: Ugly ass switch, but whats the alternative? */
	switch mtype {
	case "bid":
		msg.Payload = &game.BidAction{}
	case "pass-bid":
		msg.Payload = &game.PassBidAction{}
	case "choose-game":
		msg.Payload = &game.PassBidAction{}
	case "choose-discard":
		msg.Payload = &game.ChooseDiscardCardsAction{}
	case "play-card":
		msg.Payload = &game.PlayCardAction{}
	case "ready":
		msg.Payload = &game.ReadyAction{}
	default:
		return InboundMessage{}, errors.New("type unrecognized")
	}

	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		return InboundMessage{}, err
	}

	msg.Client = client
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
		return "", errors.New("received message has no type attribute")
	}

	return mtype, nil
}

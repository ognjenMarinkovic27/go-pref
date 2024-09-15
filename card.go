package main

import (
	"strconv"
)

type CardSuit int

var suitMap = map[CardSuit]string{
	Spades:   "♠",
	Diamonds: "♢",
	Hearts:   "♡",
	Clubs:    "♣",
}

var valMap = map[CardValue]string{
	Ten:   "T",
	Jack:  "J",
	Queen: "Q",
	King:  "K",
	Ace:   "A",
}

func cardToString(card Card) string {
	str := ""

	if card.value < 10 {
		str += strconv.Itoa(int(card.value))
	} else {
		str += valMap[card.value]
	}

	str += suitMap[card.suit]

	return str
}

const (
	Spades   CardSuit = 0
	Diamonds CardSuit = 1
	Clubs    CardSuit = 2
	Hearts   CardSuit = 3
)

type CardValue int

const (
	Seven CardValue = 7
	Eight CardValue = 8
	Nine  CardValue = 9
	Ten   CardValue = 10
	Jack  CardValue = 11
	Queen CardValue = 12
	King  CardValue = 13
	Ace   CardValue = 14
)

type Card struct {
	suit  CardSuit
	value CardValue
}

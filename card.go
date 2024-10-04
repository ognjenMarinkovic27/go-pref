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

var strToSuitMap = map[string]CardSuit{
	"♠": Spades,
	"♢": Diamonds,
	"♡": Hearts,
	"♣": Clubs,
	"2": Spades,
	"3": Diamonds,
	"4": Hearts,
	"5": Clubs,
}

var valMap = map[CardValue]string{
	Ten:   "T",
	Jack:  "J",
	Queen: "Q",
	King:  "K",
	Ace:   "A",
}

var strToValueMap = map[string]CardValue{
	"7": Seven,
	"8": Eight,
	"9": Nine,
	"T": Ten,
	"J": Jack,
	"Q": Queen,
	"K": King,
	"A": Ace,
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

func stringToCard(str string) (c Card) {
	c.value = strToValueMap[str[0:1]]
	c.suit = strToSuitMap[str[1:2]]
	return
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

func cardCompare(a, b Card) int {
	if a.suit-b.suit == 0 {
		return int(b.value - a.value)
	} else {
		return int(a.suit - b.suit)
	}
}

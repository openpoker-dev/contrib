package card

import (
	"strconv"
	"strings"
)

type (
	// Card 单张牌的定义
	Card struct {
		Rank Rank
		Suit Suit
	}

	// Rank 牌的类型
	Rank int

	// Suit 牌的花色
	Suit int
)

const (
	RankUnknown  Rank = iota // default value
	RankAceAsOne             // 1
	RankTwo                  // 2
	RankThree                // 3
	RankFour                 // 4
	RankFive                 // 5
	RankSix                  // 6
	RankSeven                // 7
	RankEight                // 8
	RankNine                 // 9
	RankTen                  // T
	RankJack                 // J
	RankQueen                // Q
	RankKing                 // K
	RankAce                  // A
	RankJoker                // Joker
)

const (
	SuitUnknown Suit = iota // 默认
	SuitHearts              // 红桃 ♥️
	SuitDiamond             // 方片 ♦️
	SuitSpades              // 黑桃 ♠️
	SuitClubs               // 梅花 ♣️
)

var (
	rankMapping = map[string]Rank{
		"2": RankTwo,
		"3": RankThree,
		"4": RankFour,
		"5": RankFive,
		"6": RankSix,
		"7": RankSeven,
		"8": RankEight,
		"9": RankNine,
		"T": RankTen,
		"J": RankJack,
		"Q": RankQueen,
		"K": RankKing,
		"A": RankAce,
	}

	suitMapping = map[string]Suit{
		"s": SuitSpades,
		"h": SuitHearts,
		"d": SuitDiamond,
		"c": SuitClubs,
	}
)

func NewCard(s string) Card {
	if len(s) != 2 {
		panic("bad card: " + s)
	}
	var card Card
	cr, exists := rankMapping[strings.ToUpper(s[:1])]
	if !exists {
		panic("unknown rank: " + s[:1])
	}
	card.Rank = cr

	suit, exists := suitMapping[strings.ToLower(s[1:])]
	if !exists {
		panic("unknown suit: " + s[1:])
	}
	card.Suit = suit
	return card
}

func (card Card) String() string {
	var output string
	switch card.Rank {
	case RankTen:
		output = "T"
	case RankJack:
		output = "J"
	case RankQueen:
		output = "Q"
	case RankKing:
		output = "K"
	case RankAce:
		output = "A"
	default:
		output = strconv.FormatInt(int64(card.Rank), 10)
	}

	switch card.Suit {
	case SuitClubs:
		output += "♣️"
	case SuitDiamond:
		output += "♦️"
	case SuitHearts:
		output += "♥️"
	case SuitSpades:
		output += "♠️"
	}
	return output
}

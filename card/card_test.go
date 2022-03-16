package card

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardOutput(t *testing.T) {
	ace := Card{Rank: RankAce, Suit: SuitClubs}
	assert.Equal(t, ace.String(), "A♣️")

	two := Card{Rank: RankTwo, Suit: SuitHearts}
	assert.Equal(t, two.String(), "2♥️")

	card := NewCard("As")
	assert.Equal(t, card.String(), "A♠️")
}

package card

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeckDeal(t *testing.T) {
	deck := NewFiftyTwoCardsDeck()
	deck.Shuffle()
	assert.Equal(t, 52, deck.Length())

	for i := 0; i <= 51; i++ {
		_, ok := deck.Deal()
		assert.True(t, ok)
	}

	_, ok := deck.Deal()
	assert.False(t, ok)
	assert.Equal(t, 0, deck.Length())
}

func TestDeckBurn(t *testing.T) {
	deck := NewFiftyTwoCardsDeck()
	deck.Shuffle()

	for i := 0; i <= 51; i++ {
		ok := deck.Burn()
		assert.True(t, ok)
	}

	ok := deck.Burn()
	assert.False(t, ok)
	assert.Equal(t, 0, deck.Length())
}

func TestDeckShuffle(t *testing.T) {
	deck := NewFiftyTwoCardsDeck()
	deck.Shuffle()

	assert.NotEqualValues(t, deck.(*fiftyTwoCardsDeck).cards, standard52CardsDeck)
}

func TestDeckCut(t *testing.T) {
	deck := NewFiftyTwoCardsDeck()
	deck.Shuffle()
	fmt.Println(deck.(*fiftyTwoCardsDeck).cards)
	deck.Cut()
	fmt.Println(deck.(*fiftyTwoCardsDeck).cards)
	assert.Equal(t, deck.Length(), 52)
}

func TestRange(t *testing.T) {
	deck := NewFiftyTwoCardsDeck()
	d2, ok := deck.(RangeableDeck)
	assert.True(t, ok)

	var looped int
	d2.Range(func(c Card) bool {
		looped++
		return true
	})

	assert.Equal(t, looped, 52)
}

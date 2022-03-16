package evaluator

import (
	"fmt"
	"testing"

	"github.com/openpoker-dev/contrib/card"
)

func TestOuts(t *testing.T) {
	deck := card.NewFiftyTwoCardsDeck()
	deck.Shuffle()
	cal := NewOutsCalculator(newDefaultEvaluatorManager())

	p1 := []card.Card{}
	for i := 0; i < 2; i++ {
		card, _ := deck.Deal()
		p1 = append(p1, card)
	}
	fmt.Printf("player-1: %s, %s\n", p1[0], p1[1])

	p2 := []card.Card{}
	for i := 0; i < 2; i++ {
		card, _ := deck.Deal()
		p2 = append(p2, card)
	}
	fmt.Printf("player-2: %s, %s\n", p2[0], p2[1])

	community := []card.Card{}
	for i := 0; i < 3; i++ {
		card, _ := deck.Deal()
		community = append(community, card)
	}
	fmt.Printf("community: %s, %s, %s\n", community[0], community[1], community[2])

	outs, rate := cal.Calculate(p1, p2, community, deck.(card.RangeableDeck))
	fmt.Println(outs)
	fmt.Println(rate)
}

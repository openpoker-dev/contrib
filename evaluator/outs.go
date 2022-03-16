package evaluator

import "github.com/openpoker-dev/contrib/card"

type (
	OutsCalculator interface {
		Calculate(a, b []card.Card, community []card.Card, deck card.RangeableDeck) ([]card.Card, float64)
	}

	calculator struct {
		em EvaluatorManager
	}
)

func NewOutsCalculator(em EvaluatorManager) OutsCalculator {
	return &calculator{em: em}
}

func (cal *calculator) Calculate(a, b, community []card.Card, deck card.RangeableDeck) ([]card.Card, float64) {
	cardsA := a[:]
	cardsA = append(cardsA, community...)
	originalHandA := cal.em.Evaluate(cardsA...)

	cardsB := b[:]
	cardsB = append(cardsB, community...)
	originalHandB := cal.em.Evaluate(cardsB...)

	if originalHandA.Compare(originalHandB) >= ResultIdentical {
		return []card.Card{}, 1.0
	}

	outs := make(map[string]card.Card)
	deck.Range(func(c card.Card) bool {
		ca := make([]card.Card, 0, len(cardsA)+1)
		ca = append(ca, cardsA...)
		ca = append(ca, c)

		cb := make([]card.Card, 0, len(cardsB)+1)
		cb = append(cb, cardsB...)
		cb = append(cb, c)

		ha := cal.em.Evaluate(ca...)
		hb := cal.em.Evaluate(cb...)
		if ha.Compare(hb) >= ResultIdentical {
			outs[c.String()] = c
		}
		return true
	})

	o1 := make([]card.Card, 0, len(outs))
	for _, v := range outs {
		o1 = append(o1, v)
	}

	rate := float64(len(outs)) / float64(deck.Length())
	rounds := 5 - len(community)
	return o1, rate * float64(rounds)
}

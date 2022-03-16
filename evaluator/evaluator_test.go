package evaluator

import (
	"fmt"
	"testing"

	"github.com/openpoker-dev/contrib/card"
	"github.com/stretchr/testify/assert"
)

func TestStraightFlushEvaluate(t *testing.T) {
	evaluator := &straightFlushEvaluator{isRoyalFlushEvaluator: false}
	bestHand, ok := evaluator.Evaluate(
		card.NewCard("Ts"),
		card.NewCard("9s"),
		card.NewCard("7s"),
		card.NewCard("6s"),
		card.NewCard("8s"),
		card.NewCard("As"),
	)
	assert.True(t, ok)
	assert.Equal(t, bestHand.Rank, RankStraightFlush)
	assert.EqualValues(t, bestHand.Cards, []card.Card{card.NewCard("Ts"), card.NewCard("9s"), card.NewCard("8s"), card.NewCard("7s"), card.NewCard("6s")})

	bestHand, ok = evaluator.Evaluate(
		card.NewCard("As"),
		card.NewCard("2s"),
		card.NewCard("5s"),
		card.NewCard("3s"),
		card.NewCard("7s"),
		card.NewCard("4s"),
	)
	assert.True(t, ok)
	assert.EqualValues(t, bestHand.Cards, []card.Card{card.NewCard("5s"), card.NewCard("4s"), card.NewCard("3s"), card.NewCard("2s"), card.NewCard("As")})

	_, ok = evaluator.Evaluate(
		card.NewCard("As"),
		card.NewCard("2s"),
		card.NewCard("5s"),
		card.NewCard("3c"),
		card.NewCard("7s"),
		card.NewCard("4s"),
	)
	assert.False(t, ok)
}

func TestRoyalFlush(t *testing.T) {
	evaluator := &straightFlushEvaluator{isRoyalFlushEvaluator: true}
	_, ok := evaluator.Evaluate(
		card.NewCard("Ts"),
		card.NewCard("9s"),
		card.NewCard("7s"),
		card.NewCard("6s"),
		card.NewCard("8s"),
		card.NewCard("As"),
	)
	assert.False(t, ok)

	_, ok = evaluator.Evaluate(
		card.NewCard("As"),
		card.NewCard("2s"),
		card.NewCard("5s"),
		card.NewCard("3s"),
		card.NewCard("7s"),
		card.NewCard("4s"),
	)
	assert.False(t, ok)

	best, ok := evaluator.Evaluate(
		card.NewCard("Ts"),
		card.NewCard("Js"),
		card.NewCard("Ks"),
		card.NewCard("9s"),
		card.NewCard("Qs"),
		card.NewCard("As"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Rank, RankRoyalFlush)
	assert.EqualValues(t, best.Cards, []card.Card{card.NewCard("As"), card.NewCard("Ks"), card.NewCard("Qs"), card.NewCard("Js"), card.NewCard("Ts")})
}

func TestFourOfAKind(t *testing.T) {
	evaluator := fourOfAKindEvaluator{}
	best, ok := evaluator.Evaluate(
		card.NewCard("Ts"),
		card.NewCard("Tc"),
		card.NewCard("Td"),
		card.NewCard("9s"),
		card.NewCard("Qs"),
		card.NewCard("As"),
		card.NewCard("3s"),
	)
	assert.False(t, ok)

	best, ok = evaluator.Evaluate(
		card.NewCard("Ts"),
		card.NewCard("Tc"),
		card.NewCard("Td"),
		card.NewCard("Th"),
		card.NewCard("Qs"),
		card.NewCard("As"),
		card.NewCard("3s"),
	)
	assert.True(t, ok)
	assert.EqualValues(t, RankFourOfAKind, best.Rank)
	for _, c := range best.Cards[:4] {
		assert.Equal(t, c.Rank, card.RankTen)
	}
	assert.EqualValues(t, best.Cards[4], card.NewCard("As"))

	best, ok = evaluator.Evaluate(
		card.NewCard("Ts"),
		card.NewCard("Tc"),
		card.NewCard("Td"),
		card.NewCard("Th"),
	)
	assert.True(t, ok)
	assert.EqualValues(t, RankFourOfAKind, best.Rank)
	for _, c := range best.Cards[:4] {
		assert.Equal(t, c.Rank, card.RankTen)
	}
}

func TestFullHouse(t *testing.T) {
	evaluator := fullHouseEvaluator{}
	best, ok := evaluator.Evaluate(
		card.NewCard("Ts"),
		card.NewCard("Tc"),
		card.NewCard("Td"),
		card.NewCard("7h"),
		card.NewCard("7d"),
		card.NewCard("7c"),
		card.NewCard("Ac"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Rank, RankFullHouse)
	assert.Equal(t, best.Cards[0].Rank, card.RankTen)
	assert.Equal(t, best.Cards[4].Rank, card.RankSeven)

	best, ok = evaluator.Evaluate(
		card.NewCard("7h"),
		card.NewCard("7d"),
		card.NewCard("7c"),
		card.NewCard("Ts"),
		card.NewCard("Tc"),
		card.NewCard("Td"),
		card.NewCard("Ac"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Rank, RankFullHouse)
	assert.Equal(t, best.Cards[0].Rank, card.RankTen)
	assert.Equal(t, best.Cards[4].Rank, card.RankSeven)

	_, ok = evaluator.Evaluate(
		card.NewCard("7h"),
		card.NewCard("2d"),
		card.NewCard("9c"),
		card.NewCard("Ts"),
		card.NewCard("Tc"),
		card.NewCard("Td"),
		card.NewCard("Ac"),
	)
	assert.False(t, ok)

	best, ok = evaluator.Evaluate(
		card.NewCard("7h"),
		card.NewCard("9d"),
		card.NewCard("9c"),
		card.NewCard("Ts"),
		card.NewCard("Tc"),
		card.NewCard("Td"),
		card.NewCard("Ac"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Cards[0].Rank, card.RankTen)
	assert.Equal(t, best.Cards[4].Rank, card.RankNine)

	best, ok = evaluator.Evaluate(
		card.NewCard("7h"),
		card.NewCard("9d"),
		card.NewCard("9c"),
		card.NewCard("Ts"),
		card.NewCard("Tc"),
		card.NewCard("Td"),
		card.NewCard("7c"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Cards[0].Rank, card.RankTen)
	assert.Equal(t, best.Cards[4].Rank, card.RankNine)
}

func TestFlush(t *testing.T) {
	evaluator := flushEvaluator{}
	best, ok := evaluator.Evaluate(
		card.NewCard("7h"),
		card.NewCard("9h"),
		card.NewCard("Kh"),
		card.NewCard("Th"),
		card.NewCard("Tc"),
		card.NewCard("Qh"),
		card.NewCard("7c"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Rank, RankFlush)
	assert.EqualValues(t, best.Cards, []card.Card{card.NewCard("Kh"), card.NewCard("Qh"), card.NewCard("th"), card.NewCard("9h"), card.NewCard("7h")})

	_, ok = evaluator.Evaluate(
		card.NewCard("7h"),
		card.NewCard("9h"),
		card.NewCard("Kh"),
		card.NewCard("Tc"),
		card.NewCard("Tc"),
		card.NewCard("Qh"),
		card.NewCard("7c"),
	)
	assert.False(t, ok)
}

func TestStraight(t *testing.T) {
	evaluator := straightEvaluator{}
	best, ok := evaluator.Evaluate(
		card.NewCard("7h"),
		card.NewCard("9h"),
		card.NewCard("8h"),
		card.NewCard("Tc"),
		card.NewCard("Jd"),
		card.NewCard("Qh"),
		card.NewCard("7c"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Rank, RankStraight)
	assert.EqualValues(t, best.Cards, []card.Card{card.NewCard("Qh"), card.NewCard("Jd"), card.NewCard("Tc"), card.NewCard("9h"), card.NewCard("8h")})

	_, ok = evaluator.Evaluate(
		card.NewCard("7h"),
		card.NewCard("9h"),
		card.NewCard("8h"),
		card.NewCard("2c"),
		card.NewCard("Jd"),
		card.NewCard("Qh"),
		card.NewCard("7c"),
	)
	assert.False(t, ok)

	best, ok = evaluator.Evaluate(
		card.NewCard("3h"),
		card.NewCard("Ah"),
		card.NewCard("4h"),
		card.NewCard("2c"),
		card.NewCard("Jd"),
		card.NewCard("Qh"),
		card.NewCard("5c"),
	)
	assert.True(t, ok)
	assert.EqualValues(t, best.Cards, []card.Card{card.NewCard("5c"), card.NewCard("4h"), card.NewCard("3h"), card.NewCard("2c"), card.NewCard("Ah")})

	best, ok = evaluator.Evaluate(
		card.NewCard("Jh"),
		card.NewCard("9h"),
		card.NewCard("7h"),
		card.NewCard("Tc"),
		card.NewCard("8d"),
		card.NewCard("Th"),
		card.NewCard("5c"),
	)
	assert.True(t, ok)
	assert.EqualValues(t, best.Cards[0].Rank, card.RankJack)
	assert.EqualValues(t, best.Cards[1].Rank, card.RankTen)
	assert.EqualValues(t, best.Cards[2].Rank, card.RankNine)
	assert.EqualValues(t, best.Cards[3].Rank, card.RankEight)
	assert.EqualValues(t, best.Cards[4].Rank, card.RankSeven)
}

func TestThreeOfAKind(t *testing.T) {
	evaluator := threeOfAKindEvaluator{}
	best, ok := evaluator.Evaluate(
		card.NewCard("Jh"),
		card.NewCard("9h"),
		card.NewCard("7h"),
		card.NewCard("Tc"),
		card.NewCard("8d"),
		card.NewCard("Th"),
		card.NewCard("5c"),
	)
	assert.False(t, ok)

	best, ok = evaluator.Evaluate(
		card.NewCard("Jh"),
		card.NewCard("7d"),
		card.NewCard("7h"),
		card.NewCard("Tc"),
		card.NewCard("8d"),
		card.NewCard("Th"),
		card.NewCard("7c"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Cards[0].Rank, card.RankSeven)
	assert.Equal(t, best.Cards[1].Rank, card.RankSeven)
	assert.Equal(t, best.Cards[2].Rank, card.RankSeven)
	assert.Equal(t, best.Cards[3].Rank, card.RankJack)
	assert.Equal(t, best.Cards[4].Rank, card.RankTen)

	best, ok = evaluator.Evaluate(
		card.NewCard("Jh"),
		card.NewCard("7d"),
		card.NewCard("7h"),
		card.NewCard("Jc"),
		card.NewCard("8d"),
		card.NewCard("Jh"),
		card.NewCard("7c"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Cards[0].Rank, card.RankJack)
	assert.Equal(t, best.Cards[1].Rank, card.RankJack)
	assert.Equal(t, best.Cards[2].Rank, card.RankJack)
}

func TestTwoParis(t *testing.T) {
	evaluator := twoPairsEvaluator{}
	best, ok := evaluator.Evaluate(
		card.NewCard("Jh"),
		card.NewCard("9d"),
		card.NewCard("7h"),
		card.NewCard("Jc"),
		card.NewCard("8d"),
		card.NewCard("Ah"),
		card.NewCard("7c"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Rank, RankTwoParis)
	assert.Equal(t, best.Cards[0].Rank, card.RankJack)
	assert.Equal(t, best.Cards[2].Rank, card.RankSeven)
	assert.Equal(t, best.Cards[4].Rank, card.RankAce)

	_, ok = evaluator.Evaluate(
		card.NewCard("Th"),
		card.NewCard("9d"),
		card.NewCard("7h"),
		card.NewCard("Jc"),
		card.NewCard("8d"),
		card.NewCard("Ah"),
		card.NewCard("7c"),
	)
	assert.False(t, ok)

	best, ok = evaluator.Evaluate(
		card.NewCard("Jh"),
		card.NewCard("7d"),
		card.NewCard("7h"),
		card.NewCard("Jc"),
		card.NewCard("8d"),
		card.NewCard("8h"),
		card.NewCard("Kc"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Rank, RankTwoParis)
	assert.Equal(t, best.Cards[0].Rank, card.RankJack)
	assert.Equal(t, best.Cards[2].Rank, card.RankEight)
	assert.Equal(t, best.Cards[4].Rank, card.RankKing)
}

func TestOnePair(t *testing.T) {
	evaluator := onePairEvaluator{}
	best, ok := evaluator.Evaluate(
		card.NewCard("Jh"),
		card.NewCard("9d"),
		card.NewCard("7h"),
		card.NewCard("Jc"),
		card.NewCard("8d"),
		card.NewCard("Ah"),
		card.NewCard("Tc"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Rank, RankOnePair)
	assert.Equal(t, best.Cards[0].Rank, card.RankJack)
	assert.Equal(t, best.Cards[2].Rank, card.RankAce)
	assert.Equal(t, best.Cards[3].Rank, card.RankTen)
	assert.Equal(t, best.Cards[4].Rank, card.RankNine)
}

func TestHighCard(t *testing.T) {
	evaluator := highCardEvaluator{}
	best, ok := evaluator.Evaluate(
		card.NewCard("Jh"),
		card.NewCard("9d"),
		card.NewCard("7h"),
		card.NewCard("Jc"),
		card.NewCard("8d"),
		card.NewCard("Ah"),
		card.NewCard("Tc"),
	)
	assert.True(t, ok)
	assert.Equal(t, best.Rank, RankHighCard)
	fmt.Println(best)
}

func BenchmarkEvaluate(b *testing.B) {
	em := newDefaultEvaluatorManager()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			d := card.NewFiftyTwoCardsDeck()
			d.Shuffle()
			var cards []card.Card
			for i := 0; i < 7; i++ {
				card, _ := d.Deal()
				cards = append(cards, card)
			}

			if best := em.Evaluate(cards...); best.Rank > RankFlush {
				fmt.Println(best)
			}
		}
	})
}

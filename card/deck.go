package card

import (
	"math/rand"
	"sync/atomic"
	"time"
)

type (
	Deck interface {
		Shuffle()           // 洗牌
		Deal() (Card, bool) // 发牌
		Burn() bool         // 烧牌
		Cut()               // 切牌，要求两堆牌数量不少于10
		Length() int        // 手牌数量
	}

	RangeableDeck interface {
		Deck
		Range(func(Card) bool) // 遍历剩余牌
	}

	fiftyTwoCardsDeck struct {
		cards    []Card
		index    int32
		shuffled int32
		cutted   int32
		rand     *rand.Rand
	}
)

var (
	// standard52CardsDeck 52张牌，没有大小王
	standard52CardsDeck = []Card{
		NewCard("As"), NewCard("Ah"), NewCard("Ad"), NewCard("Ac"),
		NewCard("Ks"), NewCard("Kh"), NewCard("Kd"), NewCard("Kc"),
		NewCard("Qs"), NewCard("Qh"), NewCard("Qd"), NewCard("Qc"),
		NewCard("Js"), NewCard("Jh"), NewCard("Jd"), NewCard("Jc"),
		NewCard("Ts"), NewCard("Th"), NewCard("Td"), NewCard("Tc"),
		NewCard("9s"), NewCard("9h"), NewCard("9d"), NewCard("9c"),
		NewCard("8s"), NewCard("8h"), NewCard("8d"), NewCard("8c"),
		NewCard("7s"), NewCard("7h"), NewCard("7d"), NewCard("7c"),
		NewCard("6s"), NewCard("6h"), NewCard("6d"), NewCard("6c"),
		NewCard("5s"), NewCard("5h"), NewCard("5d"), NewCard("5c"),
		NewCard("4s"), NewCard("4h"), NewCard("4d"), NewCard("4c"),
		NewCard("3s"), NewCard("3h"), NewCard("3d"), NewCard("3c"),
		NewCard("2s"), NewCard("2h"), NewCard("2d"), NewCard("2c"),
	}

	_ RangeableDeck = (*fiftyTwoCardsDeck)(nil)
)

func NewFiftyTwoCardsDeck() Deck {
	deck := &fiftyTwoCardsDeck{
		cards: make([]Card, len(standard52CardsDeck)),
		index: -1,
	}
	copy(deck.cards, standard52CardsDeck)
	return deck
}

func (ft *fiftyTwoCardsDeck) Shuffle() {
	if atomic.CompareAndSwapInt32(&ft.shuffled, 0, 1) {
		source := rand.NewSource(time.Now().UnixNano())
		ft.rand = rand.New(source)
		ft.rand.Shuffle(len(ft.cards), func(i, j int) {
			ft.cards[i], ft.cards[j] = ft.cards[j], ft.cards[i]
		})
	}
}

func (ft *fiftyTwoCardsDeck) Deal() (Card, bool) {
	index := int(atomic.AddInt32(&ft.index, 1))
	if index > len(ft.cards)-1 {
		return Card{}, false
	}
	return ft.cards[index], true
}

func (ft *fiftyTwoCardsDeck) Burn() bool {
	index := int(atomic.AddInt32(&ft.index, 1))
	return index <= len(ft.cards)-1
}

func (ft *fiftyTwoCardsDeck) Length() int {
	current := int(atomic.LoadInt32(&ft.index))
	if current >= len(ft.cards)-1 {
		return 0
	}
	return len(ft.cards) - 1 - current
}

func (ft *fiftyTwoCardsDeck) Cut() {
	if atomic.CompareAndSwapInt32(&ft.cutted, 0, 1) {
		p := ft.rand.Int31n(52-10-10) + 10
		cards := make([]Card, 0, 52)
		cards = append(cards, ft.cards[p:]...)
		cards = append(cards, ft.cards[:p]...)
		ft.cards = cards
	}
}

func (ft *fiftyTwoCardsDeck) Range(fn func(Card) bool) {
	next := int(atomic.LoadInt32(&ft.index) + 1)
	for next < len(ft.cards) {
		if !fn(ft.cards[next]) {
			break
		}
		next++
	}
}

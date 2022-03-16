package evaluator

import (
	"errors"
	"sort"

	"github.com/openpoker-dev/contrib/card"
)

type (
	// Evaluator 手牌评估器
	Evaluator interface {
		MinimalCardCounts() int
		Rank() HandRank
		Evaluate(...card.Card) (PokerHand, bool)
	}

	EvaluatorManager interface {
		Register(Evaluator) error
		Find(HandRank) Evaluator
		Evaluate(...card.Card) PokerHand
	}

	// HandRank 手牌等级
	HandRank int

	PokerHand struct {
		Rank  HandRank
		Cards []card.Card
	}

	CompareResult int

	// straightFlushEvaluator 同花顺（包括皇家同花顺）
	straightFlushEvaluator struct {
		isRoyalFlushEvaluator bool
	}

	// fourOfAKindEvaluator 炸弹的评估算法
	fourOfAKindEvaluator struct{}

	// fullHouseEvaluator 葫芦的评估算法
	fullHouseEvaluator struct{}

	// flushEvaluator 同花的评估算法
	flushEvaluator struct{}

	// straightEvaluator 顺子的评估算法
	straightEvaluator struct{}

	// threeOfAKindEvaluator 三条的评估算法
	threeOfAKindEvaluator struct{}

	// twoPairsEvaluator 两对的评估算法
	twoPairsEvaluator struct{}

	// onePairEvaluator 对子的嗯评估算法
	onePairEvaluator struct{}

	// highCardEvaluator 高牌的评估算法
	highCardEvaluator struct{}

	simpleEvaluatorManager struct {
		evaluators []Evaluator
		registered map[HandRank]Evaluator
	}
)

const (
	RankHighCard      HandRank = iota // 高牌 eg. 45679
	RankOnePair                       // 一对 eg. 2♥️2♦️ 3459
	RankTwoParis                      // 两对 eg. 2♥️2♦️4♠️4♣️8
	RankThreeOfAKind                  // 三条 eg. 2♥️2♦️2♠️47
	RankStraight                      // 顺子 eg. A2345
	RankFlush                         // 同花 eg. 2♥️ T♥️ J♥️ K♥️ A♥️
	RankFullHouse                     // 葫芦 eg. 2♥️2♦️2♠️4♣️4♥️
	RankFourOfAKind                   // 炸弹 eg. 2♥️2♦️2♠️2♣️4♥️
	RankStraightFlush                 // 同花顺 eg. 9♥️T♥️J♥️Q♥️K♥️
	RankRoyalFlush                    // 皇家同花顺 eg. T♥️J♥️Q♥️K♥️A♥️
)

const (
	ResultHigher    CompareResult = 1
	ResultLower     CompareResult = -1
	ResultIdentical CompareResult = 0
)

var (
	handRankDescriptions = map[HandRank]string{
		RankHighCard:      "High Card",
		RankOnePair:       "One Pair",
		RankTwoParis:      "Two Pair",
		RankThreeOfAKind:  "Three of A Kind",
		RankStraight:      "Straight",
		RankFlush:         "Flush",
		RankFullHouse:     "Full House",
		RankFourOfAKind:   "Four of A Kind",
		RankStraightFlush: "Straight Flush",
		RankRoyalFlush:    "Royal Flush",
	}
)

var (
	_ Evaluator = (*straightFlushEvaluator)(nil)
	_ Evaluator = fourOfAKindEvaluator{}
	_ Evaluator = fullHouseEvaluator{}
	_ Evaluator = flushEvaluator{}
	_ Evaluator = straightEvaluator{}
	_ Evaluator = threeOfAKindEvaluator{}
	_ Evaluator = twoPairsEvaluator{}
	_ Evaluator = onePairEvaluator{}
	_ Evaluator = highCardEvaluator{}

	_ EvaluatorManager = (*simpleEvaluatorManager)(nil)

	defaultEvaluatorManager EvaluatorManager
)

func (bh PokerHand) String() string {
	output := handRankDescriptions[bh.Rank] + ": "
	for _, card := range bh.Cards {
		output += card.String()
	}
	return output
}

func (bh PokerHand) Compare(another PokerHand) CompareResult {
	if bh.Rank > another.Rank {
		return ResultHigher
	}
	if bh.Rank < another.Rank {
		return ResultLower
	}

	switch bh.Rank {
	case RankHighCard, RankFlush:
		for index, card := range bh.Cards {
			if ret := compareTwoCards(card, another.Cards[index]); ret != 0 {
				return ret
			}
		}
		return ResultIdentical

	case RankOnePair:
		if ret := compareTwoCards(bh.Cards[0], another.Cards[0]); ret != 0 {
			return ret
		}
		for i := 2; i < len(bh.Cards); i++ {
			if ret := compareTwoCards(bh.Cards[i], another.Cards[i]); ret != 0 {
				return ret
			}
		}
		return ResultIdentical

	case RankTwoParis, RankFourOfAKind:
		if ret := compareTwoCards(bh.Cards[0], another.Cards[0]); ret != 0 {
			return ret
		}
		if ret := compareTwoCards(bh.Cards[2], another.Cards[2]); ret != 0 {
			return ret
		}
		if len(bh.Cards) > 4 {
			return compareTwoCards(bh.Cards[4], another.Cards[4])
		}
		return ResultIdentical

	case RankThreeOfAKind:
		if ret := compareTwoCards(bh.Cards[0], another.Cards[0]); ret != 0 {
			return ret
		}
		for i := 3; i < len(bh.Cards); i++ {
			if ret := compareTwoCards(bh.Cards[i], another.Cards[i]); ret != 0 {
				return ret
			}
		}
		return ResultIdentical

	case RankStraight, RankStraightFlush, RankRoyalFlush:
		return compareTwoCards(bh.Cards[0], another.Cards[0])

	default:
		return ResultIdentical
	}
}

func compareTwoCards(a, b card.Card) CompareResult {
	switch {
	case a.Rank > b.Rank:
		return 1
	case a.Rank < b.Rank:
		return -1
	default:
		return 0
	}
}

func newDefaultEvaluatorManager() *simpleEvaluatorManager {
	em := &simpleEvaluatorManager{
		evaluators: []Evaluator{},
		registered: make(map[HandRank]Evaluator),
	}

	em.Register(&straightFlushEvaluator{isRoyalFlushEvaluator: true})
	em.Register(&straightFlushEvaluator{isRoyalFlushEvaluator: false})
	em.Register(fourOfAKindEvaluator{})
	em.Register(fullHouseEvaluator{})
	em.Register(flushEvaluator{})
	em.Register(straightEvaluator{})
	em.Register(threeOfAKindEvaluator{})
	em.Register(twoPairsEvaluator{})
	em.Register(onePairEvaluator{})
	em.Register(highCardEvaluator{})
	return em
}

func (em *simpleEvaluatorManager) Register(ev Evaluator) error {
	if ev == nil {
		return errors.New("nil pointer")
	}

	em.registered[ev.Rank()] = ev
	evaluators := make([]Evaluator, 0, len(em.registered))
	for _, reg := range em.registered {
		evaluators = append(evaluators, reg)
	}

	sort.Slice(evaluators, func(i, j int) bool {
		return evaluators[i].Rank() > evaluators[j].Rank()
	})
	em.evaluators = evaluators
	return nil
}

func (em *simpleEvaluatorManager) Evaluate(cards ...card.Card) PokerHand {
	for _, evaluator := range em.evaluators {
		if best, ok := evaluator.Evaluate(cards...); ok {
			return best
		}
	}
	return PokerHand{}
}

func (em *simpleEvaluatorManager) Find(rank HandRank) Evaluator {
	return em.evaluators[rank]
}

func (sf *straightFlushEvaluator) MinimalCardCounts() int {
	return 5
}

func (sf *straightFlushEvaluator) Rank() HandRank {
	if sf.isRoyalFlushEvaluator {
		return RankRoyalFlush
	}
	return RankStraightFlush
}

func (sf *straightFlushEvaluator) Evaluate(cards ...card.Card) (PokerHand, bool) {
	if len(cards) < 5 {
		return PokerHand{}, false
	}

	// 是否有同花
	var suitOfFlush card.Suit
	suits := make(map[card.Suit][]card.Card)
	for _, card := range cards {
		suits[card.Suit] = append(suits[card.Suit], card)
		if len(suits[card.Suit]) >= 5 {
			suitOfFlush = card.Suit
		}
	}
	if suitOfFlush == card.SuitUnknown {
		return PokerHand{}, false
	}

	flushCards := suits[suitOfFlush]
	// 排序
	sort.Slice(flushCards, func(i, j int) bool {
		return flushCards[i].Rank > flushCards[j].Rank
	})

	// 是否是顺子
	var hasAce bool
	if flushCards[0].Rank == card.RankAce {
		flushCards = append(flushCards, card.Card{Suit: flushCards[0].Suit, Rank: card.RankAceAsOne}) // 最后插入一张牌
		hasAce = true
	}

	startAt := 0
	straightCount := 0
	for i := 0; i < len(flushCards); i++ {
		if i == startAt {
			straightCount = 1
			continue
		}
		current := flushCards[i]
		if current.Rank == flushCards[startAt].Rank+card.Rank((-1)*(i-startAt)) {
			straightCount++
			if straightCount >= 5 { // 命中同花顺
				bestFiveHand := &PokerHand{
					Rank:  RankStraightFlush,
					Cards: []card.Card{flushCards[i-4], flushCards[i-3], flushCards[i-2], flushCards[i-1], flushCards[i]},
				}
				if sf.isRoyalFlushEvaluator {
					if bestFiveHand.Cards[0].Rank == card.RankAce {
						bestFiveHand.Rank = RankRoyalFlush
						return *bestFiveHand, true
					} else {
						return PokerHand{}, false
					}
				}

				if hasAce && bestFiveHand.Cards[4].Rank == card.RankAceAsOne {
					bestFiveHand.Cards[4].Rank = card.RankAce
				}
				return *bestFiveHand, true
			}
			continue
		}
		// 顺子中断
		startAt = i
		straightCount = 1
		if len(flushCards)-startAt < 5 { // 都不够顺子的5张牌要求了
			break
		}
	}
	return PokerHand{}, false
}

func (fourOfAKindEvaluator) MinimalCardCounts() int {
	return 4
}

func (four fourOfAKindEvaluator) Rank() HandRank {
	return RankFourOfAKind
}

func (four fourOfAKindEvaluator) Evaluate(cards ...card.Card) (PokerHand, bool) {
	if len(cards) < 4 {
		return PokerHand{}, false
	}

	var hitRank card.Rank
	Ranks := make(map[card.Rank][]card.Card)
	for _, card := range cards {
		Ranks[card.Rank] = append(Ranks[card.Rank], card)
		if len(Ranks[card.Rank]) >= 4 {
			hitRank = card.Rank
		}
	}

	if hitRank == 0 {
		return PokerHand{}, false
	}
	best := &PokerHand{
		Rank: RankFourOfAKind,
		Cards: []card.Card{
			{Rank: hitRank, Suit: card.SuitClubs},
			{Rank: hitRank, Suit: card.SuitDiamond},
			{Rank: hitRank, Suit: card.SuitHearts},
			{Rank: hitRank, Suit: card.SuitSpades},
		},
	}

	var highCard card.Rank
	for cr := range Ranks {
		if cr != hitRank && cr > highCard {
			highCard = cr
		}
	}
	if highCard > 0 {
		best.Cards = append(best.Cards, Ranks[highCard][0])
	}
	return *best, true
}

func (fh fullHouseEvaluator) MinimalCardCounts() int {
	return 5
}

func (fh fullHouseEvaluator) Rank() HandRank {
	return RankFullHouse
}

func (fh fullHouseEvaluator) Evaluate(cards ...card.Card) (PokerHand, bool) {
	if len(cards) < 5 {
		return PokerHand{}, false
	}

	Ranks := make(map[card.Rank][]card.Card)
	for _, card := range cards {
		Ranks[card.Rank] = append(Ranks[card.Rank], card)
	}

	var hasSet int
	var biggerSet card.Rank
	var hasPair int
	var biggerPair card.Rank
	for cr, cs := range Ranks {
		switch len(cs) {
		case 3:
			hasSet++
			if hasSet == 1 {
				biggerSet = cr
				continue
			}

			// 如果有两个set，则能形成一个对子
			hasPair++ // 对子数量
			var downgradRank card.Rank
			if cr > biggerSet { // 本set更大，之前的set降为对子
				downgradRank = biggerSet
				biggerSet = cr
			} else {
				downgradRank = cr
			}
			if downgradRank > biggerPair {
				biggerPair = downgradRank
			}

		case 2:
			hasPair++
			if cr > biggerPair {
				biggerPair = cr
			}
		default:
		}
	}

	if !(hasSet > 0 && hasPair > 0) {
		return PokerHand{}, false
	}

	return PokerHand{
		Rank: RankFullHouse,
		Cards: []card.Card{
			Ranks[biggerSet][0],
			Ranks[biggerSet][1],
			Ranks[biggerSet][2],
			Ranks[biggerPair][0],
			Ranks[biggerPair][1],
		},
	}, true
}

func (flush flushEvaluator) MinimalCardCounts() int {
	return 5
}

func (flush flushEvaluator) Rank() HandRank {
	return RankFlush
}

func (flush flushEvaluator) Evaluate(cards ...card.Card) (PokerHand, bool) {
	if len(cards) < 5 {
		return PokerHand{}, false
	}

	suits := make(map[card.Suit][]card.Card)
	var suitOfFlush card.Suit
	for _, card := range cards {
		suits[card.Suit] = append(suits[card.Suit], card)
		if len(suits[card.Suit]) >= 5 {
			suitOfFlush = card.Suit
		}
	}
	if suitOfFlush == card.SuitUnknown {
		return PokerHand{}, false
	}

	flushCards := suits[suitOfFlush]
	sort.Slice(flushCards, func(i, j int) bool {
		return flushCards[i].Rank > flushCards[j].Rank
	})
	return PokerHand{
		Rank: RankFlush,
		Cards: []card.Card{
			flushCards[0],
			flushCards[1],
			flushCards[2],
			flushCards[3],
			flushCards[4],
		},
	}, true
}

func (straight straightEvaluator) MinimalCardCounts() int {
	return 5
}

func (straight straightEvaluator) Rank() HandRank {
	return RankStraight
}

func (straight straightEvaluator) Evaluate(cards ...card.Card) (PokerHand, bool) {
	if len(cards) < 5 {
		return PokerHand{}, false
	}

	// 排序
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank > cards[j].Rank
	})
	var hasAce bool
	if cards[0].Rank == card.RankAce {
		hasAce = true
		cards = append(cards, card.Card{Suit: cards[0].Suit, Rank: card.RankAceAsOne})
	}

	startAt := 0
	counts := 1
	previousRank := card.RankUnknown
	for i := 0; i < len(cards); i++ {
		if i == startAt {
			previousRank = cards[i].Rank
			continue
		}

		// eg. J T T 9 8 7 6
		expectedRank := previousRank - 1
		if cards[i].Rank == expectedRank {
			previousRank = expectedRank
			counts++
			if counts >= 5 {
				best := &PokerHand{Rank: RankStraight, Cards: []card.Card{cards[startAt]}}
				founded := 1
				for j := startAt + 1; j <= i; j++ {
					if cards[j].Rank == best.Cards[founded-1].Rank { // 遇到对子跳过
						continue
					}

					best.Cards = append(best.Cards, cards[j])
					founded++
					if founded < 5 {
						continue
					}

					if hasAce && best.Cards[4].Rank == card.RankAceAsOne {
						best.Cards[4].Rank = card.RankAce
					}
					return *best, true
				}
			}
		} else if cards[i].Rank != previousRank { // 没有对子
			if i+5 > len(cards) { // 剩余的牌形成不了对子
				break
			}
			startAt = i // 重置
			counts = 1
			previousRank = cards[i].Rank
		}
		// 有对子，不处理
	}
	return PokerHand{}, false
}

func (set threeOfAKindEvaluator) MinimalCardCounts() int {
	return 3
}

func (set threeOfAKindEvaluator) Rank() HandRank {
	return RankThreeOfAKind
}

func (set threeOfAKindEvaluator) Evaluate(cards ...card.Card) (PokerHand, bool) {
	if len(cards) < 3 {
		return PokerHand{}, false
	}

	rankGroup := make(map[card.Rank][]card.Card)
	var biggerSet card.Rank
	for _, card := range cards {
		rankGroup[card.Rank] = append(rankGroup[card.Rank], card)
		if len(rankGroup[card.Rank]) >= 3 {
			if card.Rank > biggerSet {
				biggerSet = card.Rank
			}
		}
	}

	if biggerSet == card.RankUnknown {
		return PokerHand{}, false
	}

	best := &PokerHand{Rank: set.Rank(), Cards: []card.Card{
		rankGroup[biggerSet][0],
		rankGroup[biggerSet][1],
		rankGroup[biggerSet][2],
	}}

	kickers := make([]card.Card, 0, len(cards)-len(rankGroup[biggerSet])) // kicker 踢脚
	for rank, cs := range rankGroup {
		if rank != biggerSet {
			kickers = append(kickers, cs[0])
		}
	}

	sort.Slice(kickers, func(i, j int) bool {
		return kickers[i].Rank > kickers[j].Rank
	})

	for i := 0; i < len(kickers); i++ {
		best.Cards = append(best.Cards, kickers[i])
		if len(best.Cards) >= 5 {
			break
		}
	}
	return *best, true
}

func (two twoPairsEvaluator) MinimalCardCounts() int {
	return 4
}

func (two twoPairsEvaluator) Rank() HandRank {
	return RankTwoParis
}

func (two twoPairsEvaluator) Evaluate(cards ...card.Card) (PokerHand, bool) {
	return evaluateParis(2, two.Rank(), cards...)
}

func (one onePairEvaluator) MinimalCardCounts() int {
	return 2
}

func (one onePairEvaluator) Rank() HandRank {
	return RankOnePair
}

func (one onePairEvaluator) Evaluate(cards ...card.Card) (PokerHand, bool) {
	return evaluateParis(1, one.Rank(), cards...)
}

func evaluateParis(c int, hr HandRank, cards ...card.Card) (PokerHand, bool) {
	if len(cards) < 2*c {
		return PokerHand{}, false
	}

	rankGroup := make(map[card.Rank][]card.Card)
	pairs := make([]card.Rank, 0)
	for _, card := range cards {
		rankGroup[card.Rank] = append(rankGroup[card.Rank], card)
		if len(rankGroup[card.Rank]) == 2 {
			pairs = append(pairs, card.Rank)
		}
	}

	if len(pairs) < c {
		return PokerHand{}, false
	}
	if len(pairs) > c { // 找到最大的对子
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i] > pairs[j]
		})
	}

	outputCards := make([]card.Card, 0)
	for i := 0; i < c; i++ {
		outputCards = append(outputCards, rankGroup[pairs[i]][0])
		outputCards = append(outputCards, rankGroup[pairs[i]][1])
		delete(rankGroup, pairs[i])
	}

	kickers := make([]card.Card, 0)
	for _, cs := range rankGroup {
		kickers = append(kickers, cs[0])
	}

	sort.Slice(kickers, func(i, j int) bool {
		return kickers[i].Rank > kickers[j].Rank
	})
	outputCards = append(outputCards, kickers[:(5-len(outputCards))]...)
	return PokerHand{Rank: hr, Cards: outputCards[:]}, true
}

func (high highCardEvaluator) MinimalCardCounts() int {
	return 1
}

func (high highCardEvaluator) Rank() HandRank {
	return RankHighCard
}

func (high highCardEvaluator) Evaluate(cards ...card.Card) (PokerHand, bool) {
	best := &PokerHand{Rank: RankHighCard, Cards: []card.Card{}}

	switch len(cards) {
	case 0:
		return *best, true

	case 1:
		best.Cards[0] = cards[0]
		return *best, true

	default:
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Rank > cards[j].Rank
		})
		end := 5
		if len(cards) < end {
			end = len(cards)
		}
		best.Cards = cards[:end]
	}
	return *best, true
}

func Register(e Evaluator) error {
	return defaultEvaluatorManager.Register(e)
}

func Evaluate(cards ...card.Card) PokerHand {
	return defaultEvaluatorManager.Evaluate(cards...)
}

func Default(em EvaluatorManager) {
	if em == nil {
		panic("bad EvaluatorManager")
	}
	defaultEvaluatorManager = em
}

func init() {
	defaultEvaluatorManager = newDefaultEvaluatorManager()
}

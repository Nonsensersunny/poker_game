package model

import (
	"encoding/json"
	"github.com/Nonsensersunny/poker_game/util"
	log "github.com/sirupsen/logrus"
	"sort"
)

type Play []Poker

func (p Play) String() string {
	if playBytes, err := json.Marshal(p); err != nil {
		log.Errorf("Failed to marshal play, err:%v", err)
		return ""
	} else {
		return string(playBytes)
	}
}

func UnmarshalPlayFromString(s string) (Play, error) {
	var result Play
	err := json.Unmarshal([]byte(s), &result)
	return result, err
}

func (p Play) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Play) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func NewPlay() Play {
	return make([]Poker, 0)
}

func (p Play) Sort() {
	sort.Slice(p, func(i, j int) bool {
		return p[i].Pips < p[j].Pips
	})
	// TODO add func for sort with color
}

func (p Play) CollectPips() []Pips {
	var result []Pips
	for _, v := range p {
		result = append(result, v.Pips)
	}
	return result
}

func (p Play) ToCountMap() PipsMap {
	result := NewPipsMap()
	for _, v := range p {
		result[v.Pips] += 1
	}
	return result
}

func (p Play) GreaterThan(ap Play) bool {
	if ap == nil {
		return true
	}

	playType := p.ValidatePlayType()
	playTypeA := ap.ValidatePlayType()
	if playType.EqualsTo(playTypeA) {
		return playType.MaxPips > playTypeA.MaxPips
	}

	return playType.Type == PlayTypeBomb || playType.Type == PlayTypeRocket
}

func (p Play) ValidatePlayType() PlayType {
	pLen := len(p)

	// single
	if pLen == 1 {
		return NewPlayType(PlayTypeSingle, p[0].Pips)
	}

	if pLen == 2 {
		// rocket
		if p[0].Pips+p[1].Pips == PipsJokerRed+PipsJokerBlack /** sum of jokers */ {
			return NewPlayType(PlayTypeRocket, p[0].Pips)
		}

		// pair
		if p[0].Pips == p[1].Pips {
			return NewPlayType(PlayTypePair, p[0].Pips)
		}
	}

	var (
		quadruplesNum, triplesNum, doublesNum, singlesNum []int
		cntMap                                            = p.ToCountMap()
		keys, minK, maxK                                  = cntMap.CollectKeys()
		_, minV, maxV                                     = cntMap.CollectVals()
	)
	for k, v := range cntMap {
		if v == 4 {
			quadruplesNum = append(quadruplesNum, int(k))
		} else if v == 3 {
			triplesNum = append(triplesNum, int(k))
		} else if v == 2 {
			doublesNum = append(doublesNum, int(k))
		} else {
			singlesNum = append(singlesNum, int(k))
		}
	}

	// straight
	if maxK < PipsTwo && len(keys) > 4 && maxV == 1 && ((int(maxK-minK) + 1) == len(keys)) {
		max, _ := util.Extremes(singlesNum)
		return NewPlayType(PlayTypeStraight, Pips(max))
	}

	// triple
	if pLen > 2 && pLen < 6 && len(keys) < 3 && maxV == 3 && len(singlesNum) < 2 {
		max, _ := util.Extremes(triplesNum)
		// triple with single
		if len(singlesNum) == 1 {
			return NewPlayType(PlayTypeTripleWithSingle, Pips(max))
		}
		//triple with pair
		if len(doublesNum) == 1 {
			return NewPlayType(PlayTypeTripleWithPair, Pips(max))
		}
		return NewPlayType(PlayTypeTriple, Pips(max))
	}

	// bomb
	if pLen == 4 && len(keys) == 1 {
		return NewPlayType(PlayTypeBomb, p[0].Pips)
	}

	// bomb with pair
	if pLen == 6 && len(keys) == 2 && len(doublesNum) > 0 && len(quadruplesNum) == 1 {
		return NewPlayType(PlayTypeBombWithPair, Pips(quadruplesNum[0]))
	}

	// pairs
	if len(doublesNum) > 2 && maxK < PipsTwo && minV == 2 && maxV == 2 && util.IsContinuous(doublesNum) {
		max, _ := util.Extremes(doublesNum)
		return NewPlayType(PlayTypePairs, Pips(max))
	}

	// plane
	if len(triplesNum) > 1 && maxV > 2 && maxK < PipsTwo && util.IsContinuous(triplesNum) {
		max, _ := util.Extremes(quadruplesNum)
		// plane with single
		if 4*len(triplesNum) == pLen {
			return NewPlayType(PlayTypePlaneWithSingle, Pips(max))
		}
		// plane with double
		if len(triplesNum) == len(doublesNum) || len(triplesNum) == len(doublesNum)+2*len(quadruplesNum) {
			return NewPlayType(PlayTypePlaneWithPair, Pips(max))
		}
	}

	return NewPlayType(PlayTypeInvalid, 0)
}

func (p Play) DealWithIndex(idx []int) {
	for _, v := range idx {
		p[v].Used = true
	}
}

func (p Play) DealWithIndexMap(idxMap map[int]bool) {
	for k, v := range idxMap {
		if v {
			p[k].Used = true
		}
	}
}

func (p Play) Extract(idx []int) Play {
	var result Play
	for _, v := range idx {
		result = append(result, p[v])
	}
	return result
}

func (p Play) MinPlayHint() (Play, []int) {
	// TODO to be optimized
	idx := []int{0}
	return p.Extract(idx), idx
}

func (p Play) ExtractWithMap(selectedMap map[int]bool) Play {
	var result Play
	for k, v := range selectedMap {
		if v {
			result = append(result, p[k])
		}
	}
	return result
}

func (p Play) RevokeWithIndex(idx []int) {
	for _, v := range idx {
		p[v].Used = false
	}
}

func (p Play) AllUnused() bool {
	for _, v := range p {
		if v.Used {
			return false
		}
	}
	return true
}

func (p Play) ExtractUnused() Play {
	play := NewPlay()
	for _, v := range p {
		if !v.Used {
			play = append(play, v)
		}
	}
	return play
}

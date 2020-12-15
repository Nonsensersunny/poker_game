package model

import "math"

type (
	PlayType struct {
		Type    PlayTypeEnum  `json:"type"`
		MaxPips Pips `json:"max_pips"`
	}
	PipsMap map[Pips]int
)

func (p PlayType) EqualsTo(ap PlayType) bool {
	return p.Type == p.Type
}

func (p PlayType) Explain() string {
	return p.Type.Explain(p.MaxPips)
}

func NewPlayType(t PlayTypeEnum, maxPips Pips) PlayType {
	return PlayType{
		Type:    t,
		MaxPips: maxPips,
	}
}

type PlayTypeEnum int

const (
	PlayTypeSingle PlayTypeEnum = iota
	PlayTypePair
	PlayTypeTriple
	PlayTypeTripleWithSingle
	PlayTypeTripleWithPair
	PlayTypePlaneWithSingle
	PlayTypePlaneWithPair
	PlayTypeStraight
	PlayTypePairs
	PlayTypeBomb
	PlayTypeBombWithPair
	PlayTypeRocket
	PlayTypeInvalid
)

func (p PlayTypeEnum) Explain(pip Pips) string {
	switch p {
	case PlayTypeSingle:
		return "单走一张" + pip.String()
	case PlayTypePair:
		return "对" + pip.String()
	case PlayTypeTriple:
		return "三个" + pip.String()
	case PlayTypeTripleWithSingle:
		return "三带一"
	case 	PlayTypeTripleWithPair:
		return "三带一对"
	case 	PlayTypePlaneWithSingle, PlayTypePlaneWithPair:
		return "飞机"
	case PlayTypeStraight:
		return "顺子"
	case PlayTypePairs:
		return "连对"
	case PlayTypeBomb:
		return "炸弹"
	case PlayTypeBombWithPair:
		return "四个" + pip.String() + "带一对"
	case PlayTypeRocket:
		return "火箭"
	default:
		return "无效出牌"
	}
}

func NewPipsMap() PipsMap {
	return make(map[Pips]int)
}

func (p PipsMap) CollectKeys() (result []Pips, min, max Pips) {
	min, max = math.MaxInt32, math.MinInt32
	for k := range p {
		result = append(result, k)
		if k < min {
			min = k
		}
		if k > max {
			max = k
		}
	}
	return
}

func (p PipsMap) CollectVals() (result []int, min, max int) {
	min, max = math.MaxInt32, math.MinInt32
	for _, v := range p {
		result = append(result, v)
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}

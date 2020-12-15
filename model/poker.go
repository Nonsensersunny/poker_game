package model

import "fmt"

type (
	Color int
	Pips  int
)

func (c Color) ToUnicode() string {
	switch c {
	case ColorSun:
		return "\u2600"
	case ColorMoon:
		return "\u263D"
	case ColorSpade:
		return "\u2660"
	case ColorHeart:
		return "\u2665"
	case ColorClub:
		return "\u2663"
	case ColorDiamond:
		return "\u2666"
	default:
		return ""
	}
}

func (p Pips) String() string {
	switch p {
	case PipsJokerRed, PipsJokerBlack:
		return "Joker"
	case PipsAce:
		return "A"
	case PipsTwo:
		return "2"
	case PipsJack:
		return "J"
	case PipsQueen:
		return "Q"
	case PipsKing:
		return "K"
	default:
		return fmt.Sprint(int(p))
	}
}

const (
	ColorSun Color = iota
	ColorMoon
	ColorSpade
	ColorHeart
	ColorClub
	ColorDiamond
)

const (
	PipsThree Pips = iota + 3
	PipsFour
	PipsFive
	PipsSix
	PipsSeven
	PipsEight
	PipsNine
	PipsTen
	PipsJack
	PipsQueen
	PipsKing
	PipsAce
	PipsTwo
	PipsJokerBlack
	PipsJokerRed
)

var (
	DefaultColors = []Color{
		ColorSun,
		ColorMoon,
		ColorSpade,
		ColorHeart,
		ColorClub,
		ColorDiamond,
	}
	DefaultPips = []Pips{
		PipsJokerRed,
		PipsJokerBlack,
		PipsAce,
		PipsTwo,
		PipsThree,
		PipsFour,
		PipsFive,
		PipsSix,
		PipsSeven,
		PipsEight,
		PipsNine,
		PipsTen,
		PipsJack,
		PipsQueen,
		PipsKing,
	}
)

type Poker struct {
	Color Color  `json:"color"`
	Pips  Pips   `json:"pips"`
	Name  string `json:"name"`
	Used  bool   `json:"used"`
}

func NewPoker(color Color, pips Pips) Poker {
	return Poker{
		Color: color,
		Pips:  pips,
		Name:  pips.String(),
	}
}

func InitializePoker(pips Pips) []Poker {
	switch pips {
	case PipsJokerBlack:
		return []Poker{NewPoker(ColorMoon, pips)}
	case PipsJokerRed:
		return []Poker{NewPoker(ColorSun, pips)}
	default:
		return []Poker{
			NewPoker(ColorSpade, pips),
			NewPoker(ColorHeart, pips),
			NewPoker(ColorClub, pips),
			NewPoker(ColorDiamond, pips),
		}
	}
}

package model

import "testing"

func TestValidatePlayType(t *testing.T) {
	p := NewPlay()
	p = append(p, Poker{
		Pips: PipsTwo,
	}, Poker{
		Pips: PipsTwo,
	}, Poker{
		Pips: PipsTwo,
	}, Poker{
		Pips: PipsTwo,
	}, Poker{
		Pips: PipsThree,
	}, Poker{
		Pips: PipsThree,
	//}, Poker{
	//	Pips: PipsTwo,
	})
	t.Log(p.ValidatePlayType())
}
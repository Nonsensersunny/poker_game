package model

import (
	"encoding/json"
	"errors"
	"github.com/Nonsensersunny/poker_game/util"
	log "github.com/sirupsen/logrus"
	"time"
)

type Player struct {
	Order     int64  `json:"order"`
	Name      string `json:"name"`
	Stop      int    `json:"stop"`
	Play      Play   `json:"play"`
	IsLord    bool   `json:"is_lord"`
	IsSpeaker bool   `json:"is_speaker"`
	IsReady   bool   `json:"is_ready"`
	Win       int    `json:"win"`
	Lose      int    `json:"lose"`
	Rest      int    `json:"remain"`
}

func (p *Player) IncrWin() {
	p.Win++
}

func (p *Player) IncrLost() {
	p.Lose++
}

func (p *Player) DecrRest(d int) {
	p.Rest -= d
}

func (p Player) String() string {
	if playerBytes, err := json.Marshal(p); err != nil {
		log.Errorf("Failed to marshal player, err:%v")
		return ""
	} else {
		return string(playerBytes)
	}
}

func UnmarshalPlayerFromString(s string) (Player, error) {
	var result Player
	err := json.Unmarshal([]byte(s), &result)
	return result, err
}

func NewPlayer(name string, play ...Play) Player {
	p := Player{
		Name: name,
	}
	if len(play) > 0 {
		p.Play = play[0]
	}
	return p
}

func (p *Player) SetCountDown(cd int) *Player {
	p.Stop = cd
	return p
}

func (p Player) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Player) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p Player) IsWinner() bool {
	for _, v := range p.Play {
		if !v.Used {
			return false
		}
	}
	return true
}

func (p *Player) CountDown(callback, zero func()) {
	go func() {
		for {
			if p.Stop > 0 {
				p.Stop--
				callback()
				if p.Stop == 0 {
					zero()
				}
			} else {
				return
			}
			time.Sleep(time.Second)
		}
	}()
}

func (p *Player) UpdatePlay(play Play) {
	p.Play = play
}

type Players []Player

func (p Players) Length() int {
	return len(p)
}

func (p Players) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Players) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p Players) SetLord(name string) bool {
	exists := false
	for i := 0; i < len(p); i++ {
		if p[i].Name == name {
			p[i].IsLord = true
			exists = true
		}
	}
	return exists
}

func (p Players) WinnerCheck(name string) {
	var lordWinner bool
	for i := 0; i < len(p); i++ {
		if p[i].Name == name {
			lordWinner = p[i].IsLord
			p[i].Win++
			break
		}
	}
	for i := 0; i < len(p); i++ {
		if p[i].Name != name {
			if lordWinner {
				p[i].Lose++
			} else {
				if p[i].IsLord {
					p[i].Lose++
				} else {
					p[i].Win++
				}
			}
		}
	}
}

func (p Players) SetSpeaker(name string, callback, zero func()) bool {
	exists := false
	for i := 0; i < len(p); i++ {
		if p[i].Name == name {
			p[i].IsSpeaker = true
			p[i].Stop = util.CountDown
			p[i].CountDown(callback, zero)
			exists = true
		}
	}
	return exists
}

func (p Players) Last() (player Player, err error) {
	lastIdx := len(p) - 1
	if lastIdx < 0 {
		err = errors.New("no record yet")
		return
	}
	player = p[lastIdx]
	return
}

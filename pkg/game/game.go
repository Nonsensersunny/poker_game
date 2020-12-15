package game

import (
	"encoding/json"
	"errors"
	"github.com/Nonsensersunny/poker_game/model"
	"github.com/Nonsensersunny/poker_game/util"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type Speakers struct {
	Names []string `json:"names"`
	Lord  string   `json:"lord"`
	Idx   int      `json:"idx"`
}

func (s Speakers) CurrentSpeaker() string {
	return s.Names[s.Idx%len(s.Names)]
}

func (s *Speakers) NextSpeaker() string {
	s.Idx++
	return s.Names[s.Idx%len(s.Names)]
}

type Game struct {
	Cards     []model.Poker `json:"cards"`
	rand      *rand.Rand
	Mode      Mode                    `json:"mode"`
	Paused    bool                    `json:"paused"`
	PlayerMap map[string]*model.Player `json:"player_map"`
	Speakers  Speakers                `json:"speakers"`
	Records   model.Players           `json:"records"`
}

func (g Game) IsReady() bool {
	readyCnt := 0
	for _, v := range g.PlayerMap {
		if v.IsReady {
			readyCnt++
		}
	}
	return readyCnt == g.Mode.PlayerNum
}

func (g *Game) PlayerReady(name string) error {
	if player, ok := g.PlayerMap[name]; !ok {
		log.Errorf("Player:%s not exists")
		return errors.New("player not exists")
	} else {
		player.IsReady = true
		g.ReplacePlayer(*player)
	}
	return nil
}

func (g *Game) ReplacePlayer(p model.Player) {
	g.PlayerMap[p.Name] = &p
}

func (g *Game) UpdateRestAndJudgeWinner(name string, r int) bool {
	if p, ok := g.PlayerMap[name]; ok {
		p.Rest += r
		if p.Rest == 0 {
			return true
		}
	}
	return false
}

func (g *Game) WinnerCheck(name string) {
	lord := g.Speakers.Lord
	lordWinner := lord == name
	if player, ok := g.PlayerMap[name]; ok {
		player.Win++
		if lordWinner {
			for k, v := range g.PlayerMap {
				if k != name {
					v.Lose++
				}
			}
		} else {
			for k, v := range g.PlayerMap {
				if k != name {
					if v.Name == lord {
						v.Lose++
					} else {
						v.Win++
					}
				}
			}
		}
	}
}

func (g Game) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

func NewGame(mode Mode) *Game {
	cards := make([]model.Poker, 0)
	for _, pip := range model.DefaultPips {
		cards = append(cards, model.InitializePoker(pip)...)
	}
	return &Game{
		Cards:     cards,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
		Mode:      mode,
		PlayerMap: make(map[string]*model.Player),
	}
}

func (g *Game) Shuffle() {
	cLen := len(g.Cards)
	for i := 0; i < cLen; i++ {
		idx := g.rand.Intn(cLen - i)
		g.Cards[i], g.Cards[idx+i] = g.Cards[idx+i], g.Cards[i]
	}
}

func (g *Game) Deal() []model.Play {
	var (
		result  []model.Play
		steps   = util.SplitStep(len(g.Cards)-g.Mode.Remain, g.Mode.PlayerNum)
		lastIdx = 0
	)

	for _, v := range steps {
		result = append(result, g.Cards[lastIdx:v])
		lastIdx = v
	}
	return result
}

func (g *Game) Remain() []model.Play {
	var result []model.Play
	result = append(result, g.Cards[(len(g.Cards)-g.Mode.Remain):(len(g.Cards)-1)])
	return result
}

func (g *Game) SeatFull() bool {
	return len(g.PlayerMap) >= g.Mode.PlayerNum
}

func (g *Game) AddPlayer(p *model.Player) {
	g.PlayerMap[p.Name] = p
}

func (g *Game) RemovePlayer(name string) {
	delete(g.PlayerMap, name)
	var speakers []string
	for _, v := range g.Speakers.Names {
		if v != name {
			speakers = append(speakers, v)
		}
	}
	g.Speakers.Names = speakers
}

func (g *Game) Pause() {
	g.Paused = true
}

func (g *Game) Unpause() {
	g.Paused = false
}

func (g *Game) Record(p model.Player) {
	g.Records = append(g.Records, p)
}

func (g *Game) ConfirmPlayers() {
	for k := range g.PlayerMap {
		g.Speakers.Names = append(g.Speakers.Names, k)
	}
}

func (g *Game) RandomLord() string {
	g.Speakers.Idx = g.rand.Intn(g.Mode.PlayerNum)
	g.Speakers.Lord = g.Speakers.Names[g.Speakers.Idx]
	return g.Speakers.Lord
}

func (g *Game) UpdatePlay(name string, p model.Play) {
	if player, ok := g.PlayerMap[name]; ok {
		player.Play = p
		g.PlayerMap[name] = player
	}
}

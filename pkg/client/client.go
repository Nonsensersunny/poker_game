package client

import (
	"bufio"
	"github.com/Nonsensersunny/poker_game/model"
	"github.com/Nonsensersunny/poker_game/util"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"sync"
	"time"
)

var (
	once sync.Once
)

type Client struct {
	conn         net.Conn
	Player       *model.Player
	OtherPlayers model.Players
	GameStart    bool
	selected     map[int]bool
	cursor       int
	Records      model.Players
}

func (c *Client) WinnerCheck(name string) {
	var lordWinner bool
	if c.Player.Name == name {
		c.Player.Win++
		if c.Player.IsLord {
			lordWinner = true
		}
	} else {
		for i := 0; i < c.OtherPlayers.Length(); i++ {
			if c.OtherPlayers[i].Name == name {
				c.OtherPlayers[i].Win++
				lordWinner = c.OtherPlayers[i].IsLord
			}
		}
		if lordWinner {
			c.Player.Lose++
		} else {
			c.Player.Win++
		}
		for i := 0; i < len(c.OtherPlayers); i++ {
			if c.OtherPlayers[i].Name != name {
				if lordWinner {
					c.OtherPlayers[i].Lose++
				} else {
					c.OtherPlayers[i].Win++
				}
			}
		}
	}
}

func (c *Client) SendMessage(prefix model.Prefix, s string) {
	if _, err := c.conn.Write([]byte(prefix.AssembleMessage(s))); err != nil {
		log.Errorf("Failed to send message:%v to server, err:%v", s, err)
	}
}

func (c *Client) entry() (string, error) {
	reader := bufio.NewReader(c.conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(msg) > 0 {
		msg = msg[:len(msg)-1]
	}

	return msg, nil
}

func (c *Client) messageHandler(msg string, serverSig chan string) error {
	prefix, body := model.Extract(msg)
	//log.Errorf("%v, %v", prefix, body)
	switch prefix {
	case model.PrefixNormal:
		//log.Info(body)
	case model.PrefixPlayer:
		player, err := model.UnmarshalPlayerFromString(body)
		if err != nil {
			log.Errorf("Failed to unmarshal player:%s from server, err:%v", msg, err)
			return nil
		}
		exists := false
		for i := 0; i < len(c.OtherPlayers); i++ {
			if c.OtherPlayers[i].Name == player.Name {
				c.OtherPlayers[i].Play = player.Play
				exists = true
			}
		}
		if !exists {
			c.OtherPlayers = append(c.OtherPlayers, player)
		}
		c.Records = append(c.Records, player)
		c.Refresh()
	case model.PrefixDeal:
		play, err := model.UnmarshalPlayFromString(body)
		if err != nil {
			log.Errorf("Failed to unmarshal deal:%s from server, err:%v", msg, err)
			return err
		}
		c.Player.Play = play
		c.Refresh()
	case model.PrefixUncoverLord:
		if c.Player.Name == body {
			c.Player.IsLord = true
		} else {
			c.OtherPlayers.SetLord(body)
		}
		c.Refresh()
	case model.PrefixYourName:
		c.Player.Name = body
	case model.PrefixPlayerName:
		c.OtherPlayers = append(c.OtherPlayers, model.NewPlayer(msg))
	case model.PrefixWinner:
		c.WinnerCheck(body)
		c.Refresh()
	case model.PrefixSpeaker:
		if c.Player.Name == body {
			c.Player.IsSpeaker = true
			c.Player.Stop = util.CountDown
			c.Player.CountDown(c.Refresh, func() {
				c.Pass(serverSig)
			})
		} else {
			c.OtherPlayers.SetSpeaker(body, c.Refresh, func() {})
		}
		//c.Refresh()
	case model.PrefixGameStart:
		c.GameStart = true
	case model.PrefixInvalidPlay:
		// TODO revoke last invalid play
		c.Revoke(serverSig)
	default:
		log.Infof("Unknown message:%s", msg)
	}

	return nil
}

func (c *Client) Revoke(serverSig chan string) {
	last, err := c.Records.Last()
	if err == nil && last.Name == c.Player.Name && len(last.Play) > 0 {
		c.Player.Play = append(c.Player.Play, last.Play...)
		c.Player.Play.Sort()
	}
	c.Player.IsSpeaker = true
	c.Player.Stop = last.Stop
	c.Player.CountDown(c.Refresh, func() {
		c.Pass(serverSig)
	})
}

func (c *Client) Pass(serverSig chan string) {
	if c.Player.IsSpeaker {
		last, err := c.Records.Last()
		if last.Name == c.Player.Name || err != nil {
			min, idx := c.Player.Play.MinPlayHint()
			c.Player.Play.DealWithIndex(idx)
			serverSig <- model.PrefixPlay.AssembleMessage(min.String())
			c.Player.Play = c.Player.Play.ExtractUnused()
			c.Records = append(c.Records, model.NewPlayer(c.Player.Name, min))
		} else {
			serverSig <- model.PrefixPlay.AssembleMessage("")
		}
	} else {
		return
	}

	c.Player.Stop = 0
	c.Player.IsSpeaker = false
	c.Refresh()
}

func InitClient(addr string) {
	log.Info("Initializing client...")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Infof("Connecting %v", addr)

	client := Client{
		conn:         conn,
		OtherPlayers: make([]model.Player, 0),
		selected:     make(map[int]bool),
		Player:       &model.Player{},
	}

	serverSignal := make(chan string, 10)
	defer close(serverSignal)

	clientSignal := make(chan string, 10)
	defer close(clientSignal)

	go client.InitUI(serverSignal)

	go func() {
		for {
			select {
			case msg := <-serverSignal:
				client.SendMessage("", msg)
			}
		}
	}()

	once.Do(func() {
		client.SendMessage(model.PrefixGameStart, "")
	})

	go func() {
		for {
			if client.GameStart {
				//client.Refresh()
			}
			time.Sleep(time.Second)
		}
	}()

	for {
		if msg, err := client.entry(); err != nil {
			log.Errorf("Something wrong with server message handler, err:%v", err)
			break
		} else {
			if err := client.messageHandler(msg, serverSignal); err != nil {
				log.Errorf("Failed to handle message, err:%v", err)
			}
		}
	}
}

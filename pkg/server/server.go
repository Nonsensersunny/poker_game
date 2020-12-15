package server

import (
	"bufio"
	"fmt"
	"github.com/Nonsensersunny/poker_game/model"
	"github.com/Nonsensersunny/poker_game/pkg/game"
	"github.com/Nonsensersunny/poker_game/util"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"time"
)

type Server struct {
	connMap map[string]net.Conn
	Game    *game.Game `json:"game"`
}

func (s Server) Broadcast(prefix model.Prefix, msg string, excludes ...string) {
	for k, v := range s.connMap {
		if !util.SliceContains(excludes, k) {
			_, err := v.Write([]byte(prefix.AssembleMessage(msg)))
			if err != nil {
				log.Errorf("Failed to notify:%s, err:%v", v.RemoteAddr().String(), err)
			}
		}
	}
}

func (s Server) Send(name string, msg string) {
	if conn, ok := s.connMap[name]; !ok {
		log.Errorf("No client named:%s", name)
	} else {
		if _, err := conn.Write([]byte(msg)); err != nil {
			log.Errorf("Failed to send message to:%s, err:%v", name, err)
		}
	}
}

func InitServer(addr string) {
	log.Info("Initializing server...")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Infof("Listening %v", addr)

	log.Info("Initializing game configuration...")
	g := game.NewGame(game.ModeChinesePoker /** TODO more game mode to be added */)
	server := Server{
		Game:    g,
		connMap: make(map[string]net.Conn),
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			return
		}

		remoteAddr := conn.RemoteAddr().String()
		log.Info(fmt.Sprintf("Client %v connected.", remoteAddr))

		if server.Game.SeatFull() {
			conn.Write([]byte("抱歉，房间满员，请稍后重试\n"))
			continue
		}

		server.connMap[remoteAddr] = conn
		server.Send(remoteAddr, model.PrefixYourName.AssembleMessage(remoteAddr))
		server.Broadcast(model.PrefixPlayerName, remoteAddr, remoteAddr)
		server.Broadcast(model.PrefixNormal, fmt.Sprintf("玩家%s加入房间", remoteAddr), remoteAddr)
		server.Game.AddPlayer(&model.Player{
			Name:  remoteAddr,
			Order: time.Now().Unix(),
		})
		if server.Game.SeatFull() {
			server.Broadcast(model.PrefixGameStart, "房间人数已满，可以开始游戏")
		}

		go server.handleConnection(conn, conn.RemoteAddr().String())
	}
}

func (s *Server) handleConnection(conn net.Conn, ip string) {
	reader := bufio.NewReader(conn)
	buffer, err := reader.ReadString('\n')
	if err != nil {
		log.Errorf("Client %s exit.", ip)
		conn.Close()
		name := conn.RemoteAddr().String()
		delete(s.connMap, name)
		s.Broadcast(model.PrefixNormal, fmt.Sprintf("玩家%s断开连接", name))
		s.Game.RemovePlayer(name)
		s.Broadcast(model.PrefixNormal, "游戏已暂停")
		s.Game.Pause()
		return
	}

	var realMsg string
	if len(buffer) > 0 {
		realMsg = buffer[:len(buffer)-1]
		log.Infof("Client message: %s", realMsg)
	}
	s.handleMessage(conn, realMsg)

	//conn.Write(buffer)
	s.handleConnection(conn, ip)
}

func (s *Server) handleMessage(conn net.Conn, msg string) {
	name := conn.RemoteAddr().String()
	prefix, msg := model.Extract(msg)
	switch prefix {
	case model.PrefixGameStart:
		if err := s.Game.PlayerReady(name); err != nil {
			return
		}
		s.Broadcast(model.PrefixNormal, fmt.Sprintf("玩家%s已准备就绪", name))
		if s.Game.IsReady() {
			log.Info("Confirming players...")
			s.Broadcast(model.PrefixNormal, "确认玩家...")
			s.Game.ConfirmPlayers()

			log.Info("Shuffling...")
			s.Broadcast(model.PrefixNormal, "开始洗牌...")
			s.Game.Shuffle()

			log.Info("Dealing...")
			s.Broadcast(model.PrefixNormal, "开始发牌...")
			deals := s.Game.Deal()
			idx := 0
			for k := range s.Game.PlayerMap {
				s.Send(k, model.PrefixDeal.AssembleMessage(deals[idx].String()))
				s.Game.UpdatePlay(name, deals[idx])
				s.Game.PlayerMap[k].Rest = len(deals[idx])
				idx++
			}

			s.ResetGame()
		}
	case model.PrefixPlay:
		log.Infof("Player:%s is speaking", name)
		if name != s.Game.Speakers.CurrentSpeaker() {
			log.Errorf("Player:%s is not current speaker:%s", name, s.Game.Speakers.CurrentSpeaker())
			s.Send(name, model.PrefixInvalidPlay.AssembleMessage(""))
			return
		}
		if msg == "" {
			log.Infof("Player:%s chosen pass", name)
			s.Broadcast(model.PrefixSpeaker, s.Game.Speakers.NextSpeaker())
			return
		}
		play, err := model.UnmarshalPlayFromString(msg)
		if err != nil {
			log.Errorf("Invalid play, err:%v", err)
			s.Send(name, model.PrefixInvalidPlay.AssembleMessage(""))
			return
		}

		last, err := s.Game.Records.Last()
		if err != nil {
			log.Infof("No record found, game start now, err:%v", err)
			s.Game.Records = append(s.Game.Records, model.NewPlayer(name, play))
			goto NEXT
		} else {
			if play.GreaterThan(last.Play) || last.Name == name {
				goto NEXT
			} else {
				log.Errorf("Invalid play, needs to be revoked")
				s.Send(name, model.PrefixInvalidPlay.AssembleMessage(""))
				return
			}
		}
	NEXT:
		log.Info("Valid play, finding next speaker...")
		s.Broadcast(model.PrefixPlayer, model.NewPlayer(name, play).String(), name)
		time.Sleep(time.Second)

		if s.Game.UpdateRestAndJudgeWinner(name, -len(play)) {
			log.Infof("Game over, winner:%v", name)
			s.Broadcast(model.PrefixWinner, name)
			s.Game.WinnerCheck(name)
			time.Sleep(time.Second)

			log.Infof("Reset game")
			s.ResetGame()
			return
		}

		nextSpeaker := s.Game.Speakers.NextSpeaker()
		log.Infof("Next player:%s", nextSpeaker)
		s.Broadcast(model.PrefixSpeaker, nextSpeaker)
	case model.PrefixCancelRemain, model.PrefixPass:
		if s.Game.Speakers.CurrentSpeaker() == name {
			nextSpeaker := s.Game.Speakers.NextSpeaker()
			s.Broadcast(model.PrefixSpeaker, nextSpeaker)
		} else {
			log.Errorf("Player:%s is not current speaker:%s", name, s.Game.Speakers.CurrentSpeaker())
			s.Send(name, model.PrefixInvalidPlay.AssembleMessage(""))
			return
		}
	}
}

func (s *Server) ResetGame() {
	log.Info("Random lord selecting...")
	s.Game.RandomLord()
	s.Broadcast(model.PrefixUncoverLord, s.Game.Speakers.Lord)
	time.Sleep(time.Second)

	log.Infof("Initialize speaker:%v", s.Game.Speakers.Lord)
	s.Broadcast(model.PrefixSpeaker, s.Game.Speakers.Lord)
	time.Sleep(time.Second)

	log.Info("Sending game start command...")
	s.Broadcast(model.PrefixGameStart, "")

	log.Info("Game started")
}

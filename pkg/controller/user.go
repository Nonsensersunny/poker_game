package controller

import (
	"errors"
	"fmt"
	"github.com/Nonsensersunny/poker_game/container"
	"github.com/Nonsensersunny/poker_game/model"
	"github.com/Nonsensersunny/poker_game/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func CheckName(ctx *gin.Context) {
	name := ctx.Query("name")
	nameS := container.DefaultContainer.Redis.Get(ctx, name).Val()
	if nameS != "" {
		util.ResponseWithErr(ctx, util.ErrNameOccupied, errors.New(fmt.Sprintf("name %s occupied", name)))
		return
	}

	if err := updatePlayerExpiration(ctx, name); err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrWriteDB, err)
		return
	}

	util.ResponseWithData(ctx, "success")
}

func InitGame(ctx *gin.Context) {
	name := ctx.Query("name")
	//nameS := container.DefaultContainer.Redis.Get(ctx, name).Val()
	//if nameS != "" {
	//	util.ResponseWithErr(ctx, util.ErrNameOccupied, errors.New(fmt.Sprintf("name %s occupied", name)))
	//	return
	//}

	if err := container.DefaultContainer.Redis.Set(ctx, name, "occupied", time.Hour).Err(); err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrWriteDB, err)
		return
	}

	userGame := util.GenUserGameKey(name)
	games := container.DefaultContainer.Redis.SMembers(ctx, userGame).Val()
	if len(games) > 0 {
		util.ResponseWithErr(ctx, util.ErrGameNotEnd, errors.New(fmt.Sprintf("game %s not end", userGame)))
		return
	}

	if err := container.DefaultContainer.Redis.SAdd(ctx, util.AvailableGameSet, userGame).Err(); err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrWriteDB, err)
		return
	}

	if err := container.DefaultContainer.Redis.SAdd(ctx, userGame, model.NewPlayer(name)).Err(); err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrWriteDB, err)
		return
	}

	util.ResponseWithData(ctx, userGame)
}

func DestroyGame(ctx *gin.Context) {
	key := ctx.Query("name")
	userGame := util.GenUserGameKey(key)
	games := container.DefaultContainer.Redis.SMembers(ctx, userGame).Val()
	if len(games) < 1 {
		util.ResponseWithErr(ctx, util.ErrGameNotExist, errors.New(fmt.Sprintf("game %s not exist", userGame)))
		return
	}

	if err := container.DefaultContainer.Redis.SRem(ctx, util.AvailableGameSet, key).Err(); err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrWriteDB, err)
		return
	}

	util.ResponseWithData(ctx, "success")
}

type availableGames struct {
	Name    string        `json:"name"`
	Players model.Players `json:"players"`
}

func GetAvailableGames(ctx *gin.Context) {
	var result []availableGames
	games := container.DefaultContainer.Redis.SMembers(ctx, util.AvailableGameSet).Val()
	for _, v := range games {
		players, err := getPlayers(ctx, v)
		if err != nil {
			log.Error(err)
			util.ResponseWithErr(ctx, util.ErrReadDB, err)
			return
		}

		var ps model.Players
		for _, p := range players {
			p.Play = nil
			ps = append(ps, p)
		}
		if ps != nil {
			result = append(result, availableGames{
				Name:    v,
				Players: ps,
			})
		}
	}

	util.ResponseWithData(ctx, result)
}

func JoinGame(ctx *gin.Context) {
	name := ctx.Query("name")
	gameKey := ctx.Query("game")
	players, err := getPlayers(ctx, gameKey)
	if err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrReadDB, err)
		return
	}

	for _, v := range players {
		if name == v.Name {
			util.ResponseWithErr(ctx, util.ErrDuplicateOperation, errors.New("already in game"))
			return
		}
	}

	if len(players) > 2 {
		util.ResponseWithErr(ctx, util.ErrGameOccupied, errors.New(fmt.Sprintf("game %v seat full", gameKey)))
		return
	}

	if err := container.DefaultContainer.Redis.SAdd(ctx, gameKey, model.NewPlayer(name)).Err(); err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrWriteDB, err)
		return
	}

	util.ResponseWithData(ctx, gameKey)
}

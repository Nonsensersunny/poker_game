package controller

import (
	"errors"
	"github.com/Nonsensersunny/poker_game/container"
	"github.com/Nonsensersunny/poker_game/model"
	game2 "github.com/Nonsensersunny/poker_game/pkg/game"
	"github.com/Nonsensersunny/poker_game/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
)

func StartGame(ctx *gin.Context) {
	var mode game2.Mode
	if ctx.Query("mode") == "" {
		mode = game2.ModeChinesePoker
	}
	game := game2.NewGame(mode)
	game.Shuffle()
	heaps := game.Deal()
	remain := game.Remain()

	gameKey := ctx.Query("game")
	players, err := getPlayers(ctx, gameKey)
	if err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrReadDB, err)
		return
	}

	for i := 0; i < len(players); i++ {
		heaps[i].Sort()
		if err := updateGame(ctx, gameKey, model.NewPlayer(players[i].Name, heaps[i])); err != nil {
			log.Error(err)
			util.ResponseWithErr(ctx, util.ErrWriteDB, err)
			return
		}
	}

	remainKey := util.GenGameRemainKey(gameKey)
	if err := container.DefaultContainer.Redis.Set(ctx, remainKey, remain, time.Hour).Err(); err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrWriteDB, err)
		return
	}

	util.ResponseWithData(ctx, "success")
}

func TakeRemain(ctx *gin.Context) {
	name := ctx.Query("name")
	game := ctx.Query("game")
	remainKey := util.GenGameRemainKey(game)

	remain, err := getRemain(ctx, game)
	if err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrReadDB, err)
		return
	}

	players, err := getPlayers(ctx, game)
	if err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrReadDB, err)
		return
	}

	for _, v := range players {
		if v.Name == name {
			if err := updateGame(ctx, game, model.NewPlayer(name, append(v.Play, remain...))); err != nil {
				log.Error(err)
				util.ResponseWithErr(ctx, util.ErrWriteDB, err)
				return
			}

			// clear cache
			if err := container.DefaultContainer.Redis.Del(ctx, remainKey).Err(); err != nil {
				log.Error(err)
			}
			util.ResponseWithData(ctx, "success")
			return
		}
	}
}

func UncoverRemain(ctx *gin.Context) {
	game := ctx.Query("game")
	remain, err := getRemain(ctx, game)
	if err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrReadDB, err)
		return
	}

	util.ResponseWithData(ctx, remain)
}

func Play(ctx *gin.Context) {
	name := ctx.Query("name")
	gameKey := ctx.Query("game")
	if ctx.Query("index") == "" {
		util.ResponseWithData(ctx, "success")
		return
	}

	idx, err := util.ExtractIntsFromQuery(ctx, "index")
	if err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrRequestFormat, err)
		return
	}

	players, err := getPlayers(ctx, gameKey)
	if err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrReadDB, err)
		return
	}

	for _, v := range players {
		if name == v.Name {
			extract := v.Play.Extract(idx)
			if !extract.AllUnused() {
				util.ResponseWithErr(ctx, util.ErrInvalidPlay, errors.New("used play"))
				return
			}

			playType := extract.ValidatePlayType()
			if playType.Type == model.PlayTypeInvalid {
				util.ResponseWithErr(ctx, util.ErrInvalidPlay, errors.New("invalid play"))
				return
			}

			lastPlay, err := getLastPlay(ctx, gameKey)
			if err != nil && err != redis.Nil {
				log.Error(err)
				util.ResponseWithErr(ctx, util.ErrReadDB, err)
				return
			}
			if lastPlay.Name == name {
				goto DEAL
			}

			if !extract.GreaterThan(lastPlay.Play) {
				util.ResponseWithErr(ctx, util.ErrInvalidPlay, errors.New("invalid play"))
				return
			}
		DEAL:
			v.Play.DealWithIndex(idx)
			if err := updateGame(ctx, gameKey, v); err != nil {
				log.Error(err)
				util.ResponseWithErr(ctx, util.ErrWriteDB, err)
				return
			}

			if err := recordPlay(ctx, gameKey, model.NewPlayer(name, extract)); err != nil {
				log.Error(err)
				util.ResponseWithErr(ctx, util.ErrWriteDB, err)
				return
			}

			util.ResponseWithData(ctx, "success")
			return
		}
	}
}

func GetPlayer(ctx *gin.Context) {
	name := ctx.Query("name")
	gameKey := ctx.Query("game")
	players, err := getPlayers(ctx, gameKey)
	if err != nil {
		log.Error(err)
		util.ResponseWithErr(ctx, util.ErrReadDB, err)
		return
	}

	for _, v := range players {
		if v.Name == name {
			util.ResponseWithData(ctx, v)
			return
		}
	}

	util.ResponseWithErr(ctx, util.ErrGameDataMissing, errors.New("game data not found"))
}

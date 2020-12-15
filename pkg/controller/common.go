package controller

import (
	"context"
	"encoding/json"
	"github.com/Nonsensersunny/poker_game/container"
	"github.com/Nonsensersunny/poker_game/model"
	"github.com/Nonsensersunny/poker_game/util"
	"time"
)

func updateGame(ctx context.Context, game string, player model.Player) error {
	gameInfo := container.DefaultContainer.Redis.SMembers(ctx, game)
	if gameInfo.Err() != nil {
		return gameInfo.Err()
	}

	for _, v := range gameInfo.Val() {
		var oriPlayer model.Player
		if err := json.Unmarshal([]byte(v), &oriPlayer); err != nil {
			return err
		}

		if oriPlayer.Name == player.Name {
			if err := container.DefaultContainer.Redis.SRem(ctx, game, v).Err(); err != nil {
				return err
			}
			return container.DefaultContainer.Redis.SAdd(ctx, game, player).Err()
		}
	}

	return container.DefaultContainer.Redis.SAdd(ctx, game, player).Err()
}

func getRemain(ctx context.Context, game string) (model.Play, error) {
	var result model.Play
	key := util.GenGameRemainKey(game)
	remainInfo := container.DefaultContainer.Redis.Get(ctx, key)
	if err := remainInfo.Err(); err != nil {
		return result, err
	}

	remainBytes, err := remainInfo.Bytes()
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(remainBytes, &result); err != nil {
		return result, err
	}

	return result, nil
}

func getPlayers(ctx context.Context, game string) (model.Players, error) {
	var result model.Players
	gameInfo := container.DefaultContainer.Redis.SMembers(ctx, game)
	if err := gameInfo.Err(); err != nil {
		return result, err
	}

	for _, v := range gameInfo.Val() {
		var player model.Player
		if err := json.Unmarshal([]byte(v), &player); err != nil {
			return result, err
		}
		result = append(result, player)
	}

	return result, nil
}

func getLastPlay(ctx context.Context, game string) (model.Player, error) {
	key := util.GenGameRecordKey(game)
	var result model.Player
	recordInfo := container.DefaultContainer.Redis.LRange(ctx, key, -1, -1)
	if err := recordInfo.Err(); err != nil {
		return result, err
	}

	for _, v := range recordInfo.Val() {
		if err := json.Unmarshal([]byte(v), &result); err != nil {
			return result, err
		}
		return result, nil
	}

	return result, nil
}

func recordPlay(ctx context.Context, game string, player model.Player) error {
	key := util.GenGameRecordKey(game)
	return container.DefaultContainer.Redis.RPush(ctx, key, player).Err()
}

func updatePlayerExpiration(ctx context.Context, user string) error {
	return container.DefaultContainer.Redis.SetEX(ctx, user, "occupied", time.Hour).Err()
}

package model

var (
	ValidPrefixMap = map[Prefix]bool{
		PrefixGameStart:     true,
		PrefixUncoverLord:   true,
		PrefixUncoverRemain: true,
		PrefixTakeRemain:    true,
		PrefixCancelRemain:  true,
		PrefixPlay:          true,
		PrefixNormal:        true,
		PrefixPlayer:        true,
		PrefixSpeaker:       true,
		PrefixDeal:          true,
		PrefixPass:          true,
		PrefixInvalidPlay:   true,
		PrefixYourName:      true,
		PrefixPlayerName:    true,
		PrefixWinner:        true,
	}
)

const (
	PrefixGameStart     Prefix = "GAME_START"
	PrefixUncoverLord   Prefix = "UNCOVER_LORD:"
	PrefixUncoverRemain Prefix = "UNCOVER_REMAIN:"
	PrefixTakeRemain    Prefix = "TAKE_REMAIN"
	PrefixCancelRemain  Prefix = "CANCEL_REMAIN"
	PrefixPlay          Prefix = "PLAY:"
	PrefixNormal        Prefix = "NORMAL:"
	PrefixPlayer        Prefix = "PLAYER:"
	PrefixSpeaker       Prefix = "SPEAKER:"
	PrefixDeal          Prefix = "DEAL:"
	PrefixPass          Prefix = "PASS"
	PrefixInvalidPlay   Prefix = "INVALID_PLAY"
	PrefixYourName      Prefix = "YOUR_NAME:"
	PrefixPlayerName    Prefix = "PLAYER_NAME:"
	PrefixWinner        Prefix = "WINNER:"
)

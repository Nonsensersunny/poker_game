package game

type Stage int

const (
	StageShuffle = iota
	StageDeal
	StageCallRemain
	StageTakeRemain
	StagePlay
	StageOver
)

func (s Stage) String() string {
	switch s {
	case StageShuffle:
		return "洗牌"
	case StageDeal:
		return "发牌"
	case StageCallRemain:
		return "叫地主"
	case StageTakeRemain:
		return "抢地主"
	case StagePlay:
		return "打牌"
	case StageOver:
		return "结算"
	default:
		return "Unknown stage"
	}
}

type Mode struct {
	Remain     int     `json:"remain"`
	PlayerNum  int     `json:"split"`
	Stages     []Stage `json:"stages"`
	StageIndex int     `json:"stage_index"`
}

func NewMode(remain, playerNum int, stages ...Stage) Mode {
	return Mode{
		Remain:    remain,
		PlayerNum: playerNum,
		Stages:    stages,
	}
}

var (
	ModeChinesePoker = NewMode(3, 3, []Stage{
		StageShuffle,
		StageDeal,
		StageCallRemain,
		StageTakeRemain,
		StagePlay,
		StageOver,
	}...)
)

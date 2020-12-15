package client

import (
	"fmt"
	"github.com/Nonsensersunny/poker_game/model"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func (c *Client) Refresh() {
	ui.Clear()

	var comps []ui.Drawable
	comps = append(comps, renderAllPlay(c.Player.Play, c.selected, c.cursor)...)

	var title string
	next := c.Player.Play.ExtractWithMap(c.selected)
	next.Sort()
	nextType := next.ValidatePlayType()
	if c.Player.IsLord {
		title = "[地主]"
	}
	if c.Player.Stop > 0 {
		title += fmt.Sprintf("该你出牌(倒计时：%vs)：%s", c.Player.Stop, nextType.Explain())
	} else {
		title += "等待其他玩家出牌..."
	}
	comps = append(comps, renderNextPlay(title, next, 1, 7)...)

	comps = append(comps, renderDashboard(c.Player.Win, c.Player.Lose))

	comps = append(comps, renderOtherPlayers(c.OtherPlayers)...)

	if len(comps) > 0 {
		ui.Render(comps...)
	}
}

func renderAllPlay(play model.Play, selectedMap map[int]bool, cursor int) []ui.Drawable {
	var (
		pokers   = make([]ui.Drawable, 0)
		defaultX = 3
		defaultY = 1
	)

	wrapper := widgets.NewParagraph()
	wrapper.Title = "剩余手牌:"
	wrapper.SetRect(defaultY, defaultY-1, len(play)*5+5, defaultY+5)

	pokers = append(pokers, wrapper)

	for i := 0; i < len(play); i++ {
		poker := model.NewPoker(play[i].Color, play[i].Pips)
		if !play[i].Used {
			pokers = append(pokers, drawPoker(i+1, poker, NewPosition(5*i+defaultX, defaultY, 5*(i+1)+defaultX, defaultY+3), i == cursor, selectedMap[i])...)
		}
	}

	return pokers
}

func renderNextPlay(title string, play model.Play, defaultX, defaultY int) []ui.Drawable {
	var (
		pokers   = make([]ui.Drawable, 0)
	)

	wrapper := widgets.NewParagraph()
	wrapper.Title = title
	wrapper.SetRect(defaultX, defaultY-1, 90, defaultY+5)

	pokers = append(pokers, wrapper)

	for i := 0; i < len(play); i++ {
		poker := model.NewPoker(play[i].Color, play[i].Pips)
		pokers = append(pokers, drawPoker(i+2, poker, NewPosition(5*i+defaultX+2, defaultY, 5*(i+1)+defaultX+2, defaultY+3), false, false)...)
	}

	return pokers
}

func renderTips() ui.Drawable {
	var (
		defaultX = 90
		defaultY = 1
	)

	wrapper := widgets.NewParagraph()
	wrapper.Title = "Tips:"
	wrapper.Text = "<q> 退出\t <p> Pass\n" +
		"<b> 开始\t <t> 抢地主\t <s> 排序\n" +
		"<space> (不)选择\n" +
		"<enter> 出牌\t <r> 重开"
	wrapper.SetRect(defaultX, defaultY-1, 116, defaultY+5)

	return wrapper
}

func renderDashboard(win, lose int) ui.Drawable {
	var (
		defaultX = 90
		defaultY = 7
	)

	wrapper := widgets.NewParagraph()
	wrapper.Title = "记分板："
	wrapper.TextStyle.Modifier = ui.ModifierBold
	wrapper.Text = fmt.Sprintf("胜场：%v\n负场：%v", win, lose)
	wrapper.SetRect(defaultX, defaultY-1, 116, defaultY+5)

	return wrapper
}

func renderOtherPlayers(players model.Players) []ui.Drawable {
	var comps []ui.Drawable
	for i, v := range players {
		pType := v.Play.ValidatePlayType()
		var title string
		if v.IsLord {
			title = "[地主]"
		}
		if v.IsSpeaker {
			title += fmt.Sprintf("玩家%s 倒计时：%vs：%s", v.Name, v.Stop, pType.Explain())
		} else {
			title += fmt.Sprintf("玩家%s等待其余玩家出牌...", v.Name)
		}
		comps = append(comps, renderNextPlay(title, v.Play, 1, 13+6*i)...)
	}
	return comps
}

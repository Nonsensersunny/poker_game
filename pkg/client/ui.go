package client

import (
	"fmt"
	"github.com/Nonsensersunny/poker_game/model"
	"github.com/Nonsensersunny/poker_game/util"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	log "github.com/sirupsen/logrus"
	"os"
)

func (c *Client) InitUI(serverSig chan string) {
	if err := ui.Init(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer ui.Close()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		//log.Printf("%v", e.ID)
		switch e.ID {
		case "q", "<C-c>":
			return
		case "h":
			ui.Clear()
			ui.Render(renderTips())
		case "<Space>":
			if c.Player.IsSpeaker && c.cursor < len(c.Player.Play) {
				c.selected[c.cursor] = !c.selected[c.cursor]
				c.Refresh()
			}
		case "<Left>":
			c.cursor -= 1
			if c.cursor < 0 {
				c.cursor = len(c.Player.Play) - 1
			}
			c.Refresh()
		case "<Right>":
			c.cursor += 1
			if c.cursor == len(c.Player.Play) {
				c.cursor = 0
			}
			c.Refresh()
		case "s":
			c.Player.Play.Sort()
			c.Refresh()
		case "p":
			c.Pass(serverSig)
		case "<Up>":
		case "<Down>":
		case "<Enter>":
			if c.Player.IsSpeaker {
				extract := c.Player.Play.ExtractWithMap(c.selected)
				if len(c.selected) < 1 {
					continue
				}
				last, _ := c.Records.Last()
				if extract.GreaterThan(last.Play) {
					extractType := extract.ValidatePlayType()
					if extractType.Type != model.PlayTypeInvalid {
						c.Player.Play.DealWithIndexMap(c.selected)
					}

					serverSig <- model.PrefixPlay.AssembleMessage(extract.String())

					c.selected = make(map[int]bool)
					c.Player.Play = c.Player.Play.ExtractUnused()
					c.cursor = 0
					player := model.NewPlayer(c.Player.Name, extract)
					player.Stop = util.CountDownRevoke
					c.Records = append(c.Records, player)
					c.Player.Stop = 0
				}
				c.Refresh()
			}
		}
	}
}

func transferColor(c model.Color) ui.Color {
	switch c {
	case model.ColorSpade, model.ColorClub:
		return ui.ColorWhite
	case model.ColorHeart, model.ColorDiamond:
		return ui.ColorMagenta
	case model.ColorMoon:
		return ui.ColorWhite
	default:
		return ui.ColorRed
	}
}

func drawPoker(sub int, poker model.Poker, position Position, prepared, selected bool) []ui.Drawable {
	var result []ui.Drawable

	wp := widgets.NewParagraph()
	wp.Title = poker.Name
	wp.Text = poker.Color.ToUnicode()
	wp.PaddingLeft = 1
	wp.TextStyle.Fg = transferColor(poker.Color)
	wp.SetRect(position.X1, position.Y1, position.X2, position.Y2)
	wp.BorderStyle.Fg = transferColor(poker.Color)

	idx := widgets.NewParagraph()
	idx.Title = fmt.Sprint(sub)
	idx.PaddingLeft = 1
	idx.TextStyle.Fg = transferColor(poker.Color)
	idx.SetRect(position.X1, position.Y1+3, 4+position.X2, position.Y1+4)
	idx.BorderStyle.Fg = ui.ColorWhite
	idx.Border = false

	if selected {
		wp.BorderStyle.Fg = ui.ColorGreen
		idx.BorderStyle.Fg = ui.ColorGreen
	}

	if prepared {
		wp.BorderStyle.Fg = ui.ColorBlue
		idx.BorderStyle.Fg = ui.ColorBlue
		wp.TextStyle.Fg = ui.ColorBlue
		idx.TitleStyle.Fg = ui.ColorBlue
	}

	return append(result, wp, idx)
}

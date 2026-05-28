package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (receiver *App) layoutBottomNav(gtx layout.Context) layout.Dimensions {
	if receiver.navContactsClick.Clicked(gtx) {
		receiver.page = PageHome
	}
	if receiver.navChatsClick.Clicked(gtx) {
		receiver.page = PageChats
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min
			defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xFA, G: 0xFA, B: 0xFA, A: 0xFF})
			return layout.Dimensions{Size: size}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Top border line.
			var lineHeight int = gtx.Dp(unit.Dp(1))
			defer clip.Rect{Max: image.Point{X: gtx.Constraints.Max.X, Y: lineHeight}}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xDD, G: 0xDD, B: 0xDD, A: 0xFF})
			return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: lineHeight}}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(1)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layoutNavTab(gtx, receiver.theme, &receiver.navContactsClick, "Contacts", receiver.page == PageHome)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layoutNavTab(gtx, receiver.theme, &receiver.navChatsClick, "Chats", receiver.page == PageChats)
					}),
				)
			})
		}),
	)
}

func layoutNavTab(gtx layout.Context, th *material.Theme, click *widget.Clickable, label string, active bool) layout.Dimensions {
	return material.Clickable(gtx, click, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Subtitle2(th, label)
				if active {
					lbl.Color = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF} // indigo
				} else {
					lbl.Color = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF} // gray
				}
				return lbl.Layout(gtx)
			})
		})
	})
}

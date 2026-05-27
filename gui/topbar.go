package gui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var topBarColor color.NRGBA = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
var topBarTextColor color.NRGBA = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

func (receiver *App) layoutTopBar(gtx layout.Context, title string, showBack bool) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min
			defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, topBarColor)
			return layout.Dimensions{Size: size}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				if !showBack {
					label := material.H6(receiver.theme, title)
					label.Color = topBarTextColor
					return label.Layout(gtx)
				}

				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(receiver.theme, &receiver.backClick, "< Back")
						btn.Color = topBarTextColor
						btn.Background = color.NRGBA{}
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.H6(receiver.theme, title)
						label.Color = topBarTextColor
						return label.Layout(gtx)
					}),
				)
			})
		}),
	)
}

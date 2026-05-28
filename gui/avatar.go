package gui

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var avatarColors []color.NRGBA = []color.NRGBA{
	{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}, // indigo
	{R: 0x00, G: 0x96, B: 0x88, A: 0xFF}, // teal
	{R: 0xE9, G: 0x1E, B: 0x63, A: 0xFF}, // pink
	{R: 0xFF, G: 0x57, B: 0x22, A: 0xFF}, // deep orange
	{R: 0x67, G: 0x3A, B: 0xB7, A: 0xFF}, // deep purple
	{R: 0x21, G: 0x96, B: 0xF3, A: 0xFF}, // blue
	{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}, // green
	{R: 0xFF, G: 0x98, B: 0x00, A: 0xFF}, // orange
}

func initials(name string) string {
	var parts []string = strings.Fields(name)
	if 0 == len(parts) {
		return "?"
	}
	if 1 == len(parts) {
		return strings.ToUpper(parts[0][:1])
	}
	return strings.ToUpper(parts[0][:1]) + strings.ToUpper(parts[len(parts)-1][:1])
}

func avatarColor(name string) color.NRGBA {
	var hash int
	for _, c := range name {
		hash = hash*31 + int(c)
	}
	if hash < 0 {
		hash = -hash
	}
	return avatarColors[hash%len(avatarColors)]
}

func layoutAvatar(gtx layout.Context, th *material.Theme, name string, size unit.Dp) layout.Dimensions {
	var sizePx int = gtx.Dp(size)
	var bgColor color.NRGBA = avatarColor(name)
	var text string = initials(name)

	gtx.Constraints.Min = image.Point{X: sizePx, Y: sizePx}
	gtx.Constraints.Max = image.Point{X: sizePx, Y: sizePx}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			var sz image.Point = image.Point{X: sizePx, Y: sizePx}
			defer clip.Ellipse{Max: sz}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, bgColor)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body1(th, text)
			lbl.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
			if size >= unit.Dp(48) {
				lbl = material.H6(th, text)
				lbl.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
			}
			return lbl.Layout(gtx)
		}),
	)
}

package gui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func layoutProfileHeader(gtx layout.Context, th *material.Theme, name string, title string, company string, fediID string) layout.Dimensions {
	var bannerHeight unit.Dp = unit.Dp(120)
	var avatarSize unit.Dp = unit.Dp(72)
	var avatarOverlap unit.Dp = unit.Dp(36)

	var baseColor color.NRGBA = avatarColor(name)

	// Darker variant for gradient bottom.
	var darkColor color.NRGBA = color.NRGBA{
		R: uint8(float64(baseColor.R) * 0.6),
		G: uint8(float64(baseColor.G) * 0.6),
		B: uint8(float64(baseColor.B) * 0.6),
		A: 0xFF,
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Banner + overlapping avatar.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.SW}.Layout(gtx,
				// Banner background.
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					var bannerPx int = gtx.Dp(bannerHeight)
					var overlapPx int = gtx.Dp(avatarOverlap)
					var totalHeight int = bannerPx + overlapPx

					// Top half: base color.
					func() {
						defer clip.Rect{Max: image.Point{X: gtx.Constraints.Max.X, Y: bannerPx / 2}}.Push(gtx.Ops).Pop()
						paint.Fill(gtx.Ops, baseColor)
					}()

					// Bottom half: darker color.
					func() {
						rect := image.Rect(0, bannerPx/2, gtx.Constraints.Max.X, bannerPx)
						defer clip.Rect{Min: rect.Min, Max: rect.Max}.Push(gtx.Ops).Pop()
						paint.Fill(gtx.Ops, darkColor)
					}()

					return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: totalHeight}}
				}),
				// Avatar overlapping the bottom of the banner.
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// White ring behind avatar.
						var ringSize unit.Dp = avatarSize + unit.Dp(6)
						var ringSizePx int = gtx.Dp(ringSize)

						return layout.Stack{Alignment: layout.Center}.Layout(gtx,
							layout.Expanded(func(gtx layout.Context) layout.Dimensions {
								var sz image.Point = image.Point{X: ringSizePx, Y: ringSizePx}
								defer clip.Ellipse{Max: sz}.Push(gtx.Ops).Pop()
								paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
								return layout.Dimensions{Size: sz}
							}),
							layout.Stacked(func(gtx layout.Context) layout.Dimensions {
								return layoutAvatar(gtx, th, name, avatarSize)
							}),
						)
					})
				}),
			)
		}),
		// Name + subtitle + FediID below.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(12), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.H6(th, name).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						var subtitle string
						if "" != company {
							subtitle = fmt.Sprintf("%s — %s", title, company)
						} else {
							subtitle = title
						}
						lbl := material.Body2(th, subtitle)
						lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if "" == fediID {
							return layout.Dimensions{}
						}
						lbl := material.Caption(th, fediID)
						lbl.Color = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
						return lbl.Layout(gtx)
					}),
				)
			})
		}),
	)
}

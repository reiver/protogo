package gui

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"protogo/cfg"
)

func (receiver *App) layoutProfilePage(gtx layout.Context) layout.Dimensions {
	var me Person = receiver.me

	var widgets []layout.Widget

	// Avatar + name header.
	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(12), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layoutAvatar(gtx, receiver.theme, me.Name, unit.Dp(64))
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.H6(receiver.theme, me.Name).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							var subtitle string
							if "" != me.Company {
								subtitle = fmt.Sprintf("%s — %s", me.Title, me.Company)
							} else {
								subtitle = me.Title
							}
							lbl := material.Body2(receiver.theme, subtitle)
							lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if "" == me.FediID {
								return layout.Dimensions{}
							}
							lbl := material.Caption(receiver.theme, me.FediID)
							lbl.Color = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
							return lbl.Layout(gtx)
						}),
					)
				}),
			)
		})
	})

	// Detail sections.
	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layoutDetailSection(gtx, receiver.theme, "Title", me.Title)
		})
	})

	if "" != me.Company {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layoutDetailSection(gtx, receiver.theme, "Company", me.Company)
			})
		})
	}

	if "" != me.FediID {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layoutDetailSection(gtx, receiver.theme, "Fediverse ID", me.FediID)
			})
		})
	}

	if "" != me.Note {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layoutDetailSection(gtx, receiver.theme, "Notes", me.Note)
			})
		})
	}

	// Resumes section.
	if 0 < len(me.Resumes) {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return receiver.layoutResumesSection(gtx, me.Resumes, PageProfile)
			})
		})
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, cfg.Name, false)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.profileList.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
				return widgets[index](gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutBottomNav(gtx)
		}),
	)
}

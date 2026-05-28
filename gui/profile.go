package gui

import (
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"codeberg.org/reiver/go-fediverseid"

	"protogo/cfg"
)

func (receiver *App) layoutProfilePage(gtx layout.Context) layout.Dimensions {
	// Handle Fedi ID submit.
	for {
		event, ok := receiver.fediIDEditor.Update(gtx)
		if !ok {
			break
		}
		if _, ok := event.(widget.SubmitEvent); ok {
			var text string = strings.TrimSpace(receiver.fediIDEditor.Text())
			if "" == text {
				receiver.profileFediIDError = "Fediverse ID cannot be empty."
			} else if _, err := fediverseid.ParseFediverseIDString(text); nil != err {
				receiver.profileFediIDError = "Not a valid Fediverse ID. Expected format: @name@host"
			} else {
				receiver.profileFediIDError = ""
				receiver.me.FediID = text
				receiver.fediIDEditor.SetText(text)
				persistProfileFediID(text)
			}
		}
	}

	var me Person = receiver.me

	var widgets []layout.Widget

	// Profile header with banner + avatar.
	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layoutProfileHeader(gtx, receiver.theme, me.Name, me.Title, me.Company, me.FediID)
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

	// Editable Fediverse ID field.
	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Caption(receiver.theme, "Fediverse ID")
					lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(receiver.theme, &receiver.fediIDEditor, "@you@example.com")
					editor.Color = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
					editor.HintColor = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}
					return editor.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if "" == receiver.profileFediIDError {
						return layout.Dimensions{}
					}
					return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(receiver.theme, receiver.profileFediIDError)
						lbl.Color = color.NRGBA{R: 0xD3, G: 0x2F, B: 0x2F, A: 0xFF}
						return lbl.Layout(gtx)
					})
				}),
			)
		})
	})

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

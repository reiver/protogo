package gui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

func (receiver *App) layoutPersonDetail(gtx layout.Context) layout.Dimensions {
	if receiver.backClick.Clicked(gtx) {
		receiver.page = PageHome
	}

	if receiver.chatClick.Clicked(gtx) {
		receiver.chatFrom = PagePersonDetail
		receiver.page = PageChat
	}

	if receiver.selectedPerson < 0 || receiver.selectedPerson >= len(receiver.people) {
		receiver.page = PageHome
		return layout.Dimensions{}
	}

	if receiver.favClick.Clicked(gtx) {
		receiver.people[receiver.selectedPerson].Favorite = !receiver.people[receiver.selectedPerson].Favorite
	}

	var person Person = receiver.people[receiver.selectedPerson]

	var widgets []layout.Widget

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		var favIconData []byte
		var favColor color.NRGBA
		if person.Favorite {
			favIconData = icons.ToggleStar
			favColor = color.NRGBA{R: 0xFF, G: 0xB3, B: 0x00, A: 0xFF} // amber
		} else {
			favIconData = icons.ToggleStarBorder
			favColor = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF} // gray
		}
		return layoutProfileHeader(gtx, receiver.theme, person.Name, person.Title, person.Company, person.FediID,
			func(gtx layout.Context) layout.Dimensions {
				return layoutIconButton(gtx, &receiver.favClick, favIconData, favColor)
			},
			func(gtx layout.Context) layout.Dimensions {
				return layoutIconButton(gtx, &receiver.chatClick, icons.CommunicationChat, color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF})
			},
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layoutDetailSection(gtx, receiver.theme, "Title", person.Title)
		})
	})

	if "" != person.Company {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layoutDetailSection(gtx, receiver.theme, "Company", person.Company)
			})
		})
	}

	if "" != person.Note {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layoutDetailSection(gtx, receiver.theme, "Notes", person.Note)
			})
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutResumesSection(gtx, person.Resumes, PagePersonDetail)
		})
	})

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, person.Name, true)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.personList.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
				return widgets[index](gtx)
			})
		}),
	)
}

func (receiver *App) layoutFavoriteButton(gtx layout.Context, favorited bool) layout.Dimensions {
	return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		var label string
		if favorited {
			label = "★ Favorited"
		} else {
			label = "☆ Add to Favorites"
		}
		btn := material.Button(receiver.theme, &receiver.favClick, label)
		if favorited {
			btn.Background = color.NRGBA{R: 0xFF, G: 0xB3, B: 0x00, A: 0xFF} // amber
		} else {
			btn.Background = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF} // gray
		}
		return btn.Layout(gtx)
	})
}

func layoutIconButton(gtx layout.Context, click *widget.Clickable, iconData []byte, bg color.NRGBA) layout.Dimensions {
	icon, _ := widget.NewIcon(iconData)

	return material.Clickable(gtx, click, func(gtx layout.Context) layout.Dimensions {
		var sizeDp unit.Dp = unit.Dp(48)
		var sizePx int = gtx.Dp(sizeDp)

		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				var sz image.Point = image.Point{X: sizePx, Y: sizePx}
				defer clip.Ellipse{Max: sz}.Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, bg)
				return layout.Dimensions{Size: sz}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				var iconPx int = gtx.Dp(unit.Dp(24))
				gtx.Constraints.Min = image.Point{X: iconPx, Y: iconPx}
				gtx.Constraints.Max = image.Point{X: iconPx, Y: iconPx}
				return icon.Layout(gtx, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
			}),
		)
	})
}

func layoutDetailSection(gtx layout.Context, th *material.Theme, label string, value string) layout.Dimensions {
	return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(th, label)
				lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th, value).Layout(gtx)
			}),
		)
	})
}

func (receiver *App) layoutResumesSection(gtx layout.Context, resumes []Resume, fromPage Page) layout.Dimensions {
	var th *material.Theme = receiver.theme

	// Ensure we have enough clickables for the resumes.
	for len(receiver.resumeClicks) < len(resumes) {
		receiver.resumeClicks = append(receiver.resumeClicks, widget.Clickable{})
	}

	// Handle resume clicks.
	for i := range resumes {
		if receiver.resumeClicks[i].Clicked(gtx) {
			receiver.selectedResume = i
			receiver.resumeFrom = fromPage
			receiver.page = PageResumeDetail
		}
	}

	if 0 == len(resumes) {
		return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Caption(th, "Resumes")
					lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, "No resumes")
					lbl.Color = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}
					return lbl.Layout(gtx)
				}),
			)
		})
	}

	var children []layout.FlexChild

	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		lbl := material.Caption(th, fmt.Sprintf("Resumes (%d)", len(resumes)))
		lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
		return lbl.Layout(gtx)
	}))

	for i := range resumes {
		var resumeIndex int = i
		var resume Resume = resumes[resumeIndex]
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &receiver.resumeClicks[resumeIndex], func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layoutCard(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Subtitle2(th, resume.Label).Layout(gtx)
					})
				})
			})
		}))
	}

	return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
}

package gui

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (receiver *App) layoutPersonDetail(gtx layout.Context) layout.Dimensions {
	if receiver.backClick.Clicked(gtx) {
		receiver.page = PageHome
	}

	if receiver.chatClick.Clicked(gtx) {
		receiver.page = PageChat
	}

	if receiver.selectedPerson < 0 || receiver.selectedPerson >= len(receiver.people) {
		receiver.page = PageHome
		return layout.Dimensions{}
	}

	var person Person = receiver.people[receiver.selectedPerson]

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, person.Name, true)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layoutDetailSection(gtx, receiver.theme, "Title", person.Title)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if "" == person.Company {
							return layout.Dimensions{}
						}
						return layoutDetailSection(gtx, receiver.theme, "Company", person.Company)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if "" == person.Note {
							return layout.Dimensions{}
						}
						return layoutDetailSection(gtx, receiver.theme, "Notes", person.Note)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return receiver.layoutResumesSection(gtx, person.Resumes)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(receiver.theme, &receiver.chatClick, "Chat")
							return btn.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
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

func (receiver *App) layoutResumesSection(gtx layout.Context, resumes []Resume) layout.Dimensions {
	var th *material.Theme = receiver.theme

	// Ensure we have enough clickables for the resumes.
	for len(receiver.resumeClicks) < len(resumes) {
		receiver.resumeClicks = append(receiver.resumeClicks, widget.Clickable{})
	}

	// Handle resume clicks.
	for i := range resumes {
		if receiver.resumeClicks[i].Clicked(gtx) {
			receiver.selectedResume = i
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

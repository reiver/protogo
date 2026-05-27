package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (receiver *App) layoutResumeDetail(gtx layout.Context) layout.Dimensions {
	if receiver.backClick.Clicked(gtx) {
		receiver.page = PagePersonDetail
	}

	if receiver.selectedPerson < 0 || receiver.selectedPerson >= len(receiver.people) {
		receiver.page = PageHome
		return layout.Dimensions{}
	}

	var person Person = receiver.people[receiver.selectedPerson]

	if receiver.selectedResume < 0 || receiver.selectedResume >= len(person.Resumes) {
		receiver.page = PagePersonDetail
		return layout.Dimensions{}
	}

	var resume Resume = person.Resumes[receiver.selectedResume]

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, resume.Label, true)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layoutDetailSection(gtx, receiver.theme, "Person", person.Name)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(receiver.theme, resume.Content).Layout(gtx)
					}),
				)
			})
		}),
	)
}

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

	"protogo/cfg"
)

func (receiver *App) layoutHome(gtx layout.Context) layout.Dimensions {

	// Handle person clicks.
	for i := range receiver.personClicks {
		if receiver.personClicks[i].Clicked(gtx) {
			receiver.selectedPerson = i
			receiver.page = PagePersonDetail
		}
	}

	// Handle group clicks.
	for i := range receiver.groupClicks {
		if receiver.groupClicks[i].Clicked(gtx) {
			receiver.selectedGroup = i
			receiver.page = PageGroupDetail
		}
	}

	var totalItems int = len(receiver.people) + 1 + len(receiver.groups) + 1 // +1 for each section header

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, cfg.Name, false)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.homeList.Layout(gtx, totalItems, func(gtx layout.Context, index int) layout.Dimensions {

				var peopleHeaderIndex int = 0
				var peopleStartIndex int = 1
				var peopleEndIndex int = len(receiver.people)
				var groupsHeaderIndex int = peopleEndIndex + 1
				var groupsStartIndex int = groupsHeaderIndex + 1

				switch {
				case index == peopleHeaderIndex:
					return layoutSectionHeader(gtx, receiver.theme, "People")
				case index >= peopleStartIndex && index <= peopleEndIndex:
					var personIndex int = index - peopleStartIndex
					return receiver.layoutPersonItem(gtx, personIndex)
				case index == groupsHeaderIndex:
					return layoutSectionHeader(gtx, receiver.theme, "Groups")
				default:
					var groupIndex int = index - groupsStartIndex
					return receiver.layoutGroupItem(gtx, groupIndex)
				}
			})
		}),
	)
}

func layoutSectionHeader(gtx layout.Context, th *material.Theme, title string) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			label := material.Subtitle1(th, title)
			label.Color = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
			return label.Layout(gtx)
		},
	)
}

func (receiver *App) layoutPersonItem(gtx layout.Context, index int) layout.Dimensions {
	if index < 0 || index >= len(receiver.people) {
		return layout.Dimensions{}
	}

	var person Person = receiver.people[index]

	return material.Clickable(gtx, &receiver.personClicks[index], func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layoutCard(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body1(receiver.theme, person.Name).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							var subtitle string
							if "" != person.Company {
								subtitle = fmt.Sprintf("%s — %s", person.Title, person.Company)
							} else {
								subtitle = person.Title
							}
							label := material.Caption(receiver.theme, subtitle)
							label.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
							return label.Layout(gtx)
						}),
					)
				})
			},
		)
	})
}

func (receiver *App) layoutGroupItem(gtx layout.Context, index int) layout.Dimensions {
	if index < 0 || index >= len(receiver.groups) {
		return layout.Dimensions{}
	}

	var group Group = receiver.groups[index]

	return material.Clickable(gtx, &receiver.groupClicks[index], func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layoutCard(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body1(receiver.theme, group.Name).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Caption(receiver.theme, fmt.Sprintf("%d members", len(group.Members)))
							label.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
							return label.Layout(gtx)
						}),
					)
				})
			},
		)
	})
}

func layoutCard(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min
			rr := gtx.Dp(unit.Dp(8))
			defer clip.RRect{
				Rect: image.Rectangle{Max: size},
				SE:   rr,
				SW:   rr,
				NW:   rr,
				NE:   rr,
			}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF5, A: 0xFF})
			return layout.Dimensions{Size: size}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, w)
		}),
	)
}

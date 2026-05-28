package gui

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"protogo/cfg"
)

type filteredItem struct {
	kind  int // 0 = people header, 1 = person, 2 = groups header, 3 = group
	index int // index into receiver.people or receiver.groups
}

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

	// Build filtered items based on search query.
	var query string = strings.ToLower(strings.TrimSpace(receiver.searchEditor.Text()))

	var items []filteredItem

	// Filter people.
	var filteredPeople bool
	for i, person := range receiver.people {
		if "" == query || strings.Contains(strings.ToLower(person.Name), query) || strings.Contains(strings.ToLower(person.Title), query) || strings.Contains(strings.ToLower(person.Company), query) {
			if !filteredPeople {
				items = append(items, filteredItem{kind: 0})
				filteredPeople = true
			}
			items = append(items, filteredItem{kind: 1, index: i})
		}
	}

	// Filter groups.
	var filteredGroups bool
	for i, group := range receiver.groups {
		if "" == query || strings.Contains(strings.ToLower(group.Name), query) {
			if !filteredGroups {
				items = append(items, filteredItem{kind: 2})
				filteredGroups = true
			}
			items = append(items, filteredItem{kind: 3, index: i})
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, cfg.Name, false)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutSearchBar(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.homeList.Layout(gtx, len(items), func(gtx layout.Context, index int) layout.Dimensions {
				var item filteredItem = items[index]

				switch item.kind {
				case 0:
					return layoutSectionHeader(gtx, receiver.theme, "People")
				case 1:
					return receiver.layoutPersonItem(gtx, item.index)
				case 2:
					return layoutSectionHeader(gtx, receiver.theme, "Groups")
				case 3:
					return receiver.layoutGroupItem(gtx, item.index)
				default:
					return layout.Dimensions{}
				}
			})
		}),
	)
}

func (receiver *App) layoutSearchBar(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(4), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layoutCard(gtx, func(gtx layout.Context) layout.Dimensions {
			editor := material.Editor(receiver.theme, &receiver.searchEditor, "Search people and groups...")
			return editor.Layout(gtx)
		})
	})
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

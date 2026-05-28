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
	"gioui.org/widget"
	"gioui.org/widget/material"

	"protogo/cfg"
)

type filteredItem struct {
	kind  int // 0 = favorites header, 1 = fav person, 2 = fav group, 3 = people header, 4 = person, 5 = groups header, 6 = group
	index int // index into receiver.people or receiver.groups
}

func (receiver *App) layoutHome(gtx layout.Context) layout.Dimensions {

	// Handle person clicks.
	for i := range receiver.personClicks {
		if receiver.personClicks[i].Clicked(gtx) {
			receiver.selectedPerson = i
			receiver.personFrom = PageHome
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

	// Handle person favorite toggles.
	for i := range receiver.personFavClicks {
		if receiver.personFavClicks[i].Clicked(gtx) {
			receiver.people[i].Favorite = !receiver.people[i].Favorite
		}
	}

	// Handle group favorite toggles.
	for i := range receiver.groupFavClicks {
		if receiver.groupFavClicks[i].Clicked(gtx) {
			receiver.groups[i].Favorite = !receiver.groups[i].Favorite
		}
	}

	// Build filtered items based on search query.
	var query string = strings.ToLower(strings.TrimSpace(receiver.searchEditor.Text()))

	var items []filteredItem

	// Favorites section.
	var hasFavorites bool
	for i, person := range receiver.people {
		if person.Favorite && ("" == query || strings.Contains(strings.ToLower(person.Name), query) || strings.Contains(strings.ToLower(person.Title), query) || strings.Contains(strings.ToLower(person.Company), query)) {
			if !hasFavorites {
				items = append(items, filteredItem{kind: 0})
				hasFavorites = true
			}
			items = append(items, filteredItem{kind: 1, index: i})
		}
	}
	for i, group := range receiver.groups {
		if group.Favorite && ("" == query || strings.Contains(strings.ToLower(group.Name), query)) {
			if !hasFavorites {
				items = append(items, filteredItem{kind: 0})
				hasFavorites = true
			}
			items = append(items, filteredItem{kind: 2, index: i})
		}
	}

	// Filter people.
	var filteredPeople bool
	for i, person := range receiver.people {
		if "" == query || strings.Contains(strings.ToLower(person.Name), query) || strings.Contains(strings.ToLower(person.Title), query) || strings.Contains(strings.ToLower(person.Company), query) {
			if !filteredPeople {
				items = append(items, filteredItem{kind: 3})
				filteredPeople = true
			}
			items = append(items, filteredItem{kind: 4, index: i})
		}
	}

	// Filter groups.
	var filteredGroups bool
	for i, group := range receiver.groups {
		if "" == query || strings.Contains(strings.ToLower(group.Name), query) {
			if !filteredGroups {
				items = append(items, filteredItem{kind: 5})
				filteredGroups = true
			}
			items = append(items, filteredItem{kind: 6, index: i})
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
					return layoutSectionHeader(gtx, receiver.theme, "Favorites")
				case 1:
					return receiver.layoutPersonItem(gtx, item.index)
				case 2:
					return receiver.layoutGroupItem(gtx, item.index)
				case 3:
					return layoutSectionHeader(gtx, receiver.theme, "People")
				case 4:
					return receiver.layoutPersonItem(gtx, item.index)
				case 5:
					return layoutSectionHeader(gtx, receiver.theme, "Groups")
				case 6:
					return receiver.layoutGroupItem(gtx, item.index)
				default:
					return layout.Dimensions{}
				}
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutBottomNav(gtx)
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
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layoutAvatar(gtx, receiver.theme, person.Name, unit.Dp(40))
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
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
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if "" == person.FediID {
										return layout.Dimensions{}
									}
									lbl := material.Caption(receiver.theme, person.FediID)
									lbl.Color = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
									return lbl.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layoutStarToggle(gtx, receiver.theme, &receiver.personFavClicks[index], person.Favorite)
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
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layoutAvatar(gtx, receiver.theme, group.Name, unit.Dp(40))
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
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
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layoutStarToggle(gtx, receiver.theme, &receiver.groupFavClicks[index], group.Favorite)
						}),
					)
				})
			},
		)
	})
}

func layoutStarToggle(gtx layout.Context, th *material.Theme, click *widget.Clickable, favorited bool) layout.Dimensions {
	return material.Clickable(gtx, click, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			var text string
			var starColor color.NRGBA
			if favorited {
				text = "★"
				starColor = color.NRGBA{R: 0xFF, G: 0xB3, B: 0x00, A: 0xFF} // amber
			} else {
				text = "☆"
				starColor = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF} // gray
			}
			lbl := material.H6(th, text)
			lbl.Color = starColor
			return lbl.Layout(gtx)
		})
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

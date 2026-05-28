package gui

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

func (receiver *App) layoutGroupDetail(gtx layout.Context) layout.Dimensions {
	if receiver.backClick.Clicked(gtx) {
		receiver.page = PageHome
	}

	if receiver.groupChatClick.Clicked(gtx) {
		receiver.chatFrom = PageGroupDetail
		receiver.page = PageGroupChat
	}

	if receiver.selectedGroup < 0 || receiver.selectedGroup >= len(receiver.groups) {
		receiver.page = PageHome
		return layout.Dimensions{}
	}

	if receiver.favClick.Clicked(gtx) {
		receiver.groups[receiver.selectedGroup].Favorite = !receiver.groups[receiver.selectedGroup].Favorite
	}

	var group Group = receiver.groups[receiver.selectedGroup]

	var widgets []layout.Widget

	// Profile header with banner, avatar, favorite + chat icons.
	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		var favIconData []byte
		var favColor color.NRGBA
		if group.Favorite {
			favIconData = icons.ToggleStar
			favColor = color.NRGBA{R: 0xFF, G: 0xB3, B: 0x00, A: 0xFF}
		} else {
			favIconData = icons.ToggleStarBorder
			favColor = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}
		}
		return layoutProfileHeader(gtx, receiver.theme, group.Name, fmt.Sprintf("%d members", len(group.Members)), "", "",
			func(gtx layout.Context) layout.Dimensions {
				return layoutIconButton(gtx, &receiver.favClick, favIconData, favColor)
			},
			func(gtx layout.Context) layout.Dimensions {
				return layoutIconButton(gtx, &receiver.groupChatClick, icons.CommunicationChat, color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF})
			},
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Caption(receiver.theme, "Members")
			lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
			return lbl.Layout(gtx)
		})
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16), Bottom: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutMembersList(gtx, group.Members)
		})
	})

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, group.Name, true)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.groupList.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
				return widgets[index](gtx)
			})
		}),
	)
}

func (receiver *App) findPersonIndex(name string) int {
	for i, person := range receiver.people {
		if person.Name == name {
			return i
		}
	}
	return -1
}

func (receiver *App) layoutMembersList(gtx layout.Context, members []string) layout.Dimensions {
	// Ensure we have enough clickables.
	for len(receiver.memberClicks) < len(members) {
		receiver.memberClicks = append(receiver.memberClicks, widget.Clickable{})
	}

	// Handle member clicks.
	for i, member := range members {
		if receiver.memberClicks[i].Clicked(gtx) {
			var personIndex int = receiver.findPersonIndex(member)
			if 0 <= personIndex {
				receiver.selectedPerson = personIndex
				receiver.personFrom = PageGroupDetail
				receiver.page = PagePersonDetail
			}
		}
	}

	var children []layout.FlexChild

	for i := range members {
		var memberIndex int = i
		var member string = members[memberIndex]
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &receiver.memberClicks[memberIndex], func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layoutCard(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layoutAvatar(gtx, receiver.theme, member, unit.Dp(36))
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return material.Body1(receiver.theme, member).Layout(gtx)
							}),
						)
					})
				})
			})
		}))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

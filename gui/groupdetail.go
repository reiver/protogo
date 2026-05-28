package gui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
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

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, group.Name, true)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return receiver.layoutFavoriteButton(gtx, group.Favorite)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(receiver.theme, "Members")
						lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return receiver.layoutMembersList(gtx, group.Members)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(receiver.theme, &receiver.groupChatClick, "Group Chat")
							return btn.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

func (receiver *App) layoutMembersList(gtx layout.Context, members []string) layout.Dimensions {
	var children []layout.FlexChild

	for i := range members {
		var member string = members[i]
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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
		}))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

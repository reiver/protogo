package gui

import (
	"fmt"
	"image/color"
	"sort"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"protogo/cfg"
)

type chatEntry struct {
	isGroup bool
	index   int
	name    string
	lastMsg ChatMessage
}

func (receiver *App) buildChatEntries() []chatEntry {
	var entries []chatEntry

	for i, person := range receiver.people {
		if 0 == len(person.Messages) {
			continue
		}
		entries = append(entries, chatEntry{
			isGroup: false,
			index:   i,
			name:    person.Name,
			lastMsg: person.Messages[len(person.Messages)-1],
		})
	}

	for i, group := range receiver.groups {
		if 0 == len(group.Messages) {
			continue
		}
		entries = append(entries, chatEntry{
			isGroup: true,
			index:   i,
			name:    group.Name,
			lastMsg: group.Messages[len(group.Messages)-1],
		})
	}

	sort.Slice(entries, func(a, b int) bool {
		return entries[a].lastMsg.Timestamp > entries[b].lastMsg.Timestamp
	})

	return entries
}

func (receiver *App) layoutChatsPage(gtx layout.Context) layout.Dimensions {
	var entries []chatEntry = receiver.buildChatEntries()

	// Ensure we have enough clickables.
	for len(receiver.chatItemClicks) < len(entries) {
		receiver.chatItemClicks = append(receiver.chatItemClicks, widget.Clickable{})
	}

	// Handle clicks.
	for i, entry := range entries {
		if receiver.chatItemClicks[i].Clicked(gtx) {
			if entry.isGroup {
				receiver.selectedGroup = entry.index
				receiver.chatFrom = PageChats
				receiver.page = PageGroupChat
			} else {
				receiver.selectedPerson = entry.index
				receiver.chatFrom = PageChats
				receiver.page = PageChat
			}
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, cfg.Name, false)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.chatsList.Layout(gtx, len(entries), func(gtx layout.Context, index int) layout.Dimensions {
				return receiver.layoutChatItem(gtx, entries[index], index)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutBottomNav(gtx)
		}),
	)
}

func (receiver *App) layoutChatItem(gtx layout.Context, entry chatEntry, clickIndex int) layout.Dimensions {
	return material.Clickable(gtx, &receiver.chatItemClicks[clickIndex], func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layoutCard(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layoutAvatar(gtx, receiver.theme, entry.name, unit.Dp(48))
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										return material.Body1(receiver.theme, entry.name).Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										lbl := material.Caption(receiver.theme, entry.lastMsg.Timestamp)
										lbl.Color = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}
										return lbl.Layout(gtx)
									}),
								)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								var preview string
								if entry.lastMsg.FromMe {
									preview = fmt.Sprintf("You: %s", entry.lastMsg.Text)
								} else if entry.isGroup && "" != entry.lastMsg.Sender {
									preview = fmt.Sprintf("%s: %s", entry.lastMsg.Sender, entry.lastMsg.Text)
								} else {
									preview = entry.lastMsg.Text
								}
								// Truncate long previews.
								if len(preview) > 60 {
									preview = preview[:57] + "..."
								}
								lbl := material.Caption(receiver.theme, preview)
								lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
								lbl.MaxLines = 1
								return lbl.Layout(gtx)
							}),
						)
					}),
				)
			})
		})
	})
}

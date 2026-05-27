package gui

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var chatMyBubbleColor color.NRGBA = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
var chatMyTextColor color.NRGBA = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
var chatTheirBubbleColor color.NRGBA = color.NRGBA{R: 0xE8, G: 0xE8, B: 0xE8, A: 0xFF}
var chatTheirTextColor color.NRGBA = color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xFF}
var chatTimestampColor color.NRGBA = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}

func (receiver *App) layoutChat(gtx layout.Context) layout.Dimensions {
	if receiver.backClick.Clicked(gtx) {
		receiver.page = PagePersonDetail
	}

	if receiver.selectedPerson < 0 || receiver.selectedPerson >= len(receiver.people) {
		receiver.page = PageHome
		return layout.Dimensions{}
	}

	var person *Person = &receiver.people[receiver.selectedPerson]

	// Handle send.
	if receiver.sendClick.Clicked(gtx) {
		receiver.sendMessage(person)
	}

	// Handle enter key in editor.
	for {
		event, ok := receiver.chatEditor.Update(gtx)
		if !ok {
			break
		}
		if _, ok := event.(widget.SubmitEvent); ok {
			receiver.sendMessage(person)
		}
	}

	var messageCount int = len(person.Messages)

	// Scroll to bottom on first render.
	if receiver.chatList.Position.Count == 0 && messageCount > 0 {
		receiver.chatList.Position.BeforeEnd = false
		receiver.chatList.Position.Offset = 0
		receiver.chatList.Position.First = messageCount
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, person.Name, true)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.chatList.Layout(gtx, messageCount, func(gtx layout.Context, index int) layout.Dimensions {
				return layoutChatBubble(gtx, receiver.theme, person.Messages[index])
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutChatInput(gtx)
		}),
	)
}

func (receiver *App) sendMessage(person *Person) {
	var text string = receiver.chatEditor.Text()
	if "" == text {
		return
	}

	person.Messages = append(person.Messages, ChatMessage{
		FromMe:    true,
		Text:      text,
		Timestamp: time.Now().Format("2006-01-02 15:04"),
	})

	receiver.chatEditor.SetText("")

	// Scroll to bottom.
	receiver.chatList.Position.BeforeEnd = false
	receiver.chatList.Position.Offset = 0
	receiver.chatList.Position.First = len(person.Messages)
}

func layoutChatBubble(gtx layout.Context, th *material.Theme, msg ChatMessage) layout.Dimensions {
	var bubbleColor color.NRGBA
	var textColor color.NRGBA
	var alignment layout.Direction

	if msg.FromMe {
		bubbleColor = chatMyBubbleColor
		textColor = chatMyTextColor
		alignment = layout.E
	} else {
		bubbleColor = chatTheirBubbleColor
		textColor = chatTheirTextColor
		alignment = layout.W
	}

	return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return alignment.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// Limit bubble width to 75% of available width.
			gtx.Constraints.Max.X = gtx.Constraints.Max.X * 3 / 4

			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Min
					rr := gtx.Dp(unit.Dp(12))
					defer clip.RRect{
						Rect: image.Rectangle{Max: size},
						SE:   rr,
						SW:   rr,
						NW:   rr,
						NE:   rr,
					}.Push(gtx.Ops).Pop()
					paint.Fill(gtx.Ops, bubbleColor)
					return layout.Dimensions{Size: size}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body2(th, msg.Text)
								lbl.Color = textColor
								return lbl.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Caption(th, msg.Timestamp)
									lbl.Color = chatTimestampColor
									if msg.FromMe {
										lbl.Color = color.NRGBA{R: 0xBB, G: 0xBB, B: 0xDD, A: 0xFF}
									}
									return lbl.Layout(gtx)
								})
							}),
						)
					})
				}),
			)
		})
	})
}

func (receiver *App) layoutChatInput(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min
			defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF})
			return layout.Dimensions{Size: size}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(receiver.theme, &receiver.chatEditor, "Type a message...")
						return editor.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(receiver.theme, &receiver.sendClick, "Send")
						return btn.Layout(gtx)
					}),
				)
			})
		}),
	)
}

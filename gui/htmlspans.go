package gui

import (
	"image/color"
	"strings"

	"golang.org/x/net/html"

	"gioui.org/font"
	"gioui.org/unit"
	"gioui.org/x/richtext"

	"github.com/microcosm-cc/bluemonday"

	"codeberg.org/reiver/go-smallmonday"
)

var sanitizePolicy *bluemonday.Policy

func init() {
	sanitizePolicy = bluemonday.NewPolicy()
	smallmonday.AllowA(sanitizePolicy)
	smallmonday.AllowB(sanitizePolicy)
	smallmonday.AllowBr(sanitizePolicy)
	smallmonday.AllowEm(sanitizePolicy)
	smallmonday.AllowI(sanitizePolicy)
	smallmonday.AllowP(sanitizePolicy)
	smallmonday.AllowStrong(sanitizePolicy)
}

func sanitizeHTML(rawHTML string) string {
	return sanitizePolicy.Sanitize(rawHTML)
}

var (
	colorText color.NRGBA = color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xFF}
	colorLink color.NRGBA = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
)

const defaultTextSize unit.Sp = 16

func htmlToSpans(rawHTML string) []richtext.SpanStyle {
	var sanitized string = sanitizeHTML(rawHTML)

	tokenizer := html.NewTokenizer(strings.NewReader(sanitized))

	var spans []richtext.SpanStyle
	var bold bool
	var italic bool
	var linkHref string
	var inLink bool

	for {
		tt := tokenizer.Next()

		switch tt {
		case html.ErrorToken:
			return spans

		case html.StartTagToken:
			tn, hasAttr := tokenizer.TagName()
			tag := string(tn)

			switch tag {
			case "b", "strong":
				bold = true
			case "i", "em":
				italic = true
			case "a":
				inLink = true
				linkHref = ""
				if hasAttr {
					for {
						key, val, more := tokenizer.TagAttr()
						if string(key) == "href" {
							linkHref = string(val)
						}
						if !more {
							break
						}
					}
				}
			case "p":
				if 0 < len(spans) {
					spans = append(spans, richtext.SpanStyle{
						Size:    defaultTextSize,
						Color:   colorText,
						Content: "\n\n",
					})
				}
			case "br":
				spans = append(spans, richtext.SpanStyle{
					Size:    defaultTextSize,
					Color:   colorText,
					Content: "\n",
				})
			}

		case html.EndTagToken:
			tn, _ := tokenizer.TagName()
			tag := string(tn)

			switch tag {
			case "b", "strong":
				bold = false
			case "i", "em":
				italic = false
			case "a":
				inLink = false
				linkHref = ""
			}

		case html.SelfClosingTagToken:
			tn, _ := tokenizer.TagName()
			if string(tn) == "br" {
				spans = append(spans, richtext.SpanStyle{
					Size:    defaultTextSize,
					Color:   colorText,
					Content: "\n",
				})
			}

		case html.TextToken:
			text := tokenizer.Text()
			content := string(text)
			if "" == content {
				continue
			}

			span := richtext.SpanStyle{
				Size:    defaultTextSize,
				Color:   colorText,
				Content: content,
			}

			if bold {
				span.Font.Weight = font.Bold
			}
			if italic {
				span.Font.Style = font.Italic
			}
			if inLink {
				span.Color = colorLink
				span.Interactive = true
				span.Set("url", linkHref)
			}

			spans = append(spans, span)
		}
	}
}

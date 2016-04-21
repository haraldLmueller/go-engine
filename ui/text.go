package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/walesey/freetype"
	"github.com/walesey/freetype/truetype"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
	"golang.org/x/image/font"
)

type TextElement struct {
	img                *ImageElement
	width, height      float64
	text               string
	textColor          color.Color
	textSize           float64
	textFont           *truetype.Font
	size               vmath.Vector2
	active             bool
	dirty              bool
	onFocusHandlers    []func()
	onBlurHandlers     []func()
	onKeyPressHandlers []func(key string, release bool)
}

func (te *TextElement) updateImage(size vmath.Vector2) {
	// Initialize the context.
	bg := image.Transparent
	c := freetype.NewContext()
	c.SetDPI(75)
	c.SetFont(te.textFont)
	c.SetFontSize(te.textSize)
	c.SetSrc(image.NewUniform(te.textColor))
	c.SetHinting(font.HintingNone)

	// Establish image dimensions and do ward wrap
	textHeight := c.PointToFixed(te.textSize)
	var width int
	var height int = int(textHeight >> 6)
	words := strings.Split(te.text, " ")
	lines := []string{""}
	lineNb := 0
	for _, word := range words {
		wordWithSpace := fmt.Sprintf("%v ", word)
		dimensions, _ := c.StringDimensions(wordWithSpace)
		width += int(dimensions.X >> 6)
		if width > int(size.X) {
			width = int(dimensions.X >> 6)
			height += int(dimensions.Y>>6) + 1
			lines = append(lines, "")
			lineNb += 1
		}
		lines[lineNb] = fmt.Sprintf("%v%v", lines[lineNb], wordWithSpace)
	}
	if te.height > 0 {
		height = int(te.height)
	}

	rgba := image.NewRGBA(image.Rect(0, 0, int(size.X), height+int(textHeight>>6)/3))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)

	// Draw the text.
	pt := freetype.Pt(0, int(textHeight>>6))
	for _, line := range lines {
		_, err := c.DrawString(line, pt)
		if err != nil {
			log.Printf("Error drawing string: %v\n", err)
			return
		}
		pt.Y += textHeight
	}

	te.img.SetImage(imaging.FlipV(rgba))
	te.img.SetWidth(float64(rgba.Bounds().Size().X))
	te.img.SetHeight(float64(rgba.Bounds().Size().Y))
}

func (te *TextElement) GetText() string {
	return te.text
}

func (te *TextElement) SetText(text string) {
	te.text = text
	te.dirty = true
}

func (te *TextElement) SetWidth(width float64) {
	te.width = width
	te.dirty = true
}

func (te *TextElement) SetHeight(height float64) {
	te.height = height
	te.dirty = true
}

func (te *TextElement) Activate() {
	if !te.active {
		for _, handler := range te.onFocusHandlers {
			handler()
		}
	}
	te.active = true
}

func (te *TextElement) Deactivate() {
	if te.active {
		for _, handler := range te.onBlurHandlers {
			handler()
		}
	}
	te.active = false
}

func (te *TextElement) Render(size, offset vmath.Vector2) vmath.Vector2 {
	useWidth := size.X
	useHeight := size.Y
	if te.width > 0 {
		useWidth = te.width
	}
	if te.height > 0 {
		useHeight = te.height
	}
	useSize := vmath.Vector2{useWidth, useHeight}
	if !useSize.ApproxEqual(te.size, 0.001) {
		te.dirty = true
	}
	te.size = useSize
	if te.dirty {
		te.dirty = false
		te.updateImage(useSize)
	}
	return te.img.Render(size, offset)
}

func (te *TextElement) Spatial() renderer.Spatial {
	return te.img.Spatial()
}

func (te *TextElement) mouseMove(position vmath.Vector2) {
	te.img.mouseMove(position)
}

func (te *TextElement) mouseClick(button int, release bool, position vmath.Vector2) {
	te.img.mouseClick(button, release, position)
}

func (te *TextElement) keyClick(key string, release bool) {
	if te.active && !release {
		textBytes := []byte(te.text)
		if key == "backspace" {
			if len(textBytes) > 0 {
				te.SetText(string(textBytes[:len(textBytes)-1]))
			}
		} else {
			te.SetText(fmt.Sprintf("%v%v", te.text, key))
		}
		for _, handler := range te.onKeyPressHandlers {
			handler(key, release)
		}
	}
}

func (te *TextElement) AddOnFocus(handler func()) {
	te.onFocusHandlers = append(te.onFocusHandlers, handler)
}

func (te *TextElement) AddOnBlur(handler func()) {
	te.onBlurHandlers = append(te.onBlurHandlers, handler)
}

func (te *TextElement) AddOnKeyPress(handler func(key string, release bool)) {
	te.onKeyPressHandlers = append(te.onKeyPressHandlers, handler)
}

func NewTextElement(text string, textColor color.Color, textSize float64, textFont *truetype.Font) *TextElement {
	textElem := &TextElement{
		img:       NewImageElement(image.NewAlpha(image.Rect(0, 0, 1, 1))),
		text:      text,
		textColor: textColor,
		textSize:  textSize,
		textFont:  textFont,
		dirty:     true,
	}
	return textElem
}
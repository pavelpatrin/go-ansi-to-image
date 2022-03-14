package ansitoimage

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/fogleman/gg"
	"github.com/leaanthony/go-ansi-parser"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

// NewConverter creates new ANSI-to-Image converter.
func NewConverter(config Config) (*Converter, error) {
	converter := &Converter{config: config}

	// Load monospaced regular font.
	monoRegularFontFace, err := loadFontFace(config.MonoRegularFontBytes, config.MonoRegularFontPoints)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load monospaced regular font")
	}
	converter.monoRegularFontFace = monoRegularFontFace

	// Load monospaced bold font.
	monoBoldFontFace, err := loadFontFace(config.MonoBoldFontBytes, config.MonoBoldFontPoints)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load monospaced bold font")
	}
	converter.monoBoldFontFace = monoBoldFontFace

	// Load monospaced oblique font.
	monoObliqueFontFace, err := loadFontFace(config.MonoObliqueFontBytes, config.MonoObliqueFontPoints)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load monospaced oblique font")
	}
	converter.monoObliqueFontFace = monoObliqueFontFace

	// Load monospaced oblique bold font.
	monoBoldObliqueFontFace, err := loadFontFace(config.MonoObliqueBoldFontBytes, config.MonoObliqueBoldFontPoints)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load monospaced oblique bold font")
	}
	converter.monoBoldObliqueFontFace = monoBoldObliqueFontFace

	// Prepare gg context.
	converter.ggContext = gg.NewContext(
		config.Padding*2+config.PageCols*config.CharWidth,
		config.Padding*2+config.PageRows*config.LineHeight,
	)

	return converter, nil
}

// Converter converts ANSI texts to image.
type Converter struct {
	// Converter config.
	config Config

	// Graphics toolkit context.
	ggContext *gg.Context

	// Graphics toolkit fonts.
	monoRegularFontFace     font.Face
	monoBoldFontFace        font.Face
	monoObliqueFontFace     font.Face
	monoBoldObliqueFontFace font.Face
}

// Parse parses ANSI texts to image.
func (c *Converter) Parse(text string) error {
	c.ggContext.SetRGB(0, 0, 0)
	c.ggContext.Clear()

	if len(text) != 0 {
		styledTexts, err := ansi.Parse(
			text,
			ansi.WithDefaultForegroundColor("37"),
			ansi.WithDefaultBackgroundColor("30"),
			ansi.WithIgnoreInvalidCodes(),
		)
		if err != nil {
			return errors.Wrap(err, "failed to parse ANSI text")
		}

		offsetXPix, offsetYPix := c.config.Padding, c.config.Padding
		nextX, nextY := 0, 0

		for _, styledText := range styledTexts {
			lines := strings.Split(styledText.Label, c.config.LineBreak)
			for lineIndex, styledTextLine := range lines {
				if lineIndex > 0 {
					nextY += 1
					nextX = 0
				}

				lineCharsCount := utf8.RuneCountInString(styledTextLine)
				lineXBasePix := float64(offsetXPix + (nextX * c.config.CharWidth))
				lineYBasePix := float64(offsetYPix + (nextY * c.config.LineHeight))

				if styledText.BgCol != nil {
					c.ggContext.SetRGB(
						float64(styledText.BgCol.Rgb.R)/255,
						float64(styledText.BgCol.Rgb.G)/255,
						float64(styledText.BgCol.Rgb.B)/255,
					)

					rectWidth := lineCharsCount
					if lineIndex < len(lines)-1 {
						rectWidth = c.config.PageCols - nextX
					}

					c.ggContext.DrawRectangle(
						lineXBasePix, lineYBasePix,
						float64(rectWidth*c.config.CharWidth), float64(c.config.LineHeight),
					)
					c.ggContext.Fill()
				}

				// Select font face.
				switch {
				case styledText.Bold() && styledText.Italic():
					c.ggContext.SetFontFace(c.monoBoldObliqueFontFace)
				case styledText.Italic():
					c.ggContext.SetFontFace(c.monoObliqueFontFace)
				case styledText.Bold():
					c.ggContext.SetFontFace(c.monoBoldFontFace)
				default:
					c.ggContext.SetFontFace(c.monoRegularFontFace)
				}

				// Additional styles unfortunately not supported by toolkit.
				// See: https://github.com/fogleman/gg/issues/82.
				switch {
				case styledText.Style&ansi.Underlined > 0:
				case styledText.Style&ansi.Strikethrough > 0:
				}

				if styledText.Style == ansi.Bold {
					c.ggContext.SetFontFace(c.monoBoldFontFace)
				} else {
					c.ggContext.SetFontFace(c.monoRegularFontFace)
				}

				if styledText.FgCol != nil {
					c.ggContext.SetRGB(
						float64(styledText.FgCol.Rgb.R)/255,
						float64(styledText.FgCol.Rgb.G)/255,
						float64(styledText.FgCol.Rgb.B)/255,
					)
				} else {
					c.ggContext.SetRGB(1.0, 1.0, 1.0)
				}

				for runeIndex, runeValue := range []rune(styledTextLine) {
					c.ggContext.DrawString(
						string(runeValue),
						lineXBasePix+(float64(runeIndex*c.config.CharWidth)),
						lineYBasePix+(float64(c.config.LineHeight-c.config.LineShift)),
					)
				}

				nextX += lineCharsCount
			}
		}
	}

	return nil
}

// ToPNG encodes parsed text to PNG image.
func (c *Converter) ToPNG() ([]byte, error) {
	buffer := bytes.Buffer{}

	if err := c.ggContext.EncodePNG(&buffer); err != nil {
		return nil, errors.Wrap(err, "failed to encode PNG")
	}

	return buffer.Bytes(), nil
}

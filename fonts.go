package ansitoimage

import _ "embed"

import (
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

//go:embed fonts/DejaVuSansMono.ttf
var monoRegularFontData []byte

//go:embed fonts/DejaVuSansMono-Bold.ttf
var monoBoldFontData []byte

//go:embed fonts/DejaVuSansMono-Oblique.ttf
var monoObliqueFontData []byte

//go:embed fonts/DejaVuSansMono-BoldOblique.ttf
var monoBoldObliqueFontData []byte

// loadFontFace loads truetype font face from specified font data.
func loadFontFace(data []byte, points float64) (font.Face, error) {
	trueTypeFont, err := truetype.Parse(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse truetype font")
	}

	return truetype.NewFace(trueTypeFont, &truetype.Options{Size: points}), nil
}

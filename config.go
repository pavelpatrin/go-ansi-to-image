package ansitoimage

// Config represents converter configuration.
type Config struct {
	// Text settings.
	CharWidth  int
	LineHeight int
	LineShift  int
	LineBreak  string

	// Page settings.
	PageCols int
	PageRows int
	Padding  int

	// Font settings.
	MonoRegularFontBytes      []byte
	MonoRegularFontPoints     float64
	MonoBoldFontBytes         []byte
	MonoBoldFontPoints        float64
	MonoObliqueFontBytes      []byte
	MonoObliqueFontPoints     float64
	MonoObliqueBoldFontBytes  []byte
	MonoObliqueBoldFontPoints float64
}

// DefaultConfig represents default configuration for converter.
var DefaultConfig = Config{
	// Text settings.
	CharWidth:  10,
	LineHeight: 19,
	LineShift:  4,
	LineBreak:  "\n",

	// Page settings.
	PageCols: 80,
	PageRows: 24,
	Padding:  10,

	// Font settings.
	MonoRegularFontBytes:      monoRegularFontData,
	MonoRegularFontPoints:     16.0,
	MonoBoldFontBytes:         monoBoldFontData,
	MonoBoldFontPoints:        16.0,
	MonoObliqueFontBytes:      monoObliqueFontData,
	MonoObliqueFontPoints:     16.0,
	MonoObliqueBoldFontBytes:  monoBoldObliqueFontData,
	MonoObliqueBoldFontPoints: 16.0,
}

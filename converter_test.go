package ansitoimage

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestConverter tests Converter type.
func TestConverter(t *testing.T) {
	assert := require.New(t)

	ansiBytes, err := ioutil.ReadFile(path.Join("tests", "ansi-mc.txt"))
	assert.NoError(err, "ReadFile should return no error")

	pngBytes, err := ioutil.ReadFile(path.Join("tests", "ansi-mc.png"))
	assert.NoError(err, "ReadFile should return no error")

	converter, err := NewConverter(DefaultConfig)
	assert.NoError(err, "NewConverter should return no error")

	err = converter.Parse(string(ansiBytes))
	assert.NoError(err, "Parse should return no error")

	newPNGBytes, err := converter.ToPNG()
	assert.NoError(err, "ToPNG should return no error")

	assert.Equal(pngBytes, newPNGBytes)
}

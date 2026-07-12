package manager

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseResolution covers APP-103 RC2: "WxH" parses to positive dims; empty or
// malformed specs report ok=false (⇒ no downscaling).
func TestParseResolution(t *testing.T) {
	cases := []struct {
		spec    string
		w, h    int
		ok      bool
	}{
		{"1920x1080", 1920, 1080, true},
		{"  1280 x 720 ", 1280, 720, true}, // trimmed
		{"1920X1080", 1920, 1080, true},    // case-insensitive separator
		{"", 0, 0, false},                  // empty = no resize
		{"1080", 0, 0, false},              // missing dimension
		{"axb", 0, 0, false},               // non-numeric
		{"0x0", 0, 0, false},               // non-positive
		{"-1x10", 0, 0, false},             // negative
	}
	for _, c := range cases {
		w, h, ok := parseResolution(c.spec)
		assert.Equal(t, c.ok, ok, "parseResolution(%q) ok", c.spec)
		if c.ok {
			assert.Equal(t, c.w, w, "parseResolution(%q) w", c.spec)
			assert.Equal(t, c.h, h, "parseResolution(%q) h", c.spec)
		}
	}
}

// TestFfmpegCompressArgs covers APP-103 RC2: the scale filter (when present) sits
// before the codec flags, extra args sit before the output, and the output is
// always last. Unknown target extensions are rejected.
func TestFfmpegCompressArgs(t *testing.T) {
	// No scale, no extra: -y -i in -q:v N out
	a, err := ffmpegCompressArgs("in.jpg", "out.jpg", "jpg", 80, "", nil)
	require.NoError(t, err)
	assert.Equal(t, "out.jpg", a[len(a)-1], "output is last")
	assert.NotContains(t, strings.Join(a, " "), "-vf", "no scale ⇒ no -vf")
	assert.Contains(t, strings.Join(a, " "), "-q:v")

	// With scale + extra: -vf before -q:v; extra before output.
	scale := scaleFilter(1920, 1080)
	a, err = ffmpegCompressArgs("in.jpg", "out.jpg", "jpg", 75, scale, []string{"-pix_fmt", "yuvj420p"})
	require.NoError(t, err)
	joined := strings.Join(a, " ")
	assert.Contains(t, joined, "-vf "+scale)
	assert.Less(t, indexOf(a, "-vf"), indexOf(a, "-q:v"), "-vf before -q:v")
	assert.Less(t, indexOf(a, "-pix_fmt"), indexOf(a, "out.jpg"), "extra before output")
	assert.Equal(t, "out.jpg", a[len(a)-1])

	// webp / png branches.
	a, err = ffmpegCompressArgs("in.png", "out.webp", "webp", 80, "", nil)
	require.NoError(t, err)
	assert.Contains(t, strings.Join(a, " "), "-quality")

	a, err = ffmpegCompressArgs("in.png", "out.png", "png", 80, "", nil)
	require.NoError(t, err)
	assert.Contains(t, strings.Join(a, " "), "-compression_level")

	// Unknown extension is rejected.
	_, err = ffmpegCompressArgs("in.gif", "out.gif", "gif", 80, "", nil)
	assert.Error(t, err)
}

func indexOf(s []string, v string) int {
	for i := range s {
		if s[i] == v {
			return i
		}
	}
	return -1
}

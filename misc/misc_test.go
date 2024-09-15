package misc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDurationSexagesimal(t *testing.T) {
	expected := "0:22:57.628452"
	actual := FormatDurationSexagesimal(time.Duration(1377628452000))
	if actual != expected {
		assert.Equal(t, expected, actual)
	}
}

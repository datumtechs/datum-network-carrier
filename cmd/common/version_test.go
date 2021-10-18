package common

import (
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"testing"
)

func TestVersion_DefaultValue(t *testing.T) {
	v := Version()
	t.Logf("version:%s", v)
	assert.Equal(t, "Moments ago", v)
}

func TestVersion_SetValue(t *testing.T) {
	buildDate = "{DATE}"
	v := Version()
	t.Logf("version:%s", v)
	require.NotEqual(t, "Moments ago", v)
}

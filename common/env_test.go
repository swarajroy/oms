package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	result := EnvString("SAMPLE", "no_sample")
	assert.Equal(t, "sample", result)
}

package common

import (
	"testing"

	_ "github.com/joho/godotenv/autoload" // this is a must to read the .env
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	result := EnvString("SAMPLE", "no_sample")
	assert.Equal(t, "sample", result)
}

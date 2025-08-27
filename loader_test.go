package goenvloader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSimpleStruct(t *testing.T) {
	os.Setenv("PORT", "8080")
	var cfg struct {
		Port int `env:"PORT"`
	}

	if err := Load(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 8080, cfg.Port)
}

func TestLoadSimpleStructWithDefaultTags(t *testing.T) {
	os.Setenv("PORT", "")
	var cfg struct {
		Port int `env:"PORT" default:"8080"`
	}

	if err := Load(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 8080, cfg.Port)
}

func TestLoadSimpleStructWithRequiredTags(t *testing.T) {
	os.Setenv("PORT", "")
	var cfg struct {
		Port int `env:"PORT" default:"8080" required:"true"`
	}

	err := Load(&cfg)
	assert.Contains(t, err.Error(), "required environment variable PORT is empty")
}

func TestInvalidTypeInt(t *testing.T) {
	os.Setenv("PORT", "-8080")
	var cfg struct {
		Port int `env:"PORT"`
	}

	err := Load(&cfg)
	assert.Contains(t, err.Error(), "integer must be greater than 0")
}

func TestInvalidTypeString(t *testing.T) {
	os.Setenv("NUM_URL", "")
	var cfg struct {
		NumberingPlanURL string `env:"NUM_URL"`
	}

	err := Load(&cfg)
	assert.NoError(t, err)
	assert.Equal(t, "", cfg.NumberingPlanURL)
}

func TestUnsupportedType(t *testing.T) {
	os.Setenv("PANIC", "true")
	var cfg struct {
		PanicOnErr bool `env:"PANIC"`
	}

	err := Load(&cfg)
	assert.Contains(t, err.Error(), "unsupported field type")
}

func TestNoTagErr(t *testing.T) {
	os.Setenv("PORT", "8080")
	var cfg struct {
		Port int
	}
	err := Load(&cfg)
	assert.Contains(t, err.Error(), "tag is empty")
}

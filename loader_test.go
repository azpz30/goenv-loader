package goenvloader

import (
	"os"
	"reflect"
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

func TestLoadMultipleFieldsStruct(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("POSTGRES_USER", "test_user")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("POSTGRES_PASSWORD", "test_password")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRESS_NOSSL", "true")
	os.Setenv("KAFKA_URL", "kafka1:9092,kafka2:9092")

	var cfg struct {
		Port             int    `env:"PORT"`
		PostgresUser     string `env:"POSTGRES_USER"`
		DbName           string `env:"DB_NAME"`
		PostgresPassword string `env:"POSTGRES_PASSWORD"`
		PostgresHost     string `env:"POSTGRES_HOST"`
		PostgresNoSSL    string `env:"POSTGRESS_NOSSL"`
		KafkaURL         string `env:"KAFKA_URL"`
	}

	if err := Load(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate the loaded values
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "test_user", cfg.PostgresUser)
	assert.Equal(t, "test_db", cfg.DbName)
	assert.Equal(t, "test_password", cfg.PostgresPassword)
	assert.Equal(t, "localhost", cfg.PostgresHost)
	assert.Equal(t, "true", cfg.PostgresNoSSL)
	assert.Equal(t, "kafka1:9092,kafka2:9092", cfg.KafkaURL)

	// Type checking
	v := reflect.ValueOf(cfg)
	tp := v.Type()

	expectedTypes := map[string]reflect.Kind{
		"Port":             reflect.Int,
		"PostgresUser":     reflect.String,
		"DbName":           reflect.String,
		"PostgresPassword": reflect.String,
		"PostgresHost":     reflect.String,
		"PostgresNoSSL":    reflect.String,
		"KafkaURL":         reflect.String,
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := tp.Field(i).Name
		expectedType, ok := expectedTypes[fieldName]
		if !ok {
			t.Errorf("unexpected field %s", fieldName)
			continue
		}
		if field.Kind() != expectedType {
			t.Errorf("expected field %s to be of type %s, got %s", fieldName, expectedType, field.Kind())
		}
	}
}

func TestErrorEnvEmpty(t *testing.T) {
	os.Setenv("PORT", "")

	var cfg struct {
		Port int `env:"PORT"`
	}

	err := Load(&cfg)
	assert.NoError(t, err)
}

func TestErrorWrongType(t *testing.T) {
	os.Setenv("PORT", "cannot be a string")

	var cfg struct {
		Port int `env:"PORT"`
	}

	err := Load(&cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parsing \"cannot be a string\": invalid syntax")
}

func TestLoadNestedStruct(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")

	var cfg struct {
		DB struct {
			Host string `env:"DB_HOST"`
			Port int    `env:"DB_PORT"`
		}
	}

	if err := Load(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, 5432, cfg.DB.Port)
}

func TestLoadNestedStructAndNormalVar(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("KAFKA_URL", "kafka1:9092,kafka2:9092")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("POSTGRES_USER", "test_user")

	var cfg struct {
		Port     int    `env:"PORT"`
		KafkaURL string `env:"KAFKA_URL"`
		DB       struct {
			Host string `env:"DB_HOST"`
			Port int    `env:"DB_PORT"`
		}
		PGUser string `env:"POSTGRES_USER"`
	}

	if err := Load(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "kafka1:9092,kafka2:9092", cfg.KafkaURL)
	assert.Equal(t, "test_user", cfg.PGUser)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, 5432, cfg.DB.Port)
}

func TestIncorrectInput(t *testing.T) {
	var cfg string
	err := Load(&cfg)
	assert.Contains(t, err.Error(), "expected pointer to struct")
}

func TestLoadNestedStructAndNormalVarType(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("KAFKA_URL", "kafka1:9092,kafka2:9092")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("POSTGRES_USER", "test_user")

	type cfg struct {
		Port     int    `env:"PORT"`
		KafkaURL string `env:"KAFKA_URL"`
		DB       struct {
			Host string `env:"DB_HOST"`
			Port int    `env:"DB_PORT"`
		}
		PGUser string `env:"POSTGRES_USER"`
	}

	var c cfg
	if err := Load(&c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 8080, c.Port)
	assert.Equal(t, "kafka1:9092,kafka2:9092", c.KafkaURL)
	assert.Equal(t, "test_user", c.PGUser)
	assert.Equal(t, "localhost", c.DB.Host)
	assert.Equal(t, 5432, c.DB.Port)
}

func TestLoadNestedStructRequiredTag(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("POSTGRES_USER", "test_user")

	type cfg struct {
		Port     int    `env:"PORT"`
		KafkaURL string `env:"KAFKA_URLBB"`
		DB       struct {
			Host string `env:"DB_HOST" required:"true"`
			Port int    `env:"DB_PORT" required:"true"`
		}
		PGUser string `env:"POSTGRES_USER" required:"true"`
	}

	var c cfg
	if err := Load(&c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that the values were loaded correctly
	assert.Equal(t, 8080, c.Port)
	assert.Equal(t, "", c.KafkaURL)
	assert.Equal(t, "test_user", c.PGUser)
	assert.Equal(t, "localhost", c.DB.Host)
	assert.Equal(t, 5432, c.DB.Port)
}

func TestLoadNestedStructRequiredTagErr1(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("POSTGRES_USER", "test_user")

	type cfg struct {
		Port     int    `env:"PORT"`
		KafkaURL string `env:"KAFKA_URL"`
		DB       struct {
			Host string `env:"DB_HOST1" required:"true"`
			Port int    `env:"DB_PORT" required:"true"`
		}
		PGUser string `env:"POSTGRES_USER" required:"true"`
	}

	var c cfg
	err := Load(&c)
	assert.Contains(t, err.Error(), "required environment variable DB_HOST1 is empty")
}

func TestLoadIncorrectTypeHandling(t *testing.T) {
	os.Setenv("PORT", "not_an_integer")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("POSTGRES_USER", "test_user")

	type cfg struct {
		Port     int    `env:"PORT" required:"true"`
		KafkaURL string `env:"KAFKA_URL"`
		DB       struct {
			Host string `env:"DB_HOST" required:"true"`
			Port int    `env:"DB_PORT" required:"true"`
		}
		PGUser string `env:"POSTGRES_USER" required:"true"`
	}

	var c cfg
	err := Load(&c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "strconv.Atoi: parsing \"not_an_integer\": invalid syntax")
}

func TestLoadNestedStructDefaultTag(t *testing.T) {
	os.Setenv("PORT2", "8080")
	os.Setenv("POSTGRES_USER", "test_user")

	type cfg struct {
		Port int `env:"PORT2" required:"true"`
		DB   struct {
			Host string `env:"DB_HOST2"`
			Port int    `env:"DB_PORT2"`
		}
		PGUser string `env:"POSTGRES_USER" required:"true"`
	}

	var c cfg
	err := Load(&c)
	assert.NoError(t, err)
	assert.Equal(t, "", c.DB.Host)
	assert.Equal(t, int(0), c.DB.Port)
}

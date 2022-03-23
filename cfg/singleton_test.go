package cfg_test

import (
	"os"
	"testing"

	"github.com/infiniteloopcloud/go/cfg"
)

func TestSingletonEnv(t *testing.T) {
	os.Setenv("TEST1", "VALUE1")

	v1 := cfg.Get("TEST1")
	check(t, "VALUE1", v1)
}

func TestSingletonFile(t *testing.T) {
	os.Setenv("data1", "salalala")
	os.Setenv("CFG_LOCATION", "./test_config.json")

	cfg.Rebuild()

	v1 := cfg.Get("data1")
	check(t, "value1", v1)
}

func TestSingletonFileErrorLog(t *testing.T) {
	os.Setenv("data31", "salalala")
	os.Setenv("CFG_DEBUG", "true")

	cfg.Rebuild()

	v1 := cfg.Get("data31")
	check(t, "salalala", v1)
}

func check(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("String should be %s, instead of %s", expected, actual)
	}
}

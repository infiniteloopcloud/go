//nolint:gosec
package cfg

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	var config = map[string]string{
		"data1": "value1",
		"data2": "value2",
	}
	d1, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	filename := "./test-data"
	err = os.WriteFile(filename, d1, 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := New(Opts{
		Path:      filename,
		HotReload: true,
		Infof: func(ctx context.Context, format string, args ...interface{}) {
			t.Logf(format+"\n", args...)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	data1 := cfg.Get("data1")
	if data1 != "value1" {
		t.Errorf("Data1 should be 'value1', instead of %s", data1)
	}

	time.Sleep(500 * time.Millisecond)

	config["data3"] = "value3"
	d1, err = json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filename, d1, 0644)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(500 * time.Millisecond)

	if data3 := cfg.Get("data3"); data3 != "value3" {
		t.Errorf("Data3 should be 'value3', instead of %s", data1)
	}
}

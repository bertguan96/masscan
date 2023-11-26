package masscan

import (
	"context"
	"os"
	"testing"
)

func TestNewScanner(t *testing.T) {
	path, err2 := os.Getwd()
	if err2 != nil {
		return
	}
	t.Log(path)
	scanner, err := NewScanner(context.Background(),
		WithRoot(),
		WithBinaryPath(path+"/bin/masscan"),
		WithPort("1-65535"),
		WithOutputJson("scanRes.json"),
		WithRetryTime(3),
		WithRate(40000),
		WithRandomizeHosts(),
		WithTarget([]string{"114.55.97.220", "43.135.11.122", "101.42.164.23", "43.153.24.244"}))
	if err != nil {
		t.Log(err)
		return
	}
	if err := scanner.Run(); err != nil {
		t.Log(err)
		return
	}
}

package logger

import "testing"

func TestDebugf(t *testing.T) {
	Errorf("Resolve TCP Address failed: %s", "xx")
}

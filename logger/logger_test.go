package logger_test

import (
	"credo/logger"
	"log"
	"testing"
)

func Test_Get_ReturnsNonNil(t *testing.T) {
	l := logger.Get()
	if l == nil {
		t.Fatal("Get() returned nil")
	}
}

func Test_Get_ReturnsSameInstance(t *testing.T) {
	a := logger.Get()
	b := logger.Get()
	if a != b {
		t.Error("Get() should return the same singleton instance")
	}
}

func Test_Get_IsLogLogger(t *testing.T) {
	l := logger.Get()
	// Verify it is a *log.Logger by asserting the type.
	if _, ok := any(l).(*log.Logger); !ok {
		t.Error("Get() should return *log.Logger")
	}
}

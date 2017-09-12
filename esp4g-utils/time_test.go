package utils

import (
	"testing"
	"time"
	"github.com/golang/protobuf/ptypes/duration"
)

func TestConvertDuration(t *testing.T) {
	sec := time.Duration(1) * time.Second
	ms := time.Duration(1) * time.Millisecond

	if actual := ConvertDuration(sec); !(actual.Seconds == 1 && actual.Nanos == 0) {
		t.Error("expected sec but actual", actual)
	}
	if actual := ConvertDuration(ms); !(actual.Seconds == 0 && actual.Nanos == 1000000) {
		t.Error("expected ms but actual", actual)
	}
}

func TestConvertProtoDuration(t *testing.T) {
	sec := &duration.Duration{Seconds: 1, Nanos: 0}
	ms := &duration.Duration{Seconds: 0, Nanos: 1000000}

	if actual := ConvertProtoDuration(sec); actual.Seconds() != 1 {
		t.Error("expected sec but actual", actual)
	}
	if actual := ConvertProtoDuration(ms); actual.Seconds() != 0.001 {
		t.Error("expected ms but actual", actual)
	}
}

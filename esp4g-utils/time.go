package utils

import (
	"github.com/golang/protobuf/ptypes/duration"
	"time"
)

func ConvertProtoDuration(d *duration.Duration) time.Duration {
	if d == nil {
		return time.Duration(-1)
	}
	return time.Duration(d.Seconds) * time.Second + time.Duration(d.Nanos)
}

func ConvertDuration(d time.Duration) *duration.Duration {
	return &duration.Duration{
		Seconds: d.Nanoseconds() / time.Second.Nanoseconds(),
		Nanos: int32(d.Nanoseconds() % time.Second.Nanoseconds()),
	}
}

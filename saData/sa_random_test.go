package saData

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRandomStrReturnsFixedLengthURLSafeValue(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandomStr()
		if len(s) != 20 {
			t.Fatalf("RandomStr length = %d, want 20: %q", len(s), s)
		}
		for _, c := range s {
			if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_", c) {
				t.Fatalf("RandomStr contains non URL-safe character %q in %q", c, s)
			}
		}
	}
}

func TestRandomStrDoesNotExposeUnixMilliSuffix(t *testing.T) {
	before := time.Now().UnixMilli()
	s := RandomStr()
	after := time.Now().UnixMilli()

	if len(s) < 13 {
		t.Fatalf("RandomStr length = %d, want at least 13: %q", len(s), s)
	}
	suffix := s[len(s)-13:]
	v, err := strconv.ParseInt(suffix, 10, 64)
	if err == nil && v >= before && v <= after {
		t.Fatalf("RandomStr exposes unix millisecond suffix %q in %q", suffix, s)
	}
}

func TestRandomInt64UsesInt64Scale(t *testing.T) {
	var max int64
	for i := 0; i < 64; i++ {
		v := RandomInt64()
		if v < 0 {
			t.Fatalf("RandomInt64 returned negative value: %d", v)
		}
		if v > max {
			max = v
		}
	}

	if max < 1<<40 {
		t.Fatalf("RandomInt64 max from 64 samples = %d, want value using wider int64 range", max)
	}
}

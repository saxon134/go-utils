package saData

import "testing"

func TestIdToCodeLength13DoesNotCollide(t *testing.T) {
	code1 := IdToCode(1, 13)
	code2 := IdToCode(33, 13)

	if code1 == "" || code2 == "" {
		t.Fatalf("codes should be representable: %q %q", code1, code2)
	}
	if code1 == code2 {
		t.Fatalf("IdToCode collision: 1 and 33 both encoded to %q", code1)
	}
	if got := CodeToId(code1, 13); got != 1 {
		t.Fatalf("CodeToId(%q) = %d, want 1", code1, got)
	}
	if got := CodeToId(code2, 13); got != 33 {
		t.Fatalf("CodeToId(%q) = %d, want 33", code2, got)
	}
}

func TestIdToCodeRoundTripsSupportedLengths(t *testing.T) {
	for _, length := range []int{3, 5, 8, 13} {
		seen := map[string]int64{}
		for id := int64(1); id <= 200; id++ {
			code := IdToCode(id, length)
			if code == "" {
				t.Fatalf("IdToCode(%d, %d) returned empty", id, length)
			}
			if prev, ok := seen[code]; ok {
				t.Fatalf("IdToCode collision with length %d: %d and %d both encoded to %q", length, prev, id, code)
			}
			seen[code] = id

			if got := CodeToId(code, length); got != id {
				t.Fatalf("CodeToId(%q, %d) = %d, want %d", code, length, got, id)
			}
		}
	}
}

func TestIdToCharWithMaxSourceDoesNotCollide(t *testing.T) {
	char1 := IdToCharWithSource(14, 3, MaxSource)
	char2 := IdToCharWithSource(868, 3, MaxSource)

	if char1 == char2 {
		t.Fatalf("IdToCharWithSource collision: 14 and 868 both encoded to %q", char1)
	}
	if got := CharToIdWithSource(char1, MaxSource); got != 14 {
		t.Fatalf("CharToIdWithSource(%q) = %d, want 14", char1, got)
	}
	if got := CharToIdWithSource(char2, MaxSource); got != 868 {
		t.Fatalf("CharToIdWithSource(%q) = %d, want 868", char2, got)
	}
}

func TestIdCharWithMaxSourceRoundTripsWithoutCollisions(t *testing.T) {
	seen := map[string]int64{}
	for id := int64(1); id <= 5000; id++ {
		char := IdToCharWithSource(id, 3, MaxSource)
		if prev, ok := seen[char]; ok {
			t.Fatalf("IdToCharWithSource collision: %d and %d both encoded to %q", prev, id, char)
		}
		seen[char] = id

		if got := CharToIdWithSource(char, MaxSource); got != id {
			t.Fatalf("CharToIdWithSource(%q) = %d, want %d", char, got, id)
		}
	}
}

package fallback_test

import (
	. "github.com/jenchik/fallback"

	"testing"
)

func TestClassicFallback(t *testing.T) {
	f := NewClassicFallback()
	testFallback(t, f)
}

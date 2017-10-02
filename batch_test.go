package fallback_test

import (
	. "github.com/jenchik/fallback"

	"testing"
)

func TestBatchFallback(t *testing.T) {
	f := NewBatchFallback(1)
	testFallback(t, f)
}

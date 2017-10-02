package fallback_test

import (
	"math/rand"
	"testing"
	"time"
)

var (
	keys        []string
	newKeys     []string
	dispersion  []string
	dispersions [][]string
)

func init() {
	keys = make([]string, 0, items)
	newKeys = make([]string, 0, items)
	var key string
	for i := 0; i < items; i++ {
		key = randString(lenKey)
		keys = append(keys, key)
		newKeys = append(newKeys, randString(lenKey))
	}

	dispersion = newDispersion(dispersion, steps)
	for i := 0; i < benchWorks; i++ {
		s := make([]string, 0, steps)
		dispersions = append(dispersions, newDispersion(s, steps))
	}
}

func newDispersion(in []string, x int) []string {
	rand.Seed(time.Now().UTC().UnixNano())
	var key string
	var n int
	for i := 0; i < x; i++ {
		n = rand.Intn(items)
		if rand.Intn(2) == 1 {
			key = newKeys[n]
		} else {
			key = keys[n]
		}
		in = append(in, key)
	}
	return in
}

func initMap() map[string]int {
	m := make(map[string]int, items*2)
	for i, key := range keys {
		m[key] = i
	}
	return m
}

func TestDescription(t *testing.T) {
	if testing.Verbose() {
		t.Log("Prepare elements:", items)
		t.Log("Length key:", lenKey)
		t.Log("Steps:", steps)
		t.Log("Bench workers:", benchWorks)
		t.Log("Test workers:", testWorks)
	}
}

package common

import "testing"

func TestEffectiveTextMemory(t *testing.T) {

	t.Run("Cache", testDefault)
	t.Run("ClearCache", testClearCache)
	t.Run("NoCache", testNoCache)
}

func testDefault(t *testing.T) {

	etm := NewEffectiveTextMemory()

	etm.Write("test/a.txt")
	etm.Write("test/a.txt")
	etm.Write("test/b.txt")

	if etm.CountAll != 2 {
		t.Errorf("Wrong number of all entries: %d.", etm.CountAll)
	}
	if etm.CountCollisions != 0 {
		t.Errorf("Wrong number of collisions: %d.", etm.CountCollisions)
	}
	assertContainsText(t, etm, "test/a.txt", true)
	assertContainsText(t, etm, "test/b.txt", true)
	assertContainsText(t, etm, "test/something/c.txt", false)
}

func testClearCache(t *testing.T) {

	etm := NewEffectiveTextMemory()

	etm.Write("test/a.txt")
	etm.ClearCache()
	etm.Write("test/a.txt")
	etm.Write("test/b.txt")

	if etm.CountAll != 3 {
		t.Errorf("Wrong number of all entries: %d.", etm.CountAll)
	}
	if etm.CountCollisions != 1 {
		t.Errorf("Wrong number of collisions: %d.", etm.CountCollisions)
	}
}

func testNoCache(t *testing.T) {

	etm := NewEffectiveTextMemory()
	etm.UseCache = false

	etm.Write("test/a.txt")
	etm.Write("test/a.txt")
	etm.Write("test/b.txt")

	if etm.CountAll != 3 {
		t.Errorf("Wrong number of all entries: %d.", etm.CountAll)
	}
	if etm.CountCollisions != 1 {
		t.Errorf("Wrong number of collisions: %d.", etm.CountCollisions)
	}
}

func assertContainsText(t *testing.T, etm *EffectiveTextMemory, text string, shouldContain bool) {

	if shouldContain && !etm.ContainsText(text) {
		t.Errorf("Should contain text: \"%s\".", text)
	} else if !shouldContain && etm.ContainsText(text) {
		t.Errorf("Should not contain text: \"%s\".", text)
	}
}

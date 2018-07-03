package dal

import (
	"testing"
)

func TestMemoryDatabase(t *testing.T) {

	t.Run("MemoryDatabase_AddFingerprint", testMemoryDatabaseAddFingerprint)
	t.Run("MemoryDatabase_AddFingerprints", testMemoryDatabaseAddFingerprints)
	t.Run("MemoryDatabase_AddNamePair", testMemoryDatabaseAddNamePair)
	t.Run("MemoryDatabase_Clear", testMemoryDatabaseClear)
	t.Run("MemoryDatabase_LoadNamesFromFingerprints", testMemoryDatabaseLoadNamesFromFingerprints)
}

func testMemoryDatabaseAddFingerprint(t *testing.T) {

	memoryDatabase := NewMemoryDatabase()
	testDatabaseAddFingerprint(t, memoryDatabase)
}

func testMemoryDatabaseAddFingerprints(t *testing.T) {

	memoryDatabase := NewMemoryDatabase()
	testDatabaseAddFingerprints(t, memoryDatabase)
}

func testMemoryDatabaseAddNamePair(t *testing.T) {

	memoryDatabase := NewMemoryDatabase()
	testDatabaseAddNamePair(t, memoryDatabase)
}

func testMemoryDatabaseClear(t *testing.T) {

	memoryDatabase := NewMemoryDatabase()
	testDatabaseClear(t, memoryDatabase)
}

func testMemoryDatabaseLoadNamesFromFingerprints(t *testing.T) {

	memoryDatabase := NewMemoryDatabase()
	testDatabaseLoadNamesFromFingerprints(t, memoryDatabase)
}

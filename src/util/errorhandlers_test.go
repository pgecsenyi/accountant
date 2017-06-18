package util

import (
	"os"
	"testing"
)

func TestCheckErr(t *testing.T) {

	testHelper.AssertPanic(t, func() {
		CheckErr(os.ErrExist, "")
	})
}

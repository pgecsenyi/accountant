package util

import "testing"

func TestCompareByteSlicesBothNil(t *testing.T) {

	var b1 []byte
	var b2 []byte

	result := CompareByteSlices(b1, b2)

	if !result {
		t.Errorf("The two byte slices should be equal.")
	}
}

func TestCompareByteSlicesOneNil(t *testing.T) {

	b1 := []byte{}
	var b2 []byte

	result := CompareByteSlices(b1, b2)

	if result {
		t.Errorf("The two byte slices should not be equal.")
	}
}

func TestCompareByteSlicesBothEmpty(t *testing.T) {

	b1 := []byte{}
	b2 := []byte{}

	result := CompareByteSlices(b1, b2)

	if !result {
		t.Errorf("The two byte slices should be equal.")
	}
}

func TestCompareByteSlicesOneEmpty(t *testing.T) {

	b1 := []byte{}
	b2 := []byte{31, 215}

	result := CompareByteSlices(b1, b2)

	if result {
		t.Errorf("The two byte slices should not be equal.")
	}
}

func TestCompareByteSlicesEqual(t *testing.T) {

	b1 := []byte{31, 215}
	b2 := []byte{31, 215}

	result := CompareByteSlices(b1, b2)

	if !result {
		t.Errorf("The two byte slices should be equal.")
	}
}

func TestCompareByteSlicesNotEqual(t *testing.T) {

	b1 := []byte{31, 125}
	b2 := []byte{31, 215}

	result := CompareByteSlices(b1, b2)

	if result {
		t.Errorf("The two byte slices should not be equal.")
	}
}

package checksum

import (
	"io"
	"os"
	"util"
)

// CalculateChecksumForFile Calculates checksum for the given file using the given algorithm.
func CalculateChecksumForFile(filename string, algorithm string) []byte {

	file, err := os.Open(filename)
	util.CheckErr(err, "Cannot read file "+filename+".")
	defer file.Close()

	calculator := CreateCalculator(algorithm)
	io.Copy(calculator, file)
	checksum := calculator.Sum(nil)[:]

	return checksum
}

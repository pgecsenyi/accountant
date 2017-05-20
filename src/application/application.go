package application

import (
	"bll"
	"checksum"
	"flag"
	"log"
	"os"
)

const modeCalculate = "calculate"
const modeCompare = "compare"
const modeVerify = "verify"

// Application Contains main application logic.
type Application struct {
	config configuration
}

type configuration struct {
	mode            string
	algorithm       string
	sourceDirectory string
	outputChecksum  string
	inputChecksum   string
	outputNames     string
	basePath        string
}

// Initialize Initializes the application.
func (app *Application) Initialize() {

	defaultConfig := configuration{modeCalculate, checksum.SHA1, "data/input", "data/checksum.csv", "", "", ""}
	app.parseCommandLineArguments(defaultConfig)
	app.verifyConfiguration()
	app.execute()
}

func (app *Application) parseCommandLineArguments(defaultConfig configuration) {

	doCalculate := flag.Bool(
		"calculate",
		false,
		"Calculate checksums for a directory and store the results in a CSV.")
	doCompare := flag.Bool(
		"compare",
		false,
		"Compare stored checksums with the checksums of the files in the given directory and store filename matches.")
	doVerify := flag.Bool(
		"verify",
		false,
		"Verify checksums for the files listed in the given CSV.")
	sourceDirectory := flag.String(
		"source",
		defaultConfig.sourceDirectory,
		"The source directory for which the checksums will be calculated (or will be compared).")
	algorithm := flag.String("alg", defaultConfig.algorithm, "The algorithm used to calculate new checksums.")
	outputChecksum := flag.String(
		"outchk",
		defaultConfig.outputChecksum,
		"The name of the output CSV file containing checksums.")
	inputChecksum := flag.String(
		"inchk",
		defaultConfig.inputChecksum,
		"The name of the input CSV containing checksums.")
	outputNames := flag.String(
		"outnames",
		defaultConfig.outputNames,
		"The name of the output containing new file name and old filename pairs.")
	basePath := flag.String(
		"bp",
		defaultConfig.basePath,
		"The first part of the path that will not be stored in the output. By default it will be set to source.")

	flag.Parse()

	mode := modeCalculate
	if *doCalculate {
		mode = modeCalculate
	} else if *doCompare {
		mode = modeCompare
	} else if *doVerify {
		mode = modeVerify
	}

	app.config = configuration{mode, *algorithm, *sourceDirectory, *outputChecksum, *inputChecksum, *outputNames, *basePath}
}

func (app *Application) verifyConfiguration() {

	if app.config.mode == modeCalculate || app.config.mode == modeCompare {
		if !checkIfDirectoryExists(app.config.sourceDirectory) {
			log.Fatalln(app.config.sourceDirectory + " does not exist.")
		}
	}

	if app.config.inputChecksum != "" && !checkIfFileExists(app.config.inputChecksum) {
		log.Fatalln("Input file does not exist.")
	}
	if app.config.basePath == "" {
		app.config.basePath = app.config.sourceDirectory + "/"
	}
}

func (app *Application) execute() {

	hasher := checksum.NewFileHasher(app.config.algorithm)
	if app.config.mode == modeCalculate {
		calculator := bll.Calculator{app.config.sourceDirectory, app.config.outputChecksum, app.config.basePath}
		calculator.RecordChecksumsForDirectory(&hasher)
	} else if app.config.mode == modeCompare {
		comparer := bll.Comparer{
			app.config.sourceDirectory, app.config.inputChecksum,
			app.config.outputNames, app.config.outputChecksum,
			app.config.basePath}
		comparer.RecordNameChangesForDirectory(&hasher)
	} else if app.config.mode == modeVerify {
		verifier := bll.Verifier{app.config.inputChecksum, app.config.basePath}
		verifier.VerifyRecords(&hasher)
	}
}

func checkIfFileExists(path string) bool {

	if stat, err := os.Stat(path); err == nil && !os.IsNotExist(err) && !stat.IsDir() {
		return true
	}

	return false
}

func checkIfDirectoryExists(path string) bool {

	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

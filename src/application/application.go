package application

import (
	"bll"
	"checksum"
	"flag"
	"log"
	"os"
	"util"
)

const taskCalculate = "calculate"
const taskCompare = "compare"
const taskVerify = "verify"

// Application Contains main application logic.
type Application struct {
	config configuration
}

type configuration struct {
	task           string
	algorithm      string
	inputDirectory string
	outputChecksum string
	inputChecksum  string
	outputNames    string
	basePath       string
}

// Initialize Initializes the application.
func (app *Application) Initialize() {

	defaultConfig := configuration{taskCalculate, checksum.SHA1, "data/input", "data/checksum.csv", "", "", ""}
	app.parseCommandLineArguments(defaultConfig)
	app.verifyConfiguration()
	app.execute()
}

func (app *Application) parseCommandLineArguments(defaultConfig configuration) {

	task := flag.String(
		"task",
		defaultConfig.task,
		"The task to execute: calculate, compare or verify. The first one calculates checksums for a directory and"+
			" stores the results in a CSV. The second compares stored checksums with the checksums of the files in"+
			" the given directory and stores filename matches. The third verifies checksums for the files listed in"+
			" the given CSV.")
	inputDirectory := flag.String(
		"indir",
		defaultConfig.inputDirectory,
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

	app.config = configuration{
		*task, *algorithm, *inputDirectory,
		*outputChecksum, *inputChecksum, *outputNames,
		*basePath}
}

func (app *Application) verifyConfiguration() {

	if app.config.task != taskCalculate && app.config.task != taskCompare && app.config.task != taskVerify {
		log.Fatalln("Unknown task.")
	}
	if app.config.task == taskCalculate || app.config.task == taskCompare {
		if !checkIfDirectoryExists(app.config.inputDirectory) {
			log.Fatalln(app.config.inputDirectory + " does not exist.")
		}
	}
	if app.config.inputChecksum != "" && !util.CheckIfFileExists(app.config.inputChecksum) {
		log.Fatalln("Input file does not exist.")
	}
	if app.config.basePath == "" {
		app.config.basePath = app.config.inputDirectory + "/"
	}
}

func (app *Application) execute() {

	hasher := checksum.NewFileHasher(app.config.algorithm)
	if app.config.task == taskCalculate {
		calculator := bll.Calculator{app.config.inputDirectory, app.config.outputChecksum, app.config.basePath}
		calculator.RecordChecksumsForDirectory(&hasher)
	} else if app.config.task == taskCompare {
		comparer := bll.Comparer{
			app.config.inputDirectory, app.config.inputChecksum,
			app.config.outputNames, app.config.outputChecksum,
			app.config.basePath}
		comparer.RecordNameChangesForDirectory(&hasher)
	} else if app.config.task == taskVerify {
		verifier := bll.Verifier{app.config.inputChecksum, app.config.basePath}
		verifier.VerifyRecords(&hasher)
	}
}

func checkIfDirectoryExists(path string) bool {

	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}
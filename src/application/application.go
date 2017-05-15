package application

import (
	"bll"
	"checksum"
	"flag"
	"log"
	"os"
)

// Application Contains main application logic.
type Application struct {
	config configuration
}

type configuration struct {
	algorithm       string
	sourceDirectory string
	outputChecksum  string
	inputChecksum   string
	outputNames     string
	basePath        string
}

// Initialize Initializes the application.
func (app *Application) Initialize() {

	defaultConfig := configuration{checksum.SHA1, "data/input", "data/checksum.csv", "", "", ""}
	app.parseCommandLineArguments(defaultConfig)
	app.verifyConfiguration()
	app.execute()
}

func (app *Application) parseCommandLineArguments(defaultConfig configuration) {

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

	app.config = configuration{*algorithm, *sourceDirectory, *outputChecksum, *inputChecksum, *outputNames, *basePath}
}

func (app *Application) verifyConfiguration() {

	if !checkIfDirectoryExists(app.config.sourceDirectory) {
		log.Fatalln(app.config.sourceDirectory + " does not exist.")
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
	if app.config.inputChecksum == "" {
		calculator := bll.Calculator{app.config.sourceDirectory, app.config.outputChecksum, app.config.basePath}
		calculator.RecordChecksumsForDirectory(&hasher)
	} else {
		comparer := bll.Comparer{
			app.config.sourceDirectory, app.config.inputChecksum,
			app.config.outputNames, app.config.outputChecksum,
			app.config.basePath}
		comparer.RecordNameChangesForDirectory(&hasher)
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

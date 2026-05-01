package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/enw860/covanalyze/internal/errors"
	"github.com/enw860/covanalyze/internal/formatter"
	"github.com/enw860/covanalyze/internal/models"
	"github.com/enw860/covanalyze/internal/parser"
	"github.com/golang/glog"
)

const (
	// Exit codes
	exitSuccess      = 0
	exitFileNotFound = 1
	exitParseError   = 2
	exitOutputError  = 3
)

//go:embed usage.txt
var usageText string

var (
	coverageFile = flag.String("f", "", "Coverage file path (required)")
	outputFile   = flag.String("o", "", "Output file path (default: stdout)")
	showHelp     = flag.Bool("help", false, "Show help message")
	module       = flag.String("m", "", "Module prefix to replace in file paths")
	modulePath   = flag.String("mpath", "", "Path prefix to replace module with")
)

func main() {
	// Initialize glog flags
	flag.Set("logtostderr", "true")

	// Parse flags
	flag.Parse()

	// Debug: log all args at high verbosity level
	glog.V(3).Infof("flag.Args() = %v", flag.Args())
	glog.V(3).Infof("flag.NArg() = %d", flag.NArg())
	glog.V(3).Infof("os.Args = %v", os.Args)
	glog.V(3).Infof("coverageFile = '%s'", *coverageFile)

	// Handle --help flag
	if *showHelp {
		fmt.Fprint(os.Stdout, usageText)
		os.Exit(exitSuccess)
	}

	// Validate required coverage file flag
	if *coverageFile == "" {
		fmt.Fprintln(os.Stderr, "Error: coverage file path is required (use -f flag)")
		fmt.Fprintln(os.Stderr, "Run 'covanalyze --help' for usage information")
		os.Exit(exitParseError)
	}

	// Parse coverage file
	glog.V(1).Infof("Parsing coverage file: %s", *coverageFile)
	profiles, err := parser.ParseCoverageFile(*coverageFile)
	if err != nil {
		handleError(err)
	}

	glog.V(1).Infof("Found %d file profiles", len(profiles))

	// Calculate coverage for each file
	fileReports := make([]models.FileReport, 0)
	for _, profile := range profiles {
		glog.V(2).Infof("Calculating coverage for: %s", profile.FileName)
		report := parser.CalculateFileCoverage(profile)
		fileReports = append(fileReports, report)
	}

	// Normalize file paths if module and path flags are provided
	if *module != "" || *modulePath != "" {
		if *module == "" || *modulePath == "" {
			fmt.Fprintln(os.Stderr, "Error: both -m and -mpath flags must be provided together")
			fmt.Fprintln(os.Stderr, "Run 'covanalyze --help' for usage information")
			os.Exit(exitParseError)
		}
		glog.V(1).Infof("Normalizing file paths: replacing '%s' with '%s'", *module, *modulePath)
		normalizeFilePaths(fileReports, *module, *modulePath)
	}

	// Enrich file reports with AST-based semantic context
	glog.V(1).Info("Enriching coverage reports with semantic context")
	parser.EnrichFileReports(fileReports)

	// Create output structure
	output := &models.Output{
		FileReports: fileReports,
	}

	// Format as JSON
	glog.V(1).Info("Formatting output as JSON")
	jsonBytes, err := formatter.FormatJSON(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitOutputError)
	}

	// Write output
	glog.V(3).Infof("outputFile pointer = %v, value = '%s'", outputFile, *outputFile)
	if *outputFile != "" {
		// Write to file
		glog.V(3).Infof("Inside if block, outputFile = '%s'", *outputFile)
		glog.V(1).Infof("Writing output to file: %s", *outputFile)
		err := os.WriteFile(*outputFile, jsonBytes, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(exitOutputError)
		}
		glog.V(1).Info("Output written successfully")
	} else {
		// Write to stdout
		glog.V(1).Info("Writing output to stdout")
		fmt.Println(string(jsonBytes))
	}

	// Flush glog before exit
	glog.Flush()
	os.Exit(exitSuccess)
}

// handleError processes errors and exits with the appropriate exit code.
func handleError(err error) {
	switch e := err.(type) {
	case *errors.FileNotFoundError:
		fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		glog.Flush()
		os.Exit(exitFileNotFound)
	case *errors.UnsupportedModeError:
		fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		glog.Flush()
		os.Exit(exitParseError)
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		glog.Flush()
		os.Exit(exitParseError)
	}
}

// normalizeFilePaths replaces the module prefix with the path prefix in all file reports.
func normalizeFilePaths(fileReports []models.FileReport, module, path string) {
	for i := range fileReports {
		if strings.HasPrefix(fileReports[i].File, module) {
			fileReports[i].File = strings.Replace(fileReports[i].File, module, path, 1)
			glog.V(2).Infof("Normalized file path: %s", fileReports[i].File)
		}
	}
}

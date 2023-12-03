package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	logFilePath   string
	logLevel      string
	outputToFile  bool
	reportPath    string
	defaultReport = "default_report.txt"
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd := &cobra.Command{Use: "log-analyzer", Short: "Log Analyzer"}
	rootCmd.PersistentFlags().StringVar(&logFilePath, "log-file", "", "Path to the log file")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "Log analysis level (ERROR, WARNING, INFO)")
	rootCmd.PersistentFlags().BoolVar(&outputToFile, "output-to-file", false, "Output results to a file")
	rootCmd.PersistentFlags().StringVar(&reportPath, "report-file", "", "Path to the report file")

	rootCmd.AddCommand(analyzeCmd)
	rootCmd.Execute()

}

func initConfig() {
	if envLogFile := os.Getenv("LOG_FILE"); envLogFile != "" && logFilePath == "" {
		logFilePath = envLogFile
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" && logLevel == "" {
		logLevel = envLogLevel
	}
	if envOutputToFile := os.Getenv("OUTPUT_TO_FILE"); envOutputToFile != "" && !outputToFile {
		outputToFile, _ = strconv.ParseBool(envOutputToFile)
	}
	if envReportPath := os.Getenv("REPORT_PATH"); envReportPath != "" && reportPath == "" {
		reportPath = envReportPath
	}
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze log file or stdin",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Analyzing...")
		if logFilePath != "" {
			analyzeFile(logFilePath)
		} else {
			fmt.Println("Reading from standard input...")
			analyzeStdin()
		}
	},
}

func analyzeFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	analyzeLog(file)
}

func analyzeStdin() {
	errorCount, warningCount, infoCount := 0, 0, 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.Contains(line, "ERROR"):
			errorCount++
		case strings.Contains(line, "WARNING"):
			warningCount++
		case strings.Contains(line, "INFO"):
			infoCount++
		}
	}

	analyseResult(errorCount, warningCount, infoCount)

	if outputToFile {
		writeResultsToFile(errorCount, warningCount, infoCount)
	}
}

func analyzeLog(reader io.Reader) {
	errorCount, warningCount, infoCount := 0, 0, 0

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.Contains(line, "ERROR"):
			errorCount++
		case strings.Contains(line, "WARNING"):
			warningCount++
		case strings.Contains(line, "INFO"):
			infoCount++
		}
	}

	analyseResult(errorCount, warningCount, infoCount)

	if outputToFile {
		writeResultsToFile(errorCount, warningCount, infoCount)
	}
}

func analyseResult(errorCount, warningCount, infoCount int) {
	fmt.Println("Analysis Results:")
	switch logLevel {
	case "ERROR":
		fmt.Printf("ERROR: %d\n", errorCount)
	case "WARNING":
		fmt.Printf("ERROR: %d\n", errorCount)
		fmt.Printf("WARNING: %d\n", warningCount)
	case "INFO":
		fmt.Printf("ERROR: %d\n", errorCount)
		fmt.Printf("WARNING: %d\n", warningCount)
		fmt.Printf("INFO: %d\n", infoCount)
	}
}

func writeResultsToFile(errorCount, warningCount, infoCount int) {
	reportFile := reportPath
	if reportFile == "" {
		reportFile = defaultReport
	}

	file, err := os.Create(reportFile)
	if err != nil {
		fmt.Println("Error creating report file:", err)
		return
	}
	defer file.Close()

	switch logLevel {
	case "ERROR":
		file.WriteString(fmt.Sprintf("ERROR: %d\n", errorCount))
	case "WARNING":
		file.WriteString(fmt.Sprintf("ERROR: %d\n", errorCount))
		file.WriteString(fmt.Sprintf("WARNING: %d\n", warningCount))
	case "INFO":
		file.WriteString(fmt.Sprintf("ERROR: %d\n", errorCount))
		file.WriteString(fmt.Sprintf("WARNING: %d\n", warningCount))
		file.WriteString(fmt.Sprintf("INFO: %d\n", infoCount))
	}

	fmt.Printf("Results written to %s\n", reportFile)
}

func main() {
}

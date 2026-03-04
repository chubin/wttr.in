package options

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ProcessLogFile reads a wttr.in log file, parses queries, and writes invalid entries to an error file.
// Each log line is expected in the format: "timestamp request_id protocol ip query user_agent".
// Invalid queries and their error messages are written to errorFilePath.
func ProcessLogFile(logFilePath, errorFilePath string, config *WttrInOptions) error {
	logFile, err := openLogFile(logFilePath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	errorFile, writer, err := openErrorFile(errorFilePath)
	if err != nil {
		return err
	}
	defer errorFile.Close()
	defer writer.Flush()

	return processLogLines(logFile, writer, config)
}

// openLogFile opens the log file for reading.
func openLogFile(logFilePath string) (*os.File, error) {
	logFile, err := os.Open(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return logFile, nil
}

// openErrorFile creates and opens the error file for writing.
func openErrorFile(errorFilePath string) (*os.File, *bufio.Writer, error) {
	errorFile, err := os.Create(errorFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create error file: %w", err)
	}
	writer := bufio.NewWriter(errorFile)
	return errorFile, writer, nil
}

// processLogLines reads and processes each line of the log file.
func processLogLines(logFile *os.File, writer *bufio.Writer, config *WttrInOptions) error {
	scanner := bufio.NewScanner(logFile)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if err := processLogLine(scanner.Text(), lineNumber, writer, config); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}
	return nil
}

// processLogLine processes a single log line and writes errors to the error file.
func processLogLine(line string, lineNumber int, writer *bufio.Writer, config *WttrInOptions) error {
	query, err := extractQueryFromLine(line, lineNumber, writer)
	if err != nil {
		return err
	}

	_, err = ParseQueryString(query, config)
	if err != nil {
		_, writeErr := fmt.Fprintf(writer, "Line %d: Query: %s, Error: %v\n", lineNumber, query, err)
		if writeErr != nil {
			return fmt.Errorf("failed to write to error file at line %d: %w", lineNumber, writeErr)
		}
	}
	return nil
}

// extractQueryFromLine extracts the query from a log line and validates its format.
func extractQueryFromLine(line string, lineNumber int, writer *bufio.Writer) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		_, err := fmt.Fprintf(writer, "Line %d: Invalid log format: %s\n", lineNumber, line)
		if err != nil {
			return "", fmt.Errorf("failed to write to error file at line %d: %w", lineNumber, err)
		}
		return "", nil
	}

	query := fields[4]
	if strings.Contains(query, "?") {
		parts := strings.SplitN(query, "?", 2)
		if len(parts) == 2 {
			return parts[1], nil
		}
		return "", nil
	}
	return "", nil
}

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_RunAllTestCases(t *testing.T) {
    // Grab all files
    var measurementInputFiles []string
    testCaseDir := "test-cases"

    err := filepath.Walk(testCaseDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(path) == ".txt" {
            measurementInputFiles = append(measurementInputFiles, path)
        }
        return nil
    })

    if err != nil {
        t.Fatalf("%s", err)
    }

    for _, path := range measurementInputFiles {
        result, err := calculate(path)
        if err != nil {
            t.Fatal(err)
        }

        outputPath := path[:len(path) - 3] + "out"
        expectedResults, err := os.ReadFile(outputPath)
        if err != nil {
            t.Fatal(err)
        }

        if result != strings.TrimRight(string(expectedResults), "\n") {
            t.Logf("Expected: %s", string(expectedResults))
            t.Logf("Actual:   %s", result)
            t.Fatalf("%s not matching", outputPath)
        }
    }
}

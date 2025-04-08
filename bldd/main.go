package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type LibraryUsage struct {
	Library     string
	Executables []string
}

type ArchitectureReport struct {
	Architecture string
	Libraries    []LibraryUsage
}

type Report struct {
	ScanDirectory string
	ArchReports   []ArchitectureReport
}

func isELF(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first 4 bytes to check ELF magic number
	magic := make([]byte, 4)
	if _, err := io.ReadFull(file, magic); err != nil {
		return false
	}

	return magic[0] == 0x7F && magic[1] == 'E' && magic[2] == 'L' && magic[3] == 'F'
}

func getArchitecture(filePath string) (string, error) {
	file, err := elf.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	switch file.Machine {
	case elf.EM_386:
		return "i386 (x86)", nil
	case elf.EM_X86_64:
		return "x86-64", nil
	case elf.EM_ARM:
		return "armv7", nil
	case elf.EM_AARCH64:
		return "aarch64", nil
	default:
		return "", fmt.Errorf("unsupported architecture")
	}
}

func scanDirectory(dir string) (*Report, error) {
	report := &Report{
		ScanDirectory: dir,
		ArchReports:   make([]ArchitectureReport, 0),
	}

	// Map to store library usage per architecture
	archMap := make(map[string]map[string][]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isELF(path) {
			arch, err := getArchitecture(path)
			if err != nil {
				return nil // Skip unsupported architectures
			}

			file, err := elf.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

			libs, err := file.ImportedLibraries()
			if err != nil {
				return nil
			}

			if _, exists := archMap[arch]; !exists {
				archMap[arch] = make(map[string][]string)
			}

			for _, lib := range libs {
				archMap[arch][lib] = append(archMap[arch][lib], path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert map to sorted report structure
	for arch, libMap := range archMap {
		archReport := ArchitectureReport{
			Architecture: arch,
			Libraries:    make([]LibraryUsage, 0),
		}

		for lib, execs := range libMap {
			archReport.Libraries = append(archReport.Libraries, LibraryUsage{
				Library:     lib,
				Executables: execs,
			})
		}

		// Sort libraries by number of executables (high to low)
		sort.Slice(archReport.Libraries, func(i, j int) bool {
			return len(archReport.Libraries[i].Executables) > len(archReport.Libraries[j].Executables)
		})

		report.ArchReports = append(report.ArchReports, archReport)
	}

	// Sort architectures alphabetically
	sort.Slice(report.ArchReports, func(i, j int) bool {
		return report.ArchReports[i].Architecture < report.ArchReports[j].Architecture
	})

	return report, nil
}

func generateReport(report *Report, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "Report on dynamic used libraries by ELF executables on %s\n", report.ScanDirectory)
	fmt.Fprintf(file, "==================================================\n\n")

	for _, archReport := range report.ArchReports {
		fmt.Fprintf(file, "---------- %s ----------\n", archReport.Architecture)
		for _, lib := range archReport.Libraries {
			fmt.Fprintf(file, "%s (%d execs)\n", lib.Library, len(lib.Executables))
			for _, exec := range lib.Executables {
				fmt.Fprintf(file, "-> %s\n", exec)
			}
		}
		fmt.Fprintf(file, "\n")
	}

	return nil
}

func main() {
	scanDir := flag.String("dir", ".", "Directory to scan for ELF files")
	outputFile := flag.String("output", "report.txt", "Output file for the report")
	flag.Parse()

	report, err := scanDirectory(*scanDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	err = generateReport(report, *outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Report generated successfully in %s\n", *outputFile)
}

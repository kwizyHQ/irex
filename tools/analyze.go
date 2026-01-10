// tools/analyze.go
// Ported from analyze.ps1 to Go
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func check(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			fmt.Fprintf(os.Stderr, "%s: %v\n", msg[0], err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}

type BinaryStats struct {
	Path      string  `json:"path"`
	SizeBytes int64   `json:"sizeBytes"`
	SizeMB    float64 `json:"sizeMB"`
}

type GoEnv struct {
	GoVersion string `json:"goVersion"`
	GoOS      string `json:"goOS"`
	GoArch    string `json:"goArch"`
	CGO       string `json:"cgo"`
}

type Symbol struct {
	Address string `json:"address"`
	Size    int    `json:"size"`
	Type    string `json:"type"`
	Symbol  string `json:"symbol"`
}

type TopSymbol struct {
	Symbol string  `json:"symbol"`
	SizeKB float64 `json:"sizeKB"`
	Type   string  `json:"type"`
}

type PackageSize struct {
	Package string  `json:"package"`
	SizeMB  float64 `json:"sizeMB"`
}

type ProjectStats struct {
	Symbols     int     `json:"symbols"`
	SizeMB      float64 `json:"sizeMB"`
	SizePercent float64 `json:"sizePercent"`
}

type FileInfo struct {
	Path   string  `json:"path"`
	SizeKB float64 `json:"sizeKB"`
}

type Report struct {
	GeneratedAt  string        `json:"generatedAt"`
	Binary       BinaryStats   `json:"binary"`
	BuildTime    string        `json:"buildTime"`
	Sha256       string        `json:"sha256"`
	GoEnv        GoEnv         `json:"goEnv"`
	Project      ProjectStats  `json:"project"`
	TopPackages  []PackageSize `json:"topPackages"`
	TopSymbols   []TopSymbol   `json:"topSymbols"`
	Module       string        `json:"module"`
	Dependencies []string      `json:"dependencies"`
	Coverage     string        `json:"coverage"`
	LargestFiles []FileInfo    `json:"largestFiles"`
	Lint         string        `json:"lint"`
	Staticcheck  string        `json:"staticcheck"`
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func Analyze() {
	mainPackage := "./cmd/irex-dev"
	binaryName := "irex.exe"
	outputFile := "output.json"

	if len(os.Args) > 1 {
		mainPackage = os.Args[1]
	}
	if len(os.Args) > 2 {
		binaryName = os.Args[2]
	}
	if len(os.Args) > 3 {
		outputFile = os.Args[3]
	}

	fmt.Println("▶ Building release binary...")
	buildCmd := exec.Command("go", "build", "-trimpath", "-o", binaryName, mainPackage)
	check(buildCmd.Run(), "Build failed")
	if _, err := os.Stat(binaryName); err != nil {
		check(err, fmt.Sprintf("Build failed: %s not found", binaryName))
	}

	// Binary info
	binInfo, err := os.Stat(binaryName)
	check(err, "Stat binary")
	binStats := BinaryStats{
		Path:      binaryName,
		SizeBytes: binInfo.Size(),
		SizeMB:    float64(binInfo.Size()) / (1024 * 1024),
	}
	buildTime := binInfo.ModTime().UTC().Format(time.RFC3339)

	// SHA256
	f, err := os.Open(binaryName)
	check(err, "Open binary for SHA256")
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	check(err, "Hash binary")
	sha := fmt.Sprintf("%x", h.Sum(nil))

	// Go env
	goVersion, _ := runCmd("go", "version")
	goOS, _ := runCmd("go", "env", "GOOS")
	goArch, _ := runCmd("go", "env", "GOARCH")
	cgo, _ := runCmd("go", "env", "CGO_ENABLED")
	goEnv := GoEnv{
		GoVersion: strings.TrimSpace(goVersion),
		GoOS:      strings.TrimSpace(goOS),
		GoArch:    strings.TrimSpace(goArch),
		CGO:       strings.TrimSpace(cgo),
	}

	// Go module and dependencies
	goModPath := "go.mod"
	goModBytes, err := os.ReadFile(goModPath)
	check(err, "Read go.mod")
	moduleName := ""
	var dependencies []string
	modRe := regexp.MustCompile(`^module\s+(.+)$`)
	depRe := regexp.MustCompile(`^\s*([^\s]+)\s+v[^\s]+`)
	for _, line := range strings.Split(string(goModBytes), "\n") {
		if m := modRe.FindStringSubmatch(line); m != nil {
			moduleName = m[1]
		} else if m := depRe.FindStringSubmatch(line); m != nil {
			dep := m[1]
			if !strings.HasPrefix(dep, "//") && dep != "require" && dep != "go" {
				dependencies = append(dependencies, dep)
			}
		}
	}

	// Symbol analysis
	fmt.Println("▶ Analyzing symbols...")
	nmRaw, _ := runCmd("go", "tool", "nm", "-size", binaryName)
	var symbols []Symbol
	for _, line := range strings.Split(nmRaw, "\n") {
		parts := strings.Fields(line)
		if len(parts) == 4 {
			var size int
			fmt.Sscanf(parts[1], "%d", &size)
			symbols = append(symbols, Symbol{
				Address: parts[0],
				Size:    size,
				Type:    parts[2],
				Symbol:  parts[3],
			})
		}
	}

	// Top symbols by size
	sort.Slice(symbols, func(i, j int) bool { return symbols[i].Size > symbols[j].Size })
	var topSymbols []TopSymbol
	for i, s := range symbols {
		if i >= 30 {
			break
		}
		topSymbols = append(topSymbols, TopSymbol{
			Symbol: s.Symbol,
			SizeKB: float64(s.Size) / 1024,
			Type:   s.Type,
		})
	}

	// Package-level aggregation
	pkgSizes := make(map[string]int)
	for _, s := range symbols {
		pkg := "<runtime>"
		if idx := strings.Index(s.Symbol, "/"); idx > 0 {
			pkg = s.Symbol[:idx]
		}
		pkgSizes[pkg] += s.Size
	}
	var packages []PackageSize
	for pkg, sz := range pkgSizes {
		packages = append(packages, PackageSize{
			Package: pkg,
			SizeMB:  float64(sz) / (1024 * 1024),
		})
	}
	sort.Slice(packages, func(i, j int) bool { return packages[i].SizeMB > packages[j].SizeMB })
	if len(packages) > 25 {
		packages = packages[:25]
	}

	// Project code contribution
	projectPrefix := "github.com/kwizyHQ/irex"
	projectSizeBytes := 0
	projectSymbols := 0
	for _, s := range symbols {
		if strings.Contains(s.Symbol, projectPrefix) {
			projectSymbols++
			projectSizeBytes += s.Size
		}
	}
	projectStats := ProjectStats{
		Symbols:     projectSymbols,
		SizeMB:      float64(projectSizeBytes) / (1024 * 1024),
		SizePercent: float64(projectSizeBytes) / float64(binInfo.Size()) * 100,
	}

	// Go test coverage
	fmt.Println("▶ Running Go test coverage...")
	coverageOutput, _ := runCmd("go", "test", "./...", "-cover")
	coveragePercent := ""
	for _, line := range strings.Split(coverageOutput, "\n") {
		if m := regexp.MustCompile(`coverage: ([0-9.]+)%`).FindStringSubmatch(line); m != nil {
			coveragePercent = m[1]
		}
	}

	// Largest source files
	fmt.Println("▶ Finding largest source files...")
	var sourceFiles []FileInfo
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
			sourceFiles = append(sourceFiles, FileInfo{
				Path:   path,
				SizeKB: float64(info.Size()) / 1024,
			})
		}
		return nil
	})
	sort.Slice(sourceFiles, func(i, j int) bool { return sourceFiles[i].SizeKB > sourceFiles[j].SizeKB })
	largestFiles := sourceFiles
	if len(largestFiles) > 10 {
		largestFiles = largestFiles[:10]
	}

	// Lint/staticcheck
	fmt.Println("▶ Running golint...")
	lintOutput, _ := runCmd("golint", "./...")
	fmt.Println("▶ Running staticcheck...")
	staticcheckOutput, _ := runCmd("staticcheck", "./...")

	// Final report
	report := Report{
		GeneratedAt:  time.Now().Format(time.RFC3339),
		Binary:       binStats,
		BuildTime:    buildTime,
		Sha256:       sha,
		GoEnv:        goEnv,
		Project:      projectStats,
		TopPackages:  packages,
		TopSymbols:   topSymbols,
		Module:       moduleName,
		Dependencies: dependencies,
		Coverage:     coveragePercent,
		LargestFiles: largestFiles,
		Lint:         lintOutput,
		Staticcheck:  staticcheckOutput,
	}
	b, err := json.MarshalIndent(report, "", "  ")
	check(err, "Marshal report")
	check(os.WriteFile(outputFile, b, 0644), "Write output file")
	fmt.Printf("✅ Analysis complete\n→ Output written to %s\n", outputFile)
}

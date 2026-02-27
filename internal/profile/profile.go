package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Profile struct {
	// Demo indicates if the server is in demo mode
	Demo    bool
	Addr    string
	Port    int
	Data    string
	DSN     string
	Driver  string
	Secret  string
	Version string
}

func (p *Profile) Validate() error {

	if !p.Demo && p.Data == "" {
		if runtime.GOOS == "windows" {
			p.Data = filepath.Join(os.Getenv("ProgramData"), "go-server")
		} else {
			p.Data = "/var/opt/go-server"
		}
	}

	if p.Data == "" {
		p.Data = "."
	}

	dataDir, err := checkDataDir(p.Data)
	if err != nil {
		return err
	}
	p.Data = dataDir

	if p.Driver == "sqlite" && p.DSN == "" {
		mode := "prod"
		if p.Demo {
			mode = "demo"
		}
		p.DSN = filepath.Join(dataDir, fmt.Sprintf("go-server_%s.db", mode))
	} else if p.Driver == "postgresql" && p.DSN == "" {
		// Default PostgreSQL DSN for development
		p.DSN = "host=localhost port=5432 user=postgres password=password dbname=goserver sslmode=disable"
	}

	return nil
}

func checkDataDir(dataDir string) (string, error) {
	if !filepath.IsAbs(dataDir) {
		// Use current working directory for relative paths instead of executable directory
		// This fixes the issue with go run where os.Args[0] points to temporary build directory
		currentDir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("unable to get current working directory: %w", err)
		}
		relativeDir := filepath.Join(currentDir, dataDir)
		absDir, _ := filepath.Abs(relativeDir)
		dataDir = absDir
	}
	dataDir = strings.TrimRight(dataDir, "\\/")
	if _, err := os.Stat(dataDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dataDir, 0755); err != nil {
				return "", fmt.Errorf("unable to create data folder %s: %w", dataDir, err)
			}
		} else {
			return "", fmt.Errorf("unable to access data folder %s: %w", dataDir, err)
		}
	}
	return dataDir, nil
}

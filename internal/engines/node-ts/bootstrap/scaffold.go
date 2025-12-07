package bootstrap

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed templates/scaffold/*
var templatesFS embed.FS

// Scaffold performs project initialization for a TypeScript Node project.
// It will:
//   - create target directory structure
//   - run `npm init` or `yarn init`
//   - install base dependencies (dotenv, axios, pino)
//   - install schema/framework deps (mongoose, fastify)
//   - install devDependencies (typescript, ts-node, @types/node, @types/dotenv, nodemon)
//   - run `npx tsc --init`
//   - create src/* folders and basic files (.env.example, README.md, src/app.ts, src/vendor/server.ts)
func Scaffold() error {
	target := os.Getenv("IREX_TARGET")
	if target == "" {
		return fmt.Errorf("IREX_TARGET must be set")
	}

	pm := os.Getenv("IREX_PKG_MANAGER")
	if pm == "" {
		pm = "npm"
	}

	useNodemon := false
	if strings.EqualFold(os.Getenv("IREX_DEV_NODEMON"), "true") {
		useNodemon = true
	}

	name := os.Getenv("IREX_NAME")
	if name == "" {
		name = filepath.Base(target)
	}

	// ensure target dir exists
	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("creating target dir: %w", err)
	}

	// initialize package.json using the selected package manager
	if err := runInit(target, pm); err != nil {
		return fmt.Errorf("package init: %w", err)
	}

	// ensure package name is set
	if err := setPackageName(target, pm, name); err != nil {
		return fmt.Errorf("setting package name: %w", err)
	}

	// install dependencies
	deps := []string{"dotenv", "axios", "pino"}
	frameworkDeps := []string{"fastify", "mongoose"} // reasonable defaults
	devDeps := []string{"typescript", "ts-node", "@types/node", "@types/dotenv"}
	if useNodemon {
		devDeps = append(devDeps, "nodemon")
	}

	if err := installDeps(target, pm, deps, false); err != nil {
		return fmt.Errorf("install deps: %w", err)
	}
	if err := installDeps(target, pm, frameworkDeps, false); err != nil {
		return fmt.Errorf("install framework deps: %w", err)
	}
	if err := installDeps(target, pm, devDeps, true); err != nil {
		return fmt.Errorf("install dev deps: %w", err)
	}

	// initialize TypeScript config
	if err := runCommand(target, "npx", "tsc", "--init"); err != nil {
		return fmt.Errorf("tsc init: %w", err)
	}

	// create directory structure and files
	if err := createFiles(target, useNodemon); err != nil {
		return fmt.Errorf("creating files: %w", err)
	}

	return nil
}

func runInit(target, pm string) error {
	if pm == "yarn" {
		// yarn init -y
		return runCommand(target, "yarn", "init", "-y")
	}
	// npm init -y
	return runCommand(target, "npm", "init", "-y")
}

func setPackageName(target, pm, name string) error {
	// try npm pkg set if npm
	if pm == "npm" {
		if err := runCommand(target, "npm", "pkg", "set", fmt.Sprintf("name=%s", strings.ToLower(name))); err == nil {
			return nil
		}
		// fall through to manual edit
	}

	// read package.json and set name
	pkgPath := filepath.Join(target, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return err
	}
	var pj map[string]interface{}
	if err := json.Unmarshal(data, &pj); err != nil {
		return err
	}
	pj["name"] = name
	out, err := json.MarshalIndent(pj, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pkgPath, out, 0644)
}

func installDeps(target, pm string, pkgs []string, dev bool) error {
	if len(pkgs) == 0 {
		return nil
	}
	args := []string{}
	if pm == "yarn" {
		args = append(args, "add")
		if dev {
			args = append(args, "-D")
		}
		args = append(args, pkgs...)
		return runCommand(target, "yarn", args...)
	}
	// npm
	args = append(args, "install")
	if dev {
		args = append(args, "-D")
	} else {
		args = append(args, "--save")
	}
	args = append(args, pkgs...)
	return runCommand(target, "npm", args...)
}

func runCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createFiles(target string, useNodemon bool) error {
	// create src and subfolders
	folders := []string{
		"src",
		filepath.Join("src", "hooks"),
		filepath.Join("src", "middlewares"),
		filepath.Join("src", "utils"),
		filepath.Join("src", "vendor"),
		filepath.Join("src", "workflows"),
	}
	for _, f := range folders {
		if err := os.MkdirAll(filepath.Join(target, f), 0755); err != nil {
			return err
		}
	}

	filesToCopy := []struct {
		src      string
		dest     string
		required bool
	}{
		{"templates/scaffold/app.ts", filepath.Join(target, "src", "app.ts"), true},
		{"templates/scaffold/server.ts", filepath.Join(target, "src", "vendor", "server.ts"), true},
		{"templates/scaffold/README.md", filepath.Join(target, "README.md"), true},
		{"templates/scaffold/.env.example", filepath.Join(target, ".env.example"), true},
		{"templates/scaffold/nodemon.json", filepath.Join(target, "nodemon.json"), useNodemon},
	}
	for _, f := range filesToCopy {
		if !f.required {
			continue
		}
		if err := copyEmbeddedFile(f.src, f.dest); err != nil {
			return fmt.Errorf("copying %s to %s: %w", f.src, f.dest, err)
		}
	}
	return nil
}

func copyEmbeddedFile(src, dest string) error {
	data, err := templatesFS.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, data, 0644)
}

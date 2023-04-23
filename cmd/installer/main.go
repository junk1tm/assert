//go:build installer

// Installer is a convenience tool that copies the assertions into a project,
// so that they can be used without introducing a direct dependency.
//
// Add the following directive to any .go file in the root of your project and run `go generate ./...`:
//
//	go:generate go run -tags=installer go-simpler.org/assert/cmd/installer <path/to/pkg>
//
// Where <path/to/pkg> is the location where you want to put the generated package,
// e.g. `.` for the project root (hint: for libraries you may want to put it in `internal`).
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go-simpler.org/assert"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("path to package is not specified")
	}

	moduleName, err := readModuleName()
	if err != nil {
		return err
	}

	path := filepath.Join(os.Args[1], "assert")
	fullpath := filepath.Join(path, "dotimport")

	if err := os.MkdirAll(fullpath, 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	if err := os.Chdir(path); err != nil {
		return err
	}

	importPath := moduleName
	if path != "." {
		importPath = filepath.Join(moduleName, path)
	}

	// update the import in the `dotimport/alias.go` file.
	supportFile := strings.Replace(assert.SupportFile, "go-simpler.org/assert", importPath, 1)

	if err := writeFile("assert.go", assert.MainFile); err != nil {
		return err
	}
	if err := writeFile("dotimport/alias.go", supportFile); err != nil {
		return err
	}

	return nil
}

func readModuleName() (string, error) {
	f, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer f.Close()

	header, err := bufio.NewReader(f).ReadString('\n')
	if err != nil {
		return "", err
	}

	name := strings.TrimPrefix(header, "module")
	return strings.TrimSpace(name), nil
}

func writeFile(name, content string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	const header = `// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Code generated by go-simpler.org/assert/cmd/installer. DO NOT EDIT.

`
	return os.WriteFile(name, []byte(header+content), 0644)
}

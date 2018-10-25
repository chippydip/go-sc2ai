package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	// Checkout a temp copy of the API files
	dir, err := ioutil.TempDir("", "s2client-proto")
	check(err)
	defer os.RemoveAll(dir)

	// Preserve the test directory to look at
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	fmt.Print("Press 'Enter' to continue...")
	// 	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// }()

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      "https://github.com/Blizzard/s2client-proto",
		Progress: os.Stdout,
	})
	check(err)

	// Get all the .proto files
	protoDir := filepath.Join(dir, "s2clientprotocol")
	files, err := ioutil.ReadDir(protoDir)
	check(err)

	protocArgs := []string{
		"--proto_path=" + protoDir,
		"--gogofaster_out=api",
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".proto" {
			continue
		}
		path := filepath.Join(protoDir, file.Name())

		// Upgrade the file to proto3 and fix the package name
		writeLines(path, upgradeProto(path))

		// Add the file to the list of command line args for protoc
		protocArgs = append(protocArgs, path)
	}

	// Generate go code from the .proto files
	out, err := exec.Command("protoc", protocArgs...).CombinedOutput()
	fmt.Print(string(out))
	check(err)
}

// Thing we want to use twice
const (
	importPrefix   = "import \"s2clientprotocol/"
	optionalPrefix = "optional "
	enumPrefix     = "enum "
)

func upgradeProto(path string) []string {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	var lines []string

	// Read line by line, making modifications as needed
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Get the line and trim whitespace to make matching easier
		line := strings.TrimSpace(scanner.Text())

		switch {
		// Upgrade to proto3 and set the go package name
		case line == "syntax = \"proto2\";":
			lines = append(lines, "syntax = \"proto3\";", "option go_package = \"api\";")

		// Remove subdirectory of the import so the output path isn't nested
		case strings.HasPrefix(line, importPrefix):
			lines = append(lines, "import \""+line[len(importPrefix):])

		// Remove "optional" prefixes (they are implicit in proto3)
		case strings.HasPrefix(line, optionalPrefix):
			lines = append(lines, line[len(optionalPrefix):])

		// Race already has a zero-value, so opt-out of the next fix
		case line == "enum Race {":
			lines = append(lines, line)

		// Enums must have a zero value in proto3 (and unfortunately they must be unique due to C++ scoping rules)
		case strings.HasPrefix(line, enumPrefix):
			lines = append(lines, line, line[len(enumPrefix):len(line)-2]+"_not_specified = 0;")

		// Everything else just gets copied to the output
		default:
			lines = append(lines, line)
		}
	}

	return lines
}

func writeLines(path string, lines []string) {
	file, err := os.Create(path)
	check(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	check(writer.Flush())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

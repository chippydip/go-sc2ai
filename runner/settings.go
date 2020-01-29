package runner

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// Set changes the default value of a command line flag.
func Set(name, value string) {
	if err := flag.Set(name, value); err != nil {
		log.Print(err)
	}
}

var hasLoaded = false

func loadSettings() bool {
	if flag.Parsed() {
		return hasLoaded
	}

	// Parse the command line arguments
	showHelp := flag.Bool("help", false, "Prints help message")
	flag.Parse()
	if *showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(processSettings.processPath) == 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "Please run StarCraft II before running this API")
		flag.CommandLine.Usage()
		return false
	}

	hasLoaded = true
	return true
}

func init() {
	// Default to the environment variable (Linux mostly)
	if sc2path := os.Getenv("SC2PATH"); len(sc2path) > 0 {
		log.Printf("SC2PATH: %v", sc2path)
		processSettings.processPath = filepath.Join(sc2path, "Versions", "dummy")
	}

	// Read value from ExecuteInfo.txt if the current user has run the game before
	file, err := getUserDirectory()
	if err != nil {
		log.Printf("Error getting user directory: %v", err)
	} else if len(file) > 0 {
		file = filepath.Join(file, "Starcraft II", "ExecuteInfo.txt")
		log.Printf("ExecuteInfo path: %v", file)
	}

	if props, err := newPropertyReader(file); err == nil {
		props.getString("executable", &processSettings.processPath)
		log.Printf("  executable = %v", processSettings.processPath)
	} else {
		log.Printf("Error reading `executable`: %v", err)
	}

	// Backout the defaulted path to the Versions directory and then find the latest Base game
	if len(processSettings.processPath) > 0 {
		// Find the highest version folder where the exe exists
		if path := sc2Path(); path != "" {
			path = filepath.Join(path, "Versions")
			subdirs := getSubdirs(path)
			for i := len(subdirs) - 1; i >= 0; i-- {
				p := filepath.Join(path, subdirs[i], getBinPath())
				if _, err := os.Stat(p); err == nil {
					processSettings.processPath = p
					break
				}
			}
		}
	}

	// Blizzard Flags
	flagStr("executable", &processSettings.processPath, "The path to StarCraft II.")
	//flagInt("port", &processSettings.portStart, "The port to make StarCraft II listen on.")
	flagBool("realtime", &processSettings.realtime, "Whether to run StarCraft II in real time or not.")
	flagStr("map", &gameSettings.mapName, "Which map to run.")
	flagInt("timeout", &processSettings.timeoutMS, "Timeout for how long the library will block for a response.")
}

// SetRealtime sets the default realtime option to enabled.
func SetRealtime() {
	Set("realtime", "1")
}

// SetMap sets the default map to use.
func SetMap(name string) {
	Set("map", name)
}

func getUserDirectory() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// Should really call SHGetFolderPathW, but I don't want to mess with cgo just for that
		const key = "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Shell Folders"
		out, err := exec.Command("reg", "query", key, "/v", "Personal").CombinedOutput()

		sout := strings.TrimSpace(string(out))
		if err != nil {
			log.Print("Documents directory lookup failed: ", sout)
			return "", err
		}

		// Parse the actual value out of the output
		const prefix = len("    Personal    REG_SZ    ")
		value := strings.Split(sout, "\r\n")[1][prefix:]
		return value, nil

	case "darwin":
		user, err := user.Current()
		if err != nil {
			log.Print("Failed to get current user:", err)
			return "", err
		}
		return filepath.Join(user.HomeDir, "Library", "Application Support", "Blizzard"), nil

	default:
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		return user.HomeDir, nil
	}
}

func getBinPath() string {
	switch runtime.GOOS {
	case "windows":
		return "SC2_x64.exe"
	case "darwin":
		return "SC2.app/Contents/MacOS/SC2"
	default:
		return "SC2_x64"
	}
}

func getSubdirs(dir string) []string {
	dirs := []string{}
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	sort.Strings(dirs)
	return dirs
}

func flagStr(name string, value *string, usage string) {
	flag.StringVar(value, name, *value, usage)
}

func flagInt(name string, value *int, usage string) {
	flag.IntVar(value, name, *value, usage)
}

func flagBool(name string, value *bool, usage string) {
	flag.BoolVar(value, name, *value, usage)
}

func flagVar(name string, value flag.Value, usage string) {
	flag.Var(value, name, usage)
}

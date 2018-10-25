package runner

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// Set ...
func Set(name, value string) {
	if err := flag.Set(name, value); err != nil {
		fmt.Println(err)
	}
}

// LoadSettings ...
func LoadSettings() bool {
	// Parse the command line arguments
	flag.Parse()

	if len(processSettings.processPath) == 0 {
		fmt.Println("Please run StarCraft II before running this API")
		flag.CommandLine.Usage()
		return false
	}

	return true
}

func init() {
	file, _ := getUserDirectory()
	if len(file) > 0 {
		file = filepath.Join(file, "Starcraft II", "ExecuteInfo.txt")
	}

	// ParseFromFile
	if props, err := newPropertyReader(file); err == nil {
		props.getString("executable", &processSettings.processPath)

		var realtime int
		if props.getInt("realtime", &realtime) && realtime != 0 {
			processSettings.realtime = true
		}

		props.getInt("port", &processSettings.portStart)
		props.getString("map", &gameSettings.mapName)
		props.getInt("timeout", &processSettings.timeoutMS)
	}

	// FindLatestExe
	if len(processSettings.processPath) > 0 {
		// Get the exe name and then back out to the Versions directory
		path, exe := filepath.Split(processSettings.processPath)
		path = filepath.Dir(path) // remove trailing slash
		for path != "." && filepath.Base(path) != "Versions" {
			path = filepath.Dir(path)
		}

		// Find the highest version folder where the exe exists
		if path != "." {
			subdirs := getSubdirs(path)
			for i := len(subdirs) - 1; i >= 0; i-- {
				p := filepath.Join(path, subdirs[i], exe)
				if _, err := os.Stat(p); err == nil {
					processSettings.processPath = p
					break
				}
			}
		}
	}

	// Blizzard Flags
	flagStr("e", "executable", &processSettings.processPath, "The path to StarCraft II.")
	flagInt("s", "step_size", &processSettings.stepSize, "How many steps to take per call.")
	//flagInt("p", "port", &processSettings.portStart, "The port to make StarCraft II listen on.")
	flagBool("r", "realtime", &processSettings.realtime, "Whether to run StarCraft II in real time or not.")
	flagStr("m", "map", &gameSettings.mapName, "Which map to run.")
	flagInt("t", "timeout", &processSettings.timeoutMS, "Timeout for how long the library will block for a response.")
}

func getUserDirectory() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// Should really call SHGetFolderPathW, but I don't want to mess with cgo just for that
		const key = "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Shell Folders"
		out, err := exec.Command("reg", "query", key, "/v", "Personal").CombinedOutput()

		sout := strings.TrimSpace(string(out))
		if err != nil {
			fmt.Println("Documents directory lookup failed", sout)
			return "", err
		}

		// Parse the actual value out of the output
		const prefix = len("    Personal    REG_SZ    ")
		value := strings.Split(sout, "\r\n")[1][prefix:]
		return value, nil

	case "darwin":
		panic("NYI")

	default:
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		return user.HomeDir, nil
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

func flagStr(short string, long string, value *string, usage string) {
	flag.StringVar(value, short, *value, usage)
	flag.StringVar(value, long, *value, usage)
}

func flagInt(short string, long string, value *int, usage string) {
	flag.IntVar(value, short, *value, usage)
	flag.IntVar(value, long, *value, usage)
}

func flagBool(short string, long string, value *bool, usage string) {
	flag.BoolVar(value, short, *value, usage)
	flag.BoolVar(value, long, *value, usage)
}

func flagVar(short string, long string, value flag.Value, usage string) {
	flag.Var(value, short, usage)
	flag.Var(value, long, usage)
}

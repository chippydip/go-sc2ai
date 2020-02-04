package runner

import (
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
	"time"

	"github.com/chippydip/go-sc2ai/api"
)

var (
	processPath             = defaultExecutable()
	processInterfaceOptions = &api.InterfaceOptions{
		Raw:   true,
		Score: true,
		// FeatureLayer: &api.SpatialCameraSetup{
		// 	Resolution: &api.Size2DI{X: 512, Y: 512},
		// },
	}
	processRealtime          = false
	processConnectTimeout, _ = time.ParseDuration("2m")
)

func init() {
	// Blizzard Flags
	flagStr("executable", &processPath, "The path to StarCraft II.")
	//flagInt("port", &processSettings.portStart, "The port to make StarCraft II listen on.")
	flagBool("realtime", &processRealtime, "Whether to run StarCraft II in real time or not.")
	flagDur("timeout", &processConnectTimeout, "Timeout for how long the library will block for a response.")
}

// SetExecutable sets the default executable path to use.
func SetExecutable(exePath string) {
	Set("executable", exePath)
}

// SetRealtime sets the default realtime option to enabled.
func SetRealtime() {
	Set("realtime", "1")
}

// SetConnectTimeout sets how long to wait for a connection to the game.
func SetConnectTimeout(timeout time.Duration) {
	Set("timeout", fmt.Sprint(timeout))
}

// SetInterfaceOptions sets the interface launch options when starting a game.
func SetInterfaceOptions(options *api.InterfaceOptions) {
	processInterfaceOptions = options
}

func defaultExecutable() string {
	path := ""

	// Default to the environment variable (Linux mostly)
	if sc2path := os.Getenv("SC2PATH"); len(sc2path) > 0 {
		log.Printf("SC2PATH: %v", sc2path)
		path = filepath.Join(sc2path, "Versions", "dummy")
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
		props.getString("executable", &path)
		log.Printf("  executable = %v", path)
	} else {
		log.Printf("Error reading `executable`: %v", err)
	}

	// Backout the defaulted path to the Versions directory and then find the latest Base game
	if pp := sc2Path(path); pp != "" {
		// Find the highest version folder where the exe exists
		pp = filepath.Join(pp, "Versions")
		subdirs := getSubdirs(pp)
		for i := len(subdirs) - 1; i >= 0; i-- {
			p := filepath.Join(pp, subdirs[i], getBinPath())
			if _, err := os.Stat(p); err == nil {
				path = p
				break
			}
		}
	}
	return path
}

func hasProcessPath() bool {
	if len(processPath) == 0 {
		return false
	}
	return true
}

func processPathForBuild(build uint32) string {
	path := processPath
	if build != 0 {
		// Get the exe name and then back out to the Versions directory
		_, exe := filepath.Split(path)
		root := sc2Path(path)
		if root == "" {
			log.Fatalf("Can't find game dir: %v", path)
		}
		dir := filepath.Join(sc2Path(path), "Versions")

		// Get the path of the correct version and make sure the exe exists
		path = filepath.Join(dir, fmt.Sprintf("Base%v", build), exe)
		if _, err := os.Stat(path); err != nil {
			log.Fatalf("Base version not found: %v", err)
		}
	}
	return path
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

func defaultSc2Path() string {
	return sc2Path(processPath)
}

func sc2Path(path string) string {
	for {
		prev := path
		path = filepath.Dir(path)

		if filepath.Base(path) == "Versions" {
			return filepath.Dir(path)
		} else if path == prev {
			return ""
		}
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

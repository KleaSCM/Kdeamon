package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sahilm/fuzzy"
	"golang.org/x/sys/execabs"
)

type DesktopEntry struct {
	Name    string
	Exec    string
	Icon    string
	Comment string
}

var (
	appList    []string
	appMap     map[string]DesktopEntry
	mu         sync.Mutex
	input      *widget.Entry
	suggestion string // Track suggestion for inline autocomplete
)

func main() {

	if os.Getenv("IS_DAEMON") != "1" { // Daemonize
		daemonize()
		return
	}

	Kapp := app.New() // Initialize the app
	Kwindow := Kapp.NewWindow("Kdeamon")
	Kwindow.Resize(fyne.NewSize(600, 50)) // Wide window

	loadApplications()

	// Create longer box
	input = widget.NewEntry()
	input.SetPlaceHolder("Search Kdeamon! ")

	// Handle input changes for inline autocomplete
	input.OnChanged = func(s string) {
		go func() { // Run filtering in a separate goroutine
			mu.Lock()
			defer mu.Unlock()

			if s == "" {
				suggestion = "" // Clear suggestion
				return
			}
			// fuzzy finder
			matches := fuzzy.Find(s, appList)
			if len(matches) > 0 {
				suggestion = appList[matches[0].Index] // Store the top match as suggestion without altering user input
			} else {
				suggestion = "" // Clear suggestion if no match
			}
		}()
	}

	// Enter key to launch
	input.OnSubmitted = func(s string) {
		mu.Lock()
		defer mu.Unlock()
		appName := suggestion // Use suggestions or exact match
		if appName == "" && len(appList) > 0 {
			for _, app := range appList {
				if strings.EqualFold(app, s) {
					appName = app
					break
				}
			}
		}

		if appName != "" {
			launchApplication(appMap[appName].Exec)
			Kwindow.Hide()
			input.SetText("")
		}
	}

	content := container.NewMax(input) // make it longer and less ugly
	Kwindow.SetContent(content)
	Kwindow.ShowAndRun()
}

// Daemonize and detach from terminal
func daemonize() {
	// Set up command to re-run ((detached process))
	cmd := exec.Command(os.Args[0])               // Rerun
	cmd.Env = append(os.Environ(), "IS_DAEMON=1") // brand as a daemon
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Detach from terminal
	}

	if err := cmd.Start(); err != nil { // Start the new process and exit
		fmt.Println("Error daemonizing process:", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func loadApplications() { // Load apps into appList and appMap
	mu.Lock()
	defer mu.Unlock()
	appMap = make(map[string]DesktopEntry)
	appList = []string{}

	desktopDirs := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
	}

	for _, dir := range desktopDirs {
		files, err := filepath.Glob(filepath.Join(dir, "*.desktop"))
		if err != nil {
			continue
		}
		for _, file := range files {
			entry, err := parseDesktopEntry(file)
			if err == nil && entry.Name != "" && entry.Exec != "" {
				appMap[entry.Name] = entry
				appList = append(appList, entry.Name)
			}
		}
	}
}

// func to Parse .desktop files to load
func parseDesktopEntry(filePath string) (DesktopEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return DesktopEntry{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	entry := DesktopEntry{}
	inDesktopEntry := false

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "[Desktop Entry]" {
			inDesktopEntry = true
			continue
		}
		if !inDesktopEntry || strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]

		switch key {
		case "Name":
			entry.Name = value
		case "Exec":
			entry.Exec = parseExecValue(value)
		case "Icon":
			entry.Icon = value
		case "Comment":
			entry.Comment = value
		}
	}

	return entry, nil
}

func parseExecValue(exec string) string {
	// Remove command-line args and dsktp entry vars
	exec = strings.Split(exec, " ")[0]
	exec = strings.ReplaceAll(exec, "%u", "")
	exec = strings.ReplaceAll(exec, "%U", "")
	exec = strings.ReplaceAll(exec, "%f", "")
	exec = strings.ReplaceAll(exec, "%F", "")
	exec = strings.ReplaceAll(exec, "%i", "")
	exec = strings.ReplaceAll(exec, "%c", "")
	exec = strings.ReplaceAll(exec, "%k", "")
	return strings.TrimSpace(exec)
}

func launchApplication(execPath string) {
	cmd := execabs.Command(execPath)
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error launching application:", err)
	}
}

//
// Copyright (c) 2018 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2018-02-09
//

package aw

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/deanishe/awgo/util"
)

// Dir returns the path to the workflow's root directory.
func Dir() string { return wf.Dir() }
func (wf *Workflow) Dir() string {

	if wf.dir == "" {
		wf.dir = findWorkflowRoot()
	}

	return wf.dir
}

// CacheDir returns the path to the workflow's cache directory.
func CacheDir() string { return wf.CacheDir() }
func (wf *Workflow) CacheDir() string {

	if wf.cacheDir == "" {
		wf.cacheDir = wf.Conf.Get(EnvVarCacheDir)
	}

	return wf.cacheDir
}

// OpenCache opens the workflow's cache directory in the default application (usually Finder).
func OpenCache() error { return wf.OpenCache() }
func (wf *Workflow) OpenCache() error {
	cmd := exec.Command("open", wf.CacheDir())
	return cmd.Run()
}

// ClearCache deletes all files from the workflow's cache directory.
func ClearCache() error { return wf.ClearCache() }
func (wf *Workflow) ClearCache() error {
	return util.ClearDirectory(wf.CacheDir())
}

// DataDir returns the path to the workflow's data directory.
func DataDir() string { return wf.DataDir() }
func (wf *Workflow) DataDir() string {

	if wf.dataDir == "" {
		wf.dataDir = wf.Conf.Get(EnvVarDataDir)
	}

	return wf.dataDir
}

// OpenData opens the workflow's data directory in the default application (usually Finder).
func OpenData() error { return wf.OpenData() }
func (wf *Workflow) OpenData() error {
	cmd := exec.Command("open", wf.DataDir())
	return cmd.Run()
}

// ClearData deletes all files from the workflow's cache directory.
func ClearData() error { return wf.ClearData() }
func (wf *Workflow) ClearData() error {
	return util.ClearDirectory(wf.DataDir())
}

// Reset deletes all workflow data (cache and data directories).
func Reset() error { return wf.Reset() }
func (wf *Workflow) Reset() error {
	errs := []error{}
	if err := wf.ClearCache(); err != nil {
		errs = append(errs, err)
	}
	if err := wf.ClearData(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// LogFile returns the path to the workflow's log file.
func LogFile() string { return wf.LogFile() }
func (wf *Workflow) LogFile() string {
	return filepath.Join(wf.CacheDir(), fmt.Sprintf("%s.log", wf.BundleID()))
}

// OpenLog opens the workflow's logfile in the default application (usually Console.app).
func OpenLog() error { return wf.OpenLog() }
func (wf *Workflow) OpenLog() error {
	if !util.PathExists(wf.LogFile()) {
		log.Println("Creating log file...")
	}
	cmd := exec.Command("open", wf.LogFile())
	return cmd.Run()
}

func OpenHelp() error { return wf.OpenHelp() }
func (wf *Workflow) OpenHelp() error {
	if wf.HelpURL == "" {
		return errors.New("Help URL is not set")
	}
	cmd := exec.Command("open", wf.HelpURL)
	return cmd.Run()
}

// Try to find workflow root based on presence of info.plist.
func findWorkflowRoot() string {

	var (
		dirs = []string{}        // directories to look in for info.plist
		seen = map[string]bool{} // avoid duplicates in dirs
	)

	// Add path and all its parents to dirs & seen
	queueTree := func(p string) {

		p = filepath.Clean(p)

		segs := strings.Split(p, "/")

		for i := len(segs) - 1; i > 0; i-- {

			p := strings.Join(segs[0:i], "/")

			if p == "" {
				p = "/"
			}

			if !seen[p] {

				seen[p] = true
				dirs = append(dirs, p)
			}
		}
	}

	// Working directory is the best place to start looking
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Add all paths from working directory upwards and from
	// directory executable is in upwards.
	queueTree(wd)
	queueTree(filepath.Dir(os.Args[0]))

	// Return path of first directory that contains an info.plist
	for _, dir := range dirs {

		p := filepath.Join(dir, "info.plist")

		if _, err := os.Stat(p); err == nil {
			return dir
		}
	}

	log.Printf("[warning] info.plist not found. Guessed workflow directory: %s", wd)
	// Fall back to working directory
	return wd
}

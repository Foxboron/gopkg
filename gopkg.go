package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/tools/go/vcs"
)

type Gopkg struct {
	VCS             *vcs.RepoRoot
	Pkgname         string
	Gopath          string // GOPATH is the tmpdir for where all deps are put
	Dir             string // Dir is the absolute path to the module
	Repo            string // Name of the golang module
	Version         string
	Revision        string
	DirectoryName   string
	Url             string
	Library         bool
	GoMod           bool
	Vendor          bool
	License         string
	LicenseFilename string
	SpecialLicense  bool
	Makedepends     map[string]bool
	Depends         map[string]bool
	Checkdepends    map[string]bool
}

func (pkg *Gopkg) Fetch() error {
	dir := filepath.Join(pkg.Gopath, "src", pkg.VCS.Root)
	pkg.Dir = dir
	pkg.Repo = pkg.VCS.Root
	pkg.Url = pkg.VCS.Repo
	log.Printf("Fetching %s to %s", pkg.VCS.Repo, dir)
	pkg.VCS.VCS.Create(dir, pkg.VCS.Repo)

	if _, err := os.Stat(filepath.Join(dir, "vendor")); !os.IsNotExist(err) {
		pkg.Vendor = true
	}

	if _, err := os.Stat(filepath.Join(dir, "go.mod")); os.IsNotExist(err) {
		pkg.GoMod = false
	}
	pkg.License, pkg.LicenseFilename = DetectLicenseFile(dir)
	if ok := SpecialLicense[pkg.License]; ok {
		pkg.SpecialLicense = true
	}
	return nil
}

func (pkg *Gopkg) CheckProgram() bool {
	mains, err := pkg.Exec("go", "list", "-f", "{{.ImportPath}} {{.Name}}", "./...")
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(mains, "\n") {
		if line == "" {
			continue
		}
		if strings.Split(line, " ")[1] == "main" {
			return true
		}
	}
	return false
}

func (pkg *Gopkg) GetVersion() (string, error) {
	if tag, err := pkg.Exec("git", "describe", "--exact-match", "--tags"); err == nil {
		version := strings.TrimSpace(tag)
		if strings.HasPrefix(version, "v") {
			version = version[1:]
		}
		return version, nil
	}

	// Format: 0.0.YYYYMMDD.HASH-1
	out, err := pkg.Exec("git", "log", "--pretty=format:%ct", "-n1")
	if err != nil {
		return "", err
	}

	lastCommitUnix, err := strconv.ParseInt(strings.TrimSpace(out), 0, 64)
	if err != nil {
		return "", err
	}

	sha, err := pkg.Exec("git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}

	pkg.Revision = strings.TrimSpace(sha)

	version := fmt.Sprintf("0.0.%s.%s",
		time.Unix(lastCommitUnix, 0).UTC().Format("20060102"),
		sha[:7])
	return version, nil
}

// Returns depends checkdepends
func (pkg *Gopkg) FindDeps() (map[string]bool, map[string]bool, error) {
	removeDeps := make(map[string]bool)
	Depends := make(map[string]bool)
	Checkdepends := make(map[string]bool)
	stdLibDeps, err := pkg.Exec("go", "list", "std")
	if err != nil {
		return Depends, Checkdepends, err
	}
	for _, stdlib := range strings.Split(strings.TrimSpace(stdLibDeps), "\n") {
		removeDeps[stdlib] = true
	}
	checkDeps := func(p string) bool {
		if p == "" {
			return false
		}
		if strings.HasPrefix(p, pkg.Repo+"/") || p == pkg.Repo {
			return false
		}
		if p == "C" {
			return false
		}
		if ok := removeDeps[p]; ok {
			return false
		}
		return true
	}
	depsList, err := pkg.Exec("go", "list", "-f", "{{join .Imports \"\\n\"}}\n", "./...")
	if err != nil {
		return Depends, Checkdepends, err
	}
	for _, p := range strings.Split(strings.TrimSpace(depsList), "\n") {
		if !checkDeps(p) {
			continue
		}
		Depends[p] = true
	}
	checkdepsList, err := pkg.Exec("go", "list", "-f", "{{join .TestImports \"\\n\"}}\n{{join .XTestImports \"\\n\"}}", "./...")
	if err != nil {
		return Depends, Checkdepends, err
	}
	for _, p := range strings.Split(strings.TrimSpace(checkdepsList), "\n") {
		if !checkDeps(p) {
			continue
		}
		Checkdepends[p] = true
	}
	return Depends, Checkdepends, nil
}

func (pkg *Gopkg) Exec(cmds ...string) (string, error) {
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Dir = pkg.Dir
	cmd.Env = append([]string{
		"GO111MODULE=off",
		"GOPATH=",
	})
	if value, ok := os.LookupEnv("HOME"); ok {
		cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", value))
	}
	if value, ok := os.LookupEnv("PATH"); ok {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", value))
	}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%v: %v", cmd.Args, err)
	}
	return string(out), nil
}

func InitGopkg(src, gopath string) *Gopkg {
	rr, err := vcs.RepoRootForImportPath(src, false)
	if err != nil {
		log.Fatalf("Not a path %v", err)
	}
	if src != rr.Root {
		log.Printf("Path %q", src)
		src = rr.Root
	}
	pkg := &Gopkg{
		Gopath:         gopath,
		VCS:            rr,
		DirectoryName:  path.Base(rr.Root),
		Depends:        make(map[string]bool),
		Makedepends:    make(map[string]bool),
		Checkdepends:   make(map[string]bool),
		SpecialLicense: false,
		GoMod:          false,
		Vendor:         false,
	}
	return pkg
}

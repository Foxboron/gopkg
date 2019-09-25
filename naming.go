package main

import (
	"log"
	"strings"
)

var Names = map[string]string{
	"github.com":        "github",
	"google.golang.org": "google",
	"code.google.com":   "googlecode",
	"cloud.google.com":  "googlecloud",
	"gopkg.in":          "gopkg",
	"golang.org":        "golang",
	"bitbucket.org":     "bitbucket",
	"go4.org":           "go4",
	"pkg.deepin.io":     "deepin",
}

func GetName(name string) string {
	parts := strings.Split(name, "/")
	if host, ok := Names[parts[0]]; ok {
		parts[0] = host
	} else {
		log.Printf("No shorthand for %v found!", name)
	}
	return strings.Trim("golang-"+strings.ToLower(strings.Replace(strings.Join(parts, "-"), "_", "-", -1)), "-")
}

// Normalize all package names and remove duplicates
func NormalizePackageNames(pkgs map[string]bool) map[string]bool {
	rootPkgs := make(map[string]bool)
	for pkg := range pkgs {
		pkg := GetName(pkg)
		if ok := rootPkgs[pkg]; !ok {
			rootPkgs[pkg] = true
		}
	}
	return rootPkgs
}

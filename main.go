package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

// Environment variables
// Gopath something
var Gopath = os.Getenv("GOPKG_GOPATH")

// Build flags
var (
	buildCommand = flag.NewFlagSet("build", flag.ExitOnError)

	// CheckArchweb Check Archweb for packages not to package
	CheckArchweb = buildCommand.Bool("check", false, "check archweb for packages to not package")

	//IsProgram determines if this is a library or program
	IsProgram = buildCommand.Bool("program", false, "package as if it's a program")

	Recursive = buildCommand.Bool("recursive", false, "package the dependencies")
)

func ExecBuild(args []string) {
	// gopath, err := ioutil.TempDir("", "gopkg")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	gopath := "./gopkg_gopath"
	path, err := filepath.Abs(gopath)
	if err != nil {
		log.Fatal(err)
	}
	gopkg := args[0]
	pkg, err := CreatePkg(gopkg, path)
	if err != nil {
		log.Fatal(err)
	}
	CreateTemplate(pkg, *IsProgram)

	for pkgDep := range pkg.Depends {
		pkg, err := CreatePkg(pkgDep, path)
		if err != nil {
			log.Fatal(err)
		}
		CreateTemplate(pkg, *IsProgram)
	}
}

func main() {
	args := os.Args[1:]
	cmd := ""

	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case "build":
		buildCommand.Parse(os.Args[2:])
		ExecBuild(buildCommand.Args())
	}
}

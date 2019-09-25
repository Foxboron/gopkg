package main

import (
	"log"
	"os"

	"golang.org/x/tools/go/vcs"
)

// Ensures we only have one version of the package
// People can referr to the subpackage of a package, and we
// only care about the top-most package.
func FilterDuplicatePackages(pkgs map[string]bool) map[string]bool {
	rootPkgs := make(map[string]bool)
	for pkg := range pkgs {
		rr, err := vcs.RepoRootForImportPath(pkg, false)
		if err != nil {
			log.Fatal(err)
			continue
		}
		if ok := rootPkgs[rr.Root]; !ok {
			rootPkgs[rr.Root] = true
		}
	}
	return rootPkgs
}

func CreatePkg(gopkg, gopath string) (*Gopkg, error) {
	pkg := InitGopkg(gopkg, gopath)
	pkg.Fetch()
	pkg.Pkgname = GetName(pkg.Repo)

	if *CheckArchweb {
		if ok := CheckIfPkgExists(pkg.Pkgname); !ok {
			log.Fatal("Package exists in repository")
			os.Exit(1)
		}
	}
	if pkg.CheckProgram() && !*IsProgram {
		log.Printf("This might be a program. Please check!")
	}
	ver, err := pkg.GetVersion()
	if err != nil {
		return &Gopkg{}, err
	}
	pkg.Version = ver
	depends, checkdepends, err := pkg.FindDeps()
	if err != nil {
		return &Gopkg{}, err
	}
	if *IsProgram {
		pkg.Makedepends = FilterDuplicatePackages(pkg.Makedepends)
	} else {
		pkg.Depends = FilterDuplicatePackages(depends)
	}
	pkg.Checkdepends = FilterDuplicatePackages(checkdepends)
	return pkg, nil
}

package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
)

type TemplatePkg struct {
	Gopkg
}

func JoinDeps(deps map[string]bool) string {
	ret := ""
	for dep := range deps {
		if ret != "" {
			ret += " "
		}
		ret += fmt.Sprintf(`%s`, dep)
	}
	return ret
}

func CreateTemplate(pkg *Gopkg, isprogram bool) {
	directory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	directory = path.Join(directory, "packages")
	os.Mkdir(directory, 0755)
	os.Mkdir(path.Join(directory, pkg.Pkgname), 0755)
	f, err := os.Create(path.Join(directory, pkg.Pkgname, "PKGBUILD"))
	if err != nil {
		log.Printf("Couldn't create file: %s", err)
	}
	// So we don't modify the original Gopkg struct
	tmplPkg := &TemplatePkg{Gopkg: *pkg}
	tmplPkg.Depends = NormalizePackageNames(pkg.Depends)
	tmplPkg.Makedepends = NormalizePackageNames(pkg.Makedepends)
	tmplPkg.Checkdepends = NormalizePackageNames(pkg.Checkdepends)

	tpl := template.Must(template.New("main").Funcs(template.FuncMap{"StringsJoin": JoinDeps}).ParseGlob("./templates/*.template"))
	if isprogram {
		err = tpl.ExecuteTemplate(f, fmt.Sprintf("pkgbuild-%s.template", "program"), &tmplPkg)
	} else {
		err = tpl.ExecuteTemplate(f, fmt.Sprintf("pkgbuild-%s.template", "library"), &tmplPkg)
	}
	if err != nil {
		log.Printf("Couldn't write template: %s", err)
	}
	log.Printf("Created package %s in %s", pkg.Pkgname, path.Join(directory, pkg.Pkgname))
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PackageJSON struct {
	Pkgname string `json:"pkgname"`
	Pkgbase string `json:"pkgbase"`
	Pkgver  string `json:"pkgver"`
}

type PkgQuery struct {
	Version int64          `json:"version"`
	Results []*PackageJSON `json:"results"`
}

func SearchArchweb(name string) (*PkgQuery, error) {
	var pkgq PkgQuery
	response, err := http.Get(fmt.Sprintf("https://www.archlinux.org/packages/search/json/?sort=&q=%s", name))
	if err != nil {
		return &pkgq, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &pkgq, err
	}
	err = json.Unmarshal(contents, &pkgq)
	if err != nil {
		return &pkgq, err
	}
	return &pkgq, nil
}

func CheckIfPkgExists(name string) bool {
	pkgs, err := SearchArchweb(name)
	if err != nil {
		return false
	}
	for _, pkg := range pkgs.Results {
		fmt.Println("Matched")
		fmt.Println(name)
		fmt.Println(pkg.Pkgname)
		fmt.Println("----")
		if name == pkg.Pkgname {
			return true
		}
	}
	return false
}

func FilterPackagedPackages(pkgs map[string]bool) map[string]bool {
	filteredPkgs := make(map[string]bool)
	for pkg := range pkgs {
		if ok := filteredPkgs[pkg]; ok {
			continue
		}
		if ok := CheckIfPkgExists(pkg); ok {
			fmt.Printf("%v is packaged in the repositories\n", pkg)
			continue
		}
		filteredPkgs[pkg] = true
	}
	return filteredPkgs
}

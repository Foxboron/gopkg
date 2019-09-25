package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/licensecheck"
)

//
var LicenseFilenames = []string{
	"LICENCE.txt",
	"LICENCE",
	"LICENSE.txt",
	"LICENSE",
	"license",
	"unLICENSE",
	"unlicence",
	"license.md",
	"LICENSE.md",
	"license.txt",
	"COPYING",
	"copyRIGHT",
	"COPYRIGHT.txt",
	"copying.txt",
	"LICENSE.php",
	"LICENCE.docs",
	"copying.image",
	"COPYRIGHT.go",
	"LICENSE-MIT",
	"LICENSE_1_0.txt",
	"COPYING-GPL",
	"COPYRIGHT-BSD",
	"MIT-LICENSE.txt",
	"mit-license-foo.md",
	"OFL.md",
	"ofl.textile",
	"ofl",
	"not-the-ofl",
	"README.txt",
}

var SpecialLicense = map[string]bool{
	"BSD":    true,
	"BSD2":   true,
	"BSD3":   true,
	"ISC":    true,
	"MIT":    true,
	"ZLIB":   true,
	"Python": true,
}

// In the case where licensecheck mangles the name
var NameTranslation = map[string]string{
	"BSD-2-Clause-Patent": "BSD2",
	"BSD-2-Clause":        "BSD2",
	"BSD-3-Clause":        "BSD3",
	"GPL-3.0":             "GPL3",
	"GPL-2.0":             "GPL",
	"LGPL-2.1":            "LGPL",
	"LGPL-3":              "LGPL3",
	"MPL-3":               "MPL3",
	"MPL-1.1":             "MPL",
	"MPL-1.0":             "MPL",
	"Apache-2.0":          "Apache",
	"Apache-1.0":          "Apache",
	"CPL-1.0":             "CPL",
	"AGPL-3.0":            "AGPL",
	"Zlib":                "ZLIB",
}

func DetectLicenseFile(dir string) (string, string) {
	licenseFile := ""
	fileName := ""
	for _, file := range LicenseFilenames {
		potentialFile := filepath.Join(dir, file)
		if _, err := os.Stat(potentialFile); !os.IsNotExist(err) {
			licenseFile = potentialFile
			fileName = file
			break
		}
	}
	if licenseFile == "" {
		return "Insert License", ""
	}
	dat, err := ioutil.ReadFile(licenseFile)
	if err != nil {
		return "Insert License", ""
	}
	return DetectLicense(dat), fileName
}

func DetectLicense(input []byte) string {
	cover, _ := licensecheck.Cover(input, licensecheck.Options{})
	if normalizedName, ok := NameTranslation[cover.Match[0].Name]; ok {
		log.Printf("Detected license %v as %v\n", cover.Match[0].Name, normalizedName)
		return normalizedName
	}
	log.Printf("License name %v not normalized. Please check license\n", cover.Match[0].Name)
	return cover.Match[0].Name
}

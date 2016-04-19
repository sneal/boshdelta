package boshdelta

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strings"
)

// PivnetRelease is a Pivotal network release, .pivotal file
type PivnetRelease struct {
	Path     string
	Releases []*Release
}

// NewPivnetReleaseFromFile loads a .pivotal release
func NewPivnetReleaseFromFile(path string) (*PivnetRelease, error) {
	if filepath.Ext(path) != ".pivotal" {
		return nil, fmt.Errorf("Expected a .pivotal file, but instead got a %s file", filepath.Ext(path))
	}
	pivnetRelease := &PivnetRelease{
		Path: path,
	}
	err := pivnetRelease.loadBoshReleases()
	return pivnetRelease, err
}

// UniqueProperties returns all the distinct release properties across all BOSH
// releases and jobs
func (p *PivnetRelease) UniqueProperties() map[string]*Property {
	uniqueProps := make(map[string]*Property)
	for _, rr := range p.Releases {
		for pname, p := range rr.UniqueProperties() {
			uniqueProps[pname] = p
		}
	}
	return uniqueProps
}

func (p *PivnetRelease) loadBoshReleases() (err error) {
	zipReader, err := zip.OpenReader(p.Path)
	if err != nil {
		return err
	}
	for _, zipFile := range zipReader.File {
		if !zipFile.FileInfo().IsDir() && strings.HasPrefix(zipFile.Name, "releases/") {
			zf, zerr := zipFile.Open()
			if zerr != nil {
				return zerr
			}
			release, rerr := NewRelease(zf, zipFile.Name)
			if rerr != nil {
				return rerr
			}
			p.Releases = append(p.Releases, release)
		}
		if err != nil {
			return err
		}
	}
	return err
}

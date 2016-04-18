package boshdelta

import (
	"fmt"
	"os"
	"path/filepath"
)

// PivnetRelease is a Pivotal network release, .pivotal file
type PivnetRelease struct {
	Path     string
	Releases []Release
}

// NewPivnetRelease loads a .pivotal release
func NewPivnetRelease(path string) (*PivnetRelease, error) {
	if filepath.Ext(path) != ".pivotal" {
		return nil, fmt.Errorf("Expected a .pivotal file, but instead got a %s file", filepath.Ext(path))
	}
	pivnetRelease := &PivnetRelease{
		Path: path,
	}
	err := pivnetRelease.loadReleases()
	return pivnetRelease, err
}

func (p *PivnetRelease) loadReleases() (err error) {
	// load the .pivnet file
	f, err := os.Open(p.Path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = cerr
		}
	}()

	// load all the releases
	// TODO: load the pivotal file and walk all tgz files under releases

	return nil
}

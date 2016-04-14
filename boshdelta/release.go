package boshdelta

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"
)

const releaseManifestFileName string = "release.MF"

// Release is a BOSH release
type Release struct {
	Path               string
	Jobs               []Job  `yaml:"jobs"`
	Name               string `yaml:"name"`
	CommitHash         string `yaml:"commit_hash"`
	UncommittedChanges bool   `yaml:"uncommitted_changes"`
	Version            string `yaml:"version"`
}

// Job is a job in a BOSH release
type Job struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Sha1        string `yaml:"sha1"`
	Fingerprint string `yaml:"fingerprint"`
}

// NewRelease creates a release reading in the BOSH release metadata
func NewRelease(releasePath string) (*Release, error) {
	r := &Release{
		Path: releasePath,
	}
	err := r.readManifest()
	return r, err
}

func (r *Release) readManifest() (err error) {
	f, err := os.Open(r.Path)
	if err != nil {
		return err
	}
	defer func() { err = f.Close() }()
	gzf, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzf)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		info := header.FileInfo()
		if !info.IsDir() && info.Name() == releaseManifestFileName {
			decoder := candiedyaml.NewDecoder(tarReader)
			err = decoder.Decode(&r)
			if err != nil {
				fmt.Println(err)
				return err
			}
			break
		}
	}

	return nil
}

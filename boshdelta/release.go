package boshdelta

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"
)

const releaseManifestFileName string = "release.MF"
const jobManifestFileName string = "job.MF"

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
	Name        string              `yaml:"name"`
	Version     string              `yaml:"version"`
	Sha1        string              `yaml:"sha1"`
	Fingerprint string              `yaml:"fingerprint"`
	Properties  map[string]Property `yaml:"properties"`
}

// Property is a Job manifest property
type Property struct {
	Description string      `yaml:"description"`
	Default     interface{} `yaml:"default"`
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
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = cerr
		}
	}()
	tgzWalk(f, func(h *tar.Header, tr *tar.Reader) error {
		info := h.FileInfo()
		if !info.IsDir() && info.Name() == releaseManifestFileName {
			decoder := candiedyaml.NewDecoder(tr)
			err = decoder.Decode(&r)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

type tgzWalker func(h *tar.Header, r *tar.Reader) error

func tgzWalk(tgz io.Reader, walkFn tgzWalker) error {
	gzf, err := gzip.NewReader(tgz)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gzf)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err = walkFn(header, tr)
		if err != nil {
			return err
		}
	}
	return nil
}

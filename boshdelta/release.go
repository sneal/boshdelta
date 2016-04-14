package boshdelta

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	Name        string               `yaml:"name"`
	Version     string               `yaml:"version"`
	Sha1        string               `yaml:"sha1"`
	Fingerprint string               `yaml:"fingerprint"`
	Properties  map[string]*Property `yaml:"properties"`
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

func (r *Release) FindJob(n string) *Job {
	for _, j := range r.Jobs {
		if j.Name == n {
			return &j
		}
	}
	return nil
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

	jobs := make(map[string]*Job)

	// read the release manifest and any of its contained job manifests
	tgzWalk(f, func(h *tar.Header, tr *tar.Reader) error {
		if h.FileInfo().Name() == releaseManifestFileName {
			decoder := candiedyaml.NewDecoder(tr)
			rerr := decoder.Decode(&r)
			if rerr != nil {
				return rerr
			}
		} else if strings.HasPrefix(h.Name, "./jobs") {
			jobName := strings.TrimSuffix(filepath.Base(h.Name), filepath.Ext(h.Name))
			jobs[jobName] = &Job{}
			tgzWalk(tr, func(jh *tar.Header, jtr *tar.Reader) error {
				if jh.FileInfo().Name() == jobManifestFileName {
					decoder := candiedyaml.NewDecoder(jtr)
					jerr := decoder.Decode(jobs[jobName])
					if jerr != nil {
						return jerr
					}
				}
				return nil
			})
		}
		return nil
	})

	// copy over all properties read out of the individual job manifests
	for i := range r.Jobs {
		r.Jobs[i].Properties = jobs[r.Jobs[i].Name].Properties
	}

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

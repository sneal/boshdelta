package boshdelta

import (
	"archive/tar"
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
	err := r.loadRelease()
	return r, err
}

// FindJob returns a job instance by name, otherwise nil if not found
func (r *Release) FindJob(n string) *Job {
	for _, j := range r.Jobs {
		if j.Name == n {
			return &j
		}
	}
	return nil
}

// UniqueProperties returns all the distinct release properties across all jobs
func (r *Release) UniqueProperties() map[string]*Property {
	uniqueProps := make(map[string]*Property)
	for ji := range r.Jobs {
		rjob := r.Jobs[ji]
		for pname, p := range rjob.Properties {
			uniqueProps[pname] = p
		}
	}
	return uniqueProps
}

func (r *Release) loadRelease() (err error) {
	f, err := os.Open(r.Path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = cerr
		}
	}()
	err = r.readReleaseAndJobManifests(f)
	return err
}

func (r *Release) readReleaseAndJobManifests(f io.Reader) error {
	jobs := make(map[string]*Job)

	// read the release manifest and any of its contained job manifests
	err := tgzWalk(f, func(h *tar.Header, tr *tar.Reader) error {
		if h.FileInfo().IsDir() {
			return nil
		} else if h.FileInfo().Name() == releaseManifestFileName {
			decoder := candiedyaml.NewDecoder(tr)
			rerr := decoder.Decode(&r)
			if rerr != nil {
				return rerr
			}
		} else if strings.HasPrefix(h.Name, "./jobs") {
			jobName := strings.TrimSuffix(filepath.Base(h.Name), filepath.Ext(h.Name))
			jobs[jobName] = &Job{}
			jwerr := tgzWalk(tr, func(jh *tar.Header, jtr *tar.Reader) error {
				if jh.FileInfo().Name() == jobManifestFileName {
					decoder := candiedyaml.NewDecoder(jtr)
					jerr := decoder.Decode(jobs[jobName])
					if jerr != nil {
						return jerr
					}
				}
				return nil
			})
			if jwerr != nil {
				return jwerr
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// copy over all properties read out of the individual job manifests
	for i := range r.Jobs {
		r.Jobs[i].Properties = jobs[r.Jobs[i].Name].Properties
	}

	return nil
}

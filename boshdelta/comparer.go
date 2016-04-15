package boshdelta

import "fmt"

// Delta is the result from comparing two BOSH releases
type Delta struct {
	DeltaProperties []DeltaProperty
}

// DeltaProperty is a new property added to an existing or new job
type DeltaProperty struct {
	Name        string
	Description string
}

// ContainsProperty returns true if the delta contains the specified property
func (d *Delta) ContainsProperty(property string) bool {
	for _, p := range d.DeltaProperties {
		if p.Name == property {
			return true
		}
	}
	return false
}

// Compare two BOSH releases
func Compare(release1Path, release2Path string) (*Delta, error) {
	release1, err := NewRelease(release1Path)
	if err != nil {
		return nil, err
	}
	release2, err := NewRelease(release2Path)
	if err != nil {
		return nil, err
	}
	return CompareReleases(release1, release2), nil
}

// CompareReleases compares two loaded BOSH releases
func CompareReleases(release1, release2 *Release) *Delta {
	d := &Delta{}
	for ji := range release1.Jobs {
		r1job := release1.Jobs[ji]
		r2job := release2.FindJob(r1job.Name)
		fmt.Println(r1job.Name)
		if r2job == nil {
			// new job, add all properties
			// TODO add properties
		} else {
			// existing job, compare properties
			for k := range r1job.Properties {
				if _, ok := r2job.Properties[k]; !ok {
					// new property
					d.DeltaProperties = append(d.DeltaProperties, DeltaProperty{
						Name:        k,
						Description: r2job.Properties[k].Description,
					})
				} else {
					// existing property
				}
			}
		}
	}
	return d
}

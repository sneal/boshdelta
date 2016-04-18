package boshdelta

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
	release1, err := NewReleaseFromFile(release1Path)
	if err != nil {
		return nil, err
	}
	release2, err := NewReleaseFromFile(release2Path)
	if err != nil {
		return nil, err
	}
	return CompareReleases(release1, release2), nil
}

// CompareReleases compares two loaded BOSH releases
func CompareReleases(release1, release2 *Release) *Delta {
	d := &Delta{}
	release1UniqueProps := release1.UniqueProperties()
	release2UniqueProps := release2.UniqueProperties()
	for n, p := range release2UniqueProps {
		if _, ok := release1UniqueProps[n]; !ok {
			d.DeltaProperties = append(d.DeltaProperties, DeltaProperty{
				Name:        n,
				Description: p.Description,
			})
		}
	}
	return d
}

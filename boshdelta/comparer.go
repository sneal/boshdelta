package boshdelta

// Delta is the result from comparing two BOSH releases
type Delta struct {
	DeltaProperties []DeltaProperty
}

// DeltaProperty is a new property added to an existing or new job
type DeltaProperty struct {
	JobName string
	Name    string
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
	d := &Delta{}

	for ji := range release1.Jobs {
		r1job := release1.Jobs[ji]
		r2job := release2.FindJob(r1job.Name)
		if r2job == nil {
			// new job, add all properties
			// TODO add properties
		} else {
			// existing job, compare properties
			for k := range r1job.Properties {
				if _, ok := r2job.Properties[k]; !ok {
					// new property
					d.DeltaProperties = append(d.DeltaProperties, DeltaProperty{
						JobName: r2job.Name,
						Name:    k,
					})
				} else {
					// existing property
				}
			}
		}
	}
	return nil, nil
}

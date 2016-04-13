package boshdelta

type release struct {
	path string
	jobs []job
}

type job struct {
}

func newRelease(releasePath string) *release {
	r := &release{
		path: releasePath,
	}
	r.readManifest()
	return r
}

func (*release) readManifest() {

}

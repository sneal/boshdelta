package boshdelta

// UniquePropertyer will returns all the unique properties of itself and
// it children
type UniquePropertyer interface {
	UniqueProperties() map[string]*Property
}

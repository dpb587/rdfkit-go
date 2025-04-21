package inmemory

type DatasetOption interface {
	apply(*Dataset)
}

type datasetOptionFunc func(*Dataset)

func (f datasetOptionFunc) apply(d *Dataset) {
	f(d)
}

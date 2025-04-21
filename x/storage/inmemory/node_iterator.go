package inmemory

// type subjectBindingIterator struct {
// 	index int
// 	nodes []*NodeBinding
// }

// var _ graph.SubjectBindingIterator = &subjectBindingIterator{}

// func (i *subjectBindingIterator) Next() bool {
// 	if i.index >= len(i.nodes)-1 {
// 		return false
// 	}

// 	i.index++

// 	return true
// }

// func (i *subjectBindingIterator) GetSubjectBinding() graph.SubjectBinding {
// 	return i.nodes[i.index]
// }

// func (i *subjectBindingIterator) Err() error {
// 	return nil
// }

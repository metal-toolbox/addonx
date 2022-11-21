package reconciler

// Summary is a simple summary struct for reconcile actions
type Summary struct {
	cList, uList, dList, eList []string
}

// Created returns a list of created items
func (s *Summary) Created() []string {
	return s.cList
}

// Updated returns a list of updated items
func (s *Summary) Updated() []string {
	return s.uList
}

// Deleted returns a list of deleted items
func (s *Summary) Deleted() []string {
	return s.dList
}

// Errors returns a list of errors processing items
func (s *Summary) Errors() []string {
	return s.eList
}

// AddCreated appends an item to a list of created items
func (s *Summary) AddCreated(c string) {
	if s.cList == nil {
		s.cList = []string{}
	}

	s.cList = append(s.cList, c)
}

// AddUpdated appends an item to a list of created items
func (s *Summary) AddUpdated(u string) {
	if s.uList == nil {
		s.uList = []string{}
	}

	s.uList = append(s.uList, u)
}

// AddDeleted appends an item to a list of created items
func (s *Summary) AddDeleted(d string) {
	if s.dList == nil {
		s.dList = []string{}
	}

	s.dList = append(s.dList, d)
}

// AddErrors appends an item to a list of created items
func (s *Summary) AddErrors(e string) {
	if s.eList == nil {
		s.eList = []string{}
	}

	s.eList = append(s.eList, e)
}

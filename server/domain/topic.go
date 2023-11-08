package domain

const (
	RU int = iota
	EN
)

type Topic struct {
	Idr  int
	Lang int
	Name string
}

// TODO: this is the best solution?
type TopicList []Topic

// functions/methods
// -----------------------------------------------------------------------

func (s Topic) Eq(other Topic) bool {
	return s.Lang == other.Lang &&
		s.Name == other.Name
}

func (s TopicList) Contains(other Topic) bool {
	if s == nil {
		return false
	}

	for i := range s {
		if s[i].Eq(other) {
			return true
		}
	}
	return false
}

func (s TopicList) Eq(other TopicList) bool {
	if s == nil || other == nil {
		return false
	}
	if len(s) != len(other) {
		return false
	}

	// ***

	for i := range other {
		if !s.Contains(other[i]) {
			return false
		}
	}
	return true
}

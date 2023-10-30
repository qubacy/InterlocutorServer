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

func Contains(topics []Topic, one Topic) bool {
	if topics == nil {
		return false
	}

	for i := range topics {
		if topics[i].Eq(one) {
			return true
		}
	}
	return false
}

func IsEqual(lhs, rhs Topic) bool {
	return lhs.Lang == rhs.Lang &&
		lhs.Name == rhs.Name
}

func (self Topic) Eq(other Topic) bool {
	return IsEqual(self, other)
}

func IsSomeEqual(lhs, rhs []Topic) bool {
	if lhs == nil || rhs == nil {
		return false
	}
	if len(lhs) != len(rhs) {
		return false
	}

	// ***

	for i := range lhs {
		if !Contains(rhs, lhs[i]) {
			return false
		}
	}
	return true
}

package domain

type Admin struct {
	Idr   int
	Login string
	Pass  string
}

type AdminList []Admin

// methods
// -----------------------------------------------------------------------

func (s Admin) Eq(other Admin) bool {
	return s.Login == other.Login &&
		s.Pass == other.Pass
}

func (s AdminList) Contains(other Admin) bool {
	if s == nil {
		return false
	}

	for i := range s {
		if s[i].Eq(other) {
			return true
		}
	}
	return true
}

func (s AdminList) Eq(other AdminList) bool {
	if s == nil || other == nil {
		return false
	}
	if len(s) != len(other) {
		return false
	}

	for i := range other {
		if !s.Contains(other[i]) {
			return false
		}
	}
	return true
}

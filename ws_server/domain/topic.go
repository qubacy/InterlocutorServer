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

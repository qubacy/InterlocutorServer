package domain

type Admin struct {
	Idr              int
	Login            string
	Pass             string
	RefreshTokenHash string
}

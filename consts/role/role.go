package role

type Role int

const (
	ROOT Role = iota
	ADMIN
	USER
	GUEST
)

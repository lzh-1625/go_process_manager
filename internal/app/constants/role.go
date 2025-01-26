package constants

type Role int

const (
	ROLE_ROOT Role = iota
	ROLE_ADMIN
	ROLE_USER
	ROLE_GUEST
)

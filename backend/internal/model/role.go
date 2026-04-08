package model

type Role string

const (
	RoleCreator  Role = "creator"
	RoleAdmin    Role = "admin"
	RoleMember   Role = "member"
	RoleObserver Role = "observer"
)

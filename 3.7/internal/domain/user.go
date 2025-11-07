package domain

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleViewer  Role = "viewer"
)

type User struct {
	ID       int64
	Username string
	Pass     string
	Role     Role
}

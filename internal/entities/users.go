package entities

type User struct {
	Id       int64  `json:"-" db:"id"`
	Login    string `json:"login" validate:"required" db:"login"`
	Password string `json:"password" validate:"required" db:"password"`
}

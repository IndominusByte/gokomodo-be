package auth

type JsonLoginSchema struct {
	Email    string `json:"email" validate:"required,email,min=3,max=100" db:"email"`
	Password string `json:"password" validate:"required,min=6,max=100" db:"password"`
}

type JsonPasswordOnlySchema struct {
	Password string `validate:"required,min=6,max=100"`
}

type User struct {
	Id       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	Address  string `json:"address" db:"address"`
}

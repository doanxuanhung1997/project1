package model

type InputCreateUser struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       int    `json:"role"`
	Department int    `json:"department"`
	Avatar     string `json:"avatar"`
}

type InputLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

package domain

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Register struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Country  string `json:"country"`
}

type LoginWithToken struct {
	Token string `json:"token"`
}

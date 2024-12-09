package e_user

import (
	"github.com/littlebluewhite/Account/entry/e_w_user"
	"time"
)

type UserCreate struct {
	Username string     `json:"username"`
	Name     *string    `json:"name"`
	Password string     `json:"password"`
	Birthday *time.Time `json:"birthday"`
	Email    *string    `json:"email"`
	Phone    *string    `json:"phone"`
	Country  *string    `json:"country"`
}

type UserUpdate struct {
	ID       int32      `json:"id"`
	Username *string    `json:"username"`
	Name     *string    `json:"name"`
	Birthday *time.Time `json:"birthday"`
	Email    *string    `json:"email"`
	Phone    *string    `json:"phone"`
	Country  *string    `json:"country"`
	LoginAt  *time.Time `json:"login_at"`
}

func (uu *UserUpdate) GetKey(key string) int {
	if key == "id" {
		return int(uu.ID)
	}
	return 0
}

type User struct {
	ID        int32            `json:"id"`
	Username  string           `json:"username"`
	Name      *string          `json:"name"`
	Password  string           `json:"password"`
	Birthday  *time.Time       `json:"birthday"`
	Email     *string          `json:"email"`
	Phone     *string          `json:"phone"`
	Country   *string          `json:"country"`
	LoginAt   *time.Time       `json:"login_at"`
	UpdatedAt *time.Time       `json:"updated_at"`
	CreatedAt *time.Time       `json:"created_at"`
	WUsers    []e_w_user.WUser `json:"w_user"`
}

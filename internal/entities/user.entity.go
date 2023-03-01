package entities

type User struct {
	IdUser   int    `json:"id_user" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Status   bool   `json:"status" validate:"required"`
	Token    string `json:"token" validate:"required"`
	IsAdmin  bool   `json:"is_admin" validate:"required"`
}

type UserCreate struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserRead struct {
	IdUser   int    `json:"id_user" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Status   bool   `json:"status" validate:"required"`
	Token    string `json:"token" validate:"required"`
	IsAdmin  bool   `json:"is_admin" validate:"required"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserForgotPassword struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type UserUpdatePassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type UserValidate struct {
	Token string `query:"token" validate:"required"`
}

package dto

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required"`
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
}

type UpdateUserProfileBody struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UpdateUserProfileRequest struct {
	ID        string `json:"id" validate:"required"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

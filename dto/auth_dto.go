package dto

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

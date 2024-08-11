package dto

type GoogleUserEmailResponse struct {
	Email string `json:"email"`
}

type VerifyGoogleLoginRequest struct {
	Code       string `json:"code"`
	ServiceUrl string `json:"service_url"`
}

type VerifyGoogleLoginSvcResponse struct {
	User User `json:"user"`
}

type VerifyGoogleLoginResponse struct {
	ServiceTicket string `json:"service_ticket"`
}

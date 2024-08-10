package dto

type GoogleUserEmailResponse struct {
	Email string `json:"email"`
}

type VerifyGoogleLoginRequest struct {
	Code       string `json:"code"`
	ServiceUrl string `json:"service_url"`
}

type VerifyGoogleLoginResponse struct {
	ServiceTicket *ServiceTicket `json:"service_ticket"`
	SessionCookie string         `json:"session_cookie"`
}

type IssueSTRequest struct {
	UserID     string `json:"user_id"`
	ServiceUrl string `json:"service_url"`
}

type IssueSTResponse struct {
	ServiceTicket *ServiceTicket `json:"service_ticket"`
}

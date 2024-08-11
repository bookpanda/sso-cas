package dto

type ServiceTicket struct {
	Token      string `json:"token"`
	ServiceUrl string `json:"service_url"`
	UserID     string `json:"user_id"`
	ExpiresAt  int    `json:"expires_at"`
}

type CreateServiceTicketRequest struct {
	ServiceUrl string `json:"service_url"`
	UserID     string `json:"user_id"`
}

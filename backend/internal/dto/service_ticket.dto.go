package dto

type ServiceTicket struct {
	ServiceUrl string `json:"service_url"`
	Token      string `json:"token"`
}

type CreateServiceTicketRequest struct {
	ServiceUrl string `json:"service_url"`
	User       User   `json:"user"`
}

package request

type TicketReq struct {
	CountryCode string `json:"countryCode"`
	Tag         string `json:"tag"`
}

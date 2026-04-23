package aviasales

type AviasalesTicketDetails struct {
	Link            string `json:"link"`
	DepartureAt     string `json:"departure_at"`
	Airline         string `json:"airline"`
	ReturnAt        string `json:"return_at"`
	TicketPrice     int    `json:"price"`
	DurationTo      int    `json:"duration_to"`
	DurationBack    int    `json:"duration_back"`
	Transfers       int    `json:"transfers"`
	ReturnTransfers int    `json:"return_transfers"`
}

type AviasalesResponse struct {
	Data []AviasalesTicketDetails `json:"data"`
}

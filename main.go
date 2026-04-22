package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AviasalesTicketDetails struct {
	Link         string `json:"link"`
	DepartureAt  string `json:"departure_at"`
	Airline      string `json:"airline"`
	ReturnAt     string `json:"return_at"`
	TicketPrice  int    `json:"price"`
	DurationTo   int    `json:"duration_to"`
	DurationBack int    `json:"duration_back"`
	Transfers    int    `json:"transfers"`
}

type AviasalesResponse struct {
	Data []AviasalesTicketDetails `json:"data"`
}

func main() {
	ctx := context.Background()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://api.travelpayouts.com/v1/prices/cheap?origin=MOW&destination=HKT&depart_date=2026-05&return_date=2026-06&token=0aafb3e482df4b8e661280173a4396a0", nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.travelpayouts.com/aviasales/v3/prices_for_dates?origin=MOW&destination=HKT&departure_at=2026-05&return_at=2026-06&sorting=price&limit=10&cy=rub&unique=false&token=0aafb3e482df4b8e661280173a4396a0&direct=true", nil)
	if err != nil {
		fmt.Println("Error on creating request", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error on sending request", err)
		return
	}
	defer resp.Body.Close()
	var aviaResult AviasalesResponse
	json.NewDecoder(resp.Body).Decode(&aviaResult)

	fmt.Println("Самые дешевые билеты на выбранные даты:")
	for _, results := range aviaResult.Data {
		days := daysDiff(results.DepartureAt, results.ReturnAt)
		fmt.Printf("Вылет: %s. Прилет: %s. Дней: %d. Цена: %d ₽ \n", formatDate(results.DepartureAt), formatDate(results.ReturnAt), days, results.TicketPrice)
	}

}

func formatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("02.01.2006 15:04")
}

func daysDiff(from, to string) int {
	fromTime, _ := time.Parse(time.RFC3339, from)
	toTime, _ := time.Parse(time.RFC3339, to)
	return int(toTime.Sub(fromTime).Hours() / 24)
}

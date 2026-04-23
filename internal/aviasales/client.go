package aviasales

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func DoSearch(from string, to string, monthTo string, monthBack string, withTransfers bool) (string, error) {
	ctx := context.Background()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	travelPayouts := os.Getenv("AVIASALES_TOKEN")
	requestUrl := fmt.Sprintf("https://api.travelpayouts.com/aviasales/v3/prices_for_dates?origin=%s&destination=%s&departure_at=%s&return_at=%s&sorting=price&limit=10&cy=rub&unique=false&token=%s&direct=%t", from, to, monthTo, monthBack, travelPayouts, withTransfers)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)

	if err != nil {
		return "", fmt.Errorf("error on creating request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error on sending request: %w", err)
	}

	defer resp.Body.Close()
	var aviaResult AviasalesResponse
	json.NewDecoder(resp.Body).Decode(&aviaResult)

	message := createMessage(aviaResult)
	return message, nil
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

func createMessage(avuaResult AviasalesResponse) string {
	message := "✈️ <b>Самые дешёвые билеты:</b>\n\n"

	for i, r := range avuaResult.Data {
		days := daysDiff(r.DepartureAt, r.ReturnAt)
		transfers := "прямой"
		if r.Transfers > 0 {
			transfers = fmt.Sprintf("%d пересадка туда, %d обратно", r.Transfers, r.ReturnTransfers)
		}
		message += fmt.Sprintf(
			"<b>%d.</b> 🗓 %s → %s (%d дн.)\n"+
				"💰 <b>%d ₽</b> • %s\n"+
				"🔀 %s\n"+
				"<a href=\"https://aviasales.ru/%s\">🔗 Купить билет</a>\n\n",
			i+1,
			formatDate(r.DepartureAt), formatDate(r.ReturnAt), days,
			r.TicketPrice, r.Airline,
			transfers,
			r.Link,
		)
	}
	return message
}

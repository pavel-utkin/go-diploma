package accrual

import (
	"encoding/json"
	"fmt"
	"go-diploma/internal/api/models"
	"log"
	"net/http"
	"time"
)

type Client struct {
	http.Client

	Addr string
}

func NewClient(addr string) *Client {
	client := http.Client{}
	client.Timeout = 1 * time.Second

	return &Client{
		Client: client,
		Addr:   addr,
	}
}

func (c *Client) GetOrderAccruals(orderNr int64) (models.OrderAccrual, error) {
	url := fmt.Sprintf("%s/api/orders/%d", c.Addr, orderNr)
	log.Printf("requesting accrual: %s", url)

	response, errReq := c.Get(url)
	if errReq != nil {
		return models.OrderAccrual{}, fmt.Errorf("cannot request accrual server: %w", errReq)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return models.OrderAccrual{}, fmt.Errorf("unexpected status: %s", response.Status)
	}
	if contentType := response.Header.Get("Content-Type"); contentType != "application/json" {
		return models.OrderAccrual{}, fmt.Errorf("unexpected content type [%s]", contentType)
	}

	accrualJSON := OrderAccrualJSON{}

	dec := json.NewDecoder(response.Body)
	if err := dec.Decode(&accrualJSON); err != nil {
		return models.OrderAccrual{}, fmt.Errorf("cannot parse accrual service response: %w", err)
	}

	orderAccrual, errConv := accrualJSON.ToOrderAccrual()
	if errConv != nil {
		return models.OrderAccrual{}, fmt.Errorf("cannot convert response to Accrual: %w", errConv)
	}

	return orderAccrual, nil
}

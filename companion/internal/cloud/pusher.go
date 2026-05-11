package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/farmsimcompanymanager/companion/internal/parser"
)

type Pusher struct {
	apiURL         string
	companionToken string
	client         *http.Client
}

func NewPusher(apiURL, companionToken string) *Pusher {
	return &Pusher{
		apiURL:         apiURL,
		companionToken: companionToken,
		client:         &http.Client{Timeout: 30 * time.Second},
	}
}

type SyncPayload struct {
	Company  parser.Company          `json:"company"`
	Finances []parser.DailyFinances  `json:"finances"`
	Fields   []parser.Field          `json:"fields"`
	Vehicles []parser.Vehicle        `json:"vehicles"`
}

func (p *Pusher) Push(company parser.Company, finances []parser.DailyFinances, fields []parser.Field, vehicles []parser.Vehicle) error {
	payload := SyncPayload{
		Company:  company,
		Finances: finances,
		Fields:   fields,
		Vehicles: vehicles,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest("POST", p.apiURL+"/sync", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.companionToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned %d", resp.StatusCode)
	}

	log.Printf("cloud: pushed slot %s", company.SlotID)
	return nil
}

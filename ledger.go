package blnkgo

import (
	"fmt"
	"net/http"
	"time"
)

type LedgerService service

type Ledger struct {
	LedgerID  string                 `json:"ledger_id"`
	Name      string                 `json:"name"`
	CreatedAt time.Time              `json:"created_at"`
	MetaData  map[string]interface{} `json:"meta_data,omitempty"`
}

type CreateLedgerRequest struct {
	Name     string                 `json:"name"`
	MetaData map[string]interface{} `json:"meta_data,omitempty"`
}

func (s *LedgerService) List() ([]Ledger, *http.Response, error) {
	req, err := s.client.NewRequest("ledgers", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}
	var ledgers []Ledger
	resp, err := s.client.CallWithRetry(req, &ledgers)
	if err != nil {
		return nil, resp, err
	}
	return ledgers, resp, nil
}

func (s *LedgerService) Get(id string) (*Ledger, *http.Response, error) {
	if id == "" {
		return nil, nil, fmt.Errorf("invalid: id is required")
	}
	u := fmt.Sprintf("ledgers/%s", id)
	req, err := s.client.NewRequest(u, http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	ledger := new(Ledger)
	resp, err := s.client.CallWithRetry(req, ledger)
	if err != nil {
		return nil, resp, err
	}

	return ledger, resp, nil
}

func (s *LedgerService) Create(body CreateLedgerRequest) (*Ledger, *http.Response, error) {
	req, err := s.client.NewRequest("ledgers", http.MethodPost, body)
	if err != nil {
		return nil, nil, err
	}

	ledger := new(Ledger)
	resp, err := s.client.CallWithRetry(req, ledger)
	if err != nil {
		return nil, resp, err
	}

	return ledger, resp, nil
}

func NewLedgerService(c ClientInterface) *LedgerService {
	return &LedgerService{client: c}
}

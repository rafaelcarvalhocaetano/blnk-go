package blnkgo

import (
	"fmt"
	"net/http"
	"time"
)

type SearchService service

type SearchParams struct {
	Q        string  `json:"q"`
	QueryBy  *string `json:"query_by,omitempty"`
	FilterBy *string `json:"filter_by,omitempty"`
	SortBy   *string `json:"sort_by,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PerPage  *int    `json:"per_page,omitempty"`
}

type SearchResponse struct {
	Found         int          `json:"found"`
	OutOf         int          `json:"out_of"`
	Page          int          `json:"page"`
	RequestParams SearchParams `json:"request_params"`
	SearchTimeMs  int          `json:"search_time_ms"`
	Hits          []SearchHit  `json:"hits"`
}

type SearchHit struct {
	Document SearchDocument `json:"document"`
}

type SearchDocument struct {
	BalanceID     string                 `json:"balance_id"`
	Balance       float64                `json:"balance"`
	CreditBalance float64                `json:"credit_balance"`
	DebitBalance  float64                `json:"debit_balance"`
	Currency      string                 `json:"currency"`
	Precision     int                    `json:"precision"`
	LedgerID      string                 `json:"ledger_id"`
	CreatedAt     time.Time              `json:"created_at"`
	MetaData      map[string]interface{} `json:"meta_data"`
}

func (s *SearchService) SearchDocument(body SearchParams, resource ResourceType) (*SearchResponse, *http.Response, error) {
	u := fmt.Sprintf("search/%s", resource)
	req, err := s.client.NewRequest(u, http.MethodPost, body)
	if err != nil {
		return nil, nil, err
	}

	searchResponse := new(SearchResponse)
	resp, err := s.client.CallWithRetry(req, &searchResponse)
	if err != nil {
		return nil, resp, err
	}

	return searchResponse, resp, nil
}

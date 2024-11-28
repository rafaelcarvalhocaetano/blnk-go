package blnkgo

import (
	"fmt"
	"net/http"
	"time"
)

type LedgerBalanceService service

type LedgerBalance struct {
	BalanceID             string      `json:"balance_id"`
	Balance               int         `json:"balance"`
	Version               int         `json:"version"`
	InflightBalance       int         `json:"inflight_balance"`
	CreditBalance         int         `json:"credit_balance"`
	InflightCreditBalance int         `json:"inflight_credit_balance"`
	DebitBalance          int         `json:"debit_balance"`
	InflightDebitBalance  int         `json:"inflight_debit_balance"`
	Precision             int         `json:"precision"`
	LedgerID              string      `json:"ledger_id"`
	IdentityID            string      `json:"identity_id"`
	Indicator             string      `json:"indicator"`
	Currency              string      `json:"currency"`
	CreatedAt             time.Time   `json:"created_at"`
	InflightExpiresAt     time.Time   `json:"inflight_expires_at"`
	MetaData              interface{} `json:"meta_data,omitempty"`
}

type CreateLedgerBalanceRequest struct {
	LedgerID   string      `json:"ledger_id"`
	IdentityID string      `json:"identity_id,omitempty"`
	Currency   string      `json:"currency"`
	MetaData   interface{} `json:"meta_data,omitempty"`
}

func (s *LedgerBalanceService) Create(body CreateLedgerBalanceRequest) (*LedgerBalance, *http.Response, error) {
	req, err := s.client.NewRequest("balances", http.MethodPost, body)
	if err != nil {
		return nil, nil, err
	}

	ledgerBalance := new(LedgerBalance)
	resp, err := s.client.CallWithRetry(req, &ledgerBalance)
	if err != nil {
		return nil, resp, err
	}

	return ledgerBalance, resp, nil
}

func (s *LedgerBalanceService) Get(balanceID string) (*LedgerBalance, *http.Response, error) {
	u := fmt.Sprintf("balances/%s", balanceID)
	req, err := s.client.NewRequest(u, http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}
	ledgerBalance := new(LedgerBalance)
	resp, err := s.client.CallWithRetry(req, &ledgerBalance)
	if err != nil {
		return nil, resp, err
	}
	return ledgerBalance, resp, nil
}

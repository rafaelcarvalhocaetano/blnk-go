package blnkgo

import "net/http"

type BalanceMonitorService service

type MonitorCondition struct {
	Field     string                    `json:"field"`
	Operator  MonitorConditionOperators `json:"operator"`
	Value     int64                     `json:"value"`
	Precision int64                     `json:"precision"`
}

// MonitorData represents the data structure for monitoring information.
type MonitorData struct {
	Condition   MonitorCondition `json:"condition"`
	Description *string          `json:"description,omitempty"`
	BalanceID   string           `json:"balance_id"`
	CallBackURL *string          `json:"call_back_url,omitempty"`
}

// MonitorDataResp extends MonitorData with additional fields for response data.
type MonitorDataResp struct {
	MonitorData
	MonitorID string `json:"monitor_id"`
	CreatedAt string `json:"created_at"` // ISO date string
}

func (s *BalanceMonitorService) Create(data MonitorData) (*MonitorDataResp, *http.Response, error) {
	req, err := s.client.NewRequest("balance-monitors", http.MethodPost, data)
	if err != nil {
		return nil, nil, err
	}

	monitorData := new(MonitorDataResp)
	resp, err := s.client.CallWithRetry(req, &monitorData)
	if err != nil {
		return nil, resp, err
	}

	return monitorData, resp, nil
}

func (s *BalanceMonitorService) Get(monitorID string) (*MonitorDataResp, *http.Response, error) {
	req, err := s.client.NewRequest("balance-monitors/"+monitorID, http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var resp MonitorDataResp
	httpResp, err := s.client.CallWithRetry(req, &resp)
	if err != nil {
		return nil, httpResp, err
	}

	return &resp, httpResp, nil
}

func (s *BalanceMonitorService) List() ([]MonitorDataResp, *http.Response, error) {
	req, err := s.client.NewRequest("balance-monitors", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var monitorData []MonitorDataResp
	resp, err := s.client.CallWithRetry(req, &monitorData)
	if err != nil {
		return nil, resp, err
	}

	return monitorData, resp, nil
}

func (s *BalanceMonitorService) Update(monitorID string, data MonitorData) (*MonitorDataResp, *http.Response, error) {
	req, err := s.client.NewRequest("balance-monitors/"+monitorID, http.MethodPut, data)
	if err != nil {
		return nil, nil, err
	}

	monitorData := new(MonitorDataResp)
	resp, err := s.client.CallWithRetry(req, &monitorData)
	if err != nil {
		return nil, resp, err
	}

	return monitorData, resp, nil
}

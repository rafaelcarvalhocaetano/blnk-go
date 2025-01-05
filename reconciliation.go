package blnkgo

import (
	"net/http"
)

type ReconciliationService service

type ReconciliationUploadResp struct {
	UploadID    string `json:"upload_id"`
	RecordCount int    `json:"record_count"`
	Source      string `json:"source"`
}

// Criteria represents the filtering criteria.
type Criteria struct {
	Field          CriteriaField          `json:"field"`
	Operator       ReconciliationOperator `json:"operator"`
	AllowableDrift float64                `json:"allowable_drift,omitempty"` // Optional field
}

// Matcher represents a matching rule with multiple criteria.
type Matcher struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Criteria    []Criteria `json:"criteria"`
}

// RunReconData represents the data required to run a reconciliation process.
type RunReconData struct {
	UploadID         string                 `json:"upload_id"`
	Strategy         ReconciliationStrategy `json:"strategy"`
	DryRun           bool                   `json:"dry_run"`
	GroupingCriteria CriteriaField          `json:"grouping_criteria"`
	MatchingRuleIDs  []string               `json:"matching_rule_ids"`
}

type RunReconResp struct {
	Matcher
	RuleID    string `json:"rule_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *ReconciliationService) CreateMatchingRule(matcher Matcher) (*RunReconResp, *http.Response, error) {
	req, err := s.client.NewRequest("reconciliation/matching-rules", http.MethodPost, matcher)
	if err != nil {
		return nil, nil, err
	}

	reconResp := new(RunReconResp)
	resp, err := s.client.CallWithRetry(req, reconResp)
	if err != nil {
		return nil, resp, err
	}

	return reconResp, resp, nil
}

func (s *ReconciliationService) Run(data RunReconData) (*RunReconResp, *http.Response, error) {
	req, err := s.client.NewRequest("reconciliation/start", http.MethodPost, data)
	if err != nil {
		return nil, nil, err
	}
	reconResp := new(RunReconResp)
	resp, err := s.client.CallWithRetry(req, reconResp)
	if err != nil {
		return nil, resp, err
	}

	return reconResp, resp, nil
}

func (s *ReconciliationService) Upload(source string, file interface{}, fileName string) (*ReconciliationUploadResp, *http.Response, error) {
	req, err := s.client.NewFileUploadRequest("reconciliation/upload", "file", file, fileName, map[string]string{
		"source": source,
	})

	if err != nil {
		return nil, nil, err
	}

	reconResp := new(ReconciliationUploadResp)
	resp, err := s.client.CallWithRetry(req, reconResp)
	if err != nil {
		return nil, resp, err
	}

	return reconResp, resp, nil
}

func NewReconciliationService(c ClientInterface) *ReconciliationService {
	return &ReconciliationService{client: c}
}

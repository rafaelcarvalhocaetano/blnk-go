package blnkgo_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	blnkgo "github.com/blnkfinance/blnk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSearchService() (*MockClient, *blnkgo.SearchService) {
	mockClient := &MockClient{}
	svc := blnkgo.NewSearchService(mockClient)
	return mockClient, svc
}

func TestSearchService_SearchDocument_Success(t *testing.T) {
	mockClient, svc := setupSearchService()

	body := blnkgo.SearchParams{
		Q:       "test query",
		QueryBy: "field",
		Page:    1,
		PerPage: 10,
	}

	expectedResponse := &blnkgo.SearchResponse{
		Found: 1,
		OutOf: 1,
		Page:  1,
		RequestParams: blnkgo.SearchParams{
			Q:       "test query",
			QueryBy: "field",
			Page:    1,
			PerPage: 10,
		},
		SearchTimeMs: 100,
		Hits: []blnkgo.SearchHit{
			{
				Document: blnkgo.SearchDocument{
					BalanceID:     "balance123",
					Balance:       100.0,
					CreditBalance: 50.0,
					DebitBalance:  50.0,
					Currency:      "USD",
					Precision:     2,
					LedgerID:      "ledger123",
					CreatedAt:     blnkgo.FlexibleTime{Time: time.Now()},
					MetaData:      map[string]interface{}{"key": "value"},
				},
			},
		},
	}

	mockClient.On("NewRequest", "search/resource", http.MethodPost, body).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
	}, nil).Run(func(args mock.Arguments) {
		searchResponse := args.Get(1).(*blnkgo.SearchResponse)
		*searchResponse = *expectedResponse
	})

	searchResponse, resp, err := svc.SearchDocument(body, "resource")

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, searchResponse)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockClient.AssertExpectations(t)
}

func TestSearchService_SearchDocument_EmptyRequest(t *testing.T) {
	mockClient, svc := setupSearchService()
	body := blnkgo.SearchParams{}

	mockClient.On("NewRequest", "search/ledgers", http.MethodPost, body).Return(nil, fmt.Errorf("invalid request"))
	searchResponse, resp, err := svc.SearchDocument(body, "ledgers")

	assert.Error(t, err)
	assert.Nil(t, searchResponse)
	assert.Nil(t, resp)
	mockClient.AssertExpectations(t)
}

func TestSearchService_SearchDocument_ServerError(t *testing.T) {
	mockClient, svc := setupSearchService()
	body := blnkgo.SearchParams{
		Q:       "test query",
		QueryBy: "field",
		Page:    1,
		PerPage: 10,
	}

	expectedResp := &http.Response{StatusCode: http.StatusInternalServerError}

	mockClient.On("NewRequest", "search/ledgers", http.MethodPost, body).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(expectedResp, fmt.Errorf("server error"))

	searchResponse, resp, err := svc.SearchDocument(body, "ledgers")

	assert.Error(t, err)
	assert.Nil(t, searchResponse)
	assert.Equal(t, expectedResp, resp)
	mockClient.AssertExpectations(t)
}

func TestSearchService_SearchDocument_InvalidResource(t *testing.T) {
	mockClient, svc := setupSearchService()
	body := blnkgo.SearchParams{
		Q:       "test query",
		QueryBy: "field",
		Page:    1,
		PerPage: 10,
	}

	mockClient.On("NewRequest", "search/invalid_resource", http.MethodPost, body).Return(nil, fmt.Errorf("invalid resource"))
	searchResponse, resp, err := svc.SearchDocument(body, "invalid_resource")

	assert.Error(t, err)
	assert.Nil(t, searchResponse)
	assert.Nil(t, resp)
	mockClient.AssertExpectations(t)
}

func TestSearchService_SearchDocument_EmptyResponse(t *testing.T) {
	mockClient, svc := setupSearchService()
	body := blnkgo.SearchParams{
		Q:       "test query",
		QueryBy: "field",
		Page:    1,
		PerPage: 10,
	}

	mockClient.On("NewRequest", "search/ledgers", http.MethodPost, body).Return(&http.Request{}, nil)
	mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
	}, nil).Run(func(args mock.Arguments) {
		searchResponse := args.Get(1).(*blnkgo.SearchResponse)
		*searchResponse = blnkgo.SearchResponse{}
	})

	searchResponse, resp, err := svc.SearchDocument(body, "ledgers")

	assert.NoError(t, err)
	assert.NotNil(t, searchResponse)
	assert.Equal(t, 0, searchResponse.Found)
	assert.Equal(t, 0, searchResponse.OutOf)
	assert.Equal(t, 0, searchResponse.Page)
	assert.Equal(t, 0, searchResponse.SearchTimeMs)
	assert.Empty(t, searchResponse.Hits)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockClient.AssertExpectations(t)
}

func TestFlexibleTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		jsonData     string
		wantErr      bool
		expectedTime time.Time
	}{
		{
			name:         "Unix timestamp as number",
			jsonData:     `{"created_at": 1672531200}`,
			wantErr:      false,
			expectedTime: time.Unix(1672531200, 0),
		},
		{
			name:         "Unix timestamp as string",
			jsonData:     `{"created_at": "1672531200"}`,
			wantErr:      false,
			expectedTime: time.Unix(1672531200, 0),
		},
		{
			name:         "RFC3339 string",
			jsonData:     `{"created_at": "2023-01-01T00:00:00Z"}`,
			wantErr:      false,
			expectedTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "RFC3339 string with timezone",
			jsonData:     `{"created_at": "2023-01-01T12:30:45-03:00"}`,
			wantErr:      false,
			expectedTime: time.Date(2023, 1, 1, 15, 30, 45, 0, time.UTC),
		},
		{
			name:     "Invalid format",
			jsonData: `{"created_at": "invalid-date"}`,
			wantErr:  true,
		},
		{
			name:     "Empty string",
			jsonData: `{"created_at": ""}`,
			wantErr:  true,
		},
		{
			name:     "Null value",
			jsonData: `{"created_at": null}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc struct {
				CreatedAt blnkgo.FlexibleTime `json:"created_at"`
			}

			err := json.Unmarshal([]byte(tt.jsonData), &doc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.False(t, doc.CreatedAt.Time.IsZero())

				// Verificar se o tempo foi parseado corretamente
				if !tt.expectedTime.IsZero() {
					assert.Equal(t, tt.expectedTime.Unix(), doc.CreatedAt.Time.Unix())
				}
			}
		})
	}
}

func TestFlexibleTime_MarshalJSON(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	flexTime := blnkgo.FlexibleTime{Time: testTime}

	data, err := json.Marshal(flexTime)
	assert.NoError(t, err)

	// Should marshal as Unix timestamp
	expected := fmt.Sprintf("%d", testTime.Unix())
	assert.Equal(t, expected, string(data))
}

func TestFlexibleTime_RoundTrip(t *testing.T) {
	originalTime := time.Date(2023, 8, 6, 15, 30, 45, 0, time.UTC)
	flexTime := blnkgo.FlexibleTime{Time: originalTime}

	// Marshal to JSON
	data, err := json.Marshal(flexTime)
	assert.NoError(t, err)

	// Unmarshal back
	var unmarshaled blnkgo.FlexibleTime
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	// Should be equal (considering Unix timestamp precision)
	assert.Equal(t, originalTime.Unix(), unmarshaled.Time.Unix())
}

func TestSearchDocument_MetaData_FlexibleTypes(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name: "MetaData as object",
			jsonData: `{
				"balance_id": "bal-123",
				"balance": 100.50,
				"meta_data": {"key": "value", "count": 42}
			}`,
			wantErr: false,
		},
		{
			name: "MetaData as string",
			jsonData: `{
				"balance_id": "bal-123",
				"balance": 100.50,
				"meta_data": "string metadata"
			}`,
			wantErr: false,
		},
		{
			name: "MetaData as null",
			jsonData: `{
				"balance_id": "bal-123",
				"balance": 100.50,
				"meta_data": null
			}`,
			wantErr: false,
		},
		{
			name: "MetaData as array",
			jsonData: `{
				"balance_id": "bal-123",
				"balance": 100.50,
				"meta_data": ["item1", "item2"]
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc blnkgo.SearchDocument
			err := json.Unmarshal([]byte(tt.jsonData), &doc)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "bal-123", doc.BalanceID)
				assert.Equal(t, 100.50, doc.Balance)
				// MetaData pode ser qualquer tipo, então apenas verificamos que não é nil se não for explicitamente null
				if tt.name != "MetaData as null" {
					assert.NotNil(t, doc.MetaData)
				}
			}
		})
	}
}

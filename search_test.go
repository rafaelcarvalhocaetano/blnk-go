package blnkgo_test

import (
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
					CreatedAt:     time.Now(),
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

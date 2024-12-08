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

type MockClient struct {
	mock.Mock
}

type Tests[T any] struct {
	name        string
	body        T
	expectError bool
	errorMsg    string
}

// NewRequest is a mock method that simulates creating a new HTTP request.
func (m *MockClient) NewRequest(endpoint string, method string, body interface{}) (*http.Request, error) {
	args := m.Called(endpoint, method, body)
	if req, ok := args.Get(0).(*http.Request); ok || args.Get(0) == nil {
		return req, args.Error(1)
	}

	return nil, args.Error(1)
}

// CallWithRetry is a mock method that simulates making an HTTP call with retry logic.
func (m *MockClient) CallWithRetry(req *http.Request, v interface{}) (*http.Response, error) {
	args := m.Called(req, v)
	if resp, ok := args.Get(0).(*http.Response); ok || args.Get(0) == nil {
		return resp, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockClient) NewFileUploadRequest(endpoint string, fileParam string, file interface{}, fileName string, fields map[string]string) (*http.Request, error) {
	args := m.Called(endpoint, fileParam, file)
	return args.Get(0).(*http.Request), args.Error(1)
}

// Helper function to setup mock client and service
func setupTransactionService() (*MockClient, *blnkgo.TransactionService) {
	mockClient := &MockClient{}
	svc := blnkgo.NewTransactionService(mockClient)
	return mockClient, svc
}

// Successfully create a new transaction with valid request body
func TestCreateTransactionSuccess(t *testing.T) {
	// Setup mock client and service
	tests := []Tests[blnkgo.CreateTransactionRequest]{
		{
			name: "Create transaction Success",
			body: blnkgo.CreateTransactionRequest{
				ParentTransaction: blnkgo.ParentTransaction{
					Amount:      1000,
					Reference:   "ref-21",
					Precision:   100,
					Currency:    "USD",
					Source:      "@bank-account",
					Destination: "@World",
					MetaData: map[string]interface{}{
						"transaction_type": "deposit",
						"customer_name":    "Alice Johnson",
						"customer_id":      "alice-5786",
					},
					Description: "Alice Funds",
				},
				Inflight: true,
			},
			expectError: false,
			errorMsg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient, svc := setupTransactionService()

			// Setup expected response with fixed time
			fixedTime := time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC)
			expectedResp := &blnkgo.Transaction{
				ParentTransaction: tt.body.ParentTransaction,
				TransactionID:     "tx-123",
				CreatedAt:         fixedTime,
			}

			// Setup mock expectations
			mockClient.On("NewRequest", "transactions", http.MethodPost, tt.body).Return(&http.Request{}, nil)
			mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(&http.Response{}, nil).Run(func(args mock.Arguments) {
				// Copy expected response to output parameter
				transaction := args.Get(1).(*blnkgo.Transaction)
				*transaction = *expectedResp
			})

			transaction, resp, err := svc.Create(tt.body)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, expectedResp, transaction)
			mockClient.AssertExpectations(t)
		})
	}
}

// Create transaction with invalid or missing required fields
func TestCreateTransactionInvalidRequest(t *testing.T) {
	tests := []Tests[blnkgo.CreateTransactionRequest]{
		{
			name: "Missing Amount and Currency",
			body: blnkgo.CreateTransactionRequest{
				ParentTransaction: blnkgo.ParentTransaction{
					Reference: "TEST-REF",
					Status:    blnkgo.PryTransactionStatusApplied,
				},
			},
			expectError: true,
			errorMsg:    "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient, svc := setupTransactionService()

			transaction, resp, err := svc.Create(tt.body)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Nil(t, resp)
			assert.Nil(t, transaction)
			mockClient.AssertNotCalled(t, "NewRequest")
			mockClient.AssertNotCalled(t, "CallWithRetry")
		})
	}
}

func TestUpdateTransaction(t *testing.T) {
	// Setup mock client and service
	tests := []Tests[blnkgo.UpdateStatus]{
		{
			name: "Update transaction Success",
			body: blnkgo.UpdateStatus{
				Status: blnkgo.InflightStatusCommit,
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Update transaction Fail",
			body: blnkgo.UpdateStatus{
				Status: blnkgo.InflightStatusVoid,
			},
			expectError: true,
			errorMsg:    "transaction not found",
		},
		{
			name: "Valid Url Format",
			body: blnkgo.UpdateStatus{
				Status: blnkgo.InflightStatusCommit,
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Invalid Url Format",
			body: blnkgo.UpdateStatus{
				Status: blnkgo.InflightStatusCommit,
			},
			expectError: true,
			errorMsg:    "invalid URL format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient, svc := setupTransactionService()

			// Setup expected response with fixed time
			fixedTime := time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC)
			expectedResp := &blnkgo.Transaction{
				ParentTransaction: blnkgo.ParentTransaction{
					Amount:      1000,
					Reference:   "ref-21",
					Precision:   100,
					Currency:    "USD",
					Source:      "@bank-account",
					Destination: "@World",
					MetaData: map[string]interface{}{
						"transaction_type": "deposit",
						"customer_name":    "Alice Johnson",
						"customer_id":      "alice-5786",
					},
					Description: "Alice Funds",
					Status:      blnkgo.PryTransactionStatus(tt.body.Status),
				},
				TransactionID: "tx-123",
				CreatedAt:     fixedTime,
			}

			// Setup mock expectations
			mockClient.On("NewRequest", "transactions/inflight/tx-123", http.MethodPut, tt.body).Return(&http.Request{}, nil)

			// Setup mock expectations for CallWithRetry
			if tt.expectError {
				mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(tt.errorMsg))
			} else {
				mockClient.On("CallWithRetry", mock.Anything, mock.Anything).Return(
					&http.Response{},
					nil,
				).Run(func(args mock.Arguments) {
					transaction := args.Get(1).(*blnkgo.Transaction)
					*transaction = *expectedResp
				})
			}

			transaction, resp, err := svc.Update("tx-123", tt.body)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, transaction)
				assert.Nil(t, resp)

			} else {
				// Assert
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, expectedResp, transaction)
				mockClient.AssertExpectations(t)
			}
		})
	}
}

func TestTransactionService_Update_RequestCreationFailure(t *testing.T) {
	mockClient, svc := setupTransactionService()

	// Setup mock expectations for NewRequest to return an error
	mockClient.On("NewRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to create request"))

	// Call the Update method
	transaction, resp, err := svc.Update("tx-123", blnkgo.UpdateStatus{
		Status: blnkgo.InflightStatusCommit,
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create request")
	assert.Nil(t, transaction)
	assert.Nil(t, resp)
	mockClient.AssertNotCalled(t, "CallWithRetry")
}

func TestTransactionService_Update_InvalidID(t *testing.T) {
	mockClient, svc := setupTransactionService()

	// Call the Update method with an invalid ID
	transaction, resp, err := svc.Update("", blnkgo.UpdateStatus{
		Status: blnkgo.InflightStatusCommit,
	})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transactionID is required")
	assert.Nil(t, transaction)
	assert.Nil(t, resp)
	mockClient.AssertNotCalled(t, "CallWithRetry")
}

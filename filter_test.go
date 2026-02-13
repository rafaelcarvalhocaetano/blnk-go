package blnkgo_test

import (
	"testing"

	blnkgo "github.com/blnkfinance/blnk-go"
	"github.com/stretchr/testify/assert"
)

func TestFilterParams_EmptyFilters(t *testing.T) {
	req := blnkgo.FilterParams{}
	assert.Empty(t, req.Filters)
}

func TestFilterParams_NumericOperators_WithInt(t *testing.T) {
	operators := []blnkgo.Operator{blnkgo.OpGreaterThan, blnkgo.OpGreaterThanOrEqual, blnkgo.OpLessThan, blnkgo.OpLessThanOrEqual}

	for _, op := range operators {
		t.Run(string(op)+"_with_int", func(t *testing.T) {
			req := blnkgo.FilterParams{
				Filters: []blnkgo.Filter{
					{Field: "amount", Operator: op, Value: 1000},
				},
			}
			assert.Len(t, req.Filters, 1)
			assert.Equal(t, op, req.Filters[0].Operator)
		})
	}
}

func TestFilterParams_NumericOperators_WithFloat(t *testing.T) {
	operators := []blnkgo.Operator{blnkgo.OpGreaterThan, blnkgo.OpGreaterThanOrEqual, blnkgo.OpLessThan, blnkgo.OpLessThanOrEqual}

	for _, op := range operators {
		t.Run(string(op)+"_with_float", func(t *testing.T) {
			req := blnkgo.FilterParams{
				Filters: []blnkgo.Filter{
					{Field: "amount", Operator: op, Value: 1000.50},
				},
			}
			assert.Len(t, req.Filters, 1)
			assert.Equal(t, op, req.Filters[0].Operator)
		})
	}
}

func TestFilterParams_NumericOperators_WithDateString(t *testing.T) {
	operators := []blnkgo.Operator{blnkgo.OpGreaterThan, blnkgo.OpGreaterThanOrEqual, blnkgo.OpLessThan, blnkgo.OpLessThanOrEqual}

	for _, op := range operators {
		t.Run(string(op)+"_with_date_string", func(t *testing.T) {
			req := blnkgo.FilterParams{
				Filters: []blnkgo.Filter{
					{Field: "created_at", Operator: op, Value: "2024-01-01"},
				},
			}
			assert.Len(t, req.Filters, 1)
			assert.Equal(t, op, req.Filters[0].Operator)
		})
	}
}

func TestFilterParams_LikeOperators_WithString(t *testing.T) {
	operators := []blnkgo.Operator{blnkgo.OpLike, blnkgo.OpILike}

	for _, op := range operators {
		t.Run(string(op)+"_with_string", func(t *testing.T) {
			req := blnkgo.FilterParams{
				Filters: []blnkgo.Filter{
					{Field: "name", Operator: op, Value: "%savings%"},
				},
			}
			assert.Len(t, req.Filters, 1)
			assert.Equal(t, op, req.Filters[0].Operator)
		})
	}
}

func TestFilterParams_EqOperator(t *testing.T) {
	req := blnkgo.FilterParams{
		Filters: []blnkgo.Filter{
			{Field: "status", Operator: blnkgo.OpEqual, Value: "APPLIED"},
		},
	}
	assert.Len(t, req.Filters, 1)
	assert.Equal(t, blnkgo.OpEqual, req.Filters[0].Operator)
}

func TestFilterParams_MultipleFilters(t *testing.T) {
	req := blnkgo.FilterParams{
		Filters: []blnkgo.Filter{
			{Field: "status", Operator: blnkgo.OpEqual, Value: "APPLIED"},
			{Field: "amount", Operator: blnkgo.OpGreaterThanOrEqual, Value: 10000},
			{Field: "created_at", Operator: blnkgo.OpLessThan, Value: "2024-04-01"},
		},
		SortBy:    "amount",
		SortOrder: "desc",
		Limit:     50,
	}
	assert.Len(t, req.Filters, 3)
	assert.Equal(t, "amount", req.SortBy)
	assert.Equal(t, "desc", req.SortOrder)
	assert.Equal(t, 50, req.Limit)
}

func TestFilterParams_InOperator(t *testing.T) {
	req := blnkgo.FilterParams{
		Filters: []blnkgo.Filter{
			{Field: "status", Operator: blnkgo.OpIn, Values: []interface{}{"APPLIED", "PENDING"}},
		},
	}
	assert.Len(t, req.Filters, 1)
	assert.Equal(t, blnkgo.OpIn, req.Filters[0].Operator)
	assert.Len(t, req.Filters[0].Values, 2)
}

func TestFilterParams_BetweenOperator(t *testing.T) {
	req := blnkgo.FilterParams{
		Filters: []blnkgo.Filter{
			{Field: "amount", Operator: blnkgo.OpBetween, Values: []interface{}{100, 1000}},
		},
	}
	assert.Len(t, req.Filters, 1)
	assert.Equal(t, blnkgo.OpBetween, req.Filters[0].Operator)
	assert.Len(t, req.Filters[0].Values, 2)
}

func TestFilterParams_NullOperators(t *testing.T) {
	t.Run("IsNull", func(t *testing.T) {
		req := blnkgo.FilterParams{
			Filters: []blnkgo.Filter{
				{Field: "description", Operator: blnkgo.OpIsNull},
			},
		}
		assert.Len(t, req.Filters, 1)
		assert.Equal(t, blnkgo.OpIsNull, req.Filters[0].Operator)
	})

	t.Run("IsNotNull", func(t *testing.T) {
		req := blnkgo.FilterParams{
			Filters: []blnkgo.Filter{
				{Field: "description", Operator: blnkgo.OpIsNotNull},
			},
		}
		assert.Len(t, req.Filters, 1)
		assert.Equal(t, blnkgo.OpIsNotNull, req.Filters[0].Operator)
	})
}

func TestFilterParams_WithPagination(t *testing.T) {
	req := blnkgo.FilterParams{
		Filters: []blnkgo.Filter{
			{Field: "status", Operator: blnkgo.OpEqual, Value: "APPLIED"},
		},
		Limit:        20,
		Offset:       40,
		IncludeCount: true,
	}
	assert.Equal(t, 20, req.Limit)
	assert.Equal(t, 40, req.Offset)
	assert.True(t, req.IncludeCount)
}

func TestFilterResponse_UnmarshalJSON_WithObject(t *testing.T) {
	jsonData := `{
		"data": [
			{"balance_id": "bln_123", "currency": "USD"}
		],
		"total_count": 100
	}`

	var resp blnkgo.FilterResponse
	err := resp.UnmarshalJSON([]byte(jsonData))

	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
	assert.NotNil(t, resp.TotalCount)
	assert.Equal(t, int64(100), *resp.TotalCount)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 1)
}

func TestFilterResponse_UnmarshalJSON_WithArray(t *testing.T) {
	jsonData := `[
		{"balance_id": "bln_123", "currency": "USD"},
		{"balance_id": "bln_456", "currency": "EUR"}
	]`

	var resp blnkgo.FilterResponse
	err := resp.UnmarshalJSON([]byte(jsonData))

	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
	assert.Nil(t, resp.TotalCount)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 2)
}

func TestFilterResponse_UnmarshalJSON_WithEmptyArray(t *testing.T) {
	jsonData := `[]`

	var resp blnkgo.FilterResponse
	err := resp.UnmarshalJSON([]byte(jsonData))

	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
	assert.Nil(t, resp.TotalCount)

	data, ok := resp.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 0)
}

func TestFilterResponse_UnmarshalJSON_WithEmptyObject(t *testing.T) {
	jsonData := `{"data": [], "total_count": 0}`

	var resp blnkgo.FilterResponse
	err := resp.UnmarshalJSON([]byte(jsonData))

	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
	assert.NotNil(t, resp.TotalCount)
	assert.Equal(t, int64(0), *resp.TotalCount)
}

func TestFilterResponse_UnmarshalJSON_InvalidJSON(t *testing.T) {
	jsonData := `invalid json`

	var resp blnkgo.FilterResponse
	err := resp.UnmarshalJSON([]byte(jsonData))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode FilterResponse")
}

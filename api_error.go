package blnkgo

import (
	"fmt"
	"io"
	"net/http"
)

//This function will take in Resp as a parameter and check the status code, for error in the range of 400 return a 400 error, for error in the range of 500 return a 500 error, for success return nil

// we create a struct for api error response
type ApiErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Body    []byte `json:"body"`
}

// implement the error interface for ApiErrorResponse
func (a *ApiErrorResponse) Error() string {
	return fmt.Sprintf("Status: %d, Message: %s, Body: %s", a.Status, a.Message, a.Body)
}
func (c *Client) CheckResponse(resp *http.Response) error {
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		//read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		//create a new api error response
		apiErrorResponse := &ApiErrorResponse{
			Status:  resp.StatusCode,
			Message: resp.Status,
			Body:    body,
		}
		//return the error
		return apiErrorResponse
	}

	if resp.StatusCode >= 500 {
		//read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		//create a new api error response
		apiErrorResponse := &ApiErrorResponse{
			Status:  resp.StatusCode,
			Message: resp.Status,
			Body:    body,
		}
		//return the error
		return apiErrorResponse
	}

	return nil
}

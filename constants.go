package blnkgo

import (
	"regexp"
	"strconv"
)

type Distribution string

var (
	percentageRegex = regexp.MustCompile(`^\d+%$`)
	numberRegex     = regexp.MustCompile(`^\d+$`)
)

func (d Distribution) IsValid() bool {
	return percentageRegex.MatchString(string(d)) ||
		numberRegex.MatchString(string(d)) ||
		d == "left"
}

func (d Distribution) IsPercentage() bool {
	return percentageRegex.MatchString(string(d))
}

func (d Distribution) IsNumber() bool {
	return numberRegex.MatchString(string(d))
}

func (d Distribution) IsLeft() bool {
	return d == "left"
}

// Convert percentage string to a float64
func (d Distribution) ToPercentage() float64 {
	if d.IsPercentage() {
		// Remove the '%' character and parse the remaining string
		percentage, err := strconv.ParseFloat(string(d[:len(d)-1]), 64)
		if err == nil {
			return percentage
		}
	}
	return 0
}

// Convert number string to a float64
func (d Distribution) ToNumber() float64 {
	if d.IsNumber() {
		number, err := strconv.ParseFloat(string(d), 64)
		if err == nil {
			return number
		}
	}
	return 0
}

// PryTransactionStatus represents the transaction status.
type PryTransactionStatus string

const (
	PryTransactionStatusQueued   PryTransactionStatus = "QUEUED"
	PryTransactionStatusApplied  PryTransactionStatus = "APPLIED"
	PryTransactionStatusRejected PryTransactionStatus = "REJECTED"
	PryTransactionStatusCommit   PryTransactionStatus = "COMMIT"
	PryTransactionStatusVoid     PryTransactionStatus = "VOID"
	PryTransactionStatusInFlight PryTransactionStatus = "INFLIGHT"
	PryTransactionStatusExpired  PryTransactionStatus = "EXPIRED"
)

// InflightStatus represents the status of inflight transactions.
type InflightStatus string

const (
	InflightStatusCommit InflightStatus = "commit"
	InflightStatusVoid   InflightStatus = "void"
)

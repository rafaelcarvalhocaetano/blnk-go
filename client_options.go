package blnkgo

import "time"

type ClientOption func(*Client)

func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.options.Logger = logger
	}
}

func WithRetry(count int) ClientOption {
	return func(c *Client) {
		c.options.RetryCount = count
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.options.Timeout = timeout
	}
}

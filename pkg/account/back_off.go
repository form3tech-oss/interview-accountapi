package account

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// Interface to execute a request and return a response (http.Client implements this interface)
type Service interface {
	Do(req *http.Request) (*http.Response, error)
}

// Strategy to limit the rate of requests and retry on failure.
type LimitRateAndRetry struct {
	MaxRetries *int
	Wait       *int
}

// Returns the maximum number of retries.
// Defaults to 3.
func (lr *LimitRateAndRetry) maxRetries() int {
	if lr == nil || lr.MaxRetries == nil {
		return 3
	}

	return *lr.MaxRetries
}

// Returns the wait time in milliseconds.
// Defaults to 500.
func (lr *LimitRateAndRetry) wait() int {
	if lr == nil || lr.Wait == nil {
		return 500
	}
	return *lr.Wait
}

// Executes the request and returns the response.
// service.Do is called with req to return the response.
// - If the response is 2xx/3xx or 409, it returns the response.
// - If the response is 429, it retries after waiting for the specified time.
// - If the response is 5xx, it retries after waiting for the specified time.
// - If the response is 4xx, it returns the response.
func (lr *LimitRateAndRetry) ExponentialBackOff(service Service, req *http.Request) (*http.Response, error) {

	// If the LimitRateAndRetry is nil, then just execute the request.
	if lr == nil {
		return service.Do(req)
	}

	retries := 0
	retry := false

	// When retrying, the body needs to be read again.
	body, err := readBody(req)
	if err != nil {
		return nil, err
	}

	for ok := true; ok; ok = (retry && (retries <= lr.maxRetries())) {

		// Before retrying, the process will wait in exponential time intervals.
		// The request needs to be cloned using the original context.
		if retry {
			time.Sleep(time.Duration(WaitForRetry(retries, lr.wait())) * time.Millisecond)
			req, err = cloneRequestAndCheckValidContext(req)
			if err != nil {
				return nil, err
			}
		}

		// The body needs to be copied again before executing the request
		copyBody(req, body)

		// back off and retry strategy. If response is 429, 5xx, or an unknown error then retry.
		// translated from form3 api doumentation suggested pseudocode.
		res, err := service.Do(req)
		if err != nil {
			retry = true
		} else if (http.StatusOK <= res.StatusCode && res.StatusCode < http.StatusBadRequest) || res.StatusCode == http.StatusConflict {
			return res, nil
		} else if http.StatusTooManyRequests == res.StatusCode {
			retry = true
		} else if http.StatusInternalServerError <= res.StatusCode {
			retry = true
		} else {
			return res, nil
		}
		retries++
	}

	return nil, fmt.Errorf("[exponential back-off] Max retries (%d) exceeded", lr.maxRetries())
}

// Returns the wait factor for the retry.
func WaitForRetry(retries, wait int) float64 {
	if retries == 0 {
		return 0
	}
	exp := math.Pow(float64(1.5), float64(retries))
	return exp * float64(wait)
}

// Reads the body of the request and to keep it for retries.
func readBody(req *http.Request) ([]byte, error) {
	if req.Body != nil {
		if body, err := io.ReadAll(req.Body); err != nil {
			return nil, err
		} else {
			return body, nil
		}
	}
	return nil, nil
}

// Copies the body to the request. Before retrying, the body needs to be copied again.
func copyBody(req *http.Request, body []byte) {
	if body != nil {
		req.Body = io.NopCloser(bytes.NewReader(body))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
	}
}

// Clone the request, but only if the context is valid.
// If the context is invalid, then fail fast.
func cloneRequestAndCheckValidContext(req *http.Request) (*http.Request, error) {
	ctx := req.Context()
	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
	}
	req = req.Clone(ctx)
	return req, nil
}

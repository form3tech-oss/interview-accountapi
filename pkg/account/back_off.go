package account

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

type Service interface {
	Do(req *http.Request) (*http.Response, error)
}

type LimitRateAndRetry struct {
	MaxRetries *int
	Wait       *int
}

func (lr *LimitRateAndRetry) maxRetries() int {
	if lr == nil || lr.MaxRetries == nil {
		return 3
	}

	return *lr.MaxRetries
}

func (lr *LimitRateAndRetry) wait() int {
	if lr == nil || lr.Wait == nil {
		return 500
	} else {
		return *lr.Wait
	}
}

func (lr *LimitRateAndRetry) ExponentialBackOff(service Service, req *http.Request) (*http.Response, error) {

	if lr == nil {
		return service.Do(req)
	}

	retries := 0
	retry := false

	body, err := readBody(req)
	if err != nil {
		return nil, err
	}

	for ok := true; ok; ok = (retry && (retries <= lr.maxRetries())) {

		if retry {
			time.Sleep(time.Duration(WaitForRetry(retries, lr.wait())) * time.Millisecond)
			req, err = cloneRequestAndCheckValidContext(req)
			if err != nil {
				return nil, err
			}
		}

		copyBody(req, body)

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

func WaitForRetry(retries, wait int) float64 {
	if retries == 0 {
		return 0
	}
	exp := math.Pow(float64(1.5), float64(retries))
	return exp * float64(wait)
}

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

func copyBody(req *http.Request, body []byte) {
	if body != nil {
		req.Body = io.NopCloser(bytes.NewReader(body))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
	}
}

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

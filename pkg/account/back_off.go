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

func (lr *LimitRateAndRetry) ExponentialBackOff(service Service, req *http.Request) (*http.Response, error) {

	if lr == nil {
		return service.Do(req)
	}

	retries := 0
	retry := false
	wait := 500
	maxRetries := 3

	if lr.Wait != nil {
		wait = *lr.Wait
	}

	if lr.MaxRetries != nil {
		maxRetries = *lr.MaxRetries
	}

	var body []byte
	if req.Body != nil {
		if b, err := io.ReadAll(req.Body); err != nil {
			return nil, err
		} else {
			body = b
		}
	}

	for ok := true; ok; ok = (retry && (retries <= maxRetries)) {

		if retry {
			time.Sleep(time.Duration(WaitForRetry(retries, wait)) * time.Millisecond)
			ctx := req.Context()
			if ctx != nil {
				if err := ctx.Err(); err != nil {
					return nil, err
				}
			}
			req = req.Clone(ctx)
		}

		if body != nil {
			req.Body = io.NopCloser(bytes.NewReader(body))
			req.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(body)), nil
			}
		}

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

	return nil, fmt.Errorf("[exponential back-off] Max retries (%d) exceeded", maxRetries)
}

func WaitForRetry(retries, wait int) float64 {
	if retries == 0 {
		return 0
	}
	exp := math.Pow(float64(1.5), float64(retries))
	return exp * float64(wait)
}

package account

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

type Service interface {
	Do(req *http.Request) (*http.Response, error)
}

type LimitRateAndRetry struct {
	MaxRetries *int
	Interval   *int
}

func (lr *LimitRateAndRetry) ExponentialBackOff(service Service, req *http.Request) (*http.Response, error) {

	if lr == nil {
		return service.Do(req)
	}

	retries := 0
	retry := false
	interval := 500
	maxRetries := 3

	if lr.Interval != nil {
		interval = *lr.Interval
	}

	if lr.MaxRetries != nil {
		maxRetries = *lr.MaxRetries
	}

	for ok := true; ok; ok = (retry && (retries <= maxRetries)) {

		if retry {
			timeInterval := time.Duration(interval) * time.Millisecond
			exp := math.Pow(float64(1.5), float64(retries))
			time.Sleep(time.Duration(exp) * timeInterval)
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

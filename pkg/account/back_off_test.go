package account

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackOff(t *testing.T) {

	type testCase struct {
		name    string
		service *mockService
		calls   int
		status  int
		err     error
	}

	testCases := []testCase{
		{
			"succeeds with 200",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 200},
				},
				requests: []*http.Request{},
			},
			1,
			200,
			nil,
		},
		{
			"succeeds with 204",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 204},
				},
				requests: []*http.Request{},
			},
			1,
			204,
			nil,
		},
		{
			"retry once with too many requests",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 429},
					{StatusCode: 200},
				},
				requests: []*http.Request{},
			},
			2,
			200,
			nil,
		},
		{
			"retry once with internal server error",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 500},
					{StatusCode: 200},
				},
				requests: []*http.Request{},
			},
			2,
			200,
			nil,
		},
		{
			"exceed max retries",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 429},
					{StatusCode: 500},
					{StatusCode: 500},
					{StatusCode: 500},
					{StatusCode: 200},
				},
				requests: []*http.Request{},
			},
			4,
			200,
			fmt.Errorf("[exponential back-off] Max retries (3) exceeded"),
		},
		{
			"unknown error",
			&mockService{
				responses: []*http.Response{},
				requests:  []*http.Request{},
			},
			4,
			0,
			fmt.Errorf("[exponential back-off] Max retries (3) exceeded"),
		},
		{
			"bad request",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 400},
				},
				requests: []*http.Request{},
			},
			1,
			400,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req := &http.Request{
				Method: http.MethodGet,
			}
			wait := 5
			maxRetries := 3

			lr := &LimitRateAndRetry{
				MaxRetries: &maxRetries,
				Wait:       &wait,
			}

			res, err := lr.ExponentialBackOff(tc.service, req)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.calls, len(tc.service.requests))
			if err == nil {
				assert.Equal(t, tc.status, res.StatusCode)
			}
		})
	}
}

func TestExponentialBackOffWithTimeOut(t *testing.T) {

	type testCase struct {
		name    string
		service *mockService
		timeout time.Duration
		cancel  bool
		calls   int
		status  int
		err     error
	}

	testCases := []testCase{
		{
			"succeed before context time out",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 429},
					{StatusCode: 500},
					{StatusCode: 200},
				},
				requests:      []*http.Request{},
				expectedBytes: 5,
			},
			500,
			false,
			3,
			200,
			nil,
		},
		{
			"context time out",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 429},
					{StatusCode: 500},
					{StatusCode: 500},
					{StatusCode: 500},
					{StatusCode: 200},
				},
				requests:      []*http.Request{},
				expectedBytes: 5,
			},
			100,
			false,
			2,
			0,
			context.DeadlineExceeded,
		},
		{
			"exceed max retries",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 429},
					{StatusCode: 500},
					{StatusCode: 500},
					{StatusCode: 500},
					{StatusCode: 200},
				},
				requests:      []*http.Request{},
				expectedBytes: 5,
			},
			1000,
			false,
			4,
			0,
			fmt.Errorf("[exponential back-off] Max retries (3) exceeded"),
		},
		{
			"context cancelled",
			&mockService{
				responses: []*http.Response{
					{StatusCode: 429},
					{StatusCode: 500},
					{StatusCode: 200},
				},
				requests:      []*http.Request{},
				expectedBytes: 5,
			},
			100,
			true,
			1,
			0,
			fmt.Errorf("context canceled"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req := &http.Request{
				Method: http.MethodGet,
				Body:   ioutil.NopCloser(bytes.NewReader([]byte("HELLO"))),
				GetBody: func() (io.ReadCloser, error) {
					return ioutil.NopCloser(bytes.NewReader([]byte("HELLO"))), nil
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout*time.Millisecond)
			if tc.cancel {
				cancel()
			} else {
				defer cancel()
			}
			req = req.WithContext(ctx)
			wait := 50
			maxRetries := 3

			lr := &LimitRateAndRetry{
				MaxRetries: &maxRetries,
				Wait:       &wait,
			}

			res, err := lr.ExponentialBackOff(tc.service, req)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.calls, len(tc.service.requests))
			if err == nil {
				assert.Equal(t, tc.status, res.StatusCode)
			}
		})
	}
}

type mockService struct {
	responses     []*http.Response
	requests      []*http.Request
	expectedBytes int
}

func (m *mockService) Do(req *http.Request) (*http.Response, error) {
	if m.expectedBytes > 0 {
		if b, err := io.ReadAll(req.Body); err != nil {
			return nil, err
		} else {
			if len(b) != m.expectedBytes {
				return nil, fmt.Errorf("expected %d bytes, got %d", m.expectedBytes, len(b))
			}
		}
	}

	m.requests = append(m.requests, req)
	if len(m.responses) > 0 {
		res := m.responses[0]
		m.responses = m.responses[1:]
		return res, nil
	}
	return nil, fmt.Errorf("unexpected error")
}

func TestNextInterval(t *testing.T) {

	type testCase struct {
		name     string
		index    int
		expected []float64
	}

	testCases := []testCase{
		{
			"empty",
			0,
			[]float64{0},
		},
		{
			"fout retries",
			4,
			[]float64{
				0,
				1.5,
				2.25,
				3.375,
				5.0625,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			result := []float64{}
			for i := 0; i <= tc.index; i++ {
				result = append(result, WaitForRetry(i, 1))
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

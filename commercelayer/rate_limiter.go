package commercelayer

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

// throttledTransport embeds the underlying transport
// and regulates requests to respect rate limits.
type throttledTransport struct {
	transport http.RoundTripper
	mutex     sync.Mutex
}

// RoundTrip executes an HTTP request, and if the response indicates rate
// limiting (HTTP 429), it waits for the retry interval specified in the
// "X-Ratelimit-Interval" header before retrying. It locks a mutex to ensure
// only one request is processed at a time.
func (t *throttledTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for i := 0; i < 5; i++ {
		resp, err := t.transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		interval, err := time.ParseDuration(resp.Header.Get("X-Ratelimit-Interval") + "s")
		if err != nil {
			return resp, nil
		}

		resp.Body.Close()

		time.Sleep(interval)
	}

	return nil, errors.New("max retries reached when trying to handle rate limiting")
}

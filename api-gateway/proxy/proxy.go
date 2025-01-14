package proxy

import (
	"io"
	"net/http"
)

func ProxyRequest(targetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Copy original header form the request
		for k, v := range r.Header {
			for _, h := range v {
				req.Header.Add(k, h)
			}
		}

		// execute request to micro-service
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Failed to reach service", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		// Copy response from the micro-serive
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

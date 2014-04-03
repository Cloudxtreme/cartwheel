package backends

import "io"
import "time"
import "net"
import "net/http"
import "strings"
import "github.com/tobz/cartwheel/config"

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

type Backend struct {
	config config.BackendConfig

	transport http.RoundTripper
}

func (b *Backend) getRequestForBackend(originalRequest *http.Request) *http.Request {
	// Get a shallow copy of the original request.
	backendRequest := &http.Request{}
	*backendRequest = *originalRequest

	// Make sure we don't tell the backend to close our connection, despite what the
	// incoming request specified.
	backendRequest.Proto = "HTTP/1.1"
	backendRequest.ProtoMajor = 1
	backendRequest.ProtoMinor = 1
	backendRequest.Close = false

	// Remove hop-by-hop headers.  Again, we're following the RFC and we're trying to
	// avoid sending headers that will disable keepalives.  The copy is there in case
	// we need to modify the map, otherwise we just use our shallow copy of the original.
	copiedHeaders := false
	for _, hopHeader := range hopHeaders {
		if backendRequest.Header.Get(hopHeader) != "" {
			if !copiedHeaders {
				backendRequest.Header = &http.Header{}
				copyHeaders(backendRequest.Header, originalRequest.Header)
				copiedHeaders = true
			}

			backendRequest.Header.Del(hopHeader)
		}
	}

	// Set X-Forwarded-For to the client's IP, but carry over the header if it already
	// exists by appending it to the value, comma-separated.
	if clientIP, _, err := net.SplitHostPort(originalRequest.RemoteAddr); err == nil {
		if priorHeader, ok := backendRequest.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(priorHeader, ", ") + ", " + clientIP
		}

		backendRequest.Header.Set("X-Forwarded-For", clientIP)
	}

	return backendRequest
}

func (b *Backend) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Get our modified request to send to the backend.
	backendRequest := b.getRequestForBackend(req)

	// Make the request.
	backendResponse, err := b.transport.RoundTrip(backendRequest)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer backendResponse.Body.Close()

	rw.WriteHeader(backendResponse.StatusCode)

	copyHeaders(rw.Header(), backendResponse.Header)
	copyResponse(rw, backendResponse.Body)
}

func copyHeaders(destination, source http.Header) {
	for k, vv := range source {
		for _, v := range vv {
			destination.Add(k, v)
		}
	}
}

func copyResponse(destination io.Writer, source io.Reader) {
	if IsFlushable(destination) {
		// Hard-coded flushing at 1 second intervals.
		wr := WrapWriterAsFlushable(destination, time.Second)
		wr.Start()
		defer wr.Stop()

		destination = wr
	}

	io.Copy(destination, source)
}

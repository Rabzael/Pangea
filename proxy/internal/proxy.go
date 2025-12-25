package internal

import (
	"io"
	"net"
	"net/http"
	"net/url"
)

/** General handler **/
func ProxyHandler(w http.ResponseWriter, req *http.Request) {
	var well_done bool
	var result any

	if req.Method == http.MethodConnect {
		well_done, result = HttpsProxyHandler(w, req)
	} else {
		well_done, result = ForwardProxyHandler(w, req)
	}

	if well_done {
		LogOk(req, result.(*http.Response))
	} else {
		LogError(req, result.(*url.Error).Err)
	}
}

/** Handles incoming HTTP requests and forwards them to the target server
 */
func ForwardProxyHandler(w http.ResponseWriter, req *http.Request) (bool, any) {
	completeUrl := req.URL.String()
	var response *http.Response
	var error error

	// Forward request based on method
	switch req.Method {
	case "GET":
		response, error = http.Get(completeUrl)
	case "POST":
		response, error = http.Post(completeUrl, req.Header.Get("Content-Type"), req.Body)
	case "HEAD":
		response, error = http.Head(completeUrl)
	default:
		error = http.ErrNotSupported
	}
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return false, error
	}

	// Write response header and body
	w.WriteHeader(response.StatusCode)
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	body, _ := io.ReadAll(response.Body)
	w.Write(body)
	response.Body.Close()

	req.Close = true
	return true, response
}

/** Handle incoming HTTPS requests and forwards them to the target server
 */
func HttpsProxyHandler(w http.ResponseWriter, req *http.Request) (bool, any) {

	// Create tunnel
	hj, done := w.(http.Hijacker)
	if !done {
		http.Error(w, "Can't hijack", http.StatusInternalServerError)
		req.Close = true
		return false, http.ErrNotSupported
	}

	client, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		req.Close = true
		go LogError(req, err)
		return false, err
	}

	defer client.Close()
	client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// Connect to target
	target, err := net.Dial("tcp", req.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		req.Close = true
		go LogError(req, err)
		return false, err
	}
	defer target.Close()

	// Tunnel
	tunnel_chan := make(chan int, 2)
	go func() {
		io.Copy(client, target)
		client.Close()
		tunnel_chan <- 1
	}()
	go func() {
		io.Copy(target, client)
		target.Close()
		tunnel_chan <- 2
	}()
	<-tunnel_chan
	return true, &http.Response{Status: "200 Connection Established", StatusCode: 200}
}

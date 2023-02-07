package main

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
	"time"

	"tinygo.org/x/drivers/net"
	"tinygo.org/x/drivers/net/tls"
)

// -----------------------------------------------------------------------------

type request struct {
	https  bool
	server string
	bytes  []byte
	body   string
	close  bool
}

type response struct {
	bytes []byte
	// code  int
}

type HttpClient struct {
	httpBuf     [4000]byte
	connections map[string]net.Conn
	timeout     time.Duration
}

// -----------------------------------------------------------------------------

func newGET(url string, headers map[string]string) request {
	r := request{}
	parts := strings.SplitN(url, "://", 2)
	proto := parts[0]
	parts = strings.SplitN(parts[1], "/", 2)
	r.server = parts[0]
	path := ""
	if len(parts) > 1 {
		path = parts[1]
	}
	r.https = (proto == "https")
	r.bytes = []byte("GET /" + path + " HTTP/1.1\r\n" +
		"Host: " + strings.Split(r.server, ":")[0] + "\r\n" +
		"User-Agent: HydroStick\r\n" +
		headersToString(headers),
	)
	return r
}

func newPOST(url string, headers map[string]string) request {
	r := request{}
	parts := strings.SplitN(url, "://", 2)
	proto := parts[0]
	parts = strings.SplitN(parts[1], "/", 2)
	r.server = parts[0]
	path := ""
	if len(parts) > 1 {
		path = parts[1]
	}
	r.https = (proto == "https")
	r.bytes = []byte("POST /" + path + " HTTP/1.1\r\n" +
		"Host: " + strings.Split(r.server, ":")[0] + "\r\n" +
		"User-Agent: VST Vibrator\r\n" +
		"Content-Type: application/json\r\n" +
		headersToString(headers),
	)
	return r
}

func headersToString(headers map[string]string) (result string) {
	if headers == nil {
		return
	}
	for name, value := range headers {
		result += name + ": " + value + "\r\n"
	}
	return
}

func dialHttp(https bool, server string) (conn net.Conn, err error) {
	trace(">dialHttp")
	retries := 3
	for {
		if https {
			conn, err = tls.Dial("tcp", server, nil)
		} else {
			conn, err = net.Dial("tcp", server)
		}
		if err == nil || retries == 0 {
			trace("<dialHttp " + strconv.FormatBool(err == nil) + ", " + strconv.FormatBool(retries == 0))
			return
		}
		time.Sleep(1 * time.Second)
		retries--
	}
}

// -----------------------------------------------------------------------------

func (hc *HttpClient) sendHttp(req request, keepAlive bool) (resp response, err error) {
	conn := hc.connections[req.server]
	connNil := strconv.FormatBool(conn == nil)
	ka := strconv.FormatBool(keepAlive)
	defer un(trace("sendHttp " + connNil + ", " + ka))
	if conn == nil {
		conn, err = dialHttp(req.https, req.server)
		if err != nil {
			trace("sendHttp->dialHttp " + err.Error())
			return
		}
		hc.connections[req.server] = conn
	}

	request := []byte{}
	request = append(request, req.bytes...)
	if keepAlive {
		request = append(request, []byte("Connection: keep-alive\r\n")...)
	} else {
		request = append(request, []byte("Connection: close\r\n")...)
	}
	if len(req.body) > 0 {
		request = append(request, []byte("Content-Length: "+strconv.Itoa(len(req.body))+"\r\n\r\n")...)
		request = append(request, []byte(req.body+"\r\n")...)
		trace("POST " + req.body)
	} else {
		request = append(request, []byte("\r\n\r\n")...)
		trace("GET")
	}

	trace("Sending HTTP request...")
	trace(string(request))
	trace("---")

	_, err = conn.Write(request)
	if err != nil {
		trace("sendHttp->Write " + err.Error())
		conn.Close()
		delete(hc.connections, req.server)
		return
	}
	n, err := hc.readHttp(conn)
	if err != nil {
		trace("sendHttp->readHttp " + err.Error())
		conn.Close()
		delete(hc.connections, req.server)
		return
	}
	resp.bytes = hc.httpBuf[:n]
	// println(n)
	//trace(string(resp.bytes[:12]))
	// println(string(resp.bytes))

	if !keepAlive {
		conn.Close()
		delete(hc.connections, req.server)
	}

	return
}

func (hc *HttpClient) readHttp(conn net.Conn) (int, error) {
	read := 0
	timeout := time.Now().Add(hc.timeout)
	for {
		n, err := conn.Read(hc.httpBuf[read:])
		if err != nil {
			return 0, err
		}
		if n == 0 && (read != 0 || time.Now().After(timeout)) {
			break
		}
		if n > 0 {
			read += n
		}
		runtime.Gosched()
	}
	if time.Now().After(timeout) {
		return 0, errors.New("Read timeout")
	}
	return read, nil
}

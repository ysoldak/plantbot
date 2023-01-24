package main

import (
	"strconv"
	"time"

	"tinygo.org/x/drivers/net"
)

type blynk struct {
	server string
	token  string
	client HttpClient
}

func newBlynk() *blynk {
	return &blynk{
		server: blynkEndpoint,
		token:  blynkToken,
		client: HttpClient{
			timeout:     time.Second,
			connections: map[string]net.Conn{},
		},
	}
}

func (b *blynk) updateInt(name string, value int) (err error) {
	url := b.server + "/external/api/update?token=" + b.token + "&" + name + "=" + strconv.Itoa(value)
	req := newGET(url, nil)
	res, err := b.client.sendHttp(req, false)
	if err != nil {
		return err
	} else {
		trace(string(res.bytes))
	}
	return nil
}

func (b *blynk) updateFloat(name string, value float64) (err error) {
	url := b.server + "/external/api/update?token=" + b.token + "&" + name + "=" + strconv.FormatFloat(value, 'f', 2, 64)
	req := newGET(url, nil)
	res, err := b.client.sendHttp(req, false)
	if err != nil {
		return err
	} else {
		trace(string(res.bytes))
	}
	return nil
}

func (b *blynk) sendEvent(name string) (err error) {
	url := b.server + "/external/api/logEvent?token=" + b.token + "&code=" + name
	req := newGET(url, nil)
	res, err := b.client.sendHttp(req, false)
	if err != nil {
		return err
	} else {
		trace(string(res.bytes))
	}
	return nil
}

package talktasy

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Client instance.
type Client struct {
	caCert     string
	clientCert string
	clientKey  string
}

// NewClient returns a new Client instance with the given certs.
func NewClient(caCert, clientCert, clientKey string) *Client {
	return &Client{
		caCert:     caCert,
		clientCert: clientCert,
		clientKey:  clientKey,
	}
}

// Do sends a HTTP request.
// And verifies the response.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	client, err := c.HTTPClient()
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// HTTPClient returns a http client configured with mTLS utilities.
func (c *Client) HTTPClient() (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(c.clientCert, c.clientKey)
	if err != nil {
		return nil, err
	}
	caCert, err := ioutil.ReadFile(c.caCert)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: conf,
		},
	}
	return client, nil
}

// Dial sends a websocket request and verifies the connection.
func (c *Client) Dial(uri string, headers http.Header) (*websocket.Conn, *http.Response, error) {
	cert, err := tls.LoadX509KeyPair(c.clientCert, c.clientKey)
	if err != nil {
		return nil, nil, err
	}
	caCert, err := ioutil.ReadFile(c.caCert)
	if err != nil {
		return nil, nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}
	client := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  conf,
	}
	return client.Dial(uri, headers)
}

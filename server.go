package talktasy

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

// Server instance.
type Server struct {
	caCert     string
	serverCert string
	serverKey  string
	port       int
	pool       []string
}

// NewServer returns a server instance with the given certs.
// Example:
// @caCert = "./ca.crt"
// @serverCert = "./server.crt"
// @serverKey = "./server.key"
func NewServer(caCert, serverCert, serverKey string) *Server {
	return &Server{
		caCert:     caCert,
		serverCert: serverCert,
		serverKey:  serverKey,
	}
}

// Listen starts a HTTP server on the given address.
// Example:
// @addr = "127.0.0.1:8080"
// @fn = nil | func(w http.ResponseWriter, r *http.Request){}
func (s *Server) Listen(addr string, fn http.HandlerFunc) error {
	caCert, err := ioutil.ReadFile(s.caCert)
	if err != nil {
		return err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	conf := &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                certPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}
	conf.BuildNameToCertificate()
	httpServer := &http.Server{
		Addr:      addr,
		TLSConfig: conf,
		Handler:   s.handler(fn),
	}
	return httpServer.ListenAndServeTLS(s.serverCert, s.serverKey)
}

// Handler does pre-processing for a request.
// And then passes to the fn.
func (s *Server) handler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if fn != nil {
			fn(w, r)
		} else {
			// Show debug info when no function passed
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Request verified!"))
		}
	}
}

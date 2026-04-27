package cryptos

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func NewTLSConfigFromFile(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}, nil
}

func NewTLSConfig(cert, key, ca string) (*tls.Config, error) {
	c, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(ca))

	return &tls.Config{
		Certificates: []tls.Certificate{c},
		RootCAs:      caCertPool,
	}, nil
}

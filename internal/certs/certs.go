package certs

import (
	"crypto/tls"
)

type Conf struct {
	// Certificate dir
	// Polling interval
}

type certHandler struct {
	// TODO mutex
}

func Create(conf Conf) (*certHandler, error) {
	panic("not impl")
	return nil, nil
}

func (ch *certHandler) GetClientCertificate(info *tls.CertificateRequestInfo) (tls.Certificate, error) {
	panic("not impl")
	return nil, nil
}

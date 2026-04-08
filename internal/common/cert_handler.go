package common

import (
	"crypto/tls"
)

type CertHandler interface {
	GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error)
}

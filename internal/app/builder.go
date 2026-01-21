package app

import (
	"errors"
    "net"

	miekg "github.com/miekg/dns"

	"git.zluudg.se/piplup/internal/common"
)

type builder struct {
	newApp *application
}

func NewBuilder() *builder {
	newBuilder := new(builder)
	newBuilder.newApp = new(application)

	return newBuilder
}

func (b *builder) Logger(l common.Logger) *builder {
	b.newApp.log = l
	return b
}

func (b *builder) Address(addr string) *builder {
	b.newApp.address = addr
	return b
}

func (b *builder) UdpPort(port string) *builder {
	b.newApp.udpPort = port
	return b
}

func (b *builder) TlsPort(port string) *builder {
	b.newApp.tlsPort = port
	return b
}

func (b *builder) Upstream(addr, port string) *builder {
	b.newApp.upstream = net.JoinHostPort(addr, port)
	return b
}

func (b *builder) MatchSuffix(suffix string) *builder {
	b.newApp.matchSuffix = suffix
	return b
}

func (b *builder) Inject(toInject string) *builder {
	toInjectRR, _ := miekg.NewRR(toInject)

	b.newApp.toInject = toInjectRR

	return b
}

func (b *builder) CertDir(path string) *builder {
	b.newApp.certDir = path

	return b
}


func (b *builder) Build() (*application, error) {
	if b.newApp == nil {
		return nil, errors.New("builder used up")
	}

	if b.newApp.log == nil {
		return nil, errors.New("nil logger")
	}

	if b.newApp.address == "" {
		return nil, errors.New("bad address when creating application")
	}

	if b.newApp.udpPort == "" {
		return nil, errors.New("bad port when creating application")
	}

    // TODO check more

	newApp := b.newApp
	b.newApp = nil

	return newApp, nil
}

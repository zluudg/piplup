package app

import (
	"context"
	"crypto/tls"
	"net"

	miekg "codeberg.org/miekg/dns"

	"git.zluudg.se/piplup/internal/common"
	"git.zluudg.se/piplup/internal/logger"
)

type Conf struct {
	Debug             bool   `json:"debug"`
	Address           string `json:"address"`
	UdpPort           string `json:"udp_port"`
	TlsPort           string `json:"tls_port"`
	UpstreamAddress   string `json:"upstream_address"`
	UpstreamPort      string `json:"upstream_port"`
	UpstreamTransport string `json:"upstream_transport"`
	Log               common.Logger
	Cert              common.CertHandler
}

type appHandle struct {
	log               common.Logger
	cert              common.CertHandler
	id                string
	address           string
	udpPort           string
	tlsPort           string
	upstreamAddress   string
	upstreamPort      string
	upstreamTransport string
}

func Create(c Conf) (*appHandle, error) {
	a := new(appHandle)
	a.id = "main app"

	if c.Log == nil {
		log := logger.New(
			logger.Conf{
				Debug: c.Debug,
			})
		a.log = log
	} else {
		a.log = c.Log
	}
	a.log.Debug("Debug logging enabled for %s", a.id)

	if c.Cert == nil {
		a.log.Error("Bad cert handler handle")
		return nil, common.ErrBadHandle
	}
	a.cert = c.Cert

	if c.Address == "" {
		a.log.Error("No listening address configured")
		return nil, common.ErrBadParam
	}
	a.address = c.Address

	if c.UdpPort == "" {
		a.log.Error("No listening UDP port configured")
		return nil, common.ErrBadParam
	}
	a.udpPort = c.UdpPort

	if c.TlsPort == "" {
		a.log.Error("No listening TLS port configured")
		return nil, common.ErrBadParam
	}
	a.tlsPort = c.TlsPort

	if c.UpstreamAddress == "" {
		a.log.Error("No upstream address configured")
		return nil, common.ErrBadParam
	}
	a.upstreamAddress = c.UpstreamAddress

	if c.UpstreamPort == "" {
		a.log.Error("No upstream port configured")
		return nil, common.ErrBadParam
	}
	a.upstreamPort = c.UpstreamPort

	if c.UpstreamTransport == "" {
		a.log.Error("No upstream transport configured")
		return nil, common.ErrBadParam
	}
	a.upstreamTransport = c.UpstreamTransport

	return a, nil
}

func (a *appHandle) Run(ctx context.Context, exitCh chan<- common.Exit) {
	udpSrv := miekg.NewServer()
	udpSrv.Addr = net.JoinHostPort(a.address, a.udpPort)
	udpSrv.Net = "udp4"
	udpSrv.Handler = a

	go func() {
		err := udpSrv.ListenAndServe()
		if err != nil {
			a.log.Error("UDP listener exited with: %s", err)
		}
	}()

	a.log.Info("Started UDP Listener")

	tlsSrv := miekg.NewServer()
	tlsSrv.Addr = net.JoinHostPort(a.address, a.tlsPort)
	tlsSrv.Net = "tcp4"
	tlsSrv.Handler = a
	tlsSrv.TLSConfig = &tls.Config{
		GetCertificate: a.cert.GetCertificate,
	}

	go func() {
		err := tlsSrv.ListenAndServe()
		if err != nil {
			a.log.Error("TLS listener exited with: %s", err)
		}
	}()

	a.log.Info("Started TLS Listener")

	<-ctx.Done()
	a.log.Info("Shutting down %s", a.id)

	udpSrv.Shutdown(ctx) // TODO Timeout when miekg lib starts to support it
	tlsSrv.Shutdown(ctx) // TODO Timeout when miekg lib starts to support it

	a.log.Info("Leaving app")

	exitCh <- common.Exit{ID: a.id, Err: nil}

	return
}

func (a *appHandle) ServeDNS(ctx context.Context, w miekg.ResponseWriter, r *miekg.Msg) {
	a.log.Debug("Query for %s incoming", r.Question[0].Header().Name)
	var err error
	resp := new(miekg.Msg)
	resp.Question = r.Question
	resp.MsgHeader.ID = r.MsgHeader.ID
	resp.MsgHeader.Response = true
	resp.MsgHeader.Opcode = miekg.OpcodeQuery
	resp.MsgHeader.Authoritative = true
	resp.MsgHeader.Rcode = miekg.RcodeRefused

	if false {
		a.log.Debug("Query did not match filter, won't proxy")
		_, err = resp.WriteTo(w)
		if err != nil {
			a.log.Error("Error responding: %s", err)
		}
		return
	}

	newQ := r
	upResp, err := miekg.Exchange(ctx, newQ, a.upstreamTransport, net.JoinHostPort(a.upstreamAddress, a.upstreamPort))
	if err != nil {
		if err != nil {
			a.log.Error("Error from upsream DNS: %s", err)
		}
		_, err = resp.WriteTo(w) // TODO return some error msg here instead?
		if err != nil {
			a.log.Error("Error responding after failed upstream: %s", err)
		}
		return
	}

	resp.Answer = upResp.Answer
	resp.Ns = upResp.Ns
	resp.Extra = upResp.Extra
	resp.MsgHeader.Rcode = upResp.MsgHeader.Rcode

	_, err = resp.WriteTo(w)
	if err != nil {
		a.log.Error("Error responding after succesful upstream: %s", err)
	}
	return
}

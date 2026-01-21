package app

import (
	"context"
    "net"
	"path/filepath"
	"strings"

	miekg "github.com/miekg/dns"

	"git.zluudg.se/piplup/internal/common"
)

type application struct {
	log         common.Logger
	id          string
	address     string
	port        string
	upstream    string
	toInject    miekg.RR
    certDir     string
	matchSuffix string
	exitCh      chan<- common.Exit
}

func (a *application) Run(ctx context.Context, exitCh chan<- common.Exit) {
	a.id = "main app"
	a.exitCh = exitCh

	go func() {
        err := miekg.ListenAndServe(
        net.JoinHostPort(a.address, a.port),
		"udp4",
		a)
        if err != nil {
            a.log.Error("UDP listener exited with error: %s", err)
        }
    }()

    a.log.Info("Started UDP Listener")

    go func() {
        err := miekg.ListenAndServeTLS(
        net.JoinHostPort(a.address, a.port),
		filepath.Join(a.certDir, "tls.crt"),
		filepath.Join(a.certDir, "tls.key"),
		a)
        if err != nil {
            a.log.Error("UDP listener exited with error: %s", err)
        }
    }()

    a.log.Info("Started TLS Listener")

    <- ctx.Done()

    a.log.Info("Leaving main routine")

	a.exitCh <- common.Exit{ID: a.id, Err: nil}

	return
}

func (a *application) ServeDNS(w miekg.ResponseWriter, r *miekg.Msg) {
	a.log.Debug("Query for %s incoming", r.Question[0].Name)
	var err error
	resp := new(miekg.Msg)
	resp.Question = r.Question
	resp.MsgHdr.Id = r.MsgHdr.Id
	resp.MsgHdr.Response = true
	resp.MsgHdr.Opcode = miekg.OpcodeQuery
	resp.MsgHdr.Authoritative = true
	resp.MsgHdr.Rcode = miekg.RcodeRefused
	if !strings.HasSuffix(r.Question[0].Name, a.matchSuffix) {
		a.log.Debug("Query for %s, won't proxy", r.Question[0].Name)
		err = w.WriteMsg(resp)
		if err != nil {
			a.log.Error("Error responding: %s", err)
		}
		return
	}

	newQ := r
	upResp, err := miekg.Exchange(newQ, a.upstream)
	if err != nil {
		if err != nil {
			a.log.Error("Error from upsream DNS: %s", err)
		}
		err = w.WriteMsg(resp)
		if err != nil {
			a.log.Error("Error responding after failed upstream: %s", err)
		}
		return
	}

	resp.Answer = upResp.Answer
	resp.Ns = upResp.Ns
	resp.Extra = upResp.Extra
	resp.MsgHdr.Rcode = upResp.MsgHdr.Rcode
    if a.toInject != nil {
	    resp.Extra = append(resp.Extra, a.toInject)
    }

	err = w.WriteMsg(resp)
	if err != nil {
		a.log.Error("Error responding after succesful upstream: %s", err)
	}
	return
}

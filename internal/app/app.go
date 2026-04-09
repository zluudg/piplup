package app

import (
	"context"
	"crypto/tls"
	"net"

	miekg "codeberg.org/miekg/dns"

	"git.zluudg.se/piplup/internal/app/action"
	"git.zluudg.se/piplup/internal/app/match"
	"git.zluudg.se/piplup/internal/common"
	"git.zluudg.se/piplup/internal/logger"
)

const c_DEFAULT_ACTION = "{{DEFAULT}}"

type Conf struct {
	Debug             bool          `json:"debug"`
	Address           string        `json:"address"`
	UdpPort           string        `json:"udp_port"`
	TlsPort           string        `json:"tls_port"`
	UpstreamAddress   string        `json:"upstream_address"`
	UpstreamPort      string        `json:"upstream_port"`
	UpstreamTransport string        `json:"upstream_transport"`
	Actions           []action.Conf `json:"actions"`
	Matches           []match.Conf  `json:"matches"`
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
	matchIncoming     []*match.Match
	matchOutgoing     []*match.Match
	actions           map[string]action.Action
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

	a.actions = make(map[string]action.Action)
	for _, aconf := range c.Actions {
		ac, err := action.Create(aconf)
		if err != nil {
			a.log.Error("Could not create action '%s': %s", aconf.ID, err)
			return nil, common.ErrBadParam
		}

		a.actions[aconf.ID] = ac
	}
	defaultAc, err := action.Create(action.Conf{
		ID:      c_DEFAULT_ACTION,
		Forward: true,
		Kind:    "noop",
	})
	if err != nil {
		a.log.Error("Could not create default action")
		return nil, common.ErrNotCompleted
	}
	a.actions[c_DEFAULT_ACTION] = defaultAc

	a.matchIncoming = make([]*match.Match, 0)
	a.matchOutgoing = make([]*match.Match, 0)
	for _, mconf := range c.Matches {
		_, ok := a.actions[mconf.ActionID]
		if !ok {
			a.log.Error("Match referenced unrecognized action '%s'", mconf.ActionID)
			return nil, common.ErrBadParam
		}

		m, err := match.Create(mconf)
		if err != nil {
			a.log.Error("Could not create match object '%s'", mconf.String())
			return nil, common.ErrBadParam
		}
		if mconf.Outgoing {
			a.matchOutgoing = append(a.matchOutgoing, m)
		} else {
			a.matchIncoming = append(a.matchIncoming, m)
		}
	}
	matchDefault, err := match.Create(match.Conf{
		ActionID: c_DEFAULT_ACTION,
	})
	if err != nil {
		a.log.Error("Could not create default match pattern")
		return nil, common.ErrNotCompleted
	}
	a.matchOutgoing = append(a.matchOutgoing, matchDefault)
	a.matchIncoming = append(a.matchIncoming, matchDefault)

	a.log.Info("%d incoming match patterns configured", len(a.matchIncoming))
	a.log.Info("%d outgoing match patterns configured", len(a.matchOutgoing))

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
	var chosenIncAction action.Action
	var chosenOutAction action.Action

	resp := new(miekg.Msg)

	for _, m := range a.matchIncoming {
		if !m.IsMatch(r) {
			continue
		}
		a.log.Debug("Matched incoming pattern %s", m.String())

		ac, ok := a.actions[m.ActionID()]
		if !ok {
			a.log.Error("Invalid Action ID for match '%s'", m.String())
			break
		}

		chosenIncAction = ac
		break
	}

	if chosenIncAction == nil {
		panic("No action found for match")
	}

	incProcessed, err := chosenIncAction.Apply(r)
	if err != nil {
		a.log.Error("Could not apply action: %s, dropping...", err)
		return
	} else {
		a.log.Debug("Succesfully applied action %s on incoming match", chosenIncAction.ID())
	}

	writeRaw := false
	if chosenIncAction.DoForward() {
		upResp, err := miekg.Exchange(ctx, incProcessed, a.upstreamTransport, net.JoinHostPort(a.upstreamAddress, a.upstreamPort))
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

		for _, m := range a.matchOutgoing {
			if m.IsMatch(upResp) {
				a.log.Info("Matched outgoing pattern %s", m.String()) // TODO remove
				break
			}
			ac, ok := a.actions[m.ActionID()]
			if !ok {
				a.log.Error("Invalid Action ID for match '%s'", m.String())
				break
			}

			chosenOutAction = ac
			break
		}

		if chosenOutAction == nil {
			panic("No action found for match")
		}

		outProcessed, err := chosenOutAction.Apply(upResp)
		if err != nil {
			a.log.Error("Could not apply action: %s, dropping...", err)
			return
		} else {
			a.log.Debug("Succesfully applied action %s on outgoing match", chosenOutAction.ID())
		}

		resp = outProcessed
		writeRaw = chosenOutAction.WriteRaw()
	} else {
		resp = incProcessed
		writeRaw = chosenIncAction.WriteRaw()
	}

	if writeRaw {
		_, err := w.Write(resp.Data)
		if err != nil {
			a.log.Error("Error sending raw response: %s", err)
		}
	} else {
		_, err = resp.WriteTo(w)
		if err != nil {
			a.log.Error("Error sending response: %s", err)
		}
	}
	return
}

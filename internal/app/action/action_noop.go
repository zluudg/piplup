package action

import (
	miekg "codeberg.org/miekg/dns"
)

type noopAction struct {
	actionBase
}

func createNoopAction(conf Conf) (Action, error) {
	ac := new(noopAction)
	ac.id = conf.ID
	ac.forward = conf.Forward
	return ac, nil
}

func (a *noopAction) Apply(msg *miekg.Msg) (*miekg.Msg, error) {
	return msg, nil
}

func (a *noopAction) DoForward() bool {
	return a.forward
}

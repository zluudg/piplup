package action

import (
	miekg "codeberg.org/miekg/dns"
)

type injectAction struct {
	actionBase
}

func createInjectAction(conf Conf) (Action, error) {
	ac := new(injectAction)
	ac.id = conf.ID
	ac.forward = conf.Forward
	return ac, nil
}

func (a *injectAction) Apply(msg *miekg.Msg) (*miekg.Msg, error) {
	panic("not impl")
	return msg, nil
}

func (a *injectAction) DoForward() bool {
	return a.forward
}

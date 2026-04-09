package action

import (
	miekg "codeberg.org/miekg/dns"
)

type garbageAction struct {
	actionBase
}

func createGarbageAction(conf Conf) (Action, error) {
	panic("not impl")
	return new(garbageAction), nil
}

func (a *garbageAction) Apply(msg *miekg.Msg) (*miekg.Msg, error) {
	panic("not impl")
	return msg, nil
}

func (a *garbageAction) DoForward() bool {
	return true
}

func (a *garbageAction) WriteRaw() bool {
	return true
}

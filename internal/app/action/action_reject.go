package action

import (
	miekg "codeberg.org/miekg/dns"
)

type refuseAction struct {
	actionBase
}

func createRejectAction(conf Conf) (Action, error) {
	ac := new(rejectAction)
	ac.id = conf.ID
	ac.isChainable = false

	return ac, nil
}

func (a *rejectAction) Apply(msg *miekg.Msg) (*miekg.Msg, error) {
	return msg, nil
}

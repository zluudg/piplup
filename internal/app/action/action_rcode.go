package action

import (
	miekg "codeberg.org/miekg/dns"
)

type rcodeAction struct {
	actionBase
}

func createRcodeAction(conf Conf) (Action, error) {
	ac := new(rcodeAction)
	ac.id = conf.ID
	ac.isChainable = true

	return ac, nil
}

func (a *rcodeAction) Apply(msg *miekg.Msg) (*miekg.Msg, error) {
	return msg, nil
}

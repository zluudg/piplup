package action

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand/v2"
	"slices"

	miekg "codeberg.org/miekg/dns"
)

type garbageAction struct {
	actionBase
}

func createGarbageAction(conf Conf) (Action, error) {
	ac := new(garbageAction)
	ac.id = conf.ID
	ac.forward = conf.Forward

	return ac, nil
}

func (a *garbageAction) Apply(msg *miekg.Msg) (*miekg.Msg, error) {
	outMsg := new(miekg.Msg)

	txID := make([]byte, 2)
	binary.BigEndian.PutUint16(txID, msg.MsgHeader.ID)

	data := make([]byte, rand.IntN(1200))

	_, err := crand.Read(data)
	if err != nil {
		return nil, err
	}

	outMsg.Data = slices.Concat(txID, data)

	return outMsg, nil
}

func (a *garbageAction) DoForward() bool {
	return a.forward
}

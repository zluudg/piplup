package action

import (
	miekg "codeberg.org/miekg/dns"

	"git.zluudg.se/piplup/internal/common"
)

const c_ACTION_NOOP = "noop"
const c_ACTION_GARBAGE = "garbage"

type Conf struct {
	ID      string           `json:"id"`
	Forward bool             `json:"forward"` /* Ignored when action is associated with outgoing match */
	Kind    string           `json:"kind"`
	Inject  []InjectableData `json:"inject_data"`
}

type InjectableData struct {
	Rdata      string `json:"rdata"`
	Section    string `json:"section"`
	SigningKey string `json:"signing_key"`
}

type createActionFn func(Conf) (Action, error)

type Action interface {
	Apply(*miekg.Msg) (*miekg.Msg, error)
	DoForward() bool
	WriteRaw() bool
	ID() string
}

type actionBase struct {
	id      string
	forward bool
}

func (ab *actionBase) ID() string {
	return ab.id
}

var c_CREATE_DISPATCH map[string]createActionFn = map[string]createActionFn{
	c_ACTION_NOOP:    createNoopAction,
	c_ACTION_GARBAGE: createGarbageAction,
}

func Create(conf Conf) (Action, error) {
	create, ok := c_CREATE_DISPATCH[conf.Kind]
	if !ok {
		return nil, common.ErrBadParam
	}

	return create(conf)
}

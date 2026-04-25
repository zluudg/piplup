package action

import (
	miekg "codeberg.org/miekg/dns"

	"git.zluudg.se/piplup/internal/common"
)

const (
	c_ACTION_NOOP    = "noop"
	c_ACTION_GARBAGE = "garbage"
	c_ACTION_INJECT  = "inject"
	c_ACTION_RCODE   = "rcode"
	c_ACTION_REJECT  = "reject"
)

type Conf struct {
	ID     string           `json:"id"`
	Kind   string           `json:"kind"`
	Inject []InjectableData `json:"inject_data"`
	Rcode  string           `json:"rcode"`
}

type InjectableData struct {
	Rdata      string `json:"rdata"`
	Section    string `json:"section"`
	SigningKey string `json:"signing_key"`
}

type createActionFn func(Conf) (Action, error)

type Action interface {
	Apply(*miekg.Msg) (*miekg.Msg, error)
	ID() string
	IsChainable() bool
}

type actionBase struct {
	id          string
	isChainable bool
}

func (ab *actionBase) ID() string {
	return ab.id
}

func (ab *actionBase) IsChainable() bool {
	return ab.isChainable
}

var c_CREATE_DISPATCH map[string]createActionFn = map[string]createActionFn{
	c_ACTION_NOOP:    createNoopAction,
	c_ACTION_GARBAGE: createGarbageAction,
	c_ACTION_INJECT:  createInjectAction,
	c_ACTION_RCODE:   createRcodeAction,
	c_ACTION_REJECT:  createRejectAction,
}

func Create(conf Conf) (Action, error) {
	create, ok := c_CREATE_DISPATCH[conf.Kind]
	if !ok {
		return nil, common.ErrBadParam
	}

	return create(conf)
}

package match

import (
	"fmt"
	"regexp"

	miekg "codeberg.org/miekg/dns"

	"git.zluudg.se/piplup/internal/common"
)

type Conf struct {
	Qname    *string `json:"qname"`
	Qtype    *string `json:"qtype"`
	Outgoing bool    `json:"match_outgoing"`
	ActionID string  `json:"action"`
}

func (c Conf) String() string {
	qnameRepr := "<NONE>"
	qtypeRepr := "<NONE>"

	if c.Qname != nil {
		qnameRepr = *c.Qname
	}

	if c.Qtype != nil {
		qtypeRepr = *c.Qtype
	}

	return fmt.Sprintf("%s/%s", qnameRepr, qtypeRepr)
}

type Match struct {
	regex    *regexp.Regexp
	qtype    *uint16
	actionID string
	str      string
}

func Create(conf Conf) (*Match, error) {
	var err error
	qnameRepr := "<NONE>"
	qtypeRepr := "<NONE>"
	m := new(Match)

	if conf.Qname != nil {
		m.regex, err = regexp.Compile(*conf.Qname)
		if err != nil {
			return nil, common.ErrBadParam
		}
		qnameRepr = *conf.Qname
	}

	if conf.Qtype != nil {
		typeToMatch, ok := miekg.StringToType[*conf.Qtype]
		if !ok {
			return nil, common.ErrBadParam
		}
		m.qtype = &typeToMatch
		qtypeRepr = *conf.Qtype
	}

	if conf.ActionID == "" {
		return nil, common.ErrBadParam
	}
	m.actionID = conf.ActionID

	m.str = fmt.Sprintf("%s/%s", qnameRepr, qtypeRepr)

	return m, nil
}

func (m *Match) String() string {
	return m.str
}

func (m *Match) IsMatch(msg *miekg.Msg) bool {
	isMatch := true

	if m.qtype == nil && m.regex == nil {
		return true
	}

	if len(msg.Question) == 0 {
		return false
	}

	query := msg.Question[0]

	if m.qtype != nil {
		typeMsg := miekg.RRToType(query)
		isMatch = isMatch && *m.qtype == typeMsg
	}

	if m.regex != nil {
		h := query.Header()
		if h == nil {
			return false
		}
		nameMsg := h.Name

		isMatch = isMatch && m.regex.MatchString(nameMsg)
	}

	return isMatch
}

func (m *Match) ActionID() string {
	return m.actionID
}

package match

import (
	"fmt"
	"regexp"

	miekg "codeberg.org/miekg/dns"

	"git.zluudg.se/piplup/internal/app/action"
	"git.zluudg.se/piplup/internal/common"
)

type Match struct {
	regex   *regexp.Regexp
	qtype   *uint16
	actions []action.Action
	str     string
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

	if len(conf.Actions) == 0 {
		return nil, common.ErrBadParam
	}
	m.actions = conf.Actions

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

func (m *Match) Actions() []string {
	return m.actions
}

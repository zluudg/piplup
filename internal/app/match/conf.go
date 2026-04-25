package match

import (
	"fmt"
)

type Conf struct {
	Qname    *string  `json:"qname"`
	Qtype    *string  `json:"qtype"`
	Outgoing bool     `json:"match_outgoing"`
	ActionIDs  []string `json:"actions"`
	Forward  bool     `json:"forward"`
    Actions []action.Conf
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

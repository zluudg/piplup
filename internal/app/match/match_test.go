package match

import (
	"testing"

	//miekg "codeberg.org/miekg/dns"

	"git.zluudg.se/piplup/internal/common"
)

func TestCreate_BadConf(t *testing.T) {
	mconf := Conf{}

	_, err := Create(mconf)

	expected := common.ErrBadParam
	if err != expected {
		t.Logf("Expected: %s, Got: %s", expected, err)
	}
}

//func TestIsMatch_ResponseMismatch(t *testing.T) {
//    msg := miekg.Msg{
//        MsgHeader: miekg.MsgHeader{Response: true},
//        Question:  []miekg.Question{},
//    }
//
//	m, err := match.Create(mconf)
//	if err != nil {
//		t.Logf("Could not create match object '%s/%s'", *mconf.Qname, *mconf.Qtype)
//	}
//
//    m := app.Match{
//        Response: false,
//    }
//
//    // Should return false if response headers do not match [1]
//    result := app.IsMatch(&msg, m)
//    if result {
//        t.Errorf("Expected isMatch to return false when Response headers mismatch [1]")
//    }
//}
//
//func TestIsMatch_EmptyQuestion(t *testing.T) {
//    msg := miekg.Msg{
//        MsgHeader: miekg.MsgHeader{Response: true},
//        Question:  []miekg.Question{},
//    }
//    m := app.Match{
//        Response: true,
//    }
//
//    // Should return false if question slice is empty [1]
//    result := app.IsMatch(&msg, m)
//    if result {
//        t.Errorf("Expected isMatch to return false when Question list is empty [1]")
//    }
//}
//
//func TestIsMatch_RegexHeaderNil(t *testing.T) {
//    msg := miekg.Msg{
//        MsgHeader: miekg.MsgHeader{Response: true},
//        Question:  []miekg.Question{
//            {Header: &miekg.QuestionHeader{Name: "example.com"}}, // Example structure
//        },
//    }
//    m := app.Match{
//        Response: true,
//        Regex:    nil, // If regex is not set, behavior depends on nil check
//    }
//    // Ensure header is not nil, but if it were nil, isMatch would return false [1]
//    // Note: Specific logic for nil regex or header check relies on h == nil [1]
//}
//
//func TestIsMatch_Success(t *testing.T) {
//    msg := miekg.Msg{
//        MsgHeader: miekg.MsgHeader{Response: true},
//        Question: []miekg.Question{
//            {Header: &miekg.QuestionHeader{Name: "example.com", Qtype: miekg.TypeA}},
//        },
//    }
//    m := app.Match{
//        Response:  true,
//        Qtype:     (*uint16)(uint16(miekg.TypeA)),
//        // Regex logic handled similarly to source [1]
//    }
//
//    // This test would ensure the function passes when conditions align [1]
//    result := app.IsMatch(&msg, m)
//    if !result {
//        t.Errorf("Expected isMatch to return true when all conditions are met [1]")
//    }
//}

package app

//const c_ACTION_NOOP = "noop"
//const c_ACTION_GARBAGE = "garbage"
//const c_ACTION_INJECT = "inject"
//var c_ACTION_KINDS = []string{c_ACTION_NOOP, c_ACTION_GARBAGE, c_ACTION_INJECT}
//
//type ActionConf {
//    ID string `json:"id"`
//    Kind string `json:"kind"`
//
//}
//
//type noopAction {
//    id string
//}
//
//type action interface {
//    Apply(dnsMsg miekg.Msg) ([]byte)
//}
//
//type createAction(c ActionConf) (action, error) {
//    if c.Kind == c_ACTION_NOOP {
//        return createNoopAction(c)
//    } else if c.Kind == c_ACTION_GARBAGE {
//        return createGarbageAction(c)
//    } else {
//        return nil, common.ErrBadParam
//    }
//}
//
//func (a noopAction) Apply(dnsMsg miekg.Msg) dnsMsg miekg.Msg {
//    panic("not impl")
//    return nil
//}

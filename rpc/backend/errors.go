package backend

type RpcBizErr struct {
	Code int32
	Msg  string
}

func NewRpcBizErr(code int32, msg string) *RpcBizErr { return &RpcBizErr{Code: code, Msg: msg} }

func (e *RpcBizErr) Error() string {
	return e.Msg
}

func (e *RpcBizErr) ErrCode() int32 {
	return e.Code
}

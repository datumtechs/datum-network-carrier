package rpc


type RpcBizErr struct {
	Msg string
}
func NewRpcBizErr(msg string) *RpcBizErr {return &RpcBizErr{Msg: msg}}
func (e *RpcBizErr) Error() string {
	return e.Msg
}


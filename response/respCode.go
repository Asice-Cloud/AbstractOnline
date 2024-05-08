package response

type ReCode int64

const (
	CodeSuccess           ReCode = 1000
	CodeInvalidParams     ReCode = 1001
	CodeUserExist         ReCode = 1002
	CodeUserNotExist      ReCode = 1003
	CodeInvalidPassword   ReCode = 1004
	CodeServerBusy        ReCode = 1005
	CodeInvalidToken      ReCode = 1006
	CodeInvalidAuthFormat ReCode = 1007
	CodeNotLogin          ReCode = 1008
)

var msgFlags = map[ReCode]string{
	CodeSuccess:           "success",
	CodeInvalidParams:     "query parameter wrong",
	CodeUserExist:         "repeat user_name",
	CodeUserNotExist:      "user not exists",
	CodeInvalidPassword:   "user_name or password wrong",
	CodeServerBusy:        "service busy",
	CodeInvalidToken:      "invalid Token",
	CodeInvalidAuthFormat: "authorization format wrong",
	CodeNotLogin:          "not login",
}

func (re ReCode) Msg() string {
	msg, ok := msgFlags[re]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}

package atmz

const (
	JobCheckKey              = "checkKey"
	JobAddUser               = "addUser"
	JobPublishCheckKeyResult = "publishCheckKeyResult"
	JobPublishAddUserResult  = "publishAddUserResult"
	JobError                 = "error"
	JobResult                = "result"
)

const (
	ErrorAddUser  = "ErrorAddUser"
	ErrorCheckKey = "ErrorCheckKey"
)

const (
	MsgUserAdded    = "MsgUserAdded"
	MsgUserFound    = "MsgUserFound"
	MsgUserNotFound = "MsgUserNotFound"
)

type processVariables struct {
	Key        string `json:"key"`
	KYC        string `json:"kyc"`
	UserID     string `json:"userId"`
	Industrial bool   `json:"industrial"`
	UserExists bool   `json:"userExists"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

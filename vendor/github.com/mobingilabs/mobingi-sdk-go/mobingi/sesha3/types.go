package sesha3

type TokenPayload struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

type ExecScriptInstanceResponse struct {
	Ip     string `json:"ip"`
	CmdOut []byte `json:"out"`
	Err    error  `json:"err"`
}

type ExecScriptStackResponse struct {
	StackId string                       `json:"stack_id"`
	Outputs []ExecScriptInstanceResponse `json:"outputs"`
}

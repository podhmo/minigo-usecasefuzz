package app

type Req struct {
	ForceReq *int   `json:"a" required:"true"`
	ForceOpt int    `json:"b" required:"false"`
	KeepReq  int    `json:"c" required:"true"`
	BadBool  int    `json:"d" required:"notabool"`
	OmitReq  *int   `json:"e,omitempty" required:"true"`
}

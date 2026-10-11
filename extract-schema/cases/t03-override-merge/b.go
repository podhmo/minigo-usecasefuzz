package app

type Bounded struct {
	Level int    `json:"level" jsonschema-override:"minimum=0;maximum=10"`
	Pat   string `json:"pat" jsonschema-override:"pattern=^[a-z]+$"`
	Re    string `json:"re" jsonschema-override:"type=integer;format=int64"`
	Req   *int   `json:"req" jsonschema-override:"required=true;minimum=1"`
}

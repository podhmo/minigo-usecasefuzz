package app

type A struct {
	*B
	AName string `json:"a_name"`
}

type B struct {
	*A
	BName string `json:"b_name"`
}

package app

type Bad struct {
	X int `json:"x" jsonschema-override:"novalue"`
}

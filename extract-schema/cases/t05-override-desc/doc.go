package app

type Doc struct {
	// X carries a doc comment the override should beat.
	X int `json:"x" jsonschema-override:"description=from tag"`
}

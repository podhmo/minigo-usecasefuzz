package app

type Tagged struct {
	JsonOnly string `json:"j"`
	YamlOnly string `yaml:"y"`
	TomlOnly string `toml:"t"`
	Conflict string `json:"jname" yaml:"yname" toml:"tname"`
	Omit     string `json:",omitempty"`
	Str      int    `json:",string"`
	Inline   string `json:",inline"`
	Renamed  string `json:"renamed,omitempty"`
	Skipped  string `json:"-"`
	Other    string `schema:"sc" xml:"x"`
}

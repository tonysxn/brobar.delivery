package models

type Setting struct {
	Key   string `db:"key" json:"key"`
	Type  string `db:"type" json:"setting_type"` // "string", "text", "json", "number"
	Value string `db:"value" json:"value"`
}

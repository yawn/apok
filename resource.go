package apok

type Resource struct {
	ARN               string   `json:"ARN"`
	ContextKeys       []string `json:"ContextKeys"`
	ContextKeyStrings []string `json:"ContextKeyStrings"`
	IsRequired        bool     `json:"isRequired"`
	Name              string   `json:"Name"`
	RegEx             string   `json:"RegEx"`
}

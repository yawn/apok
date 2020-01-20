package apok

type Action struct {
	// Allowed               bool         `json:"allowed"`
	ContextKeys []ContextKey `json:"contextKeys"`
	// DeniedByOrganization  bool         `json:"deniedByOrganization"`
	// Evaluated             bool         `json:"evaluated"`
	Name                  string     `json:"name"`
	NoResourceARN         string     `json:"noResourceARN"`
	RequiredResourceNames []string   `json:"requiredResourceNames"`
	ResourceEnabled       bool       `json:"resourceEnabled"`
	ServiceAware          bool       `json:"serviceAware"`
	SupportedResources    []Resource `json:"supportedResources"`
}

func (a Action) Var() string {
	return a.Name
}

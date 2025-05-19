package opksshplugingoogleworkspace

type (
	PolicyPrincipal struct {
		User  []string `json:"users,omitempty"  yaml:"users,omitempty"`
		Group []string `json:"groups,omitempty" yaml:"groups,omitempty"`
	}

	Policy map[string]*PolicyPrincipal
)

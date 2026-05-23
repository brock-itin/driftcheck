package drift

// RuleSet defines which drift checks are enabled.
type RuleSet struct {
	CheckImage bool
	CheckEnv   bool
	CheckPorts bool
}

// DefaultRuleSet returns a RuleSet with all checks enabled.
func DefaultRuleSet() RuleSet {
	return RuleSet{
		CheckImage: true,
		CheckEnv:   true,
		CheckPorts: true,
	}
}

// RuleSetFromFlags builds a RuleSet from individual flag values.
// If no flags are explicitly enabled, all checks default to true.
func RuleSetFromFlags(image, env, ports bool, anySet bool) RuleSet {
	if !anySet {
		return DefaultRuleSet()
	}
	return RuleSet{
		CheckImage: image,
		CheckEnv:   env,
		CheckPorts: ports,
	}
}

// Enabled returns true if at least one check is enabled.
func (r RuleSet) Enabled() bool {
	return r.CheckImage || r.CheckEnv || r.CheckPorts
}

// ActiveChecks returns the list of FindingType values that are active
// under this RuleSet.
func (r RuleSet) ActiveChecks() []FindingType {
	var types []FindingType
	if r.CheckImage {
		types = append(types, FindingTypeImage)
	}
	if r.CheckEnv {
		types = append(types, FindingTypeEnv)
	}
	if r.CheckPorts {
		types = append(types, FindingTypePorts)
	}
	return types
}

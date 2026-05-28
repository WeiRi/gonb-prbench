package v1

import "sort"

// FIX (PR #136685): copy the caller's verbs slice into a new slice owned by
// the builder so sort.Strings doesn't touch shared backing array.

type PolicyRule struct {
	Verbs     []string
	Resources []string
	APIGroups []string
}

type PolicyRuleBuilder struct {
	verbs     []string
	resources []string
	groups    []string
}

func NewRule(verbs ...string) *PolicyRuleBuilder {
	// FIX: copy into a new slice so sort doesn't write caller's array.
	v := make([]string, len(verbs))
	copy(v, verbs)
	return &PolicyRuleBuilder{verbs: v}
}

func (b *PolicyRuleBuilder) Resources(res ...string) *PolicyRuleBuilder {
	b.resources = append(b.resources, res...)
	return b
}

func (b *PolicyRuleBuilder) Groups(g ...string) *PolicyRuleBuilder {
	b.groups = append(b.groups, g...)
	return b
}

func (b *PolicyRuleBuilder) Rule() (PolicyRule, error) {
	sort.Strings(b.verbs)
	return PolicyRule{Verbs: b.verbs, Resources: b.resources, APIGroups: b.groups}, nil
}

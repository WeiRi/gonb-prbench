package v1

import "sort"

// Stripped reproduction of pkg/apis/rbac/v1/helpers.go pre-PR #136685.
// BUG: PolicyRuleBuilder.Rule() calls sort.Strings on Verbs, which writes
// to the underlying array; multiple builders sharing the same backing array
// race on element writes.

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
	return &PolicyRuleBuilder{verbs: verbs}   // BUG: keeps caller's slice header — shared array
}

func (b *PolicyRuleBuilder) Resources(res ...string) *PolicyRuleBuilder {
	b.resources = append(b.resources, res...)
	return b
}

func (b *PolicyRuleBuilder) Groups(g ...string) *PolicyRuleBuilder {
	b.groups = append(b.groups, g...)
	return b
}

// Rule — BUG: sort.Strings writes the shared backing array.
func (b *PolicyRuleBuilder) Rule() (PolicyRule, error) {
	sort.Strings(b.verbs)                     // line 38 — racing write
	return PolicyRule{Verbs: b.verbs, Resources: b.resources, APIGroups: b.groups}, nil
}

// Pre-fix audit.go from PR #106045 (apiserver/admission audit).
// BUG: auditHandler.logAnnotations writes to ae.Annotations map without a
// lock; concurrent admission handlers race on the map.
package admission

type Event struct {
	Annotations map[string]string
}

type auditHandler struct {
	ae *Event
}

func NewAuditHandler() *auditHandler {
	return &auditHandler{ae: &Event{Annotations: map[string]string{}}}
}

// logAnnotations — audit.go:93/94 in pre-fix (no mutex).
func (h *auditHandler) logAnnotations(key, value string) {
	if h.ae == nil {
		return
	}
	// PRE-FIX: writes ae.Annotations without lock.
	h.ae.Annotations[key] = value // racy WRITE
}

// Admit — audit.go:57/58 in pre-fix calls logAnnotations.
func (h *auditHandler) Admit(key, value string) {
	h.logAnnotations(key, value)
}

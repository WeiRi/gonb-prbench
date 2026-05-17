package hugolib

type pageMaps struct {
	id int
}

type HugoSites struct {
	content *pageMaps
}

func newPageMaps() *pageMaps { return &pageMaps{} }

// hugolib/site.go:1306
func (h *HugoSites) readAndProcessContent() {
	_ = h.content
	h.content = newPageMaps()
}

// hugolib/hugo_sites.go:256
func (h *HugoSites) GetContentPage() {
	_ = h.content
}

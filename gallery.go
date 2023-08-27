package static

type Gallery struct {
	Meta  `json:"meta,omitempty"`
	Media []*Media `json:"media"`
}

type Meta struct {
	Title   string   `json:"title"`
	Creator string   `json:"creator"`
	Subject []string `json:"subject"`
}

func NewGallery(dir string) *Gallery {
	return &Gallery{
		Meta: Meta{},
	}
}

func (g *Gallery) WithTitle(title string) *Gallery {
	g.Title = title
	return g
}

func (g *Gallery) WithCreator(creator string) *Gallery {
	g.Creator = creator
	return g
}

func (g *Gallery) WithSubject(subs ...string) *Gallery {
	g.Subject = subs
	return g
}

func (g *Gallery) WithThumbs() *Gallery {
	if len(g.Media) > 0 {
		for _, m := range g.Media {
			m.WithThumb()
		}
	}
}

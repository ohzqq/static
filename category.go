package static

type Category struct {
	Html
	Ext      []string `toml:"ext"`
	Mime     string   `toml:"mime"`
	Template string   `toml:"template"`
}

func DefaultCategory() Category {
	return Category{
		Html: DefaultHtml(),
	}
}

func (c Category) RecursiveWrite(pages ...*Page) error {
	for _, p := range pages {
		err := Write(p.Path, c.RenderPage(p))
		if err != nil {
			return err
		}

		if p.HasChildren() {
			err := c.RecursiveWrite(p.Children...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Categories() map[string]Category {
	col := Default.Category
	for n, c := range Opts.Category {
		col[n] = c
	}
	return col
}

func GetCategory(cat string) Category {
	if c, ok := Opts.Category[cat]; ok {
		return c
	}

	if c, ok := Default.Category[cat]; ok {
		return c
	}

	return Category{
		Html: DefaultHtml(),
	}
}

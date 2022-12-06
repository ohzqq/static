package static

type Category struct {
	Html
	Ext      []string `toml:"ext"`
	Mime     string   `toml:"mime"`
	Template string   `toml:"template"`
	Index
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

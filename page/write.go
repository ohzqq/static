package page

import (
	"fmt"
	"os"
	"path/filepath"
)

func Write(path string, page []byte) error {
	out := filepath.Join(path, "index.html")

	err := os.WriteFile(out, page, 0666)
	if err != nil {
		return fmt.Errorf("Rendering %s failed with error %s\n", out, err)
	}
	fmt.Printf("Rendered %s\n", out)
	return nil
}

func RecursiveWrite(pages ...*Page) error {
	for _, page := range pages {
		err := Write(page.Path, page.Render())
		if err != nil {
			return err
		}

		if page.HasChildren() {
			err := RecursiveWrite(page.Children...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

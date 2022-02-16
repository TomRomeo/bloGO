package errors

import "fmt"

type FrontMatterMissingError struct {
	FileName string
}

func (e *FrontMatterMissingError) Error() string {
	return fmt.Sprintf("Missing YML Frontmatter: %s", e.FileName)
}

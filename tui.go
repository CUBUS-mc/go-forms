package go_forms

import "github.com/charmbracelet/huh"

func fieldsToHuh(fields []Field, form *Form) *huh.Form  {
	var formItems []*huh.Group

	for _, field := range fields {
		switch field := field.(type) {
		case *FieldBaseType:
			// Do nothing
		case *TextField:
			huhField := huh.NewInput().
				Title(field.Prompt).
				Prompt("?").
				Value(&field.Value).
				Validate(func(_ string) error {if field.IsValid() {return nil} else {return field.GetError()}})
			formItems = append(formItems, huhField)
		}
	}

	return huh.NewForm(...formItems)
}

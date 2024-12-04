package go_forms

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func refreshForm(form *Form, box *fyne.Container, fyneForm *widget.Form) {
	fyneForm.Items = nil
	fyneForm.Refresh()
	fyneForm.Items = fieldsToFyneForm(form.GetFieldsToDisplay(), form, box, fyneForm)
	fyneForm.Refresh()
	box.Refresh()
}

func fieldsToFyneForm(fields []Field, form *Form, box *fyne.Container, fyneForm *widget.Form) []*widget.FormItem {
	var formItems []*widget.FormItem

	for _, field := range fields {
		switch field := field.(type) {
		case *FieldBaseType:
			// Do nothing
		case *TextField:
			entry := widget.NewEntry()
			entry.SetText(field.GetValue())
			entry.SetPlaceHolder(field.GetPlaceholder())
			entry.OnChanged = func(text string) {
				field.SetValue(text)
				refreshForm(form, box, fyneForm)
			}
			entry.Validator = func(text string) error {
				if !field.IsValid() {
					return field.GetError()
				}
				return nil
			}
			formItems = append(formItems, widget.NewFormItem(field.GetPrompt(), entry))
		case *MultipleChoiceField:
			labelsToKeys := make(map[string]string)
			options := make([]string, 0, len(field.GetOptions()))
			for key, option := range field.GetOptions() {
				options = append(options, option.Label)
				labelsToKeys[option.Label] = key
			}
			selectWidget := widget.NewSelect(options, func(value string) {
				key := labelsToKeys[value]
				field.SetValue(key)
			})
			selectWidget.SetSelected(field.Options[field.GetValue()].Label)
			selectWidget.OnChanged = func(value string) {
				key := labelsToKeys[value]
				field.SetValue(key)
				refreshForm(form, box, fyneForm)
			}
			formItems = append(formItems, widget.NewFormItem(field.GetPrompt(), selectWidget))
		case *Message:
			formItems = append(formItems, widget.NewFormItem(field.GetValue(), widget.NewLabel("")))
		case *NumberField:
			entry := widget.NewEntry()
			entry.SetText(field.GetValue())
			entry.SetPlaceHolder(field.GetPlaceholder())
			entry.OnChanged = func(text string) {
				field.SetValue(text)
				refreshForm(form, box, fyneForm)
			}
			entry.Validator = func(text string) error {
				if !field.IsValid() {
					return field.GetError()
				}
				return nil
			}
			formItems = append(formItems, widget.NewFormItem(field.GetPrompt(), entry))
		case *FieldGroup:
			if field.GetHeading() != "" {
				formItems = append(formItems, widget.NewFormItem(field.GetHeading(), widget.NewLabel("")))
			}
			formItems = append(formItems, fieldsToFyneForm(field.GetFieldsToDisplay(), form, box, fyneForm)...)
		default:
			panic("Unknown field type")
		}
	}
	return formItems
}

func FormToFyneForm(
	form *Form,
	box *fyne.Container,
	parentDialog *dialog.CustomDialog,
	window fyne.Window,
	onSubmit func(values map[string]string),
) {
	fields := form.GetFieldsToDisplay()
	fyneForm := widget.NewForm()
	fyneForm.Items = fieldsToFyneForm(fields, form, box, fyneForm)
	fyneForm.OnSubmit = func() {
		if form.IsValid() {
			onSubmit(
				form.GetFieldValues(),
			)
			parentDialog.Hide()
		} else {
			dialog.ShowError(form.GetError(), window)
		}
	}
	fyneForm.OnCancel = func() {
		parentDialog.Hide()
	}
	fyneForm.Resize(fyne.NewSize(700, 400))
	box.RemoveAll()
	box.Add(fyneForm)
	box.Refresh()
}

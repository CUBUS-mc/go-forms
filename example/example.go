package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	forms "go-forms"
)

func main() {
	// Define a simple form
	exampleForm := forms.NewForm(
		forms.NewMessage(
			"message", // Id
			[]forms.DisplayCondition{&forms.AlwaysDisplay{}}, // Display conditions (always display)
			"This is a simple form.",                         // Message
		),
		forms.NewTextField(
			"name", // Id
			[]forms.DisplayCondition{&forms.AlwaysDisplay{}}, // Display conditions (always display)
			[]forms.Validator{&forms.NotEmptyValidator{}},    // Validators (not empty)
			"John Doe", // Placeholder
			"Name: ",   // Prompt
			"",         // Default value
		),
		forms.NewNumberField(
			"age", // Id
			[]forms.DisplayCondition{&forms.DisplayAfter{FieldId: "name"}},                                             // Display conditions (display after the name field is filled out)
			[]forms.Validator{&forms.MinValidator{Min: 0}, &forms.MaxValidator{Max: 150}, &forms.IsIntegerValidator{}}, // Validators (min 0, max 150, integer)
			"e.g. 42", // Placeholder
			"Age: ",   // Prompt
			-1,        // Default value
		),
		forms.NewMultipleChoiceField(
			"color", // Id
			[]forms.DisplayCondition{&forms.DisplayAfter{FieldId: "age"}}, // Display conditions (display after the age field is filled out)
			[]forms.Validator{&forms.ChoiceValidator{}},                   // Validators (choice validator to force user to choose a color)
			"Choose a color", // Placeholder
			"Color: ",        // Prompt
			map[string]forms.Option{
				"red":   {"Red", "The color red."},
				"green": {"Green", "The color green."},
				"blue":  {"Blue", "The color blue."},
			}, // Options
			"", // Default value
		),
	)

	// Create a new form
	formApp := app.New()
	win := formApp.NewWindow("Form Example")
	win.CenterOnScreen() // Configure the window
	win.Resize(fyne.Size{Width: 900, Height: 600})
	win.Show() // Show the window

	formSubmitCallback := func( // Callback function to handle form submission
		values map[string]string,
	) {
		for key, value := range values {
			println(key + ": " + value)
		}
	}

	forms.FormToFynePopup( // Convert the form to a Fyne popup
		"Example Form",                     // Title
		fyne.Size{Width: 400, Height: 400}, // Size
		exampleForm,                        // Form
		win,                                // Parent window
		formSubmitCallback,                 // Submit callback
		func() {},                          // Cancel callback (empty)
	)

	formApp.Run() // Run the app
}

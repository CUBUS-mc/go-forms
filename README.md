# go-forms

A simple form builder for Go.

## Features
- Simple form definition
- Form validation
- Form rendering (currently fyne only)
- Display conditions for form fields
- Many form field types

## Usage
Look at the `example` directory for a simple example.

## Field types

### FieldBase
The base field type that all other fields inherit from. DO NOT USE DIRECTLY.

Properties:
- `Id` (string): Unique identifier for the field.
- `DisplayCondition` (string): Condition for displaying the field. If the condition is not met, the field will not be displayed.
- `Validators` ([]Validator): List of validators for the field.
- `Value` (string): Value of the field.
- `form` (*Form): Reference to the form that the field belongs to.
- `error` (error): Error message for the field if validation fails.

Available display conditions
- `AlwaysDisplay`: Always display the field.
- `CustomDisplayCondition`: Custom display condition. Takes a function with the prop `field any` that returns a boolean as the `Condition`.
- `IsValidDisplayCondition`: Display the field if the fields given in `FieldIds` are valid.
- `IsInvalidDisplayCondition`: Display the field if the fields given in `FieldIds` are invalid.
- `AllFieldsValidDisplayCondition`: Display the field if all fields in the form are valid.
- `HasValueDisplayCondition`: Display the field if the field with the given `FieldId` has the given `Value`.
- `DisplayAfter`: Display the field after the field with the given `FieldId` is visible and valid (entry finished).
- `OrDisplayCondition`: Display the field if any of the given `Conditions` are met.
- `AndDisplayCondition`: Display the field if all the given `Conditions` are met.

Available validators
- `CustomValidator`: Custom validator. Takes a function with the prop `field any` that returns (bool, error) as the `Validator`.
- `AllFieldsVaild`: Validate the field if all fields in the form are valid.
- `IsValidValidator`: Validate the field if the fields given in `FieldIds` are valid.

### Message
A simple message field that can be used to display text. No user input is possible.

### Text
A simple text input field.

Properties:
- `Placeholder` (string): Placeholder text for the input field.
- `Prompt` (string): Prompt text for the input field.

Available validators
- `NotEmptyValidator`: Validate that the field is not empty.
- `MinLengthValidator`: Validate that the field has a `MinLength` number of characters.
- `MaxLengthValidator`: Validate that the field has a `MaxLength` number of characters.
- `IpValidator`: Validate that the field is a valid IPv4 address.
- `RegexValidator`: Validate that the field matches the given `RegexPattern` (string).
- `UrlValidator`: Validate that the field is a valid URL. (matches the regex `^https?://.`)

### Number
A number input field that only accepts numbers.

Available validators
- `MinValidator`: Validate that the field is greater than or equal to `Min`.
- `MaxValidator`: Validate that the field is less than or equal to `Max`.
- `IsIntegerValidator`: Validate that the field is an integer.

### MultipleChoice
A multiple choice field that allows the user to select one option from a list.

Properties:
- `Options` (map[string]Option): List of options to choose from.

Option has the following properties:
- `Label` (string): Label for the option.
- `Description` (string): Description for the option.

Available validators
- `ChoiceValidator`: Validate that the field is one of the given `Choices` (added by default).

### FieldGroup
A group of fields that can be displayed conditionally.

Properties:
- `Fields` ([]Field): List of fields in the group.
- `heading` (string): Heading for the group.

## TODOs

- [ ] Add more field types
- [ ] Add way to customize error messages (e.g., for translations)
- [ ] Add more UI frameworks for rendering (e.g., charm.sh for terminal UIs)
- [ ] Add more display conditions
- [ ] Add more validators
- [ ] Add tests
- [ ] Add more examples
- [ ] Clen up code
- [ ] Add more documentation

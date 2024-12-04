package go_forms

import (
	"encoding/json"
	"net"
	"regexp"
	"strconv"
)

type Field interface {
	GetId() string
	ShouldDisplay() bool
	IsValid() bool
	GetValue() string
	SetValue(value string)
	GetError() error
}

type Validator interface {
	Validate(field any) bool
}

type DisplayCondition interface {
	DisplayCondition(field any) bool
}

// Defining the Base Field Type

type FieldBaseType struct {
	Id                string
	DisplayConditions []DisplayCondition
	Validators        []Validator
	Value             string
	form              *Form
	error             error
}

func (f *FieldBaseType) GetId() string {
	return f.Id
}

func (f *FieldBaseType) ShouldDisplay() bool {
	for _, displayCondition := range f.DisplayConditions {
		if !displayCondition.DisplayCondition(f) {
			return false
		}
	}
	return true
}

func (f *FieldBaseType) IsValid() bool {
	if !f.ShouldDisplay() {
		return true
	}
	for _, validator := range f.Validators {
		if !validator.Validate(f) {
			return false
		}
	}
	f.error = nil
	return true
}

func (f *FieldBaseType) GetValue() string {
	return f.Value
}

func (f *FieldBaseType) SetValue(value string) {
	f.Value = value
	if f.form != nil {
		f.form.onChange()
	}
}

func (f *FieldBaseType) GetError() error {
	return f.error
}

type CustomValidator struct {
	Validator func(field any) (bool, error)
}

func (v *CustomValidator) Validate(field any) bool {
	valid, err := v.Validator(field)
	if err != nil {
		field.(*FieldBaseType).error = err
	}
	return valid
}

type AllFieldsValid struct{}

func (v *AllFieldsValid) Validate(field any) bool {
	fields := field.(*FieldBaseType).form.GetAllFields()
	for _, f := range fields {
		if !f.IsValid() {
			field.(*FieldBaseType).error = &CustomError{Message: "Not all fields are valid (invalid field: " + f.GetId() + ")"}
			return false
		}
	}
	return true
}

type IsValidValidator struct {
	FieldIds []string
}

func (v *IsValidValidator) Validate(field any) bool {
	fields := field.(*FieldBaseType).form.GetAllFields()
	for _, f := range fields {
		for _, id := range v.FieldIds {
			if f.GetId() == id && !f.IsValid() {
				field.(*FieldBaseType).error = &CustomError{Message: "Not all fields that should be valid are valid (invalid field: " + f.GetId() + ")"}
				return false
			}
		}
	}
	return true
}

type AlwaysDisplay struct{}

func (d *AlwaysDisplay) DisplayCondition(_ any) bool {
	return true
}

type CustomDisplayCondition struct {
	Condition func(field any) bool
}

func (d *CustomDisplayCondition) DisplayCondition(field any) bool {
	return d.Condition(field)
}

type IsValidDisplayCondition struct {
	FieldIds []string
}

func (d *IsValidDisplayCondition) DisplayCondition(field any) bool {
	fields := field.(*FieldBaseType).form.GetAllFields()
	for _, f := range fields {
		for _, id := range d.FieldIds {
			if f.GetId() == id && !f.IsValid() {
				return false
			}
		}
	}
	return true
}

type IsInvalidDisplayCondition struct {
	FieldIds []string
}

func (d *IsInvalidDisplayCondition) DisplayCondition(field any) bool {
	fields := field.(*FieldBaseType).form.GetAllFields()
	for _, f := range fields {
		for _, id := range d.FieldIds {
			if f.GetId() == id && f.IsValid() {
				return false
			}
		}
	}
	return true
}

type AllFieldsValidDisplayCondition struct{}

func (d *AllFieldsValidDisplayCondition) DisplayCondition(field any) bool {
	fields := field.(*FieldBaseType).form.GetAllFields()
	for _, f := range fields {
		if !f.IsValid() {
			return false
		}
	}
	return true
}

type HasValueDisplayCondition struct {
	FieldId string
	Value   string
}

func (d *HasValueDisplayCondition) DisplayCondition(field any) bool {
	fields := field.(*FieldBaseType).form.GetAllFields()
	for _, f := range fields {
		if f.GetId() == d.FieldId && f.GetValue() == d.Value {
			return true
		}
	}
	return false

}

type DisplayAfter struct {
	FieldId string
}

func (d *DisplayAfter) DisplayCondition(field any) bool {
	fields := field.(*FieldBaseType).form.GetAllFields()
	for _, f := range fields {
		if f.GetId() == d.FieldId {
			return f.IsValid() && f.ShouldDisplay()
		}
	}
	return false
}

type OrDisplayCondition struct {
	Conditions []DisplayCondition
}

func (d *OrDisplayCondition) DisplayCondition(field any) bool {
	for _, condition := range d.Conditions {
		if condition.DisplayCondition(field) {
			return true
		}
	}
	return false
}

type AndDisplayCondition struct {
	Conditions []DisplayCondition
}

func (d *AndDisplayCondition) DisplayCondition(field any) bool {
	for _, condition := range d.Conditions {
		if !condition.DisplayCondition(field) {
			return false
		}
	}
	return true
}

// Defining the Field Types based on the Base Field Type

type Message struct {
	*FieldBaseType
}

// Defining the Text Field Type based on the Base Field Type

type TextField struct {
	*FieldBaseType
	Placeholder string
	Prompt      string
}

type NotEmptyValidator struct{}

func (v *NotEmptyValidator) Validate(field any) bool {
	value := field.(*FieldBaseType).Value
	valid := value != ""
	if !valid {
		field.(*FieldBaseType).error = &CustomError{Message: "Field cannot be empty"}
	}
	return valid
}

type MaxLengthValidator struct {
	MaxLength int
}

func (v *MaxLengthValidator) Validate(field any) bool {
	value := field.(*FieldBaseType).Value
	valid := len(value) <= v.MaxLength
	if !valid {
		field.(*FieldBaseType).error = &CustomError{Message: "Field is too long (length: " + strconv.Itoa(len(value)) + ", max length: " + strconv.Itoa(v.MaxLength) + ")"}
	}
	return valid
}

type MinLengthValidator struct {
	MinLength int
}

func (v *MinLengthValidator) Validate(field any) bool {
	value := field.(*FieldBaseType).Value
	valid := len(value) >= v.MinLength
	if !valid {
		field.(*FieldBaseType).error = &CustomError{Message: "Field is too short (length: " + strconv.Itoa(len(value)) + ", min length: " + strconv.Itoa(v.MinLength) + ")"}
	}
	return valid
}

type IpValidator struct{}

func (v *IpValidator) Validate(field any) bool {
	value := field.(*FieldBaseType).Value
	valid := net.ParseIP(value) != nil
	if !valid {
		field.(*FieldBaseType).error = &CustomError{Message: "Field is not a valid IP address"}
	}
	return valid
}

type RegexValidator struct {
	RegexPattern string
}

func (v *RegexValidator) Validate(field any) bool {
	value := field.(*FieldBaseType).Value
	valid := true
	if value != "" {
		valid = regexp.MustCompile(v.RegexPattern).MatchString(value)
	}
	if !valid {
		field.(*FieldBaseType).error = &CustomError{Message: "Field does not match the required pattern (" + v.RegexPattern + ")"}
	}
	return valid
}

type UrlValidator struct{}

func (v *UrlValidator) Validate(field any) bool {
	value := field.(*FieldBaseType).Value
	valid := regexp.MustCompile(`^https?://.`).MatchString(value)
	if !valid {
		field.(*FieldBaseType).error = &CustomError{Message: "Field is not a valid URL"}
	}
	return valid
}

func (t *TextField) GetPlaceholder() string {
	return t.Placeholder
}

func (t *TextField) GetPrompt() string {
	return t.Prompt
}

// Defining the Number Field Type based on the Base Field Type

type NumberField struct {
	*TextField
}

type MinValidator struct {
	Min int
}

func (v *MinValidator) Validate(field any) bool {
	value := field.(*FieldBaseType).Value
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		field.(*FieldBaseType).error = &CustomError{Message: "Field value is not a integer"}
		return false
	}
	valid := v.Min <= valueAsInt
	if !valid {
		field.(*FieldBaseType).error = &CustomError{Message: "Field value is too small (value: " + value + ", min value: " + strconv.Itoa(v.Min) + ")"}
	}
	return valid
}

type MaxValidator struct {
	Max int
}

func (v *MaxValidator) Validate(field any) bool {
	value := field.(*FieldBaseType).Value
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		field.(*FieldBaseType).error = &CustomError{Message: "Field value is not a integer"}
		return false
	}
	valid := valueAsInt <= v.Max
	if !valid {
		field.(*FieldBaseType).error = &CustomError{Message: "Field value is too big (value: " + value + ", max value: " + strconv.Itoa(v.Max) + ")"}
	}
	return valueAsInt <= v.Max
}

type IsIntegerValidator struct{}

func (v *IsIntegerValidator) Validate(field any) bool {
	value := field.(*FieldBaseType).Value
	_, err := strconv.Atoi(value)
	if err != nil {
		field.(*FieldBaseType).error = &CustomError{Message: "Field value is not a integer"}
	}
	return err == nil
}

// Defining the Multiple Choice Field Type based on the Text Field Type

type MultipleChoiceField struct {
	*TextField
	Options map[string]Option
}

type Option struct {
	Label       string
	Description string
}

type ChoiceValidator struct{}

func (v *ChoiceValidator) Validate(field any) bool {
	multipleChoiceField, ok := field.(*MultipleChoiceField)
	if !ok {
		multipleChoiceField.error = &CustomError{Message: "Field is not a multiple choice field but ChoiceValidator was used"}
		return false
	}
	_, ok = multipleChoiceField.Options[multipleChoiceField.Value]
	if !ok {
		multipleChoiceField.error = &CustomError{Message: "Field value is not a valid option"}
	}
	return ok
}

func (m *MultipleChoiceField) GetOptions() map[string]Option {
	return m.Options
}

func (m *MultipleChoiceField) IsValid() bool {
	if !m.ShouldDisplay() {
		return true
	}
	for _, validator := range m.Validators {
		if !validator.Validate(m) {
			return false
		}
	}
	m.error = nil
	return true
}

// Defining the Field Group Type based on the Base Field Type

type FieldGroup struct {
	*FieldBaseType
	Fields  []Field
	heading string
}

func (f *FieldGroup) GetFieldsToDisplay() []Field {
	var fieldsToDisplay []Field
	for _, field := range f.Fields {
		if field.ShouldDisplay() {
			fieldsToDisplay = append(fieldsToDisplay, field)
		}
	}
	return fieldsToDisplay
}

func (f *FieldGroup) GetFieldById(id string) Field {
	for _, field := range f.Fields {
		if field.GetId() == id {
			return field
		}
	}
	return nil
}

func (f *FieldGroup) GetValue() string {
	fieldValues := make(map[string]string)
	for _, field := range f.Fields {
		fieldValues[field.GetId()] = field.GetValue()
	}
	jsonFieldValues, _ := json.Marshal(fieldValues)
	return string(jsonFieldValues)
}

func (f *FieldGroup) GetHeading() string {
	return f.heading
}

func (f *FieldGroup) SetValue(value string) {
	var fieldValues map[string]string
	err := json.Unmarshal([]byte(value), &fieldValues)
	if err != nil {
		return
	}
	for _, field := range f.Fields {
		field.SetValue(fieldValues[field.GetId()])
	}
}

func (f *FieldGroup) SetHeading(heading string) {
	f.heading = heading
}

// Defining the Form Type

type Form struct {
	Fields   []Field
	onChange func()
}

func (f *Form) GetAllFields() []Field {
	fields := make([]Field, 0)
	for _, field := range f.Fields {
		fields = append(fields, field)
		if group, ok := field.(*FieldGroup); ok {
			fields = append(fields, group.Fields...)
		}
	}
	return fields
}

func (f *Form) IsValid() bool {
	for _, field := range f.Fields {
		if !field.IsValid() {
			return false
		}
	}
	return true
}

func (f *Form) GetFieldById(id string) Field {
	for _, field := range f.Fields {
		if field.GetId() == id {
			return field
		}
	}
	return nil
}

func (f *Form) GetFieldsToDisplay() []Field {
	var fieldsToDisplay []Field
	for _, field := range f.Fields {
		if field.ShouldDisplay() {
			fieldsToDisplay = append(fieldsToDisplay, field)
		}
	}
	return fieldsToDisplay
}

func (f *Form) GetFieldValues() map[string]string {
	fieldValues := make(map[string]string)
	for _, field := range f.Fields {
		fieldValues[field.GetId()] = field.GetValue()
	}
	return fieldValues
}

func (f *Form) SetOnChangeCallback(onChange func()) {
	f.onChange = onChange
}

func (f *Form) GetError() error {
	for _, field := range f.Fields {
		if !field.IsValid() {
			return &CustomError{Message: field.GetId() + "is not valid (" + field.GetError().Error() + ")"}
		}
	}
	return nil
}

// Defining the Form Builder Functions

func NewForm(fields ...Field) *Form {
	form := &Form{Fields: fields, onChange: func() {}}
	for _, field := range fields {
		switch v := field.(type) {
		case *FieldBaseType:
			v.form = form
		case *TextField:
			v.FieldBaseType.form = form
		case *NumberField:
			v.TextField.FieldBaseType.form = form
		case *MultipleChoiceField:
			v.TextField.FieldBaseType.form = form
		case *FieldGroup:
			v.FieldBaseType.form = form
		}
	}
	return form
}

func NewFieldGroup(id string, displayConditions []DisplayCondition, validators []Validator, heading string, fields ...Field) *FieldGroup {
	return &FieldGroup{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: validators}, Fields: fields, heading: heading}
}

func NewTextField(id string, displayConditions []DisplayCondition, validators []Validator, placeholder string, prompt string, defaultValue string) *TextField {
	return &TextField{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: validators, Value: defaultValue}, Placeholder: placeholder, Prompt: prompt}
}

func NewNumberField(id string, displayConditions []DisplayCondition, validators []Validator, placeholder string, prompt string, defaultValue int) *NumberField {
	return &NumberField{TextField: &TextField{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: validators, Value: strconv.Itoa(defaultValue)}, Placeholder: placeholder, Prompt: prompt}}
}

func NewMultipleChoiceField(id string, displayConditions []DisplayCondition, validators []Validator, placeholder string, prompt string, options map[string]Option, defaultValue string) *MultipleChoiceField {
	return &MultipleChoiceField{TextField: &TextField{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: validators, Value: defaultValue}, Placeholder: placeholder, Prompt: prompt}, Options: options}
}

func NewMessage(id string, displayConditions []DisplayCondition, message string) *Message {
	return &Message{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: []Validator{}, Value: message}}
}

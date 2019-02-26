package goweb

import (
	"reflect"
	"strings"
)

const (
	// FormErrMissingRequired missing required param
	FormErrMissingRequired = 1
	// FormErrNotAForm the form object to store data need to be a struct instance
	FormErrNotAForm = 2
	// FormErrTypeCannotCast type casting error
	FormErrTypeCannotCast = 3
	// FormErrNotAValidForm form validation failed
	FormErrNotAValidForm = 4
)

var (
	formCache     = make(map[string]*FormMeta, 0)
	formErrMsgMap = map[int]string{
		FormErrMissingRequired: "Missing Required Field",
		FormErrNotAForm:        "Form must be a struct",
		FormErrTypeCannotCast:  "Can not cast to type",
		FormErrNotAValidForm:   "Form validation failed",
	}
)

// FormError validation error of a form field
type FormError struct {
	FieldName string
	Type      int
	Err       error
}

func newFormError(errType int, fieldName string, err error) error {
	return &FormError{
		Type:      errType,
		FieldName: fieldName,
		Err:       err,
	}
}

func (fe *FormError) Error() string {
	msg, found := formErrMsgMap[fe.Type]
	if !found {
		msg = "Form Error"
	}
	return msg + ":" + fe.FieldName
}

// FormMeta meta info of a Form Object
type FormMeta struct {
	Type           *reflect.Type
	RequiredFields []string
	FieldMap       map[string]*FieldConfig
}

// FieldConfig Represent settings of a form field, decoded from the struct field tag named form
// For Example:
// type CreatePostForm struct {
//      title   string  `form:"title,required"`
//      user_id int     `form:"X-User-Id,required,header"`
//      content string  `form:"content"`
// }
// The first value of form tag is the name of parameter.
// header, cookie, path indecate where should the parameter be fetched. Parameter fetched from QueryString, body forms by default.
// required means the form field value must be a none empty string
type FieldConfig struct {
	ParamName  string
	Type       *reflect.Type
	FromHeader bool
	FromCookie bool
	FromPath   bool
	IsRequired bool
}

func decodeFormTag(tagValue string) *FieldConfig {
	if tagValue == "" {
		return nil
	}

	parts := strings.Split(tagValue, ",")
	c := &FieldConfig{
		ParamName:  parts[0],
		FromHeader: false,
		FromCookie: false,
		IsRequired: false,
	}

	for _, name := range parts[1:] {
		if name == "required" {
			c.IsRequired = true
		} else if name == "header" {
			c.FromHeader = true
		} else if name == "cookie" {
			c.FromCookie = true
		} else if name == "path" {
			c.FromPath = true
		}
	}

	return c
}

// InspectForm Get meta info of a form object
func InspectForm(m interface{}) (*FormMeta, error) {
	typo := reflect.Indirect(reflect.ValueOf(m)).Type()
	fullName := typo.PkgPath() + "." + typo.Name()

	if found, exsits := formCache[fullName]; exsits {
		return found, nil
	}

	if typo.Kind() != reflect.Struct {
		return nil, newFormError(FormErrNotAForm, "", nil)
	}

	fieldMap := make(map[string]*FieldConfig, 0)
	requiredFields := make([]string, 0)

	for i := 0; i < typo.NumField(); i++ {
		structField := typo.Field(i)
		fieldConf := decodeFormTag(structField.Tag.Get("form"))
		if fieldConf != nil {
			fieldConf.Type = &structField.Type
			if fieldConf.IsRequired {
				requiredFields = append(requiredFields, structField.Name)
			}
			fieldMap[structField.Name] = fieldConf
		}
	}

	formCache[fullName] = &FormMeta{
		Type:           &typo,
		FieldMap:       fieldMap,
		RequiredFields: requiredFields,
	}

	return formCache[fullName], nil
}

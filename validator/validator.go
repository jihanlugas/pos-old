package validator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (v *CustomValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func NewValidator() *CustomValidator {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	//_ = validate.RegisterValidation("notexists", notExistsOnDbTable)
	//_ = validate.RegisterValidation("existsdata", existsDataOnDbTable)
	//_ = validate.RegisterValidation("no_hp", validNoHp)
	//_ = validate.RegisterValidation("passwdComplex", checkPasswordComplexity)
	//_ = validate.RegisterValidation("photo", photoCheck, true)
	//_ = validate.RegisterValidation("hiragana", hiragana)
	//_ = validate.RegisterValidation("katakana", katakana)
	//_ = validate.RegisterValidation("kana", kana)
	//_ = validate.RegisterValidation("kanji", kanji)
	//_ = validate.RegisterValidation("electionTypeProvince", electionTypeProvince)
	//_ = validate.RegisterValidation("electionTypeRegency", electionTypeRegency)
	//_ = validate.RegisterValidation("electionTypeDistrictdapil", electionTypeDistrictdapil)

	return &CustomValidator{
		validator: validate,
	}
}

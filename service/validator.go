package service

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/videocoin/cloud-api/rpc"
	enLocale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
	enTrans "gopkg.in/go-playground/validator.v9/translations/en"
)

type requestValidator struct {
	validator  *validator.Validate
	translator *ut.Translator
}

func newRequestValidator() *requestValidator {
	lt := enLocale.New()
	en := &lt

	uniTranslator := ut.New(*en, *en)
	uniEn, _ := uniTranslator.GetTranslator("en")
	translator := &uniEn

	validate := validator.New()
	enTrans.RegisterDefaultTranslations(validate, *translator)

	validate.RegisterTranslation(
		"email",
		*translator,
		RegisterEmailTranslation,
		EmailTranslation)
	validate.RegisterTranslation(
		"secure-password",
		*translator,
		RegisterSecurePasswordTranslation,
		SecurePasswordTranslation)
	validate.RegisterTranslation(
		"confirm-password",
		*translator,
		RegisterConfirmPasswordTranslation,
		ConfirmPasswordTranslation)

	validate.RegisterValidation("confirm-password", ValidateConfirmPassword)
	validate.RegisterValidation("secure-password", ValidateSecurePassword)

	return &requestValidator{
		validator:  validate,
		translator: translator,
	}

}

func (rv *requestValidator) validate(r interface{}) *rpc.MultiValidationError {
	trans := *rv.translator
	verrs := &rpc.MultiValidationError{}

	serr := rv.validator.Struct(r)
	if serr != nil {
		verrs.Errors = []*rpc.ValidationError{}

		for _, err := range serr.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(r).Elem().FieldByName(err.Field())
			jsonField := extractValueFromTag(field.Tag.Get("json"))
			verr := &rpc.ValidationError{
				Field:   jsonField,
				Message: err.Translate(trans),
			}
			verrs.Errors = append(verrs.Errors, verr)
		}

		return verrs
	}

	return nil
}

func RegisterEmailTranslation(ut ut.Translator) error {
	return ut.Add("email", "Enter a valid email address", true)
}

func EmailTranslation(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("email", fe.Field())
	return t
}

func RegisterSecurePasswordTranslation(ut ut.Translator) error {
	return ut.Add("secure-password", "Password must be more than 8 characters and contain both numbers and letters", true)
}

func SecurePasswordTranslation(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("secure-password", fe.Field())
	return t
}

func RegisterConfirmPasswordTranslation(ut ut.Translator) error {
	return ut.Add("confirm-password", "Passwords does not match", true)
}

func ConfirmPasswordTranslation(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("confirm-password", fe.Field())
	return t
}

func ValidateConfirmPassword(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}

	return field.String() == currentField.String()
}

func ValidateSecurePassword(fl validator.FieldLevel) bool {
	field := fl.Field()
	password := field.String()

	if password == "" {
		return false
	}

	var (
		hasMinLen = false
		hasNumber = false
		hasLetter = false
		// hasUpper   = false
		// hasSpecial = false
		// hasLower   = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsNumber(char):
			hasNumber = true
			//     hasUpper = true
			// case unicode.IsLower(char):
			//     hasLower = true
			// case unicode.IsPunct(char) || unicode.IsSymbol(char):
			//     hasSpecial = true
		}
	}

	return hasMinLen && hasNumber && hasLetter
}

func extractValueFromTag(tag string) string {
	values := strings.Split(tag, ",")
	return values[0]
}

package controllers

import (
	"encoding/json"
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ApiError struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func RegisterErrors(instance *validator.Validate) (ut.Translator, error) {

	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Println("translator not found")
	}

	if err := en_translations.RegisterDefaultTranslations(instance, trans); err != nil {
		log.Println(err)
		return nil, err
	}

	_ = instance.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = instance.RegisterTranslation("uuid4", trans, func(ut ut.Translator) error {
		return ut.Add("uuid4", "{0} must be in valid UUID v4 format", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("uuid4", fe.Field())
		return t
	})

	return trans, nil
}

func ValidateRequest(body []byte, obj interface{}) (*ApiError, error) {
	err := json.Unmarshal(body, obj)
	if err != nil {
		log.Println("Error when trying to unmarshal request body.")
		return nil, err
	}
	valErrs := validate.Struct(obj)
	if valErrs != nil {
		errors := []string{}
		for _, err := range valErrs.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(translator))
		}
		log.Println("Request body validation error.")
		return &ApiError{Message: "Error when validating payload", Errors: errors}, nil
	}
	return nil, nil
}

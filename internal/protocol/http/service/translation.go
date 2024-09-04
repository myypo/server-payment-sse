package service

import (
	"errors"
	"fmt"
	httpErr "payment-sse/internal/protocol/http/error"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type TranslationService struct {
	enTrans ut.Translator
}

func NewTranslationService() TranslationService {
	serv := TranslationService{}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		uni := ut.New(en, en)

		serv.enTrans, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v, serv.enTrans)

	}

	return serv
}

const FailedToValidate string = "unhandled input validation error has occured"

func (s *TranslationService) translate(
	errs error,
	ut ut.Translator,
) *httpErr.HttpError {
	valErrs, ok := errs.(validator.ValidationErrors)
	if !ok {
		return httpErr.NewBadRequest(ErrFailedToParse)
	}

	transErrs := valErrs.Translate(ut)
	transErrsArr := make([]string, 0, len(transErrs))
	for _, v := range transErrs {
		transErrsArr = append(transErrsArr, v)
	}

	return httpErr.NewBadRequest(fmt.Errorf("%s", strings.Join(transErrsArr, ". ")))
}

func (s *TranslationService) TranslateEN(
	errs error,
) *httpErr.HttpError {
	return s.translate(errs, s.enTrans)
}

var ErrFailedToParse error = errors.New("failed to parse incoming request, invalid data provided")

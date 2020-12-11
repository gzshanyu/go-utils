package utils

import (
	"context"
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	vali "github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

type Validator struct {
	ctx      context.Context
	validate *vali.Validate
	trans    ut.Translator
}

func (p *Validator) Validate(s interface{}, name ...string) []error {
	if p.validate == nil {
		return nil
	}

	// 注册一个函数，获取struct tag里自定义的label作为字段名
	var tagName = "label"
	if name != nil {
		tagName = name[0]
	}

	p.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get(tagName)
	})

	if err := p.validate.Struct(s); err != nil {
		return p.translate(err)
	}

	return nil
}

func (p *Validator) ValidateField(field interface{}, tag string, name ...string) error {
	if p.validate == nil {
		return nil
	}

	var fieldName string
	if name != nil {
		fieldName = name[0]
	}

	if err := p.validate.Var(field, tag); err != nil {
		errs := p.translate(err)
		if errs != nil && len(errs) > 0 {
			return errors.New(fieldName + errs[0].Error())
		}
	}

	return nil
}

func (p *Validator) translate(err error) []error {
	var errs []error
	for _, e := range err.(vali.ValidationErrors) {
		errs = append(errs, errors.New(e.Translate(p.trans)))
	}
	return errs
}

func (p *Validator) load() {
	var validate = vali.New()
	var uni = ut.New(zh.New())
	var trans, _ = uni.GetTranslator("zh")
	// 验证器注册翻译器
	if err := zh_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return
	}

	p.trans = trans
	p.validate = validate
}

func NewValidator(ctx ...context.Context) *Validator {
	obj := &Validator{}
	if ctx != nil && len(ctx) > 0 {
		obj.ctx = ctx[0]
	}
	obj.load()
	return obj
}

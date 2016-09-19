package schema

import (
	"gopkg.in/asaskevich/govalidator.v4"
	"testing"
)

func TestValidateBasic(t *testing.T) {
	var p = Property{
		Path: "hello",
		Validator: []struct {
			Name     string
			Params   []interface{}
			IsCustom bool
			IsParams bool
		}{{Name: "email"}},
	}
	var data = "123"
	_, errs := p.validateTag(data)
	t.Logf("%s", errs)
}

func TestValidateParam(t *testing.T) {
	var p = Property{
		Path: "hello",
		Validator: []struct {
			Name     string
			Params   []interface{}
			IsCustom bool
			IsParams bool
		}{{Name: "stringlength", IsParams: true, Params: []interface{}{3, 20}}},
	}
	var data = 1231
	_, errs := p.validateTag(data)
	t.Logf("%s", errs)
}

func TestValidateCustom(t *testing.T) {
	govalidator.TagMap["cus"] = govalidator.Validator(func(str string) bool {
		return str == "duck"
	})

	var p = Property{
		Path: "hello",
		Validator: []struct {
			Name     string
			Params   []interface{}
			IsCustom bool
			IsParams bool
		}{{Name: "cus"}},
	}
	var data = "duck"
	_, errs := p.validateTag(data)
	t.Logf("%s", errs)
}

func TestSchema(t *testing.T) {
	var sche Schema
	var p1 = Property{
		Path: "hello",
		Validator: []struct {
			Name     string
			Params   []interface{}
			IsCustom bool
			IsParams bool
		}{{Name: "stringlength", IsParams: true, Params: []interface{}{3, 20}}},
	}
	sche.AddProperty(p1)
	errs := sche.Validate(map[string]interface{}{"hello": 1231, "name": "what"})
	t.Logf("%s", &errs)
}
func TestSchemaEmpty(t *testing.T) {
	var sche Schema
	var p1 = Property{
		Path:    "hello",
		Default: "hello world",
	}
	var p2 = Property{
		Path: "name.last",
	}
	var p3 = Property{
		Path: "name.first",
	}
	sche.AddProperty(p1, p2, p3)
	result := sche.EmptyMap()
	t.Logf("%s", result)
}

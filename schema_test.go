package schema

import (
	"encoding/json"
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
	var p2 = Property{
		Path:    "second.$",
		Default: []int{123123},
		Validator: []struct {
			Name     string
			Params   []interface{}
			IsCustom bool
			IsParams bool
		}{{Name: "alphanum", IsParams: false, Params: []interface{}{13, 20}}},
	}
	sche.AddProperty(p1, p2)
	var data = map[string]interface{}{
		"hello":  1231,
		"name":   "what",
		"first":  []map[string]interface{}{map[string]interface{}{"name": 123}},
		"second": []int64{123, 1232},
		"last":   map[string]interface{}{"t": "tt"}}

	ps, _ := sche.Validate(data)
	t.Logf("Validate: %s", data)
	data2, _ := Flatten(data)
	t.Logf("Flatten: %s", data2)

	data3 := sche.CleanFlatMap(data, ps)
	t.Logf("CleanFlatMap: %s", data3)

	data4 := sche.CleanMap(data3, ps)
	data5, _ := json.Marshal(data4)
	t.Logf("CleanMap: %s", data5)
}

func TestSchemaEmpty(t *testing.T) {
	var sche Schema
	var p1 = Property{
		Path:    "hello",
		Default: "hello world",
	}
	var p2 = Property{
		Path:    "name.$.first.$.last",
		Default: "hello",
	}
	var p3 = Property{
		Path:    "hah.$",
		Default: []int{123, 1232},
	}
	var p4 = Property{
		Path:    "x.$.what",
		Default: []int{123, 1232},
	}
	sche.AddProperty(p1, p2, p3, p4)
	result := sche.EmptyMap()
	data, _ := json.Marshal(result)
	t.Logf("%s", data)
}

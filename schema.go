package schema

import (
	"fmt"
	"github.com/astaxie/flatmap"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/asaskevich/govalidator.v4"
	"reflect"
	"regexp"
	"strings"
)

type Schema struct {
	Properties map[string]Property
}

func (s *Schema) EmptyMap() map[string]interface{} {
	var result = make(map[string]interface{})

	for path, p := range s.Properties {
		if !strings.Contains(path, ".") {
			result[path] = p.Default
		} else {
			paths := strings.Split(path, ".")
			setValue(result, paths, p.Default)
		}
	}
	return result
}

func (s *Schema) RegistValidator(tag string, f func(string) bool) {
	govalidator.TagMap[tag] = govalidator.Validator(f)
}

func (s *Schema) AddProperty(ps ...Property) {
	if s.Properties == nil {
		s.Properties = make(map[string]Property)
	}
	for _, p := range ps {
		s.Properties[p.Path] = p
	}
}

func (s *Schema) Validate(data map[string]interface{}) (paths []string, errs SchemaValidateResult) {
	flatData, err := flatmap.Flatten(data)
	if err != nil {
		errs.Add(".", err)
		return
	}

	for orgPath, v := range flatData {
		path := regDollar.ReplaceAllString(orgPath, ".$")
		if p, found := s.Properties[path]; found {
			orgValue := v
			ok, e := p.validateTag(orgValue)
			if !ok {
				errs.Add(path, e)
			} else {
				paths = append(paths, orgPath)
			}
		}
	}
	return paths, errs
}

func (s *Schema) CleanFlatMap(data map[string]interface{}, paths []string) map[string]interface{} {
	data2, _ := Flatten(data)
	var data3 = make(map[string]interface{})
	for _, v := range paths {
		data3[v] = data2[v]
	}
	return data3
}

func (s *Schema) CleanMap(data map[string]interface{}, paths []string) map[string]interface{} {
	var data2 = make(map[string]interface{})
	for _, path := range paths {
		if !strings.Contains(path, ".") {
			data2[path] = data[path]
		} else {
			newPath := regDollar.ReplaceAllString(path, ".$")
			paths := strings.Split(newPath, ".")
			setValue2(data2, paths, data[path])
		}
	}
	return data2
}

// TODO: 类型没有定义
type Property struct {
	Path      string
	Type      string
	Default   interface{}
	Validator []struct {
		Name     string
		Params   []interface{}
		IsCustom bool
		IsParams bool
	}
}

func (p *Property) validateTag(data interface{}) (bool, error) {
	var errs error
	for _, v := range p.Validator {
		if v.IsCustom {
			if validateFunc, found := govalidator.CustomTypeTagMap.Get(v.Name); found {
				ok := validateFunc(data, v.Params[0])
				if !ok {
					errs = multierror.Append(errs, fmt.Errorf("Validate Error: %s  for %s", v.Name, fmt.Sprint(data)))
				}
			}
		} else if v.IsParams {
			if validateFunc, found := govalidator.ParamTagMap[v.Name]; found {
				var params []string
				for _, param := range v.Params {
					params = append(params, fmt.Sprint(param))
				}
				ok := validateFunc(fmt.Sprint(data), params...)
				if !ok {
					errs = multierror.Append(errs, fmt.Errorf("Validate Error: %s  for %s", v.Name, fmt.Sprint(data)))
				}
			}
		} else if validateFunc, found := govalidator.TagMap[v.Name]; found {
			ok := validateFunc(fmt.Sprint(data))
			if !ok {
				errs = multierror.Append(errs, fmt.Errorf("Validate Error: %s  for %s", v.Name, fmt.Sprint(data)))
			}
		} else {

		}
	}
	return errs == nil, errs
}

type SchemaValidateResult struct {
	validates map[string]error
}

func (s *SchemaValidateResult) Add(path string, errs error) {
	if s.validates == nil {
		s.validates = make(map[string]error)
	}
	if errs != nil {
		s.validates[path] = errs
	}
}

func (s SchemaValidateResult) String() string {
	var r = ""
	for _, v := range s.validates {
		r += v.Error() + "\n"
	}
	return r
}

func (s SchemaValidateResult) Error() string {
	return s.String()
}

func (s *SchemaValidateResult) Ok() bool {
	return len(s.validates) == 0
}

var (
	regDollar, _ = regexp.Compile("(\\.[0-9]+)")
)

func setValue(data map[string]interface{}, paths []string, value interface{}) {
	if len(paths) > 2 {
		if paths[1] == "$" {
			if _, found := data[paths[0]]; !found {
				data[paths[0]] = make([]map[string]interface{}, 1)
				data[paths[0]].([]map[string]interface{})[0] = make(map[string]interface{})
			}
			setValue(data[paths[0]].([]map[string]interface{})[0], paths[2:], value)
		} else {
			data[paths[0]] = make(map[string]interface{})
			setValue(data[paths[0]].(map[string]interface{}), paths[1:], value)
		}
	} else {
		if orgP, found := data[paths[0]]; !found {
			if data == nil {
				data = make(map[string]interface{})
			}
			data[paths[0]] = value
		} else {
			v := reflect.ValueOf(orgP)
			if v.Kind() == reflect.Interface {
				v = v.Elem()
			}
			switch v.Kind() {
			case reflect.Bool,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Float64, reflect.Float32,
				reflect.String:
				data[paths[0]] = []interface{}{orgP, value}
			case reflect.Slice, reflect.Array:
				data[paths[0]] = append(data[paths[0]].([]interface{}), value)
			}

		}

	}
}

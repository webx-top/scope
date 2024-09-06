package scopes

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"reflect"
	"strings"
)

var (
	_ json.Marshaler = (*ScopedJSON)(nil)
)

func New(scope string, data interface{}) *ScopedJSON {
	return &ScopedJSON{
		scope:       scope,
		parsedScope: strings.Split(scope, `:`),
		buffer:      bytes.NewBuffer(nil),
		data:        data,
	}
}

type ScopedJSON struct {
	scope       string
	parsedScope []string
	buffer      *bytes.Buffer
	data        interface{}
}

func (j *ScopedJSON) MarshalJSON() ([]byte, error) {
	err := j.scopeAll(`json`)
	return j.buffer.Bytes(), err
}

func (j *ScopedJSON) JSON() []byte {
	r, _ := j.MarshalJSON()
	return r
}

func (j *ScopedJSON) Read(p []byte) (n int, err error) {
	return j.buffer.Read(p)
}

func (j *ScopedJSON) Write(p []byte) (n int, err error) {
	return j.buffer.Write(p)
}

func (j *ScopedJSON) scopeAll(encodingType string) (err error) {
	val := reflect.Indirect(reflect.ValueOf(j.data))
	if val.Kind() == reflect.Slice {
		slicer := make([]map[string]interface{}, 0)
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i).Interface()
			slicer = append(slicer, j.extractData(v))
		}
		err = j.marshalWrite(slicer, encodingType)
	} else {
		err = j.marshalWrite(j.extractData(j.data), encodingType)
	}
	return
}

func (j *ScopedJSON) marshalWrite(data interface{}, encodingType string) error {
	var b []byte
	var err error
	if encodingType == `xml` {
		b, err = xml.Marshal(data)
	} else {
		b, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}
	_, err = j.Write(b)
	return err
}

func (j *ScopedJSON) extractData(data interface{}) map[string]interface{} {
	val := reflect.Indirect(reflect.ValueOf(data))
	thisData := make(map[string]interface{})
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tagVal := typeField.Tag

		jTags := tagVal.Get("json")
		jsonTag := strings.Split(jTags, ",")
		if len(jsonTag) > 0 {
			switch jsonTag[0] {
			case ``:
				jsonTag[0] = typeField.Name
			case `-`:
				continue
			}
			if len(jsonTag) == 2 {
				if jsonTag[1] == "omitempty" && valueField.Interface() == "" {
					continue
				}
			}
		} else {
			jsonTag = append(jsonTag, typeField.Name)
		}

		tag := tagVal.Get("scope")
		if tag == "" {
			thisData[jsonTag[0]] = valueField.Interface()
			continue
		}

		tags := strings.Split(tag, ",")
		if inSlice(j.parsedScope, tags) {
			thisData[jsonTag[0]] = valueField.Interface()
		}
	}
	return thisData
}

func inSlice(parsedScope []string, allowedScopes []string) bool {
	for _, allowedScope := range allowedScopes {
		scopeBSplit := strings.Split(allowedScope, ":")
		if matchParsedScopes(parsedScope, scopeBSplit) {
			return true
		}
	}
	return false
}

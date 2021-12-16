/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package json_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	runtimetesting "k8s.io/apimachinery/pkg/runtime/testing"
	"k8s.io/apimachinery/pkg/util/diff"
)

type testDecodable struct {
	metav1.TypeMeta `json:",inline"`

	Other     string
	Value     int           `json:"value"`
	Spec      DecodableSpec `json:"spec"`
	Interface interface{}   `json:"interface"`
}

// DecodableSpec has 15 fields.
type DecodableSpec struct {
	A int `json:"A"`
	B int `json:"B"`
	C int `json:"C"`
	D int `json:"D"`
	E int `json:"E"`
	F int `json:"F"`
	G int `json:"G"`
	H int `json:"h"`
	I int `json:"i"`
	J int `json:"j"`
	K int `json:"k"`
	L int `json:"l"`
	M int `json:"m"`
	N int `json:"n"`
	O int `json:"o"`
}

func (d *testDecodable) DeepCopyObject() runtime.Object {
	if d == nil {
		return nil
	}
	out := new(testDecodable)
	d.DeepCopyInto(out)
	return out
}
func (d *testDecodable) DeepCopyInto(out *testDecodable) {
	*out = *d
	out.Other = d.Other
	out.Value = d.Value
	out.Spec = d.Spec
	out.Interface = d.Interface
	return
}

type testDecodeCoercion struct {
	metav1.TypeMeta `json:",inline"`

	Bool bool `json:"bool"`

	Int   int `json:"int"`
	Int32 int `json:"int32"`
	Int64 int `json:"int64"`

	Float32 float32 `json:"float32"`
	Float64 float64 `json:"float64"`

	String string `json:"string"`

	Struct testDecodable `json:"struct"`

	Array []string          `json:"array"`
	Map   map[string]string `json:"map"`
}

func (d *testDecodeCoercion) DeepCopyObject() runtime.Object {
	if d == nil {
		return nil
	}
	out := new(testDecodeCoercion)
	d.DeepCopyInto(out)
	return out
}
func (d *testDecodeCoercion) DeepCopyInto(out *testDecodeCoercion) {
	*out = *d
	return
}

func TestDecode(t *testing.T) {
	testCases := []struct {
		creater runtime.ObjectCreater
		typer   runtime.ObjectTyper
		yaml    bool
		pretty  bool
		strict  bool

		data       []byte
		defaultGVK *schema.GroupVersionKind
		into       runtime.Object

		errFn          func(error) bool
		expectedObject runtime.Object
		expectedGVK    *schema.GroupVersionKind
	}{
		{
			data: []byte("{}"),

			expectedGVK: &schema.GroupVersionKind{},
			errFn:       func(err error) bool { return strings.Contains(err.Error(), "Object 'Kind' is missing in") },
		},
		{
			data:       []byte("{}"),
			defaultGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			creater:    &mockCreater{err: fmt.Errorf("fake error")},

			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			errFn:       func(err error) bool { return err.Error() == "fake error" },
		},
		{
			data:           []byte("{}"),
			defaultGVK:     &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			creater:        &mockCreater{obj: &testDecodable{}},
			expectedObject: &testDecodable{},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
		},

		// version without group is not defaulted
		{
			data:           []byte(`{"apiVersion":"blah"}`),
			defaultGVK:     &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			creater:        &mockCreater{obj: &testDecodable{}},
			expectedObject: &testDecodable{TypeMeta: metav1.TypeMeta{APIVersion: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "", Version: "blah"},
		},
		// group without version is defaulted
		{
			data:           []byte(`{"apiVersion":"other/"}`),
			defaultGVK:     &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			creater:        &mockCreater{obj: &testDecodable{}},
			expectedObject: &testDecodable{TypeMeta: metav1.TypeMeta{APIVersion: "other/"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
		},
		// group version, kind is defaulted
		{
			data:           []byte(`{"apiVersion":"other1/blah1"}`),
			defaultGVK:     &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			creater:        &mockCreater{obj: &testDecodable{}},
			expectedObject: &testDecodable{TypeMeta: metav1.TypeMeta{APIVersion: "other1/blah1"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other1", Version: "blah1"},
		},
		// gvk all provided then not defaulted at all
		{
			data:           []byte(`{"kind":"Test","apiVersion":"other/blah"}`),
			defaultGVK:     &schema.GroupVersionKind{Kind: "Test1", Group: "other1", Version: "blah1"},
			creater:        &mockCreater{obj: &testDecodable{}},
			expectedObject: &testDecodable{TypeMeta: metav1.TypeMeta{APIVersion: "other/blah", Kind: "Test"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
		},
		//gvk defaulting if kind not provided in data and defaultGVK use into's kind
		{
			data:           []byte(`{"apiVersion":"b1/c1"}`),
			into:           &testDecodable{TypeMeta: metav1.TypeMeta{Kind: "a3", APIVersion: "b1/c1"}},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "a3", Group: "b1", Version: "c1"}},
			defaultGVK:     nil,
			creater:        &mockCreater{obj: &testDecodable{}},
			expectedObject: &testDecodable{TypeMeta: metav1.TypeMeta{Kind: "a3", APIVersion: "b1/c1"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "a3", Group: "b1", Version: "c1"},
		},

		// accept runtime.Unknown as into and bypass creator
		{
			data: []byte(`{}`),
			into: &runtime.Unknown{},

			expectedGVK: &schema.GroupVersionKind{},
			expectedObject: &runtime.Unknown{
				Raw:         []byte(`{}`),
				ContentType: runtime.ContentTypeJSON,
			},
		},
		{
			data: []byte(`{"test":"object"}`),
			into: &runtime.Unknown{},

			expectedGVK: &schema.GroupVersionKind{},
			expectedObject: &runtime.Unknown{
				Raw:         []byte(`{"test":"object"}`),
				ContentType: runtime.ContentTypeJSON,
			},
		},
		{
			data:        []byte(`{"test":"object"}`),
			into:        &runtime.Unknown{},
			defaultGVK:  &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &runtime.Unknown{
				TypeMeta:    runtime.TypeMeta{APIVersion: "other/blah", Kind: "Test"},
				Raw:         []byte(`{"test":"object"}`),
				ContentType: runtime.ContentTypeJSON,
			},
		},

		// unregistered objects can be decoded into directly
		{
			data:        []byte(`{"kind":"Test","apiVersion":"other/blah","value":1,"Other":"test"}`),
			into:        &testDecodable{},
			typer:       &mockTyper{err: runtime.NewNotRegisteredErrForKind("mock", schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"})},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodable{
				TypeMeta: metav1.TypeMeta{APIVersion: "other/blah", Kind: "Test"},
				Other:    "test",
				Value:    1,
			},
		},
		// registered types get defaulted by the into object kind
		{
			data:        []byte(`{"value":1,"Other":"test"}`),
			into:        &testDecodable{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodable{
				Other: "test",
				Value: 1,
			},
		},
		// registered types get defaulted by the into object kind even without version, but return an error
		{
			data:        []byte(`{"value":1,"Other":"test"}`),
			into:        &testDecodable{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: ""}},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: ""},
			errFn:       func(err error) bool { return strings.Contains(err.Error(), "Object 'apiVersion' is missing in") },
			expectedObject: &testDecodable{
				Other: "test",
				Value: 1,
			},
		},
		// Error on invalid number
		{
			data:        []byte(`{"kind":"Test","apiVersion":"other/blah","interface":1e1000}`),
			creater:     &mockCreater{obj: &testDecodable{}},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			errFn: func(err error) bool {
				return strings.Contains(err.Error(), `json: cannot unmarshal number 1e1000 into Go struct field testDecodable.interface of type float64`)
			},
		},
		// Unmarshalling is case-sensitive
		{
			// "VaLue" should have been "value"
			data:        []byte(`{"kind":"Test","apiVersion":"other/blah","VaLue":1,"Other":"test"}`),
			into:        &testDecodable{},
			typer:       &mockTyper{err: runtime.NewNotRegisteredErrForKind("mock", schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"})},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodable{
				TypeMeta: metav1.TypeMeta{APIVersion: "other/blah", Kind: "Test"},
				Other:    "test",
			},
		},
		// Unmarshalling is case-sensitive for big struct.
		{
			// "b" should have been "B", "I" should have been "i"
			data:        []byte(`{"kind":"Test","apiVersion":"other/blah","spec": {"A": 1, "b": 2, "h": 3, "I": 4}}`),
			into:        &testDecodable{},
			typer:       &mockTyper{err: runtime.NewNotRegisteredErrForKind("mock", schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"})},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodable{
				TypeMeta: metav1.TypeMeta{APIVersion: "other/blah", Kind: "Test"},
				Spec:     DecodableSpec{A: 1, H: 3},
			},
		},
		// Unknown fields should return an error from the strict JSON deserializer.
		{
			data:        []byte(`{"unknown": 1}`),
			into:        &testDecodable{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			errFn: func(err error) bool {
				return strings.Contains(err.Error(), `unknown field "unknown"`)
			},
			strict: true,
		},
		// Unknown fields should return an error from the strict YAML deserializer.
		{
			data:        []byte("unknown: 1\n"),
			into:        &testDecodable{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			errFn: func(err error) bool {
				return strings.Contains(err.Error(), `unknown field "unknown"`)
			},
			yaml:   true,
			strict: true,
		},
		// Duplicate fields should return an error from the strict JSON deserializer.
		{
			data:        []byte(`{"value":1,"value":1}`),
			into:        &testDecodable{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			errFn: func(err error) bool {
				return strings.Contains(err.Error(), `duplicate field "value"`)
			},
			strict: true,
		},
		// Duplicate fields should return an error from the strict YAML deserializer.
		{
			data: []byte("value: 1\n" +
				"value: 1\n"),
			into:        &testDecodable{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			errFn: func(err error) bool {
				return strings.Contains(err.Error(), `"value" already set in map`)
			},
			yaml:   true,
			strict: true,
		},
		// Duplicate fields should return an error from the strict JSON deserializer for unstructured.
		{
			data:        []byte(`{"value":1,"value":1}`),
			into:        &unstructured.Unstructured{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK: &schema.GroupVersionKind{},
			errFn: func(err error) bool {
				return strings.Contains(err.Error(), `duplicate field "value"`)
			},
			strict: true,
		},
		// Duplicate fields should return an error from the strict YAML deserializer for unstructured.
		{
			data: []byte("value: 1\n" +
				"value: 1\n"),
			into:        &unstructured.Unstructured{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK: &schema.GroupVersionKind{},
			errFn: func(err error) bool {
				return strings.Contains(err.Error(), `"value" already set in map`)
			},
			yaml:   true,
			strict: true,
		},
		// Strict JSON decode into unregistered objects directly.
		{
			data:        []byte(`{"kind":"Test","apiVersion":"other/blah","value":1,"Other":"test"}`),
			into:        &testDecodable{},
			typer:       &mockTyper{err: runtime.NewNotRegisteredErrForKind("mock", schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"})},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodable{
				TypeMeta: metav1.TypeMeta{APIVersion: "other/blah", Kind: "Test"},
				Other:    "test",
				Value:    1,
			},
			strict: true,
		},
		// Strict YAML decode into unregistered objects directly.
		{
			data: []byte("kind: Test\n" +
				"apiVersion: other/blah\n" +
				"value: 1\n" +
				"Other: test\n"),
			into:        &testDecodable{},
			typer:       &mockTyper{err: runtime.NewNotRegisteredErrForKind("mock", schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"})},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodable{
				TypeMeta: metav1.TypeMeta{APIVersion: "other/blah", Kind: "Test"},
				Other:    "test",
				Value:    1,
			},
			yaml:   true,
			strict: true,
		},
		// Valid strict JSON decode without GVK.
		{
			data:        []byte(`{"value":1234}`),
			into:        &testDecodable{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodable{
				Value: 1234,
			},
			strict: true,
		},
		// Valid strict YAML decode without GVK.
		{
			data:        []byte("value: 1234\n"),
			into:        &testDecodable{},
			typer:       &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodable{
				Value: 1234,
			},
			yaml:   true,
			strict: true,
		},

		// coerce from null
		{
			data:           []byte(`{"bool":null,"int":null,"int32":null,"int64":null,"float32":null,"float64":null,"string":null,"array":null,"map":null,"struct":null}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{},
			strict:         true,
		},
		{
			data:           []byte(`{"bool":null,"int":null,"int32":null,"int64":null,"float32":null,"float64":null,"string":null,"array":null,"map":null,"struct":null}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{},
			yaml:           true,
			strict:         true,
		},
		// coerce from string
		{
			data:           []byte(`{"string":""}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{},
			strict:         true,
		},
		{
			data:           []byte(`{"string":""}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{},
			yaml:           true,
			strict:         true,
		},
		// coerce from array
		{
			data:           []byte(`{"array":[]}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{Array: []string{}},
			strict:         true,
		},
		{
			data:           []byte(`{"array":[]}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{Array: []string{}},
			yaml:           true,
			strict:         true,
		},
		// coerce from map
		{
			data:           []byte(`{"map":{},"struct":{}}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{Map: map[string]string{}},
			strict:         true,
		},
		{
			data:           []byte(`{"map":{},"struct":{}}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{Map: map[string]string{}},
			yaml:           true,
			strict:         true,
		},
		// coerce from int
		{
			data:           []byte(`{"int":1,"int32":1,"int64":1,"float32":1,"float64":1}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{Int: 1, Int32: 1, Int64: 1, Float32: 1, Float64: 1},
			strict:         true,
		},
		{
			data:           []byte(`{"int":1,"int32":1,"int64":1,"float32":1,"float64":1}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{Int: 1, Int32: 1, Int64: 1, Float32: 1, Float64: 1},
			yaml:           true,
			strict:         true,
		},
		// coerce from float
		{
			data:           []byte(`{"float32":1.0,"float64":1.0}`),
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{Float32: 1, Float64: 1},
			strict:         true,
		},
		{
			data:           []byte(`{"int":1.0,"int32":1.0,"int64":1.0,"float32":1.0,"float64":1.0}`), // floating point gets dropped in yaml -> json step
			into:           &testDecodeCoercion{},
			typer:          &mockTyper{gvk: &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}},
			expectedGVK:    &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"},
			expectedObject: &testDecodeCoercion{Int: 1, Int32: 1, Int64: 1, Float32: 1, Float64: 1},
			yaml:           true,
			strict:         true,
		},
	}

	for i, test := range testCases {
		var s runtime.Serializer
		if test.yaml {
			s = json.NewSerializerWithOptions(json.DefaultMetaFactory, test.creater, test.typer, json.SerializerOptions{Yaml: test.yaml, Pretty: false, Strict: test.strict})
		} else {
			s = json.NewSerializerWithOptions(json.DefaultMetaFactory, test.creater, test.typer, json.SerializerOptions{Yaml: test.yaml, Pretty: test.pretty, Strict: test.strict})
		}
		obj, gvk, err := s.Decode([]byte(test.data), test.defaultGVK, test.into)

		if !reflect.DeepEqual(test.expectedGVK, gvk) {
			t.Errorf("%d: unexpected GVK: %v", i, gvk)
		}

		switch {
		case err == nil && test.errFn != nil:
			t.Errorf("%d: failed: not getting the expected error", i)
			continue
		case err != nil && test.errFn == nil:
			t.Errorf("%d: failed: %v", i, err)
			continue
		case err != nil:
			if !test.errFn(err) {
				t.Errorf("%d: failed: %v", i, err)
			}
			if !runtime.IsStrictDecodingError(err) && obj != nil {
				t.Errorf("%d: should have returned nil object", i)
			}
			continue
		}

		if test.into != nil && test.into != obj {
			t.Errorf("%d: expected into to be returned: %v", i, obj)
			continue
		}

		if !reflect.DeepEqual(test.expectedObject, obj) {
			t.Errorf("%d: unexpected object:\n%s", i, diff.ObjectGoPrintSideBySide(test.expectedObject, obj))
		}
	}
}

func TestCacheableObject(t *testing.T) {
	gvk := schema.GroupVersionKind{Group: "group", Version: "version", Kind: "MockCacheableObject"}
	creater := &mockCreater{obj: &runtimetesting.MockCacheableObject{}}
	typer := &mockTyper{gvk: &gvk}
	serializer := json.NewSerializer(json.DefaultMetaFactory, creater, typer, false)

	runtimetesting.CacheableObjectTest(t, serializer)
}

type mockCreater struct {
	apiVersion string
	kind       string
	err        error
	obj        runtime.Object
}

func (c *mockCreater) New(kind schema.GroupVersionKind) (runtime.Object, error) {
	c.apiVersion, c.kind = kind.GroupVersion().String(), kind.Kind
	return c.obj, c.err
}

type mockTyper struct {
	gvk *schema.GroupVersionKind
	err error
}

func (t *mockTyper) ObjectKinds(obj runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	if t.gvk == nil {
		return nil, false, t.err
	}
	return []schema.GroupVersionKind{*t.gvk}, false, t.err
}

func (t *mockTyper) Recognizes(_ schema.GroupVersionKind) bool {
	return false
}

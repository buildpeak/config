package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/buildpeak/config"
)

const PASSWORD = "123456"

func TestMain(m *testing.M) {
	os.Setenv("PASSWORD", PASSWORD)

	os.Exit(m.Run())
}

func TestFind(t *testing.T) {
	y := `version: 1.0.0

string: test string
int: 100
float: 123.456
bool: yes

list:
  - one
  - two
  - three

object:
  attr_one: 1
  attr_two:
    - 1
    - 2
  attr_three:
    one: alpha
    two: bravo
    three: charlie
`
	var tests = []struct {
		key string
		ok  bool
		exp interface{}
	}{
		{"string", true, "test string"},
		{"int", true, 100},
		{"float", true, 123.456},
		{"bool", true, true},
		{"list", true, []interface{}{"one", "two", "three"}},
		{"object.attr_one", true, 1},
		{"object.attr_two", true, []interface{}{1, 2}},
		{"object.attr_three.one", true, "alpha"},
		{"object.attr_three.two", true, "bravo"},
		{"object.attr_three.three", true, "charlie"},
		{"notfound", false, nil},
		{"object.attr_four", false, nil},
		{"object.attr_three.four", false, nil},
	}
	c, err := Decode([]byte(y))
	if err != nil {
		t.Fatalf("Decoding testing yaml failed: %v", err)
	}

	for _, test := range tests {
		v, ok := c.Find(test.key)
		assert.Equalf(t, test.ok, ok, "key %s okay: want %t, got %t", test.key, test.ok, ok)
		assert.Equalf(t, test.exp, v, "key %s value: want %v, got %v", test.key, test.exp, v)
	}
}

func TestDecode(t *testing.T) {
	var tests = []struct {
		yaml string
		json string
		keys []string
		exps []interface{}
	}{
		{
			yaml: `name: Johnson`,
			json: `{"name":"Johnson"}`,
			keys: []string{`name`},
			exps: []interface{}{"Johnson"},
		},
		{
			yaml: `person:
  name: Johnson
  age: 80
  home.address: 11 Shenton Rd`,
			json: `{"person":{"name":"Johnson","age":80,"home.address":"11 Shenton Rd"}}`,
			keys: []string{`person.name`, `person.age`, `person.home\.address`},
			exps: []interface{}{"Johnson", 80, "11 Shenton Rd"},
		},
	}
	for _, test := range tests {
		cj, err := Decode([]byte(test.json))
		if err != nil {
			t.Errorf("Decoding %s failed: %s", test.json, err)
			continue
		}
		cy, err := Decode([]byte(test.yaml))
		if err != nil {
			t.Errorf("Decode %s failed: %s", test.yaml, err)
			continue
		}

		for i, key := range test.keys {
			jv, ok := cj.Find(key)
			if !ok {
				t.Errorf("key %s not found", key)
			}
			if jv != test.exps[i] {
				t.Errorf("Expectation unmet: want %v(%T) got %v(%T)\nJSON: %s",
					test.exps[i],
					test.exps[i],
					jv,
					jv,
					test.json)
			}

			yv, ok := cy.Find(key)
			if !ok {
				t.Errorf("key %s not found", key)
			}
			if yv != test.exps[i] {
				t.Errorf("Expectation unmet: want %v(%T) got %v(%T)\nYAML: %s",
					test.exps[i],
					test.exps[i],
					yv,
					yv,
					test.yaml)
			}
		}
	}
}

func TestGetEnv(t *testing.T) {
	var tests = []struct {
		yaml string
		json string
		keys []string
		exps []string
	}{
		{
			yaml: `password: ${PASSWORD}`,
			json: `{"password":"${PASSWORD}"}`,
			keys: []string{`password`},
			exps: []string{PASSWORD},
		},
	}
	for _, test := range tests {
		cj, err := Decode([]byte(test.json))
		if err != nil {
			t.Errorf("Decoding %s failed: %s", test.json, err)
			continue
		}
		cy, err := Decode([]byte(test.yaml))
		if err != nil {
			t.Errorf("Decode %s failed: %s", test.yaml, err)
			continue
		}

		for i, key := range test.keys {
			jv := cj.GetEnv(key)
			if jv != test.exps[i] {
				t.Errorf("Expectation unmet: want %v(%T) got %v(%T)\nJSON: %s",
					test.exps[i],
					test.exps[i],
					jv,
					jv,
					test.json)
			}

			yv := cy.GetEnv(key)
			if yv != test.exps[i] {
				t.Errorf("Expectation unmet: want %v(%T) got %v(%T)\nYAML: %s",
					test.exps[i],
					test.exps[i],
					yv,
					yv,
					test.yaml)
			}
		}
	}
}

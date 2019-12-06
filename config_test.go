package config_test

import (
	"os"
	"testing"

	. "github.com/buildpeak/config"
)

const PASSWORD = "123456"

func TestMain(m *testing.M) {
	os.Setenv("PASSWORD", PASSWORD)

	os.Exit(m.Run())
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
		cj := Decode([]byte(test.json))
		cy := Decode([]byte(test.yaml))

		for i, key := range test.keys {
			jv := cj.Get(key)
			if jv != test.exps[i] {
				t.Errorf("Expectation unmet: want %v(%T) got %v(%T)\nJSON: %s",
					test.exps[i],
					test.exps[i],
					jv,
					jv,
					test.json)
			}

			yv := cy.Get(key)
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
		cj := Decode([]byte(test.json))
		cy := Decode([]byte(test.yaml))
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

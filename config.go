package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"regexp"

	"gopkg.in/yaml.v2"
	//"github.com/goccy/go-yaml"
)

//Configor ...
type Configor interface {
	GetString(key string) string
	GetInt64(key string) int64
	GetFloat64(key string) float32
	GetBool(key string) bool
}

//Config ...
type Config struct {
	mp map[string]interface{}
}

//Load ...
func Load(filename string) *Config {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var c *Config
	if regexp.MustCompile(`(?i:\.json|\.ya?ml)$`).MatchString(filename) {
		c = Decode(content)
	} else {
		panic(errors.New("unsupported file type: " + path.Ext(filename)))
	}
	return c
}

//Decode ...
func Decode(b []byte) *Config {
	c := new(Config)
	if err := yaml.Unmarshal(b, &c.mp); err != nil {
		panic(err)
	}
	return c
}

//GetString ...
func (c Config) GetString(key string) string {
	return c.Get(key).(string)
}

//GetInt64 ...
func (c Config) GetInt64(key string) int64 {
	return c.Get(key).(int64)
}

//GetFloat64 returns a float64 value by key
func (c Config) GetFloat64(key string) float64 {
	return c.Get(key).(float64)
}

//GetBool returns a bool value under key
func (c Config) GetBool(key string) bool {
	return c.Get(key).(bool)
}

//GetEnv returns a environment variable with name under key
func (c Config) GetEnv(key string) string {
	vn := c.GetString(key)
	if vn[0] == '$' {
		vn = vn[1:]
	}
	if vn[0] == '{' && vn[len(vn)-1] == '}' {
		vn = vn[1 : len(vn)-1]
	}
	return os.Getenv(vn)
}

//Get a value from Config by key (a . separated string)
func (c Config) Get(key string) interface{} {
	keys, err := tokenizeString(key, '.', '\\')
	if err != nil {
		panic(err)
	}

	var mp = c.mp
	for _, ky := range keys {
		val, ok := mp[ky]
		if !ok {
			return val
		}

		if reflect.TypeOf(val).Kind() != reflect.Map {
			return val
		}
		if m, ok := val.(map[string]interface{}); ok {
			mp = m
			continue
		}

		if m, ok := val.(map[interface{}]interface{}); ok {
			for k, v := range m {
				mp[k.(string)] = v
			}
		}
	}

	return nil
}

//GetenvOr get an environment variable or returns the default value
func GetenvOr(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultVal
}

func tokenizeString(s string, sep, escape rune) (tokens []string, err error) {
	var runes []rune
	inEscape := false
	for _, r := range s {
		switch {
		case inEscape:
			inEscape = false
			fallthrough
		default:
			runes = append(runes, r)
		case r == escape:
			inEscape = true
		case r == sep:
			tokens = append(tokens, string(runes))
			runes = runes[:0]
		}
	}
	tokens = append(tokens, string(runes))
	if inEscape {
		err = errors.New("invalid terminal escape")
	}
	return tokens, err
}

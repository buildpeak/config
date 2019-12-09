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
func Load(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c *Config
	if regexp.MustCompile(`(?i:\.json|\.ya?ml)$`).MatchString(filename) {
		c, err = Decode(content)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unsupported file type: " + path.Ext(filename))
	}

	return c, nil
}

//Decode ...
func Decode(b []byte) (*Config, error) {
	c := new(Config)
	if err := yaml.Unmarshal(b, &c.mp); err != nil {
		return nil, err
	}
	return c, nil
}

//GetString ...
func (c Config) GetString(key string) string {
	if v, ok := c.Lookup(key); ok {
		return os.ExpandEnv(v.(string))
	}
	return ""
}

//GetInt64 ...
func (c Config) GetInt64(key string) int64 {
	if v, ok := c.Lookup(key); ok {
		return v.(int64)
	}
	return 0
}

//GetFloat64 returns a float64 value by key
func (c Config) GetFloat64(key string) float64 {
	if v, ok := c.Lookup(key); ok {
		return v.(float64)
	}
	return 0.0
}

//GetBool returns a bool value under key
func (c Config) GetBool(key string) bool {
	if v, ok := c.Lookup(key); ok {
		return v.(bool)
	}
	return false
}

//Lookup a value from Config with a key (. separated string)
func (c Config) Lookup(key string) (interface{}, bool) {
	keys, _ := tokenizeString(key, '.', '\\') // if key ends with \, just ignore

	var mp = c.mp
	for _, ky := range keys {
		val, ok := mp[ky]
		if !ok {
			return nil, false
		}

		if reflect.TypeOf(val).Kind() != reflect.Map {
			return val, true
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

	return nil, false
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

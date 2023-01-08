/*
github.com/steowens/yamlcfg (c) 2023 by Stephen Owens>

github.com/steowens/datastructures is licensed under a
Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License.

You should have received a copy of the license (LICENSE) along with this
work. If not, see <http://creativecommons.org/licenses/by-nc-sa/4.0/>.
*/
package yamlcfg

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type YamlType int64

const (
	String YamlType = iota
	Integer
	Float
	Bool
	Array
	Map
	Nil
)

var yamlTypeStrings = [...]string{"String", "Integer", "Float", "Bool", "Array", "Map", "Nil"}

func (typ YamlType) String() string {
	return yamlTypeStrings[typ]
}

type Config struct {
	values map[interface{}]interface{}
}

func (cfg *Config) GetString(path string) (string, error) {
	val, typ := cfg.Fetch(path)
	if typ == Nil {
		return "", nil
	}
	if typ == String {
		return val.(string), nil
	}
	if typ == Integer {
		return fmt.Sprintf("%v", val.(int)), nil
	}
	if typ == Float {
		return fmt.Sprintf("%v", val.(float64)), nil
	}
	if typ == Bool {
		return fmt.Sprintf("%v", val.(bool)), nil
	}
	return "", fmt.Errorf("Cannot convert %s into string", typ.String())
}

func (cfg *Config) GetInt(path string) (int64, error) {
	val, typ := cfg.Fetch(path)
	if typ == Nil {
		return 0, nil
	}
	if typ == String {
		intVar, err := strconv.ParseInt(val.(string), 0, 10)
		return intVar, err
	}
	if typ == Integer {
		return int64(val.(int)), nil
	}
	return 0, fmt.Errorf("Cannot convert %s into int64", typ.String())
}

func (cfg *Config) GetFloat(path string) (float64, error) {
	val, typ := cfg.Fetch(path)
	if typ == Nil {
		return 0.0, nil
	}
	if typ == Float {
		return val.(float64), nil
	}
	if typ == String {
		fltVar, err := strconv.ParseFloat(val.(string), 64)
		return fltVar, err
	}
	if typ == Integer {
		return float64(val.(int)), nil
	}
	return 0, fmt.Errorf("Cannot convert %s into float64", typ.String())
}

func (cfg *Config) GetBool(path string) (bool, error) {
	val, typ := cfg.Fetch(path)
	if typ == Bool {
		return val.(bool), nil
	}
	if typ == String {
		boolVar, err := strconv.ParseBool(val.(string))
		return boolVar, err
	}
	return false, fmt.Errorf("Cannot convert %s into float64", typ.String())
}

func (cfg *Config) Fetch(path string) (result interface{}, typ YamlType) {
	raw, typ := cfg.fetchRaw(path)
	if typ == Array {
		a, _ := raw.([]interface{})
		rg := make([]interface{}, len(a))
		for i := 0; i < len(a); i++ {
			v := a[i]
			vt := yamlTypeOf(v)
			if vt == Map {
				m, _ := v.(map[string]interface{})
				rg = append(rg, m)
			} else {
				rg = append(rg, v)
			}
		}
	} else if typ == Map {
		m, _ := raw.(map[interface{}]interface{})
		result = &Config{
			values: m,
		}
	} else {
		result = raw
	}
	return
}

func (cfg *Config) fetchRaw(path string) (result interface{}, typ YamlType) {
	splits := strings.Split(path, ".")
	result = fetch(splits, cfg.values)
	typ = yamlTypeOf(result)
	return
}

func yamlTypeOf(value interface{}) (typ YamlType) {
	if value == nil {
		typ = Nil
		return
	}
	_, ok := value.(string)
	if ok {
		typ = String
		return
	}
	_, ok = value.(float64)
	if ok {
		typ = Float
		return
	}
	_, ok = value.(int64)
	if ok {
		typ = Integer
		return
	}
	_, ok = value.(int)
	if ok {
		typ = Integer
		return
	}
	_, ok = value.(bool)
	if ok {
		typ = Bool
		return
	}
	_, ok = value.([]interface{})
	if ok {
		typ = Array
		return
	}
	_, ok = value.(map[interface{}]interface{})
	if ok {
		typ = Map
		return
	}
	return
}

func LoadFile(path string) (config *Config, err error) {
	input, e := ioutil.ReadFile(path)
	if e != nil {
		err = fmt.Errorf("Unable to load config file %s. Reason: %s", path, e.Error())
		return
	}
	return loadBytes(input)
}

func loadBytes(input []byte) (config *Config, err error) {
	structured := make(map[interface{}]interface{})
	e := yaml.Unmarshal(input, &structured)
	if e != nil {
		err = fmt.Errorf("Unable to unmarshal input stream. Reason: %s", e.Error())
		return
	}
	config = &Config{
		values: structured,
	}
	return
}

func fetch(path []string, values map[interface{}]interface{}) (result interface{}) {
	if values == nil {
		return
	}
	if path == nil || len(path) == 0 {
		return
	}
	if len(path) == 1 {
		result = values[path[0]]
		return
	}
	temp := values[path[0]]
	if temp == nil {
		return
	}
	m, ok := temp.(map[interface{}]interface{})
	if ok {
		result = fetch(path[1:], m)
	}
	return
}

/*
github.com/steowens/yamlcfg (c) 2023 by Stephen Owens>

github.com/steowens/datastructures is licensed under a
Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License.

You should have received a copy of the license (LICENSE) along with this
work. If not, see <http://creativecommons.org/licenses/by-nc-sa/4.0/>.
*/
package yamlcfg

import (
	"testing"
)

func TestYamlLoad(t *testing.T) {
	cfg, err := LoadFile("testconfig.yaml")
	if err != nil {
		t.Fatalf("Error loading input file: %s", err.Error())
	}
	if cfg == nil {
		t.Fatal("Loadfile returns nil cfg")
	}
	l := len(cfg.values)
	if l == 0 {
		t.Fatal("No values in returned config")
	}
}

func TestYamlTypes(t *testing.T) {
	cfg, err := LoadFile("testconfig.yaml")
	if err != nil {
		t.Fatalf("Error loading input file: %s", err.Error())
	}
	if cfg == nil {
		t.Fatal("Loadfile returns nil cfg")
	}
	testForType("rootstring", cfg, String, t)
	testForType("rootint", cfg, Integer, t)
	testForType("rootfloat", cfg, Float, t)
	testForType("rootbool", cfg, Bool, t)
	testForType("rootarray", cfg, Array, t)
	testForType("rootobj", cfg, Map, t)

	testForType("rootobj.subobj1", cfg, Map, t)
	testForType("rootobj.subarray1", cfg, Array, t)

	testForType("rootobj.subobj1.aString", cfg, String, t)
	testForType("rootobj.subobj1.anInt", cfg, Integer, t)
	testForType("rootobj.subobj1.aFloat", cfg, Float, t)
	testForType("rootobj.subobj1.aBool", cfg, Bool, t)

	testForType("rootobj.subobj1.aMartian", cfg, Nil, t)
}

func TestMapsBecomeConfigs(t *testing.T) {
	cfg, err := LoadFile("testconfig.yaml")
	if err != nil {
		t.Fatalf("Error loading input file: %s", err.Error())
	}
	if cfg == nil {
		t.Fatal("Loadfile returns nil cfg")
	}

	val, _ := cfg.Fetch("rootobj")
	_, ok := val.(*Config)
	if !ok {
		t.Fatal("rootobj does not come back as a *Config")
	}

	val, _ = cfg.Fetch("rootobj.subobj1")
	_, ok = val.(*Config)
	if !ok {
		t.Fatal("rootobj.subobj1 does not come back as a *Config")
	}
}

func testForType(path string, cfg *Config, expectedType YamlType, t *testing.T) {
	val, typ := cfg.fetchRaw(path)
	if val == nil && expectedType != Nil {
		t.Fatalf("Fetch returned nil at %s when fetching %s value", path, expectedType.String())
	}
	if typ != expectedType {
		t.Fatalf("Fetch returned %s type at %s when fetching value. Ex[ected %s", typ.String(), path, expectedType.String())
	}
}

func TestGetters(t *testing.T) {
	cfg, err := LoadFile("testconfig.yaml")
	if err != nil {
		t.Fatalf("Error loading input file: %s", err.Error())
	}
	if cfg == nil {
		t.Fatal("Loadfile returns nil cfg")
	}

	testGetString(cfg, t, "rootstring", "this is a string", false)
	testGetString(cfg, t, "rootint", "23466", false)
	testGetString(cfg, t, "rootfloat", "234.66", false)
	testGetString(cfg, t, "rootbool", "false", false)
	testGetString(cfg, t, "rootarray", "", true)
	testGetString(cfg, t, "rootobj", "", true)

	testGetInt(cfg, t, "rootstring", 0, true)
	testGetInt(cfg, t, "rootint", 23466, false)
	testGetInt(cfg, t, "rootfloat", 0, true)
	testGetInt(cfg, t, "rootbool", 0, true)
	testGetInt(cfg, t, "rootarray", 0, true)
	testGetInt(cfg, t, "rootobj", 0, true)

	testGetFloat(cfg, t, "rootstring", 0, true)
	testGetFloat(cfg, t, "rootint", 23466.0, false)
	testGetFloat(cfg, t, "rootfloat", 234.66, false)
	testGetFloat(cfg, t, "rootbool", 0, true)
	testGetFloat(cfg, t, "rootarray", 0, true)
	testGetFloat(cfg, t, "rootobj", 0, true)

	testGetBool(cfg, t, "rootstring", false, true)
	testGetBool(cfg, t, "rootint", false, true)
	testGetBool(cfg, t, "rootfloat", false, true)
	testGetBool(cfg, t, "rootbool", false, false)
	testGetBool(cfg, t, "rootarray", false, true)
	testGetBool(cfg, t, "rootobj", false, true)
	testGetBool(cfg, t, "rootobj.subobj1.aBool", true, false)
	testGetBool(cfg, t, "rootobj.subobj1.tfTrue", true, false)
	testGetBool(cfg, t, "rootobj.subobj1.tfFalse", false, false)
}

func testGetString(cfg *Config, t *testing.T, path string, expVal string, isErr bool) {
	val, err := cfg.GetString(path)
	if !isErr && err != nil {
		t.Fatalf("Got unexpected error %s at path %s", err.Error(), path)
		return
	}
	if isErr && err == nil {
		t.Fatalf("Expected GetString('%s') to error", path)
		return
	}
	if val != expVal {
		t.Fatalf("Value '%s' does not match expected value '%s'", val, expVal)
	}
	return
}

func testGetInt(cfg *Config, t *testing.T, path string, expVal int64, isErr bool) {
	val, err := cfg.GetInt(path)
	if !isErr && err != nil {
		t.Fatalf("Got unexpected error %s at path %s", err.Error(), path)
		return
	}
	if isErr && err == nil {
		t.Fatalf("Expected GetInt('%s') to error", path)
		return
	}
	if val != expVal {
		t.Fatalf("Value '%d' does not match expected value '%d'", val, expVal)
	}
	return
}

func testGetFloat(cfg *Config, t *testing.T, path string, expVal float64, isErr bool) {
	val, err := cfg.GetFloat(path)
	if !isErr && err != nil {
		t.Fatalf("Got unexpected error %s at path %s", err.Error(), path)
		return
	}
	if isErr && err == nil {
		t.Fatalf("Expected GetFloat('%s') to error", path)
		return
	}
	if val != expVal {
		t.Fatalf("Value '%f' does not match expected value '%f'", val, expVal)
	}
	return
}

func testGetBool(cfg *Config, t *testing.T, path string, expVal bool, isErr bool) {
	val, err := cfg.GetBool(path)
	if !isErr && err != nil {
		t.Fatalf("Got unexpected error %s at path %s", err.Error(), path)
		return
	}
	if isErr && err == nil {
		t.Fatalf("Expected GetFloat('%s') to error", path)
		return
	}
	if val != expVal {
		t.Fatalf("Value '%t' does not match expected value '%t'", val, expVal)
	}
	return
}

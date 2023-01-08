# github.com/steowens/yamlcfg

A handy yaml parser for generalized configuration without having to define config objecs.

Config values can be accessed by dotted notation.

Example usage:

    import (
        "fmt"
        "github.com/steowens/yamlcfg"
    )

    cfg, err := yamlcfg.LoadFile("testconfig.yaml")
    if err != nil {
        fmt.Fatalf("Got error loading config file %s", err.Error())
    }
    fmt.Printf("rootstring == %s", cfg.GetString("rootstring"))
    fmt.Printf("rootobj.subobj1.aString: %s", cfg.GetString("rootobj.subobj1.aString"))
    fmt.Printf("rootobj.subobj1.aFloat: %f", cfg.GetFloat("rootobj.subobj1.aFloat"))
    
---
github.com/steowens/yamlcfg (c) by Stephen Owens>

github.com/steowens/datastructures is licensed under a
Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License.

You should have received a copy of the license (LICENSE) along with this
work. If not, see <http://creativecommons.org/licenses/by-nc-sa/4.0/>.


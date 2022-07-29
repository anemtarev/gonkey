package parser

import "time"

// Fixture is a representation of the test data, that is preloaded to a redis database before the test starts
/* Example (yaml):
```yaml
inherits:
	- parent_template
	- child_template
	- other_fixture
databases:
	1:
		keys:
			$name: keys1
			values:
				a:
					value: 1
					expiration: 10s
				b:
					value: 2
	2:
		keys:
			$name: keys2
			values:
				c: 3
				d: 4
```
*/
type Fixture struct {
	Inherits  []string           `yaml:"inherits"`
	Templates Templates          `yaml:"templates"`
	Databases map[int]Database `yaml:"databases"`
}

type Templates struct {
	Keys []*Keys           `yaml:"keys"`
	Maps []*MapRecordValue `yaml:"maps"`
	Sets []*SetRecordValue `yaml:"sets"`
}

// Database contains data to load into Redis database
type Database struct {
	Maps *Maps `yaml:"maps"`
	Sets *Sets `yaml:"sets"`
	Keys *Keys `yaml:"keys"`
}

// Keys represent key/value pairs, that will be loaded into Redis database
type Keys struct {
	Name   string               `yaml:"$name"`
	Extend string               `yaml:"$extend"`
	Values map[string]*KeyValue `yaml:"values"`
}

// KeyValue represent a redis key/value pair
type KeyValue struct {
	Value      string        `yaml:"value"`
	Expiration time.Duration `yaml:"expiration"`
}

// Maps represent hash data structures, that will be loaded into Redis database
type Maps struct {
	Values map[string]*MapRecordValue `yaml:"values"`
}

// MapRecordValue represent a single hash data structure
type MapRecordValue struct {
	Name   string            `yaml:"$name"`
	Extend string            `yaml:"$extend"`
	Values map[string]string `yaml:"values"`
}

// Sets represent sets data structures, that will be loaded into Redis database
type Sets struct {
	Values map[string]*SetRecordValue `yaml:"values" toml:"" json:"values"`
}

// SetRecordValue represent a single set data structure
type SetRecordValue struct {
	Name   string               `yaml:"$name"`
	Extend string               `yaml:"$extend"`
	Values map[string]*SetValue `yaml:"values"`
}

// SetValue represent a set value object
type SetValue struct {
	Expiration time.Duration `yaml:"expiration" toml:"expiration" json:"expiration"`
}

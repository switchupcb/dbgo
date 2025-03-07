package extract

import "reflect"

// Symbols are extracted from the internal types (compiled at runtime).
var Symbols = make(map[string]map[string]reflect.Value)

//go:generate yaegi extract github.com/switchupcb/jet/v2/postgres

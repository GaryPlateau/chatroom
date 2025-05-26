package driver

import (
	"sync"
)

var lock = &sync.Mutex{}

package constant

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// ProjectRoot folder of this project
	ProjectRoot = filepath.Join(filepath.Dir(b), "..")
)

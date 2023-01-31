module github.com/lmpizarro/gorofex

go 1.18

replace github.com/lmpizarro/gorofex/pkg/lib => ../pkg/lib/

require (
	github.com/lmpizarro/gorofex/pkg/lib v0.0.0-00010101000000-000000000000
	gonum.org/v1/gonum v0.12.0
)

require golang.org/x/exp v0.0.0-20191002040644-a1355ae1e2c3 // indirect

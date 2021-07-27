module github.com/microsoft/BladeMonRT/tests/workflows

go 1.16

replace github.com/microsoft/BladeMonRT/workflows => ../../workflows

replace github.com/microsoft/BladeMonRT/nodes => ../../nodes

require (
	github.com/microsoft/BladeMonRT/nodes v0.0.0-00010101000000-000000000000
	github.com/microsoft/BladeMonRT/workflows v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	gotest.tools v2.2.0+incompatible
)
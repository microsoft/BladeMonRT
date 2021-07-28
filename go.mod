module github.com/microsoft/BladeMonRT/main

go 1.16

replace github.com/microsoft/BladeMonRT/workflows => ./workflows

replace github.com/microsoft/BladeMonRT/nodes => ./nodes

require (
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/microsoft/BladeMonRT v0.0.0-20210721172623-33610c6f9245 // indirect
	github.com/microsoft/BladeMonRT/nodes v0.0.0-20210721172623-33610c6f9245
	github.com/microsoft/BladeMonRT/tests/mocks v0.0.0-00010101000000-000000000000 // indirect
	github.com/microsoft/BladeMonRT/tests/mocks/nodes v0.0.0-00010101000000-000000000000 // indirect
	github.com/microsoft/BladeMonRT/workflows v0.0.0-20210721172623-33610c6f9245
	github.com/pkg/errors v0.9.1 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)

replace github.com/microsoft/BladeMonRT/tests/mocks => ./tests/mocks

replace github.com/microsoft/BladeMonRT/tests/main => ./

replace github.com/microsoft/BladeMonRT => ./

replace github.com/microsoft/BladeMonRT/tests/mocks/nodes => ./tests/mocks/nodes

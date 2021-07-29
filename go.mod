module github.com/microsoft/BladeMonRT/main

go 1.16

replace github.com/microsoft/BladeMonRT/workflows => ./workflows

replace github.com/microsoft/BladeMonRT/nodes => ./nodes

require (
	github.com/golang/mock v1.6.0
	github.com/microsoft/BladeMonRT/nodes v0.0.0-20210721172623-33610c6f9245
	github.com/microsoft/BladeMonRT/workflows v0.0.0-20210721172623-33610c6f9245
	gotest.tools v2.2.0+incompatible
)

replace github.com/microsoft/BladeMonRT/nodes/dummy_node => ../nodes/dummy_node

module github.com/microsoft/BladeMonRT/main

go 1.16

replace github.com/microsoft/BladeMonRT/workflows => ./workflows

replace github.com/microsoft/BladeMonRT/nodes => ./nodes

require (
	github.com/microsoft/BladeMonRT/nodes v0.0.0-20210721172623-33610c6f9245
	github.com/microsoft/BladeMonRT/workflows v0.0.0-20210721172623-33610c6f9245
)

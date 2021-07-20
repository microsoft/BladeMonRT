module example.com/node_tests

go 1.16

replace example.com/nodes => ../../nodes

require (
	example.com/nodes v0.0.0-00010101000000-000000000000 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)

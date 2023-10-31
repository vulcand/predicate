module github.com/vulcand/predicate

go 1.19

// we use a pseudo version for github.com/gravitational/trace
// because the a bump of GRPC has been made in this package and can influence predicate clients.
// https://github.com/gravitational/trace/compare/14a9a7dd6aaf...v1.1.17
require (
	github.com/gravitational/trace v1.1.16-0.20220114165159-14a9a7dd6aaf
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/term v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

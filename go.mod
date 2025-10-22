module github.com/keep94/itertools

go 1.23.0

require github.com/stretchr/testify v1.10.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Regression where Take consumes 1 extra value from its source iterator
retract v0.7.0

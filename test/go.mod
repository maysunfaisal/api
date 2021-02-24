module github.com/maysunfaisal/api/test

go 1.13

require (
	github.com/devfile/api/v2 v2.0.0
	github.com/ghodss/yaml v1.0.0
	github.com/santhosh-tekuri/jsonschema v1.2.4
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/devfile/api/v2 v2.0.0 => ../

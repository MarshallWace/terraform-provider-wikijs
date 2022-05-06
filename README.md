# Terraform Provider Wikijs 
Terraform provider to sync config for [Wiki.js](https://github.com/requarks/wiki) 

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install # builds binary to  $GOPATH/bin and does some package caching stuff
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```sh
$ go get github.com/author/dependency
$ go mod tidy # ensures go.mod file is in sync with source code dependencies
$ go mod vendor # downloads dependencies to 'vendor' folder
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider
### Directory layout
* Resources/Data Sources go into /wikijs
* Schema such as wikijs types go into /wikijs/schema. These schema are needed to make graphql requests. 

### Testing changes
1. Set WIKIJS_HOST and WIKIJS_TOKEN in env
2. In root directory, `make testacc` to run acceptance tests in all \*/\*_test.go files. 
3. To try out on actual terraform:
   1. In root directory, `make` to build binary and move it to correct location
   2. in examples/test, run `terraform init && terraform plan` to view output. Change the `test.tf` file accordingly.

### Notes on defining graphql schema
This provider syncs config via the Wikijs graphql API. The majority of the provider is 
powered by a [graphql library](https://github.com/hasura/go-graphql-client) which converts
Go structs into graphql requests. 

Note that when defining structs where it's variables will be used to form the request, you can use JSON tags
to convert the variables from uppercase to lowercase. The necessity for this is because the variables need to be 
uppercase in order for the JSON encoding library to work (public vs private in Go), but need to be lowercase when 
sent in the graphql request. For example:
```
type PageRuleInput struct {
	Id      gqlc.String   `json:"id"`
	Deny    gqlc.Boolean  `json:"deny"`
}
```


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources. In this case, API calls will be made to the WIKIJS_HOST provided. 
The created resources will still ultimately be deleted if the tests exit successfully.
![PhoenixNAP Logo](https://phoenixnap.com/wp-content/themes/phoenixnap-v2/img/v2/logo.svg)
Terraform phoenixNAP provider
==================


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.2+
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Building the provider
---------------------

Clone repository to: `$GOPATH/src/github.com/phoenixnap/terraform-provider-pnap`

```sh
$ mkdir -p $GOPATH/src/github.com/phoenixnap; cd $GOPATH/src/github.com/phoenixnap
$ git clone git@github.com:phoenixnap/terraform-provider-pnap
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/phoenixnap/terraform-provider-pnap
$ make build
```

Using the provider
----------------------

The PhoenixNAP provider will be installed on `terraform init` of a template using the `pnap_server` resource.


Developing the provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13.7+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-pnap
...
```


Testing provider code
---------------------------

We have acceptance tests in the provider. 

To run an acceptance test, find the relevant test function in `*_test.go` (for example TestAccPnapServer_basic), and run it as

```sh
TF_ACC=1 go test -v -timeout=20m -run=TestAccPnapServer_basic
```

If you want to see HTTP traffic, set `TF_LOG=DEBUG`, i.e.

```sh
TF_LOG=DEBUG TF_ACC=1 go test -v -timeout=20m -run=TestAccPnapServer_basic
```



Testing the provider with Terraform
---------------------------------------

Once you've built the plugin binary (see [Developing the provider](#developing-the-provider) above), it can be incorporated within your Terraform environment using the `-plugin-dir` option. Subsequent runs of Terraform will then use the plugin from your development environment.

```sh
$ terraform init -plugin-dir $GOPATH/bin
```


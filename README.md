<h1 align="center">
  <br>
  <a href="https://phoenixnap.com/bare-metal-cloud"><img src="https://user-images.githubusercontent.com/78744488/109779287-16da8600-7c06-11eb-81a1-97bf44983d33.png" alt="phoenixnap Bare Metal Cloud" width="300"></a>
  <br>
  Terraform phoenixNAP Provider
  <br>
</h1>

<p align="center">
Terraform is a powerful infrastructure as code tool for provisioning and managing cloud resources programmatically. phoenixNAP's <a href="https://phoenixnap.com/bare-metal-cloud">Bare Metal Cloud</a> server platform comes with a custom-built Terraform provider <i><b>pnap</b></i> which allows you to easily deploy and destroy your Bare Metal Cloud servers with code.
</p>

<p align="center">
  <a href="https://phoenixnap.com/bare-metal-cloud">Bare Metal Cloud</a> •
  <a href="https://registry.terraform.io/providers/phoenixnap/pnap/latest">Terraform Provider</a> •
  <a href="https://developers.phoenixnap.com/">Developers Portal</a> •
  <a href="http://phoenixnap.com/kb">Knowledge Base</a> •
  <a href="https://developers.phoenixnap.com/support">Support</a>
</p>

## Requirements
-	[Bare Metal Cloud](https://bmc.phoenixnap.com) account
-	[Terraform](https://www.terraform.io/downloads.html) 0.12.2+
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Creating a Bare Metal Cloud account
You need to have a Bare Metal Cloud account in order to use the ***pnap*** Terraform provider with Bare Metal Cloud. 

1. Go to the [Bare Metal Cloud signup page](https://support.phoenixnap.com/wap-jpost3/bmcSignup).
2. Follow the prompts to set up your account.
3. Use your credentials to [log in to Bare Metal Cloud portal](https://bmc.phoenixnap.com).

:arrow_forward: **Video tutorial:** [How to Create a Bare Metal Cloud Account](https://www.youtube.com/watch?v=RLRQOisEB-k)
<br>
:arrow_forward: **Video tutorial:** [Introduction to Bare Metal Cloud](https://www.youtube.com/watch?v=8TLsqgLDMN4)

## Installing Terraform locally
Follow this helpful tutorial to learn how to install Terraform on your local machine. 

-   [How To Install Terraform On CentOS 7/Ubuntu 18.04](https://phoenixnap.com/kb/how-to-install-terraform-centos-ubuntu)

## Building the provider

Clone the repository to: `$GOPATH/src/github.com/phoenixnap/terraform-provider-pnap`.

```sh
$ mkdir -p $GOPATH/src/github.com/phoenixnap; cd $GOPATH/src/github.com/phoenixnap
$ git clone git@github.com:phoenixnap/terraform-provider-pnap
```

Navigate to the provider directory and build the provider with `make build`.

```sh
$ cd $GOPATH/src/github.com/phoenixnap/terraform-provider-pnap
$ make build
```

## Using the provider

The *pnap* provider will be installed on `terraform init` as a template of the `pnap_server` resource.

[Terraform provider documentation](https://registry.terraform.io/providers/phoenixnap/pnap/latest/docs)

## Developing the provider

If you want to work on developing the provider, you need to have [Go](http://www.golang.org) installed on your machine. Go version 1.13.7+ is *required*. You will also need to properly set up a [GOPATH](http://golang.org/doc/code.html#GOPATH) and add `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-pnap
...
```

## Testing provider code

You can run acceptance tests with the provider. Find the relevant test function in `*_test.go` (e.g., TestAccPnapServer_basic) and run it as:

```sh
TF_ACC=1 go test -v -timeout=20m -run=TestAccPnapServer_basic
```

If you want to see HTTP traffic, set `TF_LOG=DEBUG`.

```sh
TF_LOG=DEBUG TF_ACC=1 go test -v -timeout=20m -run=TestAccPnapServer_basic
```

## Testing the provider with Terraform

Once you've built the plugin binary (see [Developing the provider](#developing-the-provider)), you can incorporate it into your Terraform environment using the `-plugin-dir` option. Subsequent runs of Terraform will use the plugin from your development environment.

```sh
$ terraform init -plugin-dir $GOPATH/bin
```

## Bare Metal Cloud community
Become part of the Bare Metal Cloud community to get updates on new features, help us improve the platform, and engage with developers and other users. 

-   Follow [@phoenixNAP on Twitter](https://twitter.com/phoenixnap)
-   Join the [official Slack channel](https://phoenixnap.slack.com)
-   Sign up for our [Developers Monthly newsletter](https://phoenixnap.com/developers-monthly-newsletter)


### Resources
-	[Product page](https://phoenixnap.com/bare-metal-cloud)
-	[Instance pricing](https://phoenixnap.com/bare-metal-cloud/instances)
-	[YouTube tutorials](https://www.youtube.com/watch?v=8TLsqgLDMN4&list=PLWcrQnFWd54WwkHM0oPpR1BrAhxlsy1Rc&ab_channel=PhoenixNAPGlobalITServices)
-	[Developers Portal](https://developers.phoenixnap.com)
-	[Knowledge Base](https://phoenixnap.com/kb)
-	[Blog](https:/phoenixnap.com/blog)

### Documentation
-	[API documentation](https://developers.phoenixnap.com/apis)

### Contact phoenixNAP
Get in touch with us if you have questions or need help with Bare Metal Cloud. 

<p align="left">
  <a href="https://twitter.com/phoenixNAP">Twitter</a> •
  <a href="https://www.facebook.com/phoenixnap">Facebook</a> •
  <a href="https://www.linkedin.com/company/phoenix-nap">LinkedIn</a> •
  <a href="https://www.instagram.com/phoenixnap">Instagram</a> •
  <a href="https://www.youtube.com/user/PhoenixNAPdatacenter">YouTube</a> •
  <a href="https://developers.phoenixnap.com/support">Email</a> 
</p>

<p align="center">
  <br>
  <a href="https://phoenixnap.com/bare-metal-cloud"><img src="https://user-images.githubusercontent.com/81640346/115243282-0c773b80-a123-11eb-9de7-59e3934a5712.jpg" alt="phoenixnap Bare Metal Cloud"></a>
</p>

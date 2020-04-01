package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

//TestMain inits sweeper
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/jesinity/terraform-provider-sigil/internal/provider"
)

var version = "dev"

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/jesinity/sigil",
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.New(os.Stderr, "", log.LstdFlags).Printf("provider serve error: %s", err)
		os.Exit(1)
	}
}

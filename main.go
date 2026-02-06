package main

import (
	"context"
	"log"
	"os"

	"github.com/jesinity/terraform-provider-cloudomen/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/jesinity/cloudomen",
	}

	if err := providerserver.Serve(context.Background(), provider.New("dev"), opts); err != nil {
		log.New(os.Stderr, "", log.LstdFlags).Printf("provider serve error: %s", err)
		os.Exit(1)
	}
}

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jesinity/terraform-provider-sigil/internal/naming"
)

func TestResolveCloudPrecedence(t *testing.T) {
	base := providerConfigModel{Cloud: types.StringValue(naming.CloudAzure)}
	override := providerConfigModel{Cloud: types.StringValue(naming.CloudAWS)}

	resolved := resolveCloud(types.StringValue(naming.CloudAzure), base, true, override, true)
	if resolved != naming.CloudAWS {
		t.Fatalf("expected override cloud to win (%q), got %q", naming.CloudAWS, resolved)
	}

	resolved = resolveCloud(types.StringNull(), base, true, providerConfigModel{}, false)
	if resolved != naming.CloudAzure {
		t.Fatalf("expected base cloud (%q), got %q", naming.CloudAzure, resolved)
	}

	resolved = resolveCloud(types.StringNull(), providerConfigModel{}, false, providerConfigModel{}, false)
	if resolved != naming.CloudAWS {
		t.Fatalf("expected default cloud (%q), got %q", naming.CloudAWS, resolved)
	}

	resolved = resolveCloud(types.StringValue(naming.CloudGCP), providerConfigModel{}, false, providerConfigModel{}, false)
	if resolved != naming.CloudGCP {
		t.Fatalf("expected explicit top-level cloud (%q), got %q", naming.CloudGCP, resolved)
	}
}

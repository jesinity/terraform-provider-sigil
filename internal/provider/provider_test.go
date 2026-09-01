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

func TestResolvePlatformPrecedence(t *testing.T) {
	base := providerConfigModel{Platform: types.StringValue(naming.PlatformDatabricks)}
	if got := resolvePlatform(types.StringNull(), base, true, providerConfigModel{}, false); got != naming.PlatformDatabricks {
		t.Fatalf("expected config platform, got %q", got)
	}
	if got := resolvePlatform(types.StringValue(naming.PlatformDatabricks), providerConfigModel{}, false, providerConfigModel{}, false); got != naming.PlatformDatabricks {
		t.Fatalf("expected top-level platform, got %q", got)
	}
	override := providerConfigModel{Platform: types.StringValue(naming.PlatformDatabricks)}
	if got := resolvePlatform(types.StringValue(""), base, true, override, true); got != naming.PlatformDatabricks {
		t.Fatalf("expected overrides platform, got %q", got)
	}
	if got := resolvePlatform(types.StringNull(), providerConfigModel{}, false, providerConfigModel{}, false); got != "" {
		t.Fatalf("expected empty platform default, got %q", got)
	}
}

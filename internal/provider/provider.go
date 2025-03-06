package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure HetznerRobotProvider satisfies various provider interfaces.
var _ provider.Provider = &HetznerRobotProvider{}

// HetznerRobotProvider defines the provider implementation.
type HetznerRobotProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// HetznerRobotProviderModel describes the provider data model.
type HetznerRobotProviderModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *HetznerRobotProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hetzner_robot"
	resp.Version = p.version
}

func (p *HetznerRobotProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Description: "The Hetzner Robot API username.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The Hetzner Robot API password.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *HetznerRobotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data HetznerRobotProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// You could add validation here

	tflog.Info(ctx, "Configured Hetzner Robot client", map[string]any{"username": data.Username.ValueString()})

	// You can create a client config struct to pass to resources and data sources
	// For now, we'll just pass the configuration model as is since we're not implementing resources yet
	resp.ResourceData = data
	resp.DataSourceData = data
}

func (p *HetznerRobotProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServerResource,
	}
}

func (p *HetznerRobotProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	// We don't have any data sources yet
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HetznerRobotProvider{
			version: version,
		}
	}
}

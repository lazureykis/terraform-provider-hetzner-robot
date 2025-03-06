package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/lazureykis/terraform-provider-hetzner-robot/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ServerResource{}
var _ resource.ResourceWithImportState = &ServerResource{}

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

// ServerResource defines the resource implementation.
type ServerResource struct {
	client *client.Client
}

// ServerResourceModel describes the resource data model.
type ServerResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ServerIP   types.String `tfsdk:"server_ip"`
	ServerIPv6 types.String `tfsdk:"server_ipv6"`
	Product    types.String `tfsdk:"product"`
	Datacenter types.String `tfsdk:"datacenter"`
	Status     types.String `tfsdk:"status"`
	RescueOS   types.String `tfsdk:"rescue_os"`
	Rescue     types.Bool   `tfsdk:"rescue"`
}

func (r *ServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Hetzner Robot dedicated server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Server ID.",
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Server name.",
				Required:    true,
			},
			"server_ip": schema.StringAttribute{
				Description: "Server main IPv4 address.",
				Computed:    true,
			},
			"server_ipv6": schema.StringAttribute{
				Description: "Server main IPv6 address.",
				Computed:    true,
			},
			"product": schema.StringAttribute{
				Description: "Server product/type.",
				Computed:    true,
			},
			"datacenter": schema.StringAttribute{
				Description: "Server datacenter location.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Server status.",
				Computed:    true,
			},
			"rescue_os": schema.StringAttribute{
				Description: "Rescue system OS type to enable. Valid values: linux64, linux32, freebsd64.",
				Optional:    true,
			},
			"rescue": schema.BoolAttribute{
				Description: "Whether rescue mode is enabled.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *ServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(HetznerRobotProviderModel)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected HetznerRobotProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client.NewClient(
		providerData.Username.ValueString(),
		providerData.Password.ValueString(),
	)
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// We don't actually create servers via the API, but we can manage existing ones
	// For this resource, Create actually means "Import the server by ID"
	server, err := r.client.GetServer(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server: %s", err))
		return
	}

	// Update server name if needed
	if data.Name.ValueString() != server.Name {
		server, err = r.client.UpdateServer(data.ID.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update server name: %s", err))
			return
		}
	}

	// Handle rescue mode
	if !data.Rescue.IsNull() && data.Rescue.ValueBool() && !data.RescueOS.IsNull() {
		err = r.client.EnableRescueMode(data.ID.ValueString(), client.ServerRescueOptions{
			OS: data.RescueOS.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to enable rescue mode: %s", err))
			return
		}
	}

	// Update model with actual server data
	data.ID = types.StringValue(server.ID)
	data.Name = types.StringValue(server.Name)
	data.ServerIP = types.StringValue(server.ServerIP)
	data.ServerIPv6 = types.StringValue(server.ServerIPv6)
	data.Product = types.StringValue(server.Product)
	data.Datacenter = types.StringValue(server.Datacenter)
	data.Status = types.StringValue(server.Status)

	if server.RescueOS != "" {
		data.RescueOS = types.StringValue(server.RescueOS)
		data.Rescue = types.BoolValue(true)
	} else {
		data.Rescue = types.BoolValue(false)
	}

	tflog.Trace(ctx, "created server resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.GetServer(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server: %s", err))
		return
	}

	// Update model with actual server data
	data.ID = types.StringValue(server.ID)
	data.Name = types.StringValue(server.Name)
	data.ServerIP = types.StringValue(server.ServerIP)
	data.ServerIPv6 = types.StringValue(server.ServerIPv6)
	data.Product = types.StringValue(server.Product)
	data.Datacenter = types.StringValue(server.Datacenter)
	data.Status = types.StringValue(server.Status)

	if server.RescueOS != "" {
		data.RescueOS = types.StringValue(server.RescueOS)
		data.Rescue = types.BoolValue(true)
	} else {
		data.Rescue = types.BoolValue(false)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read current state data
	var state ServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update server name if changed
	if !data.Name.Equal(state.Name) {
		server, err := r.client.UpdateServer(data.ID.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update server name: %s", err))
			return
		}

		// Update with returned data
		data.Name = types.StringValue(server.Name)
	}

	// Handle rescue mode changes
	rescueRequested := !data.Rescue.IsNull() && data.Rescue.ValueBool()
	rescueCurrent := !state.Rescue.IsNull() && state.Rescue.ValueBool()

	if rescueRequested != rescueCurrent {
		if rescueRequested {
			// Enable rescue mode
			if data.RescueOS.IsNull() {
				resp.Diagnostics.AddError(
					"Missing Parameter",
					"rescue_os must be specified when enabling rescue mode",
				)
				return
			}

			err := r.client.EnableRescueMode(data.ID.ValueString(), client.ServerRescueOptions{
				OS: data.RescueOS.ValueString(),
			})
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to enable rescue mode: %s", err))
				return
			}

			data.Rescue = types.BoolValue(true)
		} else {
			// Disable rescue mode
			err := r.client.DisableRescueMode(data.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to disable rescue mode: %s", err))
				return
			}

			data.Rescue = types.BoolValue(false)
			data.RescueOS = types.StringNull()
		}
	} else if rescueRequested && !data.RescueOS.Equal(state.RescueOS) {
		// Rescue mode is already enabled but OS type has changed
		// Need to disable and then re-enable with new OS
		err := r.client.DisableRescueMode(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to disable rescue mode: %s", err))
			return
		}

		err = r.client.EnableRescueMode(data.ID.ValueString(), client.ServerRescueOptions{
			OS: data.RescueOS.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to enable rescue mode: %s", err))
			return
		}

		data.Rescue = types.BoolValue(true)
	}

	// Sync with server to get current state
	server, err := r.client.GetServer(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server: %s", err))
		return
	}

	// Update computed attributes
	data.ServerIP = types.StringValue(server.ServerIP)
	data.ServerIPv6 = types.StringValue(server.ServerIPv6)
	data.Product = types.StringValue(server.Product)
	data.Datacenter = types.StringValue(server.Datacenter)
	data.Status = types.StringValue(server.Status)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Dedicated servers can't be deleted via the API
	// This just removes the server from Terraform state
	tflog.Trace(ctx, "removed server resource from state")
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set the ID from the import request
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

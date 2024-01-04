package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ctfer-io/go-ctfd/api"
	"github.com/ctfer-io/terraform-provider-ctfd/internal/provider/challenge"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = (*challengeDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*challengeDataSource)(nil)
)

func NewChallengeDataSource() datasource.DataSource {
	return &challengeDataSource{}
}

type challengeDataSource struct {
	client *api.Client
}

type challengesDataSourceModel struct {
	ID         types.String             `tfsdk:"id"`
	Challenges []challengeResourceModel `tfsdk:"challenges"`
}

func (ch *challengeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_challenges"
}

func (ch *challengeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"challenges": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Identifier of the challenge.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the challenge, displayed as it.",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							MarkdownDescription: "Category of the challenge that CTFd groups by on the web UI.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the challenge, consider using multiline descriptions for better style.",
							Computed:            true,
						},
						"connection_info": schema.StringAttribute{
							MarkdownDescription: "Connection Information to connect to the challenge instance, usefull for pwn or web pentest.",
							Computed:            true,
						},
						"max_attempts": schema.Int64Attribute{
							MarkdownDescription: "Maximum amount of attempts before being unable to flag the challenge.",
							Computed:            true,
						},
						"function": schema.StringAttribute{
							MarkdownDescription: "Decay function to define how the challenge value evolve through solves, either linear or logarithmic.",
							Computed:            true,
						},
						"value": schema.Int64Attribute{
							Computed: true,
						},
						"initial": schema.Int64Attribute{
							Computed: true,
						},
						"decay": schema.Int64Attribute{
							Computed: true,
						},
						"minimum": schema.Int64Attribute{
							Computed: true,
						},
						// TODO add support of next
						"state": schema.StringAttribute{
							MarkdownDescription: "State of the challenge, either hidden or visible.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the challenge defining its layout, either standard or dynamic.",
							Computed:            true,
						},
						"requirements": schema.SingleNestedAttribute{
							MarkdownDescription: "List of required challenges that needs to get flagged before this one being accessible. Usefull for skill-trees-like strategy CTF.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"behavior": schema.StringAttribute{
									MarkdownDescription: "Behavior if not unlocked, either hidden or anonymized.",
									Computed:            true,
								},
								"prerequisites": schema.ListAttribute{

									MarkdownDescription: "List of the challenges ID.",
									Computed:            true,
									ElementType:         types.StringType,
								},
							},
						},
						"flags": schema.ListNestedAttribute{
							MarkdownDescription: "List of challenge flags that solves it.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: challenge.FlagSubdatasourceAttributes(),
							},
							Computed: true,
						},
						"tags": schema.ListAttribute{
							MarkdownDescription: "List of challenge tags that will be displayed to the end-user. You could use them to give some quick insights of what a challenge involves.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"topics": schema.ListAttribute{
							MarkdownDescription: "List of challenge topics that are displayed to the administrators for maintenance and planification.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"hints": schema.ListNestedAttribute{
							MarkdownDescription: "List of hints about the challenge displayed to the end-user.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: challenge.HintSubdatasourceAttributes(),
							},
							Computed: true,
						},
						"files": schema.ListNestedAttribute{
							MarkdownDescription: "List of files given to players to flag the challenge.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: challenge.FileSubdatasourceAttributes(),
							},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (ch *challengeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *github.com/ctfer-io/go-ctfd/api.Client, got: %T. Please open an issue at https://github.com/ctfer-io/terraform-provider-ctfd", req.ProviderData),
		)
		return
	}

	ch.client = client
}

func (ch *challengeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state challengesDataSourceModel

	challs, err := ch.client.GetChallenges(&api.GetChallengesParams{}, api.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read CTFd Challenges",
			err.Error(),
		)
		return
	}

	state.Challenges = make([]challengeResourceModel, 0, len(challs))
	for _, c := range challs {
		chall := challengeResourceModel{
			ID: types.StringValue(strconv.Itoa(c.ID)),
		}
		chall.Read(ctx, resp.Diagnostics, ch.client)

		state.Challenges = append(state.Challenges, chall)
	}

	state.ID = types.StringValue("placeholder")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

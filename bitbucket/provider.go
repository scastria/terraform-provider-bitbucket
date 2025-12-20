package bitbucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-bitbucket/bitbucket/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BB_WORKSPACE", nil),
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("BB_ACCESS_TOKEN", nil),
				ConflictsWith: []string{"client_id", "client_secret"},
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("BB_CLIENT_ID", nil),
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"client_secret"},
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("BB_CLIENT_SECRET", nil),
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"client_id"},
			},
			"num_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BB_NUM_RETRIES", 3),
			},
			"retry_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BB_RETRY_DELAY", 30),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bitbucket_repository":                 resourceRepository(),
			"bitbucket_pipelines_config":           resourcePipelinesConfig(),
			"bitbucket_dynamic_pipelines_provider": resourceDynamicPipelinesProvider(),
			"bitbucket_webhook":                    resourceWebhook(),
			"bitbucket_environment":                resourceEnvironment(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bitbucket_project": dataSourceProject(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	workspace := d.Get("workspace").(string)
	accessToken := d.Get("access_token").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	numRetries := d.Get("num_retries").(int)
	retryDelay := d.Get("retry_delay").(int)

	//Check for valid authentication
	if (clientId == "") && (clientSecret == "") && (accessToken == "") {
		return nil, diag.Errorf("You must specify either client_id/client_secret for Client Credentials Authentication or access_token")
	}

	var diags diag.Diagnostics
	c, err := client.NewClient(ctx, workspace, accessToken, clientId, clientSecret, numRetries, retryDelay)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}

package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-bitbucket/bitbucket/client"
)

func resourceRepositoryPipelinesConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryPipelinesConfigCreate,
		ReadContext:   resourceRepositoryPipelinesConfigRead,
		UpdateContext: resourceRepositoryPipelinesConfigUpdate,
		DeleteContext: resourceRepositoryPipelinesConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func fillRepositoryPipelinesConfig(c *client.RepositoryPipelinesConfig, d *schema.ResourceData) {
	c.RepositoryId = d.Get("repository_id").(string)
	c.Enabled = d.Get("is_enabled").(bool)
}

func fillResourceDataFromRepositoryPipelinesConfig(c *client.RepositoryPipelinesConfig, d *schema.ResourceData) {
	d.Set("repository_id", c.RepositoryId)
	d.Set("is_enabled", c.Enabled)
}

func resourceRepositoryPipelinesConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId := d.Get("repository_id").(string)
	buf := bytes.Buffer{}
	newRepositoryPipelinesConfig := client.RepositoryPipelinesConfig{}
	fillRepositoryPipelinesConfig(&newRepositoryPipelinesConfig, d)
	err := json.NewEncoder(&buf).Encode(newRepositoryPipelinesConfig)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RepositoryPipelinesConfigPath, c.Workspace, repositoryId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.RepositoryPipelinesConfig{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = repositoryId
	fillResourceDataFromRepositoryPipelinesConfig(retVal, d)
	d.SetId(repositoryId)
	return diags
}

func resourceRepositoryPipelinesConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RepositoryPipelinesConfigPath, c.Workspace, d.Id())
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.RepositoryPipelinesConfig{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = d.Id()
	fillResourceDataFromRepositoryPipelinesConfig(retVal, d)
	return diags
}

func resourceRepositoryPipelinesConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upRepositoryPipelinesConfig := client.RepositoryPipelinesConfig{}
	fillRepositoryPipelinesConfig(&upRepositoryPipelinesConfig, d)
	err := json.NewEncoder(&buf).Encode(upRepositoryPipelinesConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RepositoryPipelinesConfigPath, c.Workspace, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.RepositoryPipelinesConfig{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.RepositoryId = d.Id()
	fillResourceDataFromRepositoryPipelinesConfig(retVal, d)
	return diags
}

func resourceRepositoryPipelinesConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

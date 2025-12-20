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

func resourcePipelinesConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelinesConfigCreate,
		ReadContext:   resourcePipelinesConfigRead,
		UpdateContext: resourcePipelinesConfigUpdate,
		DeleteContext: resourcePipelinesConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func fillPipelinesConfig(c *client.PipelinesConfig, d *schema.ResourceData) {
	c.RepositoryId = d.Get("repository_id").(string)
	c.Enabled = d.Get("is_enabled").(bool)
}

func fillResourceDataFromPipelinesConfig(c *client.PipelinesConfig, d *schema.ResourceData) {
	d.Set("repository_id", c.RepositoryId)
	d.Set("is_enabled", c.Enabled)
}

func resourcePipelinesConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId := d.Get("repository_id").(string)
	buf := bytes.Buffer{}
	newPipelinesConfig := client.PipelinesConfig{}
	fillPipelinesConfig(&newPipelinesConfig, d)
	err := json.NewEncoder(&buf).Encode(newPipelinesConfig)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.PipelinesConfigPath, c.Workspace, repositoryId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.PipelinesConfig{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = repositoryId
	fillResourceDataFromPipelinesConfig(retVal, d)
	d.SetId(repositoryId)
	return diags
}

func resourcePipelinesConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.PipelinesConfigPath, c.Workspace, d.Id())
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.PipelinesConfig{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = d.Id()
	fillResourceDataFromPipelinesConfig(retVal, d)
	return diags
}

func resourcePipelinesConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upPipelinesConfig := client.PipelinesConfig{}
	fillPipelinesConfig(&upPipelinesConfig, d)
	err := json.NewEncoder(&buf).Encode(upPipelinesConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.PipelinesConfigPath, c.Workspace, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.PipelinesConfig{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.RepositoryId = d.Id()
	fillResourceDataFromPipelinesConfig(retVal, d)
	return diags
}

func resourcePipelinesConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	// Set to false to simulate deletion
	upPipelinesConfig := client.PipelinesConfig{
		Enabled: false,
	}
	err := json.NewEncoder(&buf).Encode(upPipelinesConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.PipelinesConfigPath, c.Workspace, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, false, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

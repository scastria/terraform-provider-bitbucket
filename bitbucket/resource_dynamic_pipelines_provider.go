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

func resourceDynamicPipelinesProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDynamicPipelinesProviderCreate,
		ReadContext:   resourceDynamicPipelinesProviderRead,
		UpdateContext: resourceDynamicPipelinesProviderUpdate,
		DeleteContext: resourceDynamicPipelinesProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func fillDynamicPipelinesProvider(c *client.DynamicPipelinesProvider, d *schema.ResourceData) {
	c.RepositoryId = d.Get("repository_id").(string)
	c.AppAri = d.Get("provider_id").(string)
}

func fillResourceDataFromDynamicPipelinesProvider(c *client.DynamicPipelinesProvider, d *schema.ResourceData) {
	d.Set("repository_id", c.RepositoryId)
	d.Set("provider_id", c.AppAri)
}

func resourceDynamicPipelinesProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId := d.Get("repository_id").(string)
	buf := bytes.Buffer{}
	newDynamicPipelinesProvider := client.DynamicPipelinesProvider{}
	fillDynamicPipelinesProvider(&newDynamicPipelinesProvider, d)
	err := json.NewEncoder(&buf).Encode(newDynamicPipelinesProvider)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DynamicPipelinesProviderPath, c.Workspace, repositoryId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.DynamicPipelinesProvider{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = repositoryId
	fillResourceDataFromDynamicPipelinesProvider(retVal, d)
	d.SetId(repositoryId)
	return diags
}

func resourceDynamicPipelinesProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.DynamicPipelinesProviderPath, c.Workspace, d.Id())
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.DynamicPipelinesProvider{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = d.Id()
	fillResourceDataFromDynamicPipelinesProvider(retVal, d)
	return diags
}

func resourceDynamicPipelinesProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upDynamicPipelinesProvider := client.DynamicPipelinesProvider{}
	fillDynamicPipelinesProvider(&upDynamicPipelinesProvider, d)
	err := json.NewEncoder(&buf).Encode(upDynamicPipelinesProvider)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DynamicPipelinesProviderPath, c.Workspace, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.DynamicPipelinesProvider{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.RepositoryId = d.Id()
	fillResourceDataFromDynamicPipelinesProvider(retVal, d)
	return diags
}

func resourceDynamicPipelinesProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	// Set to false to simulate deletion
	upDynamicPipelinesProvider := client.DynamicPipelinesProvider{
		AppAri: "",
	}
	err := json.NewEncoder(&buf).Encode(upDynamicPipelinesProvider)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DynamicPipelinesProviderPath, c.Workspace, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

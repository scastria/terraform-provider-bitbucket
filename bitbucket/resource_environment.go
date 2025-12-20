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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-bitbucket/bitbucket/client"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Test", "Staging", "Production"}, false),
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func fillEnvironment(c *client.Environment, d *schema.ResourceData) {
	c.RepositoryId = d.Get("repository_id").(string)
	c.Name = d.Get("name").(string)
	c.Type.Name = d.Get("type").(string)
}

func fillResourceDataFromEnvironment(c *client.Environment, d *schema.ResourceData) {
	d.Set("repository_id", c.RepositoryId)
	d.Set("name", c.Name)
	d.Set("type", c.Type.Name)
	d.Set("uuid", c.Uuid)
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newEnvironment := client.Environment{}
	fillEnvironment(&newEnvironment, d)
	var retVal *client.Environment = nil
	if retVal == nil {
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(newEnvironment)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		requestPath := fmt.Sprintf(client.EnvironmentPath, c.Workspace, newEnvironment.RepositoryId)
		requestHeaders := http.Header{
			headers.ContentType: []string{client.ApplicationJson},
		}
		body, err := c.HttpRequest(ctx, false, http.MethodPost, requestPath, nil, requestHeaders, &buf)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		retVal = &client.Environment{}
		err = json.NewDecoder(body).Decode(retVal)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}
	retVal.RepositoryId = newEnvironment.RepositoryId
	fillResourceDataFromEnvironment(retVal, d)
	d.SetId(retVal.EnvironmentEncodeId())
	return diags
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId, id := client.EnvironmentDecodeId(d.Id())
	requestPath := fmt.Sprintf(client.EnvironmentPathGet, c.Workspace, repositoryId, id)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Environment{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = repositoryId
	fillResourceDataFromEnvironment(retVal, d)
	return diags
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId, id := client.EnvironmentDecodeId(d.Id())
	buf := bytes.Buffer{}
	upEnvironment := client.EnvironmentChanges{}
	upEnvironment.Change.Name = d.Get("name").(string)
	err := json.NewEncoder(&buf).Encode(upEnvironment)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.EnvironmentPathUpdate, c.Workspace, repositoryId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(ctx, false, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	// API does not return updated object, so we need to read it again
	requestPath = fmt.Sprintf(client.EnvironmentPathGet, c.Workspace, repositoryId, id)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Environment{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = repositoryId
	fillResourceDataFromEnvironment(retVal, d)
	return diags
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId, id := client.EnvironmentDecodeId(d.Id())
	requestPath := fmt.Sprintf(client.EnvironmentPathGet, c.Workspace, repositoryId, id)
	_, err := c.HttpRequest(ctx, false, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

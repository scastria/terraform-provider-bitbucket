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

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		UpdateContext: resourceRepositoryUpdate,
		DeleteContext: resourceRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"use_existing": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return d.Id() != "" },
			},
		},
	}
}

func fillRepository(c *client.Repository, d *schema.ResourceData) {
	c.Project.Uuid = d.Get("project_id").(string)
	c.Slug = d.Get("key").(string)
	name, ok := d.GetOk("name")
	if ok {
		c.Name = name.(string)
	}
	c.IsPrivate = d.Get("is_private").(bool)
}

func fillResourceDataFromRepository(c *client.Repository, d *schema.ResourceData) {
	d.Set("project_id", c.Project.Uuid)
	d.Set("key", c.Slug)
	d.Set("name", c.Name)
	d.Set("is_private", c.IsPrivate)
}

func resourceRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	key := d.Get("key").(string)
	useExisting := d.Get("use_existing").(bool)
	var body *bytes.Buffer = nil
	var err error
	if useExisting {
		// Try to read an existing repo with the given key and return it if found
		requestPath := fmt.Sprintf(client.RepositoryPath, c.Workspace, key)
		body, err = c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			re := err.(*client.RequestError)
			if re.StatusCode != http.StatusNotFound {
				return diag.FromErr(err)
			}
			body = nil
		}
	}
	if body == nil {
		buf := bytes.Buffer{}
		newRepository := client.Repository{}
		fillRepository(&newRepository, d)
		err := json.NewEncoder(&buf).Encode(newRepository)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		requestPath := fmt.Sprintf(client.RepositoryPath, c.Workspace, key)
		requestHeaders := http.Header{
			headers.ContentType: []string{client.ApplicationJson},
		}
		body, err = c.HttpRequest(ctx, false, http.MethodPost, requestPath, nil, requestHeaders, &buf)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}
	retVal := &client.Repository{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromRepository(retVal, d)
	d.SetId(retVal.Uuid)
	return diags
}

func resourceRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RepositoryPath, c.Workspace, d.Id())
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Repository{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromRepository(retVal, d)
	return diags
}

func resourceRepositoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upRepository := client.Repository{}
	fillRepository(&upRepository, d)
	err := json.NewEncoder(&buf).Encode(upRepository)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RepositoryPath, c.Workspace, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Repository{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	fillResourceDataFromRepository(retVal, d)
	return diags
}

func resourceRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RepositoryPath, c.Workspace, d.Id())
	_, err := c.HttpRequest(ctx, false, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

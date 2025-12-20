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

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebhookCreate,
		ReadContext:   resourceWebhookRead,
		UpdateContext: resourceWebhookUpdate,
		DeleteContext: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"events": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_active": {
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
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func fillWebhook(c *client.Webhook, d *schema.ResourceData) {
	c.RepositoryId = d.Get("repository_id").(string)
	c.Url = d.Get("url").(string)
	c.Events = convertSetToArray(d.Get("events").(*schema.Set))
	title, ok := d.GetOk("title")
	if ok {
		c.Description = title.(string)
	}
	c.Active = d.Get("is_active").(bool)
	c.UseExisting = d.Get("use_existing").(bool)
}

func fillResourceDataFromWebhook(c *client.Webhook, d *schema.ResourceData) {
	d.Set("repository_id", c.RepositoryId)
	d.Set("url", c.Url)
	d.Set("events", c.Events)
	d.Set("title", c.Description)
	d.Set("is_active", c.Active)
	d.Set("uuid", c.Uuid)
	d.Set("use_existing", c.UseExisting)
}

func resourceWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newWebhook := client.Webhook{}
	fillWebhook(&newWebhook, d)
	var retVal *client.Webhook = nil
	if newWebhook.UseExisting {
		// Try to find an existing webhook with the given url and return it if found
		// TODO: Paginate through results if many webhooks exist
		requestPath := fmt.Sprintf(client.WebhookPath, c.Workspace, newWebhook.RepositoryId)
		body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		retVals := &client.WebhookCollection{}
		err = json.NewDecoder(body).Decode(retVals)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		for _, wh := range retVals.Values {
			if wh.Url == newWebhook.Url {
				retVal = &wh
				break
			}
		}
	}
	if retVal == nil {
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(newWebhook)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		requestPath := fmt.Sprintf(client.WebhookPath, c.Workspace, newWebhook.RepositoryId)
		requestHeaders := http.Header{
			headers.ContentType: []string{client.ApplicationJson},
		}
		body, err := c.HttpRequest(ctx, false, http.MethodPost, requestPath, nil, requestHeaders, &buf)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		retVal = &client.Webhook{}
		err = json.NewDecoder(body).Decode(retVal)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}
	retVal.RepositoryId = newWebhook.RepositoryId
	fillResourceDataFromWebhook(retVal, d)
	d.SetId(retVal.WebhookEncodeId())
	return diags
}

func resourceWebhookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId, id := client.WebhookDecodeId(d.Id())
	requestPath := fmt.Sprintf(client.WebhookPathGet, c.Workspace, repositoryId, id)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Webhook{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = repositoryId
	fillResourceDataFromWebhook(retVal, d)
	return diags
}

func resourceWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId, id := client.WebhookDecodeId(d.Id())
	buf := bytes.Buffer{}
	upWebhook := client.Webhook{}
	fillWebhook(&upWebhook, d)
	err := json.NewEncoder(&buf).Encode(upWebhook)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.WebhookPathGet, c.Workspace, repositoryId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Webhook{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.RepositoryId = repositoryId
	fillResourceDataFromWebhook(retVal, d)
	return diags
}

func resourceWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId, id := client.WebhookDecodeId(d.Id())
	requestPath := fmt.Sprintf(client.WebhookPathGet, c.Workspace, repositoryId, id)
	_, err := c.HttpRequest(ctx, false, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

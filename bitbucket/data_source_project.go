package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-bitbucket/bitbucket/client"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"contains_repository_name"},
				AtLeastOneOf:  []string{"key", "contains_repository_name"},
			},
			"contains_repository_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key"},
				AtLeastOneOf:  []string{"key", "contains_repository_name"},
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	key := d.Get("key").(string)
	containsRepositoryName := d.Get("contains_repository_name").(string)
	if containsRepositoryName != "" {
		slug := convertNameToSlug(containsRepositoryName)
		requestPath := fmt.Sprintf(client.RepositoryPath, c.Workspace, slug)
		requestQuery := url.Values{}
		body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, requestQuery, nil, &bytes.Buffer{})
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
		key = retVal.Project.Key
	}
	requestPath := fmt.Sprintf(client.ProjectPath, c.Workspace, key)
	requestQuery := url.Values{}
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, requestQuery, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Project{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("key", retVal.Key)
	d.SetId(retVal.Uuid)
	return diags
}

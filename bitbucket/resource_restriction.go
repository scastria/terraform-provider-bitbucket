package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-bitbucket/bitbucket/client"
)

func resourceRestriction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRestrictionCreate,
		ReadContext:   resourceRestrictionRead,
		UpdateContext: resourceRestrictionUpdate,
		DeleteContext: resourceRestrictionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: resourceRestrictionDiff,
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: Implement all kinds
				//ValidateFunc: validation.StringInSlice([]string{"push", "delete", "force", "restrict_merges", "require_tasks_to_be_completed", "require_approvals_to_merge", "require_review_group_approvals_to_merge", "require_default_reviewer_approvals_to_merge", "require_no_changes_requested", "require_passing_builds_to_merge", "require_commits_behind", "reset_pullrequest_approvals_on_change", "smart_reset_pullrequest_approvals", "reset_pullrequest_changes_requested_on_change", "require_all_dependencies_merged", "enforce_merge_checks", "allow_auto_merge_when_builds_pass", "require_all_comments_resolved"}, false),
				ValidateFunc: validation.StringInSlice([]string{"require_approvals_to_merge", "require_default_reviewer_approvals_to_merge", "require_passing_builds_to_merge", "require_commits_behind"}, false),
			},
			"branch_match_kind": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"branching_model", "glob"}, false),
			},
			"branch_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"feature", "bugfix", "release", "hotfix", "development", "production"}, false),
			},
			"pattern": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"use_existing": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return d.Id() != "" },
			},
			"restriction_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceRestrictionDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	requiresBranchType := []string{"branching_model"}
	_, ok := d.GetOk("branch_type")
	if slices.Contains(requiresBranchType, d.Get("branch_match_kind").(string)) && !ok {
		return fmt.Errorf("branch_type must be set when branch_match_kind is one of: %v", requiresBranchType)
	}
	requiresPattern := []string{"glob"}
	_, ok = d.GetOk("pattern")
	if slices.Contains(requiresPattern, d.Get("branch_match_kind").(string)) && !ok {
		return fmt.Errorf("pattern must be set when branch_match_kind is one of: %v", requiresPattern)
	}
	requiresValue := []string{"require_approvals_to_merge", "require_default_reviewer_approvals_to_merge", "require_passing_builds_to_merge", "require_commits_behind"}
	_, ok = d.GetOk("value")
	if slices.Contains(requiresValue, d.Get("kind").(string)) && !ok {
		return fmt.Errorf("value must be set when kind is one of: %v", requiresValue)
	}
	return nil
}

func fillRestriction(c *client.Restriction, d *schema.ResourceData) {
	c.RepositoryId = d.Get("repository_id").(string)
	c.Kind = d.Get("kind").(string)
	c.BranchMatchKind = d.Get("branch_match_kind").(string)
	branchType, ok := d.GetOk("branch_type")
	if ok {
		c.BranchType = branchType.(string)
	}
	pattern, ok := d.GetOk("pattern")
	if ok {
		c.Pattern = pattern.(string)
	}
	value, ok := d.GetOk("value")
	if ok {
		c.Value = value.(int)
	}
	c.UseExisting = d.Get("use_existing").(bool)
}

func fillResourceDataFromRestriction(c *client.Restriction, d *schema.ResourceData) {
	d.Set("repository_id", c.RepositoryId)
	d.Set("kind", c.Kind)
	d.Set("branch_match_kind", c.BranchMatchKind)
	d.Set("branch_type", c.BranchType)
	d.Set("pattern", c.Pattern)
	d.Set("value", c.Value)
	d.Set("restriction_id", c.Id)
	d.Set("use_existing", c.UseExisting)
}

func resourceRestrictionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newRestriction := client.Restriction{}
	fillRestriction(&newRestriction, d)
	var retVal *client.Restriction = nil
	if newRestriction.UseExisting {
		// Try to find an existing restriction with the given kind, branch_match_kind, branch_type, and pattern and return it if found
		// TODO: Paginate through results if many restrictions exist
		requestPath := fmt.Sprintf(client.RestrictionPath, c.Workspace, newRestriction.RepositoryId)
		body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		retVals := &client.RestrictionCollection{}
		err = json.NewDecoder(body).Decode(retVals)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		for _, r := range retVals.Values {
			if (r.Kind == newRestriction.Kind) && (r.BranchMatchKind == newRestriction.BranchMatchKind) && (r.BranchType == newRestriction.BranchType) && (r.Pattern == newRestriction.Pattern) {
				retVal = &r
				break
			}
		}
	}
	if retVal == nil {
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(newRestriction)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		requestPath := fmt.Sprintf(client.RestrictionPath, c.Workspace, newRestriction.RepositoryId)
		requestHeaders := http.Header{
			headers.ContentType: []string{client.ApplicationJson},
		}
		body, err := c.HttpRequest(ctx, false, http.MethodPost, requestPath, nil, requestHeaders, &buf)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		retVal = &client.Restriction{}
		err = json.NewDecoder(body).Decode(retVal)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}
	retVal.RepositoryId = newRestriction.RepositoryId
	fillResourceDataFromRestriction(retVal, d)
	d.SetId(retVal.RestrictionEncodeId())
	return diags
}

func resourceRestrictionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId, id := client.RestrictionDecodeId(d.Id())
	requestPath := fmt.Sprintf(client.RestrictionPathGet, c.Workspace, repositoryId, id)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Restriction{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = repositoryId
	fillResourceDataFromRestriction(retVal, d)
	return diags
}

func resourceRestrictionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId, id := client.RestrictionDecodeId(d.Id())
	buf := bytes.Buffer{}
	upRestriction := client.Restriction{}
	fillRestriction(&upRestriction, d)
	err := json.NewEncoder(&buf).Encode(upRestriction)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RestrictionPathGet, c.Workspace, repositoryId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Restriction{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.RepositoryId = repositoryId
	fillResourceDataFromRestriction(retVal, d)
	return diags
}

func resourceRestrictionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	repositoryId, id := client.RestrictionDecodeId(d.Id())
	requestPath := fmt.Sprintf(client.RestrictionPathGet, c.Workspace, repositoryId, id)
	_, err := c.HttpRequest(ctx, false, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

package kibana

import (
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
	"github.com/ewilde/go-kibana"
)

func resourceDir() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaSearchCreate,
		Read:   resourceKibanaSearchRead,
		Update: resourceKibanaSearchUpdate,
		Delete: resourceKibanaSearchDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the kibana saved search",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the kibana saved search",
				Optional:    true,
			},
			"display_columns": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sort_by_columns": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sort_ascending": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Default:  false,
			},
			"search": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeString,
							Required: true,
						},
						"filters": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"query": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"type": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceKibanaSearchCreate(d *schema.ResourceData, meta interface{}) error {
	searchRequest, err := createKibanaSearchCreateRequestFromResourceData(d)
	if err != nil {
		return fmt.Errorf("failed to create kibana search api: %v error: %v", searchRequest, err)
	}

	api, err := meta.(*kibana.KibanaClient).Search().Create(searchRequest)

	if err != nil {
		return fmt.Errorf("failed to create kibana saved search: %v error: %v", searchRequest, err)
	}

	d.SetId(api.Id)

	return resourceKibanaSearchRead(d, meta)
}

func resourceKibanaSearchRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKibanaSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKibanaSearchDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func createKibanaSearchCreateRequestFromResourceData(d *schema.ResourceData) (*kibana.SearchRequest, error) {

	sortOrder := kibana.Descending;
	if readBoolFromResource(d, "sort_ascending") {
		sortOrder = kibana.Ascending
	}

	searchBuilder := kibana.NewSearchSourceBuilder()

	if v, _ := d.GetOk("search"); v != nil {
		searchSet := v.(*schema.Set).List()
		if len(searchSet) == 1 {
			searchMap := searchSet[0].(map[string]interface{})
			searchBuilder.WithIndexId(searchMap["index"].(string))
			filters := searchMap["filters"].(interface{})

			for _, x := range filters.([]interface{}) {
				matchSet := x.(map[string]interface{})["match"].(*schema.Set).List()
				match := matchSet[0].(map[string]interface{})
				searchBuilder.WithFilter(&kibana.SearchFilter{
					Query: &kibana.SearchFilterQuery{
						Match: map[string]*kibana.SearchFilterQueryAttributes{
							"@tags": {
								Query: match["query"].(string),
								Type:  match["type"].(string),
							},
						},
					},
				})
			}
		}
	}

	searchSource, err := searchBuilder.Build()
	if err != nil {
		return nil, err
	}

	return kibana.NewRequestBuilder().
		WithTitle(readStringFromResource(d, "name")).
		WithDescription(readStringFromResource(d, "description")).
		WithDisplayColumns(readArrayFromResource(d, "display_columns")).
		WithSortColumns(readArrayFromResource(d, "sort_by_columns"), sortOrder).
		WithSearchSource(searchSource).
		Build()
}

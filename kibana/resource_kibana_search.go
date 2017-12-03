package kibana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKibanaSearch() *schema.Resource {
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
						"filters": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"query": {
													Type:     schema.TypeString,
													Required: true,
												},
												"type": {
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
	response, err := meta.(*kibana.KibanaClient).Search().GetById(d.Id())

	if err != nil {
		return fmt.Errorf("could not find kibana search: %v", err)
	}

	d.Set("name", response.Attributes.Title)
	d.Set("description", response.Attributes.Description)
	d.Set("display_columns", response.Attributes.Columns)
	d.Set("sort_by_columns", response.Attributes.Sort[:len(response.Attributes.Sort)-1])

	sortAscending := false
	if response.Attributes.Sort[1] == "ASC" {
		sortAscending = true
	}

	d.Set("sort_ascending", sortAscending)

	responseSearch := &kibana.SearchSource{}
	if err := json.Unmarshal([]byte(response.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch); err != nil {
		return err
	}

	filters := make([]interface{}, 0, len(responseSearch.Filter))
	for _, x := range responseSearch.Filter {
		filters = append(filters, map[string]interface{}{"match": flattenMatches(x.Query)})
	}

	search := []interface{}{map[string]interface{}{
		"index":   responseSearch.IndexId,
		"filters": filters,
	}}

	if err := d.Set("search", search); err != nil {
		return err
	}

	return nil
}

func resourceKibanaSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	searchRequest, err := createKibanaSearchCreateRequestFromResourceData(d)
	if err != nil {
		return fmt.Errorf("failed to update kibana search api: %v error: %v", searchRequest, err)
	}

	_, err = meta.(*kibana.KibanaClient).Search().Update(d.Id(), &kibana.UpdateSearchRequest{Attributes: searchRequest.Attributes})

	if err != nil {
		return fmt.Errorf("failed to update kibana saved search: %v error: %v", searchRequest, err)
	}

	return resourceKibanaSearchRead(d, meta)
}

func resourceKibanaSearchDelete(d *schema.ResourceData, meta interface{}) error {
	err := meta.(*kibana.KibanaClient).Search().Delete(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kibana search: %v", err)
	}

	return nil
}

func createKibanaSearchCreateRequestFromResourceData(d *schema.ResourceData) (*kibana.CreateSearchRequest, error) {

	sortOrder := kibana.Descending
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
							match["field_name"].(string): {
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

func flattenMatches(searchFilterQuery *kibana.SearchFilterQuery) *schema.Set {

	s := schema.NewSet(matchHash, []interface{}{})
	for k, v := range searchFilterQuery.Match {
		s.Add(flattenMatch(k, v))
	}
	return s
}

func flattenMatch(field string, value *kibana.SearchFilterQueryAttributes) map[string]interface{} {
	m := map[string]interface{}{}
	m["field_name"] = field
	m["query"] = value.Query
	m["type"] = value.Type

	return m
}

func matchHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["field_name"].(string)))
	buf.WriteString(fmt.Sprintf("%s", m["query"].(string)))
	buf.WriteString(fmt.Sprintf("%s", m["type"].(string)))
	return hashcode.String(buf.String())
}

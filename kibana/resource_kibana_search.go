package kibana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
							Optional: true,
						},
						"index_ref_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"query": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"filters": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exists": {
										Type:     schema.TypeString,
										Optional: true,
									},
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
									"meta": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"index": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"index_ref_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"negate": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"disabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"alias": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"key": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"params": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
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
									},
								},
							},
							Optional: true,
						},
					},
				},
			},
			"references": {
				Type:        schema.TypeSet,
				Description: "A list of references",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceKibanaSearchCreate(d *schema.ResourceData, meta interface{}) error {
	searchClient := meta.(*kibana.KibanaClient).Search()
	searchRequest, err := createKibanaSearchCreateRequestFromResourceData(d, searchClient)
	if err != nil {
		return fmt.Errorf("failed to create kibana search api: %v error: %v", searchRequest, err)
	}

	log.Printf("[INFO] Creating Kibana search %s", searchRequest.Attributes.Title)

	api, err := searchClient.Create(searchRequest)

	if err != nil {
		return fmt.Errorf("failed to create kibana saved search: %v error: %v", searchRequest, err)
	}

	d.SetId(api.Id)
	return resourceKibanaSearchRead(d, meta)
}

func resourceKibanaSearchRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading Kibana search %s", d.Id())

	response, err := meta.(*kibana.KibanaClient).Search().GetById(d.Id())

	if err != nil {
		return handleNotFoundError(err, d)
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
		existsField := ""
		if x.Exists != nil {
			existsField = x.Exists.Field
		}

		filters = append(filters, map[string]interface{}{
			"exists": existsField,
			"match":  flattenMatches(x.Query),
			"meta":   flattenMeta(x.Meta),
		})
	}

	search := []interface{}{map[string]interface{}{
		"index":          responseSearch.IndexId,
		"index_ref_name": responseSearch.IndexRefName,
		"query":          extractQueryAsString(responseSearch.Query),
		"filters":        filters,
	}}

	if err := d.Set("search", search); err != nil {
		return err
	}

	d.Set("references", flattenSearchReferences(response.References))

	return nil
}
func resourceKibanaSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	searchClient := meta.(*kibana.KibanaClient).Search()
	searchRequest, err := createKibanaSearchCreateRequestFromResourceData(d, searchClient)
	if err != nil {
		return fmt.Errorf("failed to update kibana search api: %v error: %v", searchRequest, err)
	}

	log.Printf("[INFO] Creating Kibana search %s", searchRequest.Attributes.Title)

	_, err = searchClient.Update(d.Id(), &kibana.UpdateSearchRequest{Attributes: searchRequest.Attributes, References: searchRequest.References})

	if err != nil {
		return fmt.Errorf("failed to update kibana saved search: %v error: %v", searchRequest, err)
	}

	return resourceKibanaSearchRead(d, meta)
}

func resourceKibanaSearchDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Kibana search %s", d.Id())

	err := meta.(*kibana.KibanaClient).Search().Delete(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kibana search: %v", err)
	}

	d.SetId("")

	return nil
}

func createKibanaSearchCreateRequestFromResourceData(d *schema.ResourceData, searchClient kibana.SearchClient) (*kibana.CreateSearchRequest, error) {

	sortOrder := kibana.Descending
	if readBoolFromResource(d, "sort_ascending") {
		sortOrder = kibana.Ascending
	}

	searchBuilder := searchClient.NewSearchSource()

	if v, _ := d.GetOk("search"); v != nil {
		searchSet := v.(*schema.Set).List()
		if len(searchSet) == 1 {
			searchMap := searchSet[0].(map[string]interface{})
			if v, ok := searchMap["index"]; ok {
				searchBuilder.WithIndexId(v.(string))
			}
			if v, ok := searchMap["index_ref_name"]; ok {
				searchBuilder.WithIndexRefName(v.(string))
			}
			stringApplyIfExists(searchMap["query"], func(value string) {
				searchBuilder.WithQuery(value)
			})

			filters := searchMap["filters"].(interface{})

			for _, filter := range filters.([]interface{}) {
				matchSet := filter.(map[string]interface{})["match"].(*schema.Set).List()
				var query *kibana.SearchFilterQuery
				var meta *kibana.SearchFilterMetaData
				var existsFilter *kibana.SearchFilterExists

				if len(matchSet) > 0 {
					match := matchSet[0].(map[string]interface{})
					query = &kibana.SearchFilterQuery{
						Match: map[string]*kibana.SearchFilterQueryAttributes{
							match["field_name"].(string): {
								Query: match["query"].(string),
								Type:  match["type"].(string),
							},
						},
					}
				}

				if metaList, ok := filter.(map[string]interface{})["meta"]; ok {
					metaListSet := metaList.(*schema.Set).List()
					var params *kibana.SearchFilterQueryAttributes
					if len(metaListSet) > 0 {
						metaMap := metaListSet[0].(map[string]interface{})
						paramsListSet := metaMap["params"].(*schema.Set).List()
						if len(paramsListSet) > 0 {
							paramsMap := paramsListSet[0].(map[string]interface{})
							params = &kibana.SearchFilterQueryAttributes{
								Query: paramsMap["query"].(string),
								Type:  paramsMap["type"].(string),
							}
						}

						meta = &kibana.SearchFilterMetaData{
							Index:    metaMap["index"].(string),
							Negate:   boolOrDefault(metaMap["negate"], false),
							Disabled: boolOrDefault(metaMap["disabled"], false),
							Alias:    stringOrDefault(metaMap["alias"], ""),
							Type:     metaMap["type"].(string),
							Key:      metaMap["key"].(string),
							Value:    metaMap["value"].(string),
							Params:   params,
						}

						if v, ok := metaMap["index"]; ok {
							meta.Index = v.(string)
						}
						if v, ok := metaMap["index_ref_name"]; ok {
							meta.IndexRefName = v.(string)
						}
					}
				}

				stringApplyIfExists(filter.(map[string]interface{})["exists"], func(value string) {
					existsFilter = &kibana.SearchFilterExists{
						Field: value,
					}
				})

				searchBuilder.WithFilter(&kibana.SearchFilter{
					Query:  query,
					Meta:   meta,
					Exists: existsFilter,
				})
			}
		}
	}

	searchSource, err := searchBuilder.Build()
	if err != nil {
		return nil, err
	}

	request := kibana.NewSearchRequestBuilder().
		WithTitle(readStringFromResource(d, "name")).
		WithDescription(readStringFromResource(d, "description")).
		WithDisplayColumns(readArrayFromResource(d, "display_columns")).
		WithSortColumns(readArrayFromResource(d, "sort_by_columns"), sortOrder).
		WithSearchSource(searchSource)

	references := readSearchReferencesFromResource(d)
	if len(references) > 0 {
		request.WithReferences(references)
	}

	return request.Build()
}

func flattenMatches(searchFilterQuery *kibana.SearchFilterQuery) *schema.Set {
	s := schema.NewSet(matchHash, []interface{}{})
	if searchFilterQuery == nil {
		return s
	}

	for k, v := range searchFilterQuery.Match {
		s.Add(flattenMatch(k, v))
	}
	return s
}

func flattenMeta(searchFilterMetaData *kibana.SearchFilterMetaData) *schema.Set {

	s := schema.NewSet(metaHash, []interface{}{})
	m := map[string]interface{}{}

	if searchFilterMetaData == nil {
		return s
	}

	m["index"] = searchFilterMetaData.Index
	m["index_ref_name"] = searchFilterMetaData.IndexRefName
	m["negate"] = searchFilterMetaData.Negate
	m["disabled"] = searchFilterMetaData.Disabled
	m["alias"] = searchFilterMetaData.Alias
	m["type"] = searchFilterMetaData.Type
	m["key"] = searchFilterMetaData.Key
	m["value"] = searchFilterMetaData.Value
	m["params"] = flattenMetaParams(searchFilterMetaData.Params)
	s.Add(m)

	return s
}

func flattenMetaParams(searchFilterMetaData *kibana.SearchFilterQueryAttributes) *schema.Set {
	s := schema.NewSet(matchParamsHash, []interface{}{})

	if searchFilterMetaData == nil {
		return s
	}

	m := map[string]interface{}{}
	m["type"] = searchFilterMetaData.Type
	m["query"] = searchFilterMetaData.Query

	s.Add(m)

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

func matchParamsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s", m["query"].(string)))
	buf.WriteString(fmt.Sprintf("%s", m["type"].(string)))
	return hashcode.String(buf.String())
}

func metaHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["index"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["index_ref_name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	buf.WriteString(fmt.Sprintf("%s-", m["index"].(string)))
	buf.WriteString(fmt.Sprintf("%v", m["negate"].(bool)))
	buf.WriteString(fmt.Sprintf("%v", m["disabled"].(bool)))
	buf.WriteString(fmt.Sprintf("%s", m["alias"].(string)))
	buf.WriteString(fmt.Sprintf("%s", m["type"].(string)))
	buf.WriteString(fmt.Sprintf("%s", m["key"].(string)))
	buf.WriteString(fmt.Sprintf("%s", m["value"].(string)))
	return hashcode.String(buf.String())
}

func extractQueryAsString(query interface{}) string {
	if queryMap, ok := query.(map[string]interface{}); ok {
		if value, ok := queryMap["query_string"]; ok {
			return value.(map[string]interface{})["query"].(string)
		}

		if value, ok := queryMap["query"]; ok {
			return value.(string)
		}
	}

	return ""
}

func flattenSearchReferences(refs []*kibana.SearchReferences) []interface{} {
	if refs == nil {
		return nil
	}

	out := make([]interface{}, 0)

	for _, ref := range refs {
		if ref == nil {
			continue
		}

		out = append(out, map[string]interface{}{
			"id":   ref.Id,
			"name": ref.Name,
			"type": ref.Type.String(),
		})
	}

	return out
}

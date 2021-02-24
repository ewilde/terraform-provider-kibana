package kibana

import (
	"fmt"
	"log"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	goversion "github.com/mcuadros/go-version"
)

func resourceKibanaVisualization() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaVisualizationCreate,
		Read:   resourceKibanaVisualizationRead,
		Update: resourceKibanaVisualizationUpdate,
		Delete: resourceKibanaVisualizationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the kibana saved visualization",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the kibana saved visualization",
				Optional:    true,
			},
			"saved_search_id": {
				Type:        schema.TypeString,
				Description: "Saved search id this visualization is based on, 'references' and 'saved_search_id' are mutually exclusive, you may set one or the other, but not both",
				Optional:    true,
			},
			"references": {
				Type:        schema.TypeSet,
				Description: "A list of references, 'references' and 'saved_search_id' are mutually exclusive, you may set one or the other, but not both",
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
			"visualization_state": {
				Type:        schema.TypeString,
				Description: "Visualization state for this resource",
				Required:    true,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newJson, _ := structure.NormalizeJsonString(new)
					oldJson, _ := structure.NormalizeJsonString(old)
					return newJson == oldJson
				},
			},
			"search_source_json": {
				Type:        schema.TypeString,
				Description: "Search source json",
				Optional:    true,
				Default:     "{}",
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newJson, _ := structure.NormalizeJsonString(new)
					oldJson, _ := structure.NormalizeJsonString(old)
					return newJson == oldJson
				},
			},
		},
	}
}

func resourceKibanaVisualizationCreate(d *schema.ResourceData, meta interface{}) error {
	version := meta.(*kibana.KibanaClient).Config.KibanaVersion
	visualizationRequest, err := createKibanaVisualizationCreateRequestFromResourceData(d, version)
	if err != nil {
		return fmt.Errorf("failed to create kibana visualization api: %v error: %v", visualizationRequest, err)
	}

	log.Printf("[INFO] Creating Kibana visualization %s", visualizationRequest.Attributes.Title)

	api, err := meta.(*kibana.KibanaClient).Visualization().Create(visualizationRequest)

	if err != nil {
		return fmt.Errorf("failed to create kibana saved visualization: %v error: %v", visualizationRequest, err)
	}

	d.SetId(api.Id)
	return resourceKibanaVisualizationRead(d, meta)
}

func resourceKibanaVisualizationRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading Kibana visualization %s", d.Id())

	response, err := meta.(*kibana.KibanaClient).Visualization().GetById(d.Id())

	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.Set("name", response.Attributes.Title)
	d.Set("description", response.Attributes.Description)
	version := meta.(*kibana.KibanaClient).Config.KibanaVersion
	if goversion.Compare(version, "7.0.0", "<") {
		d.Set("saved_search_id", response.Attributes.SavedSearchId)
	} else {
		if len(response.References) == 1 &&
			response.References[0].Type == kibana.VisualizationReferencesTypeSearch {
			d.Set("saved_search_id", response.References[0].Id)
		}

	}
	err = d.Set("references", flattenVisualizationReferences(response.References))
	if err != nil {
		return err
	}
	if response.Attributes.KibanaSavedObjectMeta != nil {
		d.Set("search_source_json", response.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)
	}
	d.Set("visualization_state", response.Attributes.VisualizationState)

	return nil
}

func resourceKibanaVisualizationUpdate(d *schema.ResourceData, meta interface{}) error {
	version := meta.(*kibana.KibanaClient).Config.KibanaVersion
	visualizationRequest, err := createKibanaVisualizationCreateRequestFromResourceData(d, version)
	if err != nil {
		return fmt.Errorf("failed to update kibana visualization api: %v error: %v", visualizationRequest, err)
	}

	log.Printf("[INFO] Creating Kibana visualization %s", visualizationRequest.Attributes.Title)

	_, err = meta.(*kibana.KibanaClient).Visualization().Update(d.Id(), &kibana.UpdateVisualizationRequest{Attributes: visualizationRequest.Attributes, References: visualizationRequest.References})

	if err != nil {
		return fmt.Errorf("failed to update kibana saved visualization: %v error: %v", visualizationRequest, err)
	}

	return resourceKibanaVisualizationRead(d, meta)
}

func resourceKibanaVisualizationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Kibana visualization %s", d.Id())

	err := meta.(*kibana.KibanaClient).Visualization().Delete(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kibana visualization: %v", err)
	}

	d.SetId("")

	return nil
}

func createKibanaVisualizationCreateRequestFromResourceData(d *schema.ResourceData, version string) (*kibana.CreateVisualizationRequest, error) {
	request := kibana.NewVisualizationRequestBuilder().
		WithTitle(readStringFromResource(d, "name")).
		WithDescription(readStringFromResource(d, "description")).
		WithSavedSearchId(readStringFromResource(d, "saved_search_id")).
		WithVisualizationState(readStringFromResource(d, "visualization_state"))

	searchMeta := readStringFromResource(d, "search_source_json")
	if len(searchMeta) > 0 {
		request.WithKibanaSavedObjectMeta(&kibana.SearchKibanaSavedObjectMeta{SearchSourceJSON: searchMeta})
	}

	references := readVisualizationReferencesFromResource(d)
	if len(references) > 0 {
		request.WithReferences(references)
	}

	return request.Build(version)
}

func flattenVisualizationReferences(refs []*kibana.VisualizationReferences) []interface{} {
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

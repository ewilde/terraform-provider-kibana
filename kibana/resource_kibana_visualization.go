package kibana

import (
	"fmt"
	"github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"log"
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
				Description: "Saved search id this visualization is based on",
				Required:    true,
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
		},
	}
}

func resourceKibanaVisualizationCreate(d *schema.ResourceData, meta interface{}) error {
	visualizationRequest, err := createKibanaVisualizationCreateRequestFromResourceData(d)
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
	d.Set("saved_search_id", response.Attributes.SavedSearchId)
	d.Set("visualization_state", response.Attributes.VisualizationState)

	return nil
}

func resourceKibanaVisualizationUpdate(d *schema.ResourceData, meta interface{}) error {
	visualizationRequest, err := createKibanaVisualizationCreateRequestFromResourceData(d)
	if err != nil {
		return fmt.Errorf("failed to update kibana visualization api: %v error: %v", visualizationRequest, err)
	}

	log.Printf("[INFO] Creating Kibana visualization %s", visualizationRequest.Attributes.Title)

	_, err = meta.(*kibana.KibanaClient).Visualization().Update(d.Id(), &kibana.UpdateVisualizationRequest{Attributes: visualizationRequest.Attributes})

	if err != nil {
		return fmt.Errorf("failed to update kibana saved visualization: %v error: %v", visualizationRequest, err)
	}

	return resourceKibanaVisualizationRead(d, meta)
}

func resourceKibanaVisualizationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Kibana visualization %s", d.Id())

	err := meta.(*kibana.KibanaClient).Visualization().Delete(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kibana visualization: %v", err)
	}

	d.SetId("")

	return nil
}

func createKibanaVisualizationCreateRequestFromResourceData(d *schema.ResourceData) (*kibana.CreateVisualizationRequest, error) {
	return kibana.NewVisualizationRequestBuilder().
		WithTitle(readStringFromResource(d, "name")).
		WithDescription(readStringFromResource(d, "description")).
		WithSavedSearchId(readStringFromResource(d, "saved_search_id")).
		WithVisualizationState(readStringFromResource(d, "visualization_state")).
		Build()
}

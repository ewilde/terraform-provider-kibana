package kibana

import (
	"fmt"
	"github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"log"
)

func resourceKibanaDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaDashboardCreate,
		Read:   resourceKibanaDashboardRead,
		Update: resourceKibanaDashboardUpdate,
		Delete: resourceKibanaDashboardDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the kibana saved dashboard",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the kibana saved dashboard",
				Optional:    true,
			},
			"panels_json": {
				Type:        schema.TypeString,
				Description: "Panels json",
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
			"options_json": {
				Type:        schema.TypeString,
				Description: "Options json",
				Optional:    true,
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
			"ui_state_json": {
				Type:        schema.TypeString,
				Description: "Ui state json",
				Optional:    true,
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
			"time_restore": {
				Type:        schema.TypeBool,
				Description: "Saved the time selection",
				Optional:    true,
			},
			"search_source_json": {
				Type:        schema.TypeBool,
				Description: "Search source json",
				Optional:    true,
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

func resourceKibanaDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	dashboardRequest, err := createKibanaDashboardCreateRequestFromResourceData(d)
	if err != nil {
		return fmt.Errorf("failed to create kibana dashboard api: %v error: %v", dashboardRequest, err)
	}

	log.Printf("[INFO] Creating Kibana dashboard %s", dashboardRequest.Attributes.Title)

	api, err := meta.(*kibana.KibanaClient).Dashboard().Create(dashboardRequest)

	if err != nil {
		return fmt.Errorf("failed to create kibana saved dashboard: %v error: %v", dashboardRequest, err)
	}

	d.SetId(api.Id)
	return resourceKibanaDashboardRead(d, meta)
}

func resourceKibanaDashboardRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading Kibana dashboard %s", d.Id())

	response, err := meta.(*kibana.KibanaClient).Dashboard().GetById(d.Id())

	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.Set("name", response.Attributes.Title)
	d.Set("description", response.Attributes.Description)
	d.Set("panels_json", response.Attributes.PanelsJson)
	d.Set("options_json", response.Attributes.OptionsJson)
	d.Set("ui_state_json", response.Attributes.UiStateJSON)
	d.Set("time_restore", response.Attributes.TimeRestore)

	if response.Attributes.KibanaSavedObjectMeta != nil {
		d.Set("search_source_json", response.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)
	}

	return nil
}

func resourceKibanaDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	dashboardRequest, err := createKibanaDashboardCreateRequestFromResourceData(d)
	if err != nil {
		return fmt.Errorf("failed to update kibana dashboard api: %v error: %v", dashboardRequest, err)
	}

	log.Printf("[INFO] Creating Kibana dashboard %s", dashboardRequest.Attributes.Title)

	_, err = meta.(*kibana.KibanaClient).Dashboard().Update(d.Id(), &kibana.UpdateDashboardRequest{Attributes: dashboardRequest.Attributes})

	if err != nil {
		return fmt.Errorf("failed to update kibana saved dashboard: %v error: %v", dashboardRequest, err)
	}

	return resourceKibanaDashboardRead(d, meta)
}

func resourceKibanaDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Kibana dashboard %s", d.Id())

	err := meta.(*kibana.KibanaClient).Dashboard().Delete(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kibana dashboard: %v", err)
	}

	d.SetId("")

	return nil
}

func createKibanaDashboardCreateRequestFromResourceData(d *schema.ResourceData) (*kibana.CreateDashboardRequest, error) {
	request := kibana.NewDashboardRequestBuilder().
		WithTitle(readStringFromResource(d, "name")).
		WithDescription(readStringFromResource(d, "description")).
		WithPanelsJson(readStringFromResource(d, "panels_json")).
		WithOptionsJson(readStringFromResource(d, "options_json")).
		WithUiStateJson(readStringFromResource(d, "ui_state_json")).
		WithTimeRestore(readBoolFromResource(d, "time_restore"))

	searchMeta := readStringFromResource(d, "search_source_json")
	if len(searchMeta) > 0 {
		request.WithKibanaSavedObjectMeta(&kibana.SearchKibanaSavedObjectMeta{SearchSourceJSON: searchMeta})
	}

	return request.Build()
}

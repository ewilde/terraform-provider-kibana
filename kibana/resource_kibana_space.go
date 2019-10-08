package kibana

import (
	"fmt"
	"log"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKibanaSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaSpaceCreate,
		Read:   resourceKibanaSpaceRead,
		Update: resourceKibanaSpaceUpdate,
		Delete: resourceKibanaSpaceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Id of the kibana space",
				Required:    true,
			},
			"title": {
				Type:        schema.TypeString,
				Description: "Full name of the kibana space",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the kibana space",
				Optional:    true,
				Required:    false,
			},
			"color": {
				Type:        schema.TypeString,
				Description: "Color of the kibana space",
				Optional:    true,
				Required:    false,
			},
			"initials": {
				Type:        schema.TypeString,
				Description: "Initials of the kibana space",
				Optional:    true,
				Required:    false,
			},
			"imageurl": {
				Type:        schema.TypeString,
				Description: "the data-url encoded image to display in the space avatar",
				Optional:    true,
				Required:    false,
			},
			"disabled_features": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Required: false,
			},
		},
	}
}

func resourceKibanaSpaceCreate(data *schema.ResourceData, meta interface{}) error {
	spaceClient := meta.(*kibana.KibanaClient).Space()
	space, err := createKibanaSpaceCreateRequestFromResourceData(data, spaceClient)
	if err != nil {
		return err
	}
	err = spaceClient.Create(space)
	if err != nil {
		return err
	}
	data.SetId(space.Id)
	return resourceKibanaSpaceRead(data, meta)
}

func resourceKibanaSpaceUpdate(data *schema.ResourceData, meta interface{}) error {
	spaceClient := meta.(*kibana.KibanaClient).Space()
	space, err := createKibanaSpaceCreateRequestFromResourceData(data, spaceClient)
	if err != nil {
		return err
	}
	err = spaceClient.Update(space)
	if err != nil {
		return err
	}
	data.SetId(space.Id)
	return resourceKibanaSpaceRead(data, meta)
}

func createKibanaSpaceCreateRequestFromResourceData(data *schema.ResourceData, searchClient kibana.SpaceClient) (*kibana.Space, error) {
	space := &kibana.Space{
		Id:               readStringFromResource(data, "name"),
		Name:             readStringFromResource(data, "title"),
		Description:      readStringFromResource(data, "description"),
		Color:            readStringFromResource(data, "color"),
		Initials:         readStringFromResource(data, "initials"),
		ImageUrl:         readStringFromResource(data, "imageurl"),
		DisabledFeatures: readArrayFromResource(data, "disabled_features"),
	}
	return space, nil
}

func resourceKibanaSpaceRead(data *schema.ResourceData, meta interface{}) error {
	spaceClient := meta.(*kibana.KibanaClient).Space()

	spaceID := data.Id()

	space, err := spaceClient.GetByID(spaceID)
	if err != nil {
		return err
	}
	data.SetId(spaceID)
	data.Set("name", space.Id)
	data.Set("title", space.Name)
	data.Set("description", space.Description)
	data.Set("color", space.Color)
	data.Set("initials", space.Initials)
	data.Set("imageurl", space.ImageUrl)
	data.Set("disabled_features", space.DisabledFeatures)

	if err != nil {
		return err
	}

	return nil
}

func resourceKibanaSpaceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Kibana space %s", d.Id())
	err := meta.(*kibana.KibanaClient).Space().Delete(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kibana space: %v", err)
	}

	d.SetId("")

	return nil
}

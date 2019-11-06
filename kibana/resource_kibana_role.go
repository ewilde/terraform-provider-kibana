package kibana

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	kibana "github.com/ewilde/go-kibana"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKibanaRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaRoleCreate,
		Read:   resourceKibanaRoleRead,
		Update: resourceKibanaRoleCreate,
		Delete: resourceKibanaRoleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the kibana role",
				Required:    true,
			},
			"elasticsearch": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							MinItems: 1,
						},
						"indices": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"names": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
									"privileges": {
										MinItems: 1,
										Required: true,
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"run_as": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
			"kibana": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"spaces": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
						"feature": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "Name of the kibana role",
										Required:    true,
									},
									"privileges": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceKibanaRoleCreate(data *schema.ResourceData, meta interface{}) error {
	roleClient := meta.(*kibana.KibanaClient).Role()
	role, err := createKibanaRoleCreateRequestFromResourceData(data, roleClient)
	if err != nil {
		return err
	}
	err = roleClient.CreateOrUpdate(role)
	if err != nil {
		return err
	}
	data.SetId(role.Name)
	return resourceKibanaRoleRead(data, meta)
}

func createKibanaRoleCreateRequestFromResourceData(data *schema.ResourceData, searchClient kibana.RoleClient) (*kibana.Role, error) {
	role := &kibana.Role{
		Name:     readStringFromResource(data, "name"),
		Metadata: make(map[string]interface{}),
	}
	if v, ok := data.GetOk("elasticsearch"); ok {
		esRoleData := v.([]interface{})[0].(map[string]interface{})

		clusterNames := esRoleData["cluster"].([]interface{})
		clusters := make([]string, 0, len(clusterNames))
		for _, clusterName := range clusterNames {
			clusters = append(clusters, clusterName.(string))
		}
		runAsNames := esRoleData["run_as"].([]interface{})
		runAs := make([]string, 0, len(runAsNames))
		for _, runAsName := range runAsNames {
			runAs = append(runAs, runAsName.(string))
		}
		indices := make([]interface{}, 0)
		if v, ok := esRoleData["indices"]; ok {

			for _, indiceConfig := range v.([]interface{}) {
				indice := indiceConfig.(map[string]interface{})

				indiceNamesData := indice["names"].([]interface{})
				indiceNames := make([]string, 0, len(indiceNamesData))
				for _, indiceName := range indiceNamesData {
					indiceNames = append(indiceNames, indiceName.(string))
				}
				indicePrivilegesData := indice["privileges"].([]interface{})
				indicePrivileges := make([]string, 0, len(indicePrivilegesData))
				for _, privilege := range indicePrivilegesData {
					indicePrivileges = append(indicePrivileges, privilege.(string))
				}

				indices = append(indices, map[string]interface{}{
					"names":      indiceNames,
					"privileges": indicePrivileges,
				})
			}

		}
		role.ElasticSearch = &kibana.RoleElasticSearch{
			Cluster: clusters,
			RunAs:   runAs,
			Indices: indices,
		}
	}

	rolesKibana := []*kibana.RoleKibana{}
	if v, ok := data.GetOk("kibana"); ok {
		rolesKibana = make([]*kibana.RoleKibana, 0, len(v.([]interface{})))

		for _, kibanaConfig := range v.([]interface{}) {
			roleKibanaConfig := kibanaConfig.(map[string]interface{})

			basesConfig := roleKibanaConfig["base"].([]interface{})
			bases := make([]string, 0, len(basesConfig))
			for _, base := range basesConfig {
				bases = append(bases, base.(string))
			}

			spacesConfig := roleKibanaConfig["spaces"].([]interface{})
			spaces := make([]string, 0, len(spacesConfig))
			for _, space := range spacesConfig {
				spaces = append(spaces, space.(string))
			}
			features := make(map[string][]string, 0)
			if v, ok := roleKibanaConfig["feature"]; ok {

				for _, featureConfig := range v.(*schema.Set).List() {
					feature := featureConfig.(map[string]interface{})
					privilegesConfig := feature["privileges"].([]interface{})
					privileges := make([]string, 0, len(privilegesConfig))
					for _, privilege := range privilegesConfig {
						privileges = append(privileges, privilege.(string))
					}

					features[feature["name"].(string)] = privileges

				}
			}

			roleKibana := &kibana.RoleKibana{
				Base:    bases,
				Spaces:  spaces,
				Feature: features,
			}
			rolesKibana = append(rolesKibana, roleKibana)
		}
	}
	role.Kibana = rolesKibana
	return role, nil
}

func setRoleData(role *kibana.Role, data *schema.ResourceData) error {
	data.Set("name", role.Name)
	if role.ElasticSearch != nil {
		if err := data.Set("elasticsearch", flattenRoleElasticSearch(role.ElasticSearch)); err != nil {
			return err
		}
	}
	if role.Kibana != nil {
		if err := data.Set("kibana", flattenRoleKibana(role.Kibana)); err != nil {
			return err
		}
	}
	return nil
}

func flattenRoleElasticSearch(in *kibana.RoleElasticSearch) []interface{} {
	m := make(map[string]interface{})
	m["cluster"] = in.Cluster
	m["run_as"] = in.RunAs
	m["indices"] = flattenRoleElasticSearchIndices(in.Indices)
	return []interface{}{m}
}

func flattenRoleElasticSearchIndices(in []interface{}) []interface{} {
	var out = make([]interface{}, 0, 0)
	for _, v := range in {
		u := v.(map[string]interface{})
		out = append(out, map[string]interface{}{
			"names":      u["names"],
			"privileges": u["privileges"],
		})
	}
	return out
}
func flattenRoleKibana(in []*kibana.RoleKibana) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})
		m["base"] = v.Base
		m["spaces"] = v.Spaces
		m["feature"] = schema.NewSet(featureNameMappingHash, flattenRoleKibanaFeature(v.Feature))
		out[i] = m
	}
	return out
}
func featureNameMappingHash(v interface{}) int {
	var buf bytes.Buffer
	x := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf(x["name"].(string)))
	return hashcode.String(buf.String())
}

func flattenRoleKibanaFeature(in map[string][]string) []interface{} {
	var out = make([]interface{}, 0, 0)
	for k, v := range in {
		out = append(out, map[string]interface{}{
			"name":       k,
			"privileges": v,
		})
	}
	return out
}

func resourceKibanaRoleRead(data *schema.ResourceData, meta interface{}) error {
	roleClient := meta.(*kibana.KibanaClient).Role()

	roleID := data.Get("name").(string)

	role, err := roleClient.GetByID(roleID)
	if err != nil {
		return err
	}
	data.SetId(roleID)
	err = setRoleData(role, data)
	if err != nil {
		return err
	}

	return nil
}

func resourceKibanaRoleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Kibana role %s", d.Id())
	err := meta.(*kibana.KibanaClient).Role().Delete(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kibana role: %v", err)
	}

	d.SetId("")

	return nil
}

func mapStringStringToMapStringInterface(in map[string]string) map[string]interface{} {
	if in == nil || len(in) == 0 {
		return make(map[string]interface{}, 0)
	}

	mapped := make(map[string]interface{}, len(in))
	for k, v := range in {
		mapped[k] = v
	}
	return mapped
}

func mapStringSliceToMap(in []string) map[string]string {
	mapped := make(map[string]string, len(in))
	for _, v := range in {
		if len(v) > 0 {
			splitted := strings.Split(v, "=")
			key := splitted[0]
			value := splitted[1]
			mapped[key] = value
		}
	}
	return mapped
}

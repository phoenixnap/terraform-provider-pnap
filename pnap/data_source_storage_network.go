package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkstorageapi/storagenetwork"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	networkstorageapiclient "github.com/phoenixnap/go-sdk-bmc/networkstorageapi"
)

func dataSourceStorageNetwork() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceStorageNetworkRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ips": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path_suffix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capacity_in_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used_capacity_in_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permissions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nfs": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"read_write": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"read_only": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"root_squash": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"no_squash": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"all_squash": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_billing_tag": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"created_by": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceStorageNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := storagenetwork.NewGetStorageNetworksCommand(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	numOfStorageNets := 0
	for _, instance := range resp {
		if instance.Name != nil && *instance.Name == d.Get("name").(string) || instance.Id != nil && *instance.Id == d.Get("id").(string) {
			numOfStorageNets++
			d.SetId(*instance.Id)
			d.Set("name", *instance.Name)
			if instance.Description != nil {
				d.Set("description", *instance.Description)
			}
			if instance.Status != nil {
				d.Set("status", *instance.Status)
			}
			if instance.Location != nil {
				d.Set("location", *instance.Location)
			}
			if instance.NetworkId != nil {
				d.Set("network_id", *instance.NetworkId)
			}
			var ips []interface{}
			for _, v := range instance.Ips {
				ips = append(ips, v)
			}
			d.Set("ips", ips)
			if instance.CreatedOn != nil {
				createdOn := *instance.CreatedOn
				d.Set("created_on", createdOn.String())
			}
			volumes := flattenDataVolumes(instance.Volumes)

			if err := d.Set("volumes", volumes); err != nil {
				return err
			}
		}
	}
	if numOfStorageNets > 1 {
		return fmt.Errorf("too many storage networks with name %s (found %d, expected 1)", d.Get("name").(string), numOfStorageNets)
	}
	return nil
}

func flattenDataVolumes(volumes []networkstorageapiclient.Volume) []interface{} {
	if volumes != nil {
		vols := make([]interface{}, len(volumes))
		for i, v := range volumes {
			vol := make(map[string]interface{})
			if v.Id != nil {
				vol["id"] = *v.Id
			}
			if v.Name != nil {
				vol["name"] = *v.Name
			}
			if v.Description != nil {
				vol["description"] = *v.Description
			}
			if v.Path != nil {
				vol["path"] = *v.Path
			}
			if v.PathSuffix != nil {
				vol["path_suffix"] = *v.PathSuffix
			}
			if v.CapacityInGb != nil {
				vol["capacity_in_gb"] = int(*v.CapacityInGb)
			}
			if v.UsedCapacityInGb != nil {
				vol["used_capacity_in_gb"] = int(*v.UsedCapacityInGb)
			}
			if v.Protocol != nil {
				vol["protocol"] = *v.Protocol
			}
			if v.Status != nil {
				vol["status"] = *v.Status
			}
			if v.CreatedOn != nil {
				createdOn := *v.CreatedOn
				vol["created_on"] = createdOn.String()
			}
			if v.Permissions != nil {
				perms := make([]interface{}, 1)
				permsItem := make(map[string]interface{})
				permissions := *v.Permissions
				if permissions.Nfs != nil {
					nfs := *permissions.Nfs
					nf := make([]interface{}, 1)
					nfItem := make(map[string]interface{})
					var readWrite, readOnly, rootSquash, allSquash, noSquash []interface{}
					for _, v := range nfs.ReadWrite {
						readWrite = append(readWrite, v)
					}
					nfItem["read_write"] = readWrite
					for _, v := range nfs.ReadOnly {
						readOnly = append(readOnly, v)
					}
					nfItem["read_only"] = readOnly
					for _, v := range nfs.RootSquash {
						rootSquash = append(rootSquash, v)
					}
					nfItem["root_squash"] = rootSquash
					for _, v := range nfs.AllSquash {
						allSquash = append(allSquash, v)
					}
					nfItem["all_squash"] = allSquash
					for _, v := range nfs.NoSquash {
						noSquash = append(noSquash, v)
					}
					nfItem["no_squash"] = noSquash
					nf[0] = nfItem
					permsItem["nfs"] = nf
				}
				perms[0] = permsItem
				vol["permissions"] = perms
			}
			if v.Tags != nil && len(v.Tags) > 0 {
				tags := make([]interface{}, len(v.Tags))
				for i, j := range v.Tags {
					tagAssignment := make(map[string]interface{})
					tagAssignment["id"] = j.Id
					tagAssignment["name"] = j.Name
					if j.Value != nil {
						tagAssignment["value"] = *j.Value
					}
					tagAssignment["is_billing_tag"] = j.IsBillingTag
					if j.CreatedBy != nil {
						tagAssignment["created_by"] = *j.CreatedBy
					}
					tags[i] = tagAssignment
				}
				vol["tags"] = tags
			}
			vols[i] = vol
		}
		return vols
	}
	return make([]interface{}, 0)
}

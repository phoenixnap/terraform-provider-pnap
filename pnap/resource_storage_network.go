package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkstorageapi/storagenetwork"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"

	networkstorageapiclient "github.com/phoenixnap/go-sdk-bmc/networkstorageapi"
)

func resourceStorageNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageNetworkCreate,
		Read:   resourceStorageNetworkRead,
		Update: resourceStorageNetworkUpdate,
		Delete: resourceStorageNetworkDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(pnapRetryTimeout),
			Update: schema.DefaultTimeout(pnapRetryTimeout),
			Delete: schema.DefaultTimeout(pnapDeleteRetryTimeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"path_suffix": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"capacity_in_gb": {
										Type:     schema.TypeInt,
										Required: true,
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceStorageNetworkCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(receiver.BMCSDK)

	request := &networkstorageapiclient.StorageNetworkCreate{}
	request.Name = d.Get("name").(string)
	request.Location = d.Get("location").(string)
	var desc = d.Get("description").(string)
	if len(desc) > 0 {
		request.Description = &desc
	}
	var volumes = d.Get("volumes").([]interface{})
	if len(volumes) > 0 {
		volumesObject := make([]networkstorageapiclient.VolumeCreate, len(volumes))
		for i, j := range volumes {
			volumesItem := j.(map[string]interface{})
			volume := volumesItem["volume"].([]interface{})[0]
			volumeItem := volume.(map[string]interface{})
			volumeObject := networkstorageapiclient.VolumeCreate{}

			volumeObject.Name = volumeItem["name"].(string)
			var volDesc = volumeItem["description"].(string)
			if len(volDesc) > 0 {
				volumeObject.Description = &volDesc
			}
			var pathSuffix = volumeItem["path_suffix"].(string)
			if len(pathSuffix) > 0 {
				volumeObject.PathSuffix = &pathSuffix
			}
			volumeObject.CapacityInGb = int32(volumeItem["capacity_in_gb"].(int))

			volumesObject[i] = volumeObject
		}

		request.Volumes = volumesObject
	}
	requestCommand := storagenetwork.NewCreateStorageNetworkCommand(client, *request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	} else if resp.Id == nil {
		return fmt.Errorf("unknown storage network identifier")
	} else {
		d.SetId(*resp.Id)
	}

	return resourceStorageNetworkRead(d, m)
}

func resourceStorageNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	storageNetworkID := d.Id()
	requestCommand := storagenetwork.NewGetStorageNetworkCommand(client, storageNetworkID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	if resp.Id == nil {
		return fmt.Errorf("unknown storage network identifier")
	}
	d.SetId(*resp.Id)
	if resp.Name != nil {
		d.Set("name", *resp.Name)
	}
	if resp.Description != nil {
		d.Set("description", *resp.Description)
	}
	if resp.Status != nil {
		d.Set("status", *resp.Status)
	}
	if resp.Location != nil {
		d.Set("location", *resp.Location)
	}
	if resp.NetworkId != nil {
		d.Set("network_id", *resp.NetworkId)
	}
	var ips []interface{}
	for _, v := range resp.Ips {
		ips = append(ips, v)
	}
	d.Set("ips", ips)
	if len(resp.CreatedOn.String()) > 0 {
		d.Set("created_on", resp.CreatedOn.String())
	}
	volumes := flattenVolumes(resp.Volumes)

	if err := d.Set("volumes", volumes); err != nil {
		return err
	}
	return nil
}

func resourceStorageNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") || d.HasChange("description") {
		client := m.(receiver.BMCSDK)
		storageNetworkID := d.Id()
		request := &networkstorageapiclient.StorageNetworkUpdate{}
		var name = d.Get("name").(string)
		request.Name = &name
		var desc = d.Get("description").(string)
		request.Description = &desc
		requestCommand := storagenetwork.NewUpdateStorageNetworkCommand(client, storageNetworkID, *request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported action")
	}
	return resourceStorageNetworkRead(d, m)
}

func resourceStorageNetworkDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	storageNetworkID := d.Id()

	requestCommand := storagenetwork.NewDeleteStorageNetworkCommand(client, storageNetworkID)
	err := requestCommand.Execute()
	if err != nil {
		return err
	}

	return nil
}

func flattenVolumes(volumes []networkstorageapiclient.Volume) []interface{} {
	if volumes != nil {
		vols := make([]interface{}, len(volumes))
		for i, v := range volumes {
			volsItem := make(map[string]interface{})
			vol := make([]interface{}, 1)
			volItem := make(map[string]interface{})
			if v.Id != nil {
				volItem["id"] = *v.Id
			}
			if v.Name != nil {
				volItem["name"] = *v.Name
			}
			if v.Description != nil {
				volItem["description"] = *v.Description
			}
			if v.Path != nil {
				volItem["path"] = *v.Path
			}
			if v.PathSuffix != nil {
				volItem["path_suffix"] = *v.PathSuffix
			}
			if v.CapacityInGb != nil {
				volItem["capacity_in_gb"] = int(*v.CapacityInGb)
			}
			if v.Protocol != nil {
				volItem["protocol"] = *v.Protocol
			}
			if v.Status != nil {
				volItem["status"] = *v.Status
			}
			if v.CreatedOn != nil && len(v.CreatedOn.String()) > 0 {
				volItem["created_on"] = v.CreatedOn.String()
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
				volItem["permissions"] = perms
			}
			vol[0] = volItem
			volsItem["volume"] = vol
			vols[i] = volsItem
		}
		return vols
	}
	return make([]interface{}, 0)
}

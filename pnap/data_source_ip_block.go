package pnap

import (
	"fmt"

	"github.com/PNAP/go-sdk-helper-bmc/command/ipapi/ipblock"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/phoenixnap/go-sdk-bmc/ipapi/v3"
)

func dataSourceIpBlock() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIpBlockRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr_block_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"cidr"},
			},
			"cidr": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assigned_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assigned_resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
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
			"is_bring_your_own": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIpBlockRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := ipblock.NewGetIpBlocksCommand(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	numOfBlocks := 0
	for _, instance := range resp {
		if instance.Cidr == d.Get("cidr").(string) || instance.Id == d.Get("id").(string) {
			numOfBlocks++
			d.SetId(instance.Id)
			d.Set("location", instance.Location)
			d.Set("cidr_block_size", instance.CidrBlockSize)
			d.Set("cidr", instance.Cidr)
			d.Set("status", instance.Status)
			if instance.AssignedResourceId != nil {
				d.Set("assigned_resource_id", *instance.AssignedResourceId)
			} else {
				d.Set("assigned_resource_id", "")
			}
			if instance.AssignedResourceType != nil {
				d.Set("assigned_resource_type", *instance.AssignedResourceType)
			} else {
				d.Set("assigned_resource_type", "")
			}
			if instance.Description != nil {
				d.Set("description", *instance.Description)
			} else {
				d.Set("description", "")
			}
			tags := flattenDataTags(instance.Tags)
			if err := d.Set("tags", tags); err != nil {
				return err
			}
			d.Set("is_bring_your_own", instance.IsBringYourOwn)
			if len(instance.CreatedOn.String()) > 0 {
				d.Set("created_on", instance.CreatedOn.String())
			}
		}
	}
	if numOfBlocks > 1 {
		return fmt.Errorf("too many IP Blocks with CIDR %s (found %d, expected 1)", d.Get("cidr").(string), numOfBlocks)
	}

	return nil
}

// Returns list of assigned tags
func flattenDataTags(tags []ipapi.TagAssignment) []interface{} {
	if tags != nil {
		readTags := tags
		tagsMake := make([]interface{}, len(readTags))
		for i, j := range readTags {
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
			tagsMake[i] = tagAssignment
		}
		return tagsMake
	}
	return make([]interface{}, 0)
}

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
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
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
	cidr := d.Get("cidr").(string)
	id := d.Get("id").(string)
	for _, instance := range resp {
		if instance.Cidr != nil && *instance.Cidr == cidr || instance.Id != nil && *instance.Id == id {
			numOfBlocks++
			if instance.Id != nil {
				d.SetId(*instance.Id)
			} else {
				d.SetId("")
			}
			if instance.Location != nil {
				d.Set("location", *instance.Location)
			} else {
				d.Set("location", "")
			}
			if instance.CidrBlockSize != nil {
				d.Set("cidr_block_size", *instance.CidrBlockSize)
			} else {
				d.Set("cidr_block_size", "")
			}
			if instance.Cidr != nil {
				d.Set("cidr", *instance.Cidr)
			} else {
				d.Set("cidr", "")
			}
			if instance.IpVersion != nil {
				d.Set("ip_version", *instance.IpVersion)
			} else {
				d.Set("ip_version", "")
			}
			if instance.Status != nil {
				d.Set("status", *instance.Status)
			} else {
				d.Set("status", "")
			}
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
			if instance.IsBringYourOwn != nil {
				d.Set("is_bring_your_own", *instance.IsBringYourOwn)
			} else {
				d.Set("is_bring_your_own", nil)
			}
			if instance.CreatedOn != nil {
				createdOn := *instance.CreatedOn
				d.Set("created_on", createdOn.String())
			}
		}
	}
	if numOfBlocks > 1 && len(cidr) > 0 {
		return fmt.Errorf("too many IP Blocks with CIDR %s (found %d, expected 1)", cidr, numOfBlocks)
	} else if numOfBlocks > 1 && len(id) > 0 {
		return fmt.Errorf("too many IP Blocks with ID %s (found %d, expected 1)", id, numOfBlocks)
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

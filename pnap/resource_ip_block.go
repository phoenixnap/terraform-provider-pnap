package pnap

import (
	"fmt"
	"log"
	"time"

	"github.com/PNAP/go-sdk-helper-bmc/command/ipapi/ipblock"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ipapiclient "github.com/phoenixnap/go-sdk-bmc/ipapi/v3"
)

const (
	pnapIpBlockRetryDelay   = 15 * time.Second
	pnapIpBlockRetryTimeout = 15 * time.Minute
)

func resourceIpBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpBlockCreate,
		Read:   resourceIpBlockRead,
		Update: resourceIpBlockUpdate,
		Delete: resourceIpBlockDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(pnapRetryTimeout),
			Update: schema.DefaultTimeout(pnapRetryTimeout),
			Delete: schema.DefaultTimeout(pnapDeleteRetryTimeout),
		},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr_block_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_assignment": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
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
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  nil,
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
			"cidr": {
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
			"is_bring_your_own": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIpBlockCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(receiver.BMCSDK)

	request := &ipapiclient.IpBlockCreate{}
	request.Location = d.Get("location").(string)
	request.CidrBlockSize = d.Get("cidr_block_size").(string)
	var desc = d.Get("description").(string)
	if len(desc) > 0 {
		request.Description = &desc
	}
	tags := d.Get("tags").([]interface{})
	if len(tags) > 0 {
		tagsObject := make([]ipapiclient.TagAssignmentRequest, len(tags))
		for i, j := range tags {
			tarObject := ipapiclient.TagAssignmentRequest{}
			tagsItem := j.(map[string]interface{})

			tagAssign := tagsItem["tag_assignment"].([]interface{})[0]
			tagAssignItem := tagAssign.(map[string]interface{})

			tarObject.Name = tagAssignItem["name"].(string)
			value := tagAssignItem["value"].(string)
			if len(value) > 0 {
				tarObject.Value = &value
			}
			tagsObject[i] = tarObject
		}
		request.Tags = tagsObject
	}

	requestCommand := ipblock.NewCreateIpBlockCommand(client, *request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	d.SetId(resp.Id)

	return resourceIpBlockRead(d, m)
}

func resourceIpBlockRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	ipBlockID := d.Id()
	requestCommand := ipblock.NewGetIpBlockCommand(client, ipBlockID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	d.SetId(resp.Id)
	d.Set("location", resp.Location)
	d.Set("cidr_block_size", resp.CidrBlockSize)
	d.Set("cidr", resp.Cidr)
	d.Set("status", resp.Status)
	if resp.AssignedResourceId != nil {
		d.Set("assigned_resource_id", *resp.AssignedResourceId)
	} else {
		d.Set("assigned_resource_id", "")
	}
	if resp.AssignedResourceType != nil {
		d.Set("assigned_resource_type", *resp.AssignedResourceType)
	} else {
		d.Set("assigned_resource_type", "")
	}
	if resp.Description != nil {
		d.Set("description", *resp.Description)
	} else {
		d.Set("description", "")
	}
	if resp.Tags != nil && len(resp.Tags) > 0 {
		var tagsInput = d.Get("tags").([]interface{})
		tags := flattenTags(resp.Tags, tagsInput)
		if err := d.Set("tags", tags); err != nil {
			return err
		}
	}
	d.Set("is_bring_your_own", resp.IsBringYourOwn)
	if len(resp.CreatedOn.String()) > 0 {
		d.Set("created_on", resp.CreatedOn.String())
	}
	return nil
}

func resourceIpBlockUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("description") {
		client := m.(receiver.BMCSDK)
		request := &ipapiclient.IpBlockPatch{}
		var desc = d.Get("description").(string)
		request.Description = &desc

		ipBlockID := d.Id()
		requestCommand := ipblock.NewPatchIpBlockCommand(client, ipBlockID, *request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else if d.HasChange("tags") {
		tags := d.Get("tags").([]interface{})
		client := m.(receiver.BMCSDK)
		ipBlockID := d.Id()

		var request []ipapiclient.TagAssignmentRequest

		if len(tags) > 0 {
			request = make([]ipapiclient.TagAssignmentRequest, len(tags))

			for i, j := range tags {
				tarObject := ipapiclient.TagAssignmentRequest{}
				tagsItem := j.(map[string]interface{})

				tagAssign := tagsItem["tag_assignment"].([]interface{})[0]
				tagAssignItem := tagAssign.(map[string]interface{})

				tarObject.Name = tagAssignItem["name"].(string)
				value := tagAssignItem["value"].(string)
				if len(value) > 0 {
					tarObject.Value = &value
				}
				request[i] = tarObject
			}
		}
		requestCommand := ipblock.NewPutTagsIpBlockCommand(client, ipBlockID, request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported action")
	}

	return resourceIpBlockRead(d, m)
}

func resourceIpBlockDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	ipBlockID := d.Id()

	waitResultError := ipBlockWaitForUnassign(ipBlockID, &client)
	if waitResultError != nil {
		return waitResultError
	}

	requestCommand := ipblock.NewDeleteIpBlockCommand(client, ipBlockID)
	_, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	return nil
}

func flattenTags(tagsRead []ipapiclient.TagAssignment, tagsInput []interface{}) []interface{} {
	if len(tagsInput) == 0 {
		tagsInput = make([]interface{}, 1)
		tagsInputItem := make(map[string]interface{})
		tagsInput[0] = tagsInputItem
	}
	if len(tagsInput) > 0 {
		tags := tagsRead
		for _, j := range tagsInput {
			tagsInputItem := j.(map[string]interface{})
			if tagsInputItem["tag_assignment"] != nil && len(tagsInputItem["tag_assignment"].([]interface{})) > 0 {
				tagAssign := tagsInputItem["tag_assignment"].([]interface{})[0]
				tagAssignItem := tagAssign.(map[string]interface{})
				nameInput := tagAssignItem["name"].(string)
				for _, l := range tags {
					if nameInput == l.Name {
						tagAssignItem["id"] = l.Id
						tagAssignItem["value"] = l.Value
						tagAssignItem["is_billing_tag"] = l.IsBillingTag
						tagAssignItem["created_by"] = l.CreatedBy
					}
				}
			}
		}
	}
	return tagsInput
}

func ipBlockWaitForUnassign(id string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for ip block %s to be unassigned...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"unassigning", "assigning"},
		Target:     []string{"unassigned", "assigned"},
		Refresh:    refreshForIpBlockStatus(client, id),
		Timeout:    pnapIpBlockRetryTimeout,
		Delay:      pnapIpBlockRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for ip block (%s) to be unassigned: %v", id, err)
	}

	return nil
}

func refreshForIpBlockStatus(client *receiver.BMCSDK, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		requestCommand := ipblock.NewGetIpBlockCommand(*client, id)

		resp, err := requestCommand.Execute()
		if err != nil {
			return 0, "", err
		} else {
			return 0, resp.Status, nil
		}
	}
}

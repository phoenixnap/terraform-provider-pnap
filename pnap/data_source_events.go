package pnap

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PNAP/go-sdk-helper-bmc/command/auditapi/event"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEvents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEventsRead,

		Schema: map[string]*schema.Schema{
			"events": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"username": {
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

func dataSourceEventsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	requestCommand := event.NewGetEventsCommand(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	qEvents := d.Get("events").([]interface{})
	if len(qEvents) != 1 {
		return fmt.Errorf("unsupported action")
	}
	qEvent := d.Get("events").([]interface{})[0]
	qEventItem := qEvent.(map[string]interface{})
	qName := qEventItem["name"].(string)

	var events []interface{}
	for _, instance := range resp {
		if instance.Name != nil {
			name := *instance.Name
			if name == qName {
				event := make(map[string]interface{})
				event["name"] = name
				event["timestamp"] = instance.Timestamp.String()

				userInfo := make([]interface{}, 1)
				userInfoItem := make(map[string]interface{})

				userInfoItem["account_id"] = instance.UserInfo.AccountId
				if instance.UserInfo.ClientId != nil {
					userInfoItem["client_id"] = *instance.UserInfo.ClientId
				}
				userInfoItem["username"] = instance.UserInfo.Username

				userInfo[0] = userInfoItem
				event["user_info"] = userInfo
				events = append(events, event)
			}
		}
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("events", events)
	return nil
}

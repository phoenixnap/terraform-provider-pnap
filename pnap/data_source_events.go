package pnap

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PNAP/go-sdk-helper-bmc/command/auditapi/event"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEvents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEventsRead,

		Schema: map[string]*schema.Schema{
			"from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verb": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"events": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
	query := dto.Query{}

	from := d.Get("from").(string)
	if from != "" {
		t1, err1 := time.Parse(time.RFC3339, from)
		if err1 != nil {
			return err1
		} else {
			query.From = t1
		}
	}
	to := d.Get("to").(string)
	if to != "" {
		t2, err2 := time.Parse(time.RFC3339, to)
		if err2 != nil {
			return err2
		} else {
			query.To = t2
		}
	}
	query.Limit = int32(d.Get("limit").(int))
	query.Order = d.Get("order").(string)
	query.Username = d.Get("username").(string)
	query.Verb = d.Get("verb").(string)
	query.Uri = d.Get("uri").(string)

	requestCommand := event.NewGetEventsCommandWithQuery(client, &query)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	qEvents := d.Get("events").([]interface{})
	var events []interface{}

	if len(qEvents) > 0 {
		if len(qEvents) != 1 {
			return fmt.Errorf("unsupported action")
		}
		qEvent := qEvents[0]
		qEventItem := qEvent.(map[string]interface{})
		qName := qEventItem["name"].(string)

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
	} else {
		events = make([]interface{}, len(resp))

		for num, instance := range resp {
			event := make(map[string]interface{})
			if instance.Name != nil {
				event["name"] = *instance.Name
			}
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
			events[num] = event
		}
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("events", events)
	return nil
}

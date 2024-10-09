package pnap

import (
	"fmt"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/bgppeergroup"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBgpPeerGroup() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceBgpPeerGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"location"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"ipv4_prefixes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_allocation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_bring_your_own_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"in_use": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"target_asn_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"is_bring_your_own": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"verification_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"verification_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"active_asn_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"is_bring_your_own": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"verification_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"verification_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"advertised_routes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rpki_roa_origin_asn": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ebgp_multi_hop": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"peering_loopbacks_v4": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"keep_alive_timer_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hold_timer_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBgpPeerGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	bgpID := d.Get("id").(string)
	if len(bgpID) > 0 {
		requestCommand := bgppeergroup.NewGetBgpPeerGroupsCommand(client)
		resp, err := requestCommand.Execute()
		if err != nil {
			return err
		}
		numOfGroups := 0
		for _, instance := range resp {
			if instance.Id == bgpID {
				numOfGroups++
				d.SetId(instance.Id)
				d.Set("status", instance.Status)
				d.Set("location", instance.Location)

				ipv4Prefixes := flattenIpv4Prefixes(instance.Ipv4Prefixes)
				if err := d.Set("ipv4_prefixes", ipv4Prefixes); err != nil {
					return err
				}
				target := instance.TargetAsnDetails
				targetAsnDetails := flattenAsnDetails(&target)
				if err := d.Set("target_asn_details", targetAsnDetails); err != nil {
					return err
				}
				activeAsnDetails := flattenAsnDetails(instance.ActiveAsnDetails)
				if err := d.Set("active_asn_details", activeAsnDetails); err != nil {
					return err
				}
				d.Set("password", instance.Password)
				d.Set("advertised_routes", instance.AdvertisedRoutes)
				d.Set("rpki_roa_origin_asn", int(instance.RpkiRoaOriginAsn))
				d.Set("ebgp_multi_hop", int(instance.EBgpMultiHop))
				var peeringLoopbacks []interface{}
				for _, v := range instance.PeeringLoopbacksV4 {
					peeringLoopbacks = append(peeringLoopbacks, v)
				}
				d.Set("peering_loopbacks_v4", peeringLoopbacks)
				d.Set("keep_alive_timer_seconds", int(instance.KeepAliveTimerSeconds))
				d.Set("hold_timer_seconds", int(instance.HoldTimerSeconds))

				if instance.CreatedOn != nil {
					createdOn := *instance.CreatedOn
					d.Set("created_on", createdOn)
				}
				if instance.LastUpdatedOn != nil {
					lastUpdatedOn := *instance.LastUpdatedOn
					d.Set("last_updated_on", lastUpdatedOn)
				}
			}
		}
		if numOfGroups > 1 {
			return fmt.Errorf("too many BGP Peer Groups with id %s (found %d, expected 1)", d.Get("id").(string), numOfGroups)
		}
		return nil
	} else {
		query := dto.Query{}
		location := d.Get("location").(string)
		query.LocationString = location
		requestCommand := bgppeergroup.NewGetBgpPeerGroupsWithQueryCommand(client, &query)
		resp, err := requestCommand.Execute()
		if err != nil {
			return err
		}
		numOfGroups := 0
		for _, instance := range resp {
			numOfGroups++
			d.SetId(instance.Id)
			d.Set("status", instance.Status)
			d.Set("location", instance.Location)

			ipv4Prefixes := flattenIpv4Prefixes(instance.Ipv4Prefixes)
			if err := d.Set("ipv4_prefixes", ipv4Prefixes); err != nil {
				return err
			}
			target := instance.TargetAsnDetails
			targetAsnDetails := flattenAsnDetails(&target)
			if err := d.Set("target_asn_details", targetAsnDetails); err != nil {
				return err
			}
			activeAsnDetails := flattenAsnDetails(instance.ActiveAsnDetails)
			if err := d.Set("active_asn_details", activeAsnDetails); err != nil {
				return err
			}
			d.Set("password", instance.Password)
			d.Set("advertised_routes", instance.AdvertisedRoutes)
			d.Set("rpki_roa_origin_asn", int(instance.RpkiRoaOriginAsn))
			d.Set("ebgp_multi_hop", int(instance.EBgpMultiHop))
			var peeringLoopbacks []interface{}
			for _, v := range instance.PeeringLoopbacksV4 {
				peeringLoopbacks = append(peeringLoopbacks, v)
			}
			d.Set("peering_loopbacks_v4", peeringLoopbacks)
			d.Set("keep_alive_timer_seconds", int(instance.KeepAliveTimerSeconds))
			d.Set("hold_timer_seconds", int(instance.HoldTimerSeconds))

			if instance.CreatedOn != nil {
				createdOn := *instance.CreatedOn
				d.Set("created_on", createdOn)
			}
			if instance.LastUpdatedOn != nil {
				lastUpdatedOn := *instance.LastUpdatedOn
				d.Set("last_updated_on", lastUpdatedOn)
			}
		}
		if numOfGroups > 1 {
			return fmt.Errorf("too many BGP Peer Groups with location %s (found %d, expected 1)", d.Get("location").(string), numOfGroups)
		}
		return nil
	}
}

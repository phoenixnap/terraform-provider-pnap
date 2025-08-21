package pnap

import (
	"fmt"

	"github.com/PNAP/go-sdk-helper-bmc/command/networkapi/bgppeergroup"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	networkapiclient "github.com/phoenixnap/go-sdk-bmc/networkapi/v4"
)

func resourceBgpPeerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceBgpPeerGroupCreate,
		Read:   resourceBgpPeerGroupRead,
		Update: resourceBgpPeerGroupUpdate,
		Delete: resourceBgpPeerGroupDelete,

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
			"asn": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"advertised_routes": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_prefixes": { // Deprecated
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
			"ip_prefixes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_allocation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
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
			"peering_loopbacks_v6": {
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

func resourceBgpPeerGroupCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(receiver.BMCSDK)

	request := &networkapiclient.BgpPeerGroupCreate{}
	request.Location = d.Get("location").(string)

	var asn = d.Get("asn").(int)
	request.Asn = int64(asn)

	var password = d.Get("password").(string)
	if len(password) > 0 {
		request.Password = &password
	}

	request.AdvertisedRoutes = d.Get("advertised_routes").(string)

	requestCommand := bgppeergroup.NewCreateBgpPeerGroupCommand(client, *request)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	d.SetId(resp.Id)

	return resourceBgpPeerGroupRead(d, m)
}

func resourceBgpPeerGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	bgpID := d.Id()
	requestCommand := bgppeergroup.NewGetBgpPeerGroupCommand(client, bgpID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	d.SetId(resp.Id)
	d.Set("status", resp.Status)
	d.Set("location", resp.Location)

	ipv4Prefixes := flattenIpv4Prefixes(resp.Ipv4Prefixes)
	if err := d.Set("ipv4_prefixes", ipv4Prefixes); err != nil {
		return err
	}
	ipPrefixes := flattenIpPrefixes(resp.IpPrefixes)
	if err := d.Set("ip_prefixes", ipPrefixes); err != nil {
		return err
	}
	target := resp.TargetAsnDetails
	targetAsnDetails := flattenAsnDetails(&target)
	if err := d.Set("target_asn_details", targetAsnDetails); err != nil {
		return err
	}
	activeAsnDetails := flattenAsnDetails(resp.ActiveAsnDetails)
	if err := d.Set("active_asn_details", activeAsnDetails); err != nil {
		return err
	}
	d.Set("password", resp.Password)
	d.Set("advertised_routes", resp.AdvertisedRoutes)
	d.Set("rpki_roa_origin_asn", int(resp.RpkiRoaOriginAsn))
	d.Set("ebgp_multi_hop", int(resp.EBgpMultiHop))
	var peeringLoopbacks []interface{}
	for _, v := range resp.PeeringLoopbacksV4 {
		peeringLoopbacks = append(peeringLoopbacks, v)
	}
	d.Set("peering_loopbacks_v4", peeringLoopbacks)
	var peeringLoopbacks6 []interface{}
	for _, v6 := range resp.PeeringLoopbacksV6 {
		peeringLoopbacks6 = append(peeringLoopbacks6, v6)
	}
	d.Set("peering_loopbacks_v6", peeringLoopbacks6)
	d.Set("keep_alive_timer_seconds", int(resp.KeepAliveTimerSeconds))
	d.Set("hold_timer_seconds", int(resp.HoldTimerSeconds))

	if resp.CreatedOn != nil {
		createdOn := *resp.CreatedOn
		d.Set("created_on", createdOn)
	}
	if resp.LastUpdatedOn != nil {
		lastUpdatedOn := *resp.LastUpdatedOn
		d.Set("last_updated_on", lastUpdatedOn)
	}

	return nil
}

func resourceBgpPeerGroupUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("asn") || d.HasChange("password") || d.HasChange("advertised_routes") {
		client := m.(receiver.BMCSDK)
		request := &networkapiclient.BgpPeerGroupPatch{}

		if d.HasChange("asn") {
			asn := d.Get("asn").(int)
			asn64 := int64(asn)
			request.Asn = &asn64
		}
		if d.HasChange("password") {
			password := d.Get("password").(string)
			request.Password = &password
		}
		if d.HasChange("advertised_routes") {
			routes := d.Get("advertised_routes").(string)
			request.AdvertisedRoutes = &routes
		}

		requestCommand := bgppeergroup.NewUpdateBgpPeerGroupCommand(client, d.Id(), *request)

		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("unsupported action")
	}
	return resourceBgpPeerGroupRead(d, m)

}

func resourceBgpPeerGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)

	bgpID := d.Id()

	requestCommand := bgppeergroup.NewDeleteBgpPeerGroupCommand(client, bgpID)
	_, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	return nil
}

func flattenIpv4Prefixes(ipv4Prefixes []networkapiclient.BgpIPv4Prefix) []interface{} {
	if ipv4Prefixes != nil {
		ss := make([]interface{}, len(ipv4Prefixes))
		for i, v := range ipv4Prefixes {
			s := make(map[string]interface{})
			s["ipv4_allocation_id"] = v.Ipv4AllocationId
			s["cidr"] = v.Cidr
			s["status"] = v.Status
			s["is_bring_your_own_ip"] = v.IsBringYourOwnIp
			s["in_use"] = v.InUse

			ss[i] = s
		}
		return ss
	}
	return make([]interface{}, 0)
}

func flattenIpPrefixes(ipPrefixes []networkapiclient.BgpIpPrefix) []interface{} {
	if ipPrefixes != nil {
		ss := make([]interface{}, len(ipPrefixes))
		for i, v := range ipPrefixes {
			s := make(map[string]interface{})
			s["ip_allocation_id"] = v.IpAllocationId
			s["cidr"] = v.Cidr
			s["ip_version"] = v.IpVersion
			s["status"] = v.Status

			ss[i] = s
		}
		return ss
	}
	return make([]interface{}, 0)
}

func flattenAsnDetails(AsnDetails *networkapiclient.AsnDetails) []interface{} {
	if AsnDetails != nil {
		ss := make([]interface{}, 1)
		s := make(map[string]interface{})
		s["asn"] = int(AsnDetails.Asn)
		s["is_bring_your_own"] = AsnDetails.IsBringYourOwn
		s["verification_status"] = AsnDetails.VerificationStatus
		if AsnDetails.VerificationReason != nil {
			s["verification_reason"] = *AsnDetails.VerificationReason
		}

		ss[0] = s
		return ss
	}
	return make([]interface{}, 0)
}

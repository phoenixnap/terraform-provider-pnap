---
layout: "pnap"
page_title: "phoenixNAP: pnap_bgp_peer_group"
sidebar_current: "docs-pnap-datasource-bgp-peer-group"
description: |-
  Provides a phoenixNAP BGP Peer Group datasource. This can be used to read BGP Peer Groups.
---

# pnap_bgp_peer_group Datasource

Provides a phoenixNAP BGP Peer Group datasource. This can be used to read BGP Peer Groups.



## Example Usage

Fetch a BGP Peer Group data by location and show it's IPv4 Peering Loopback addresses.

```hcl
# Fetch a BGP Peer Group
data "pnap_bgp_peer_group" "BGP-Peer-Group-1" {
    location = "PHX"
}

# Show IPv4 prefixes
output "BGPPeerGroup1" {
    value = data.pnap_bgp_peer_group.BGP-Peer-Group-1.peering_loopbacks_v4
}
```

## Argument Reference

The following arguments are supported:

* `location` - The BGP Peer Group location. Supported values are `PHX`, `ASH`, `SGP`, `NLD`, `CHI`, `SEA` and `AUS`.
* `id` - The unique identifier of the BGP Peer Group.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the BGP Peer Group.
* `status` - The BGP Peer Group status.
* `location` - The BGP Peer Group location.
* `ipv4_prefixes` - The list of BGP Peer Group IPv4 prefixes.
    * `ipv4_allocation_id` - IPv4 allocation ID.
    * `cidr` - The IP block in CIDR format.
    * `status`- The BGP IPv4 Prefix status.
    * `is_bring_your_own_ip` - Identifies IP as a "bring your own" IP block.
    * `in_use` - The boolean value of the BGP IPv4 Prefix is in use.
* `target_asn_details ` - BGP Peer Group ASN details.
    * `asn` - The BGP Peer Group ASN.
    * `is_bring_your_own` - True if the BGP Peer Group ASN is a "bring your own" ASN.
    * `verification_status` - The BGP Peer Group ASN verification status.
    * `verification_reason` - The BGP Peer Group ASN verification reason for the respective status.
* `active_asn_details ` - BGP Peer Group ASN details.
    * `asn` - The BGP Peer Group ASN.
    * `is_bring_your_own` - True if the BGP Peer Group ASN is a "bring your own" ASN.
    * `verification_status` - The BGP Peer Group ASN verification status.
    * `verification_reason` - The BGP Peer Group ASN verification reason for the respective status.
* `password`- The BGP Peer Group password.
* `advertised_routes` - The Advertised routes for the BGP Peer Group.
* `rpki_roa_origin_asn` - The RPKI ROA Origin ASN of the BGP Peer Group based on location.
* `ebgp_multi_hop` - The eBGP Multi-hop of the BGP Peer Group.
* `peering_loopbacks_v4` - The IPv4 Peering Loopback addresses of the BGP Peer Group. Valid IP formats are IPv4 addresses.
* `keep_alive_timer_seconds` - The Keep Alive Timer in seconds, of the BGP Peer Group.
* `hold_timer_seconds` - The Hold Timer in seconds, of the BGP Peer Group.
* `created_on` - Date and time of creation.
* `last_updated_on` - Date and time of last update.
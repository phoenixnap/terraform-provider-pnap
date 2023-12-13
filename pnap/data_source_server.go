package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/phoenixnap/go-sdk-bmc/bmcapi/v2"

	"github.com/PNAP/go-sdk-helper-bmc/command/bmcapi/server"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
)

func dataSourceServer() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceServerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"hostname"},
			},
			"hostname": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip_addresses": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_ip_addresses": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"netris_controller": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_os": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"netris_softgate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_os": {
							Type:     schema.TypeString,
							Computed: true,
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
			"network_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_network_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"configuration_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_networks": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ips": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"dhcp": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"status_description": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"ip_blocks_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"configuration_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_blocks": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"vlan_id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"public_network_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_networks": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ips": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"status_description": {
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
				},
			},
			"storage_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"root_partition": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"raid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"superseded_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supersedes": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	//serverID := d.Id()
	requestCommand := server.NewGetServersCommand(client)
	//requestCommand.SetRequester(client)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}
	/* code := resp.StatusCode
	if code != 200 {
		response := &dto.ErrorMessage{}
		response.FromBytes(resp)
		return fmt.Errorf("API Returned Code: %v, Message: %v, Validation Errors: %v", code, response.Message, response.ValidationErrors)
	}
	response := &dto.Servers{}
	response.FromBytes(resp) */

	numOfServers := 0
	for _, instance := range resp {
		if instance.Hostname == d.Get("hostname").(string) || instance.Id == d.Get("id").(string) {
			numOfServers++
			d.SetId(instance.Id)
			d.Set("status", instance.Status)
			d.Set("hostname", instance.Hostname)

			d.Set("os", instance.Os)
			d.Set("type", instance.Type)
			d.Set("location", instance.Location)
			d.Set("network_type", instance.NetworkType)

			var privateIPs []interface{}
			for _, v := range instance.PrivateIpAddresses {
				privateIPs = append(privateIPs, v)
			}
			d.Set("private_ip_addresses", privateIPs)
			var publicIPs []interface{}
			for _, k := range instance.PublicIpAddresses {
				publicIPs = append(publicIPs, k)
			}
			d.Set("public_ip_addresses", publicIPs)
			if len(instance.PublicIpAddresses) > 0 {
				d.Set("primary_ip_address", instance.PublicIpAddresses[0])
			}

			if instance.OsConfiguration != nil {
				if instance.OsConfiguration.NetrisController != nil {
					netrisController := make([]interface{}, 1)
					netrisControllerItem := make(map[string]interface{})
					if instance.OsConfiguration.NetrisController.HostOs != nil {
						netrisControllerItem["host_os"] = *instance.OsConfiguration.NetrisController.HostOs
					}
					netrisController[0] = netrisControllerItem
					d.Set("netris_controller", netrisController)
				}
				if instance.OsConfiguration.NetrisSoftgate != nil {
					netrisSoftgate := make([]interface{}, 1)
					netrisSoftgateItem := make(map[string]interface{})
					if instance.OsConfiguration.NetrisSoftgate.HostOs != nil {
						netrisSoftgateItem["host_os"] = *instance.OsConfiguration.NetrisSoftgate.HostOs
					}
					netrisSoftgate[0] = netrisSoftgateItem
					d.Set("netris_softgate", netrisSoftgate)
				}
			}

			tags := flattenServerDataTags(instance.Tags)
			if err := d.Set("tags", tags); err != nil {
				return err
			}
			netConf := flattenServerDataNetworkConfiguration(instance.NetworkConfiguration)
			if err := d.Set("network_configuration", netConf); err != nil {
				return err
			}
			if instance.StorageConfiguration.RootPartition != nil {
				storageConfiguration := make([]interface{}, 1)
				storageConfigurationItem := make(map[string]interface{})
				rootPartition := make([]interface{}, 1)
				rootPartitionItem := make(map[string]interface{})
				if instance.StorageConfiguration.RootPartition.Raid != nil {
					rootPartitionItem["raid"] = *instance.StorageConfiguration.RootPartition.Raid
				}
				if instance.StorageConfiguration.RootPartition.Size != nil {
					rootPartitionItem["size"] = int(*instance.StorageConfiguration.RootPartition.Size)
				}
				rootPartition[0] = rootPartitionItem
				storageConfigurationItem["root_partition"] = rootPartition
				storageConfiguration[0] = storageConfigurationItem
				d.Set("storage_configuration", storageConfiguration)
			}

			d.Set("superseded_by", instance.SupersededBy)
			d.Set("supersedes", instance.Supersedes)
		}
	}

	if numOfServers > 1 {
		return fmt.Errorf("too many devices found with hostname %s (found %d, expected 1)", d.Get("hostname").(string), numOfServers)
	}

	return nil
}

// Returns list of assigned tags
func flattenServerDataTags(tags []bmcapi.TagAssignment) []interface{} {
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

// Returns entire network details
func flattenServerDataNetworkConfiguration(networkConfiguration bmcapi.NetworkConfiguration) []interface{} {
	netConf := make([]interface{}, 1)
	nc := make(map[string]interface{})

	if networkConfiguration.GatewayAddress != nil {
		nc["gateway_address"] = *networkConfiguration.GatewayAddress
	}
	if networkConfiguration.PrivateNetworkConfiguration != nil {
		privateNetworkConfiguration := *networkConfiguration.PrivateNetworkConfiguration
		prNetConf := make([]interface{}, 1)
		prnc := make(map[string]interface{})
		if privateNetworkConfiguration.ConfigurationType != nil {
			prnc["configuration_type"] = *privateNetworkConfiguration.ConfigurationType
		}
		if privateNetworkConfiguration.PrivateNetworks != nil {
			privateNetworks := privateNetworkConfiguration.PrivateNetworks
			prNet := make([]interface{}, len(privateNetworks))
			for i, j := range privateNetworks {
				spn := make(map[string]interface{})
				spn["id"] = j.Id
				if j.Ips != nil {
					ips := make([]interface{}, len(j.Ips))
					for k, l := range j.Ips {
						ips[k] = l
					}
					spn["ips"] = ips
				}
				if j.Dhcp != nil {
					spn["dhcp"] = *j.Dhcp
				}
				if j.StatusDescription != nil {
					spn["status_description"] = *j.StatusDescription
				}
				prNet[i] = spn
			}
			prnc["private_networks"] = prNet
		}
		prNetConf[0] = prnc
		nc["private_network_configuration"] = prNetConf
	}
	if networkConfiguration.IpBlocksConfiguration != nil {
		ipBlocksConfiguration := *networkConfiguration.IpBlocksConfiguration
		ipBlocksConf := make([]interface{}, 1)
		ibc := make(map[string]interface{})
		if ipBlocksConfiguration.ConfigurationType != nil {
			ibc["configuration_type"] = *ipBlocksConfiguration.ConfigurationType
		}
		if ipBlocksConfiguration.IpBlocks != nil {
			ipBlocks := ipBlocksConfiguration.IpBlocks
			ib := make([]interface{}, len(ipBlocks))
			for i, j := range ipBlocks {
				sib := make(map[string]interface{})
				sib["id"] = j.Id
				if j.VlanId != nil {
					sib["vlan_id"] = int(*j.VlanId)
				}
				ib[i] = sib
			}
			ibc["ip_blocks"] = ib
		}
		ipBlocksConf[0] = ibc
		nc["ip_blocks_configuration"] = ipBlocksConf
	}
	if networkConfiguration.PublicNetworkConfiguration != nil {
		publicNetworkConfiguration := *networkConfiguration.PublicNetworkConfiguration
		puNetConf := make([]interface{}, 1)
		punc := make(map[string]interface{})
		if publicNetworkConfiguration.PublicNetworks != nil {
			publicNetworks := publicNetworkConfiguration.PublicNetworks
			puNet := make([]interface{}, len(publicNetworks))
			for i, j := range publicNetworks {
				spn := make(map[string]interface{})
				spn["id"] = j.Id
				ips := make([]interface{}, len(j.Ips))
				for k, l := range j.Ips {
					ips[k] = l
				}
				spn["ips"] = ips
				if j.StatusDescription != nil {
					spn["status_description"] = *j.StatusDescription
				}
				puNet[i] = spn
			}
			punc["public_networks"] = puNet
		}
		puNetConf[0] = punc
		nc["public_network_configuration"] = puNetConf
	}
	netConf[0] = nc
	return netConf
}

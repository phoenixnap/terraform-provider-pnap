package pnap

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net/netip"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/command/bmcapi/server"
	"github.com/PNAP/go-sdk-helper-bmc/dto"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"

	bmcapiclient "github.com/phoenixnap/go-sdk-bmc/bmcapi/v3"
)

const (
	pnapRetryTimeout       = 100 * time.Minute
	pnapDeleteRetryTimeout = 15 * time.Minute
	pnapRetryDelay         = 5 * time.Second
	pnapRetryMinTimeout    = 3 * time.Second
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(pnapRetryTimeout),
			Update: schema.DefaultTimeout(pnapRetryTimeout),
			Delete: schema.DefaultTimeout(pnapDeleteRetryTimeout),
		},

		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_keys": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cpu": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cores_per_cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_frequency_in_ghz": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_type": {
				Type:                  schema.TypeString,
				Optional:              true,
				Computed:              true,
				DiffSuppressFunc:      supressUserDefinedNetworkType,
				DiffSuppressOnRefresh: true,
			},
			"install_default_ssh_keys": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ssh_key_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"reservation_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pricing_model": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rdp_allowed_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bring_your_own_license": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"management_ui_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_password": {
				Type:     schema.TypeString,
				Computed: true,
				//Sensitive: true,
			},
			"management_access_allowed_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"install_os_to_ram": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"esxi": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datastore_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datastore_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"cloud_init": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_data": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ipxe": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"native_vlan_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vlan_id": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"static_dhcp_address_v4": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
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
						"netris_web_console_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"netris_user_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"netris_softgate": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_os": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"controller_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"controller_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"controller_auth_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"provisioned_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_ip_blocks": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"transfer_reservation_to": {
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
			"network_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_address": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"private_network_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_address": { //Deprecated
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"configuration_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"private_networks": {
										Type:       schema.TypeList,
										Optional:   true,
										Computed:   true,
										ConfigMode: schema.SchemaConfigModeAttr,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_private_network": {
													Type:       schema.TypeList,
													Required:   true,
													MaxItems:   1,
													ConfigMode: schema.SchemaConfigModeAttr,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:     schema.TypeString,
																Required: true,
															},
															"ips": {
																Type:     schema.TypeSet,
																Optional: true,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
															"dhcp": {
																Type:     schema.TypeBool,
																Optional: true,
																Computed: true,
																Default:  nil,
															},
															"status_description": {
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
								},
							},
						},
						"ip_blocks_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"configuration_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ip_blocks": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_ip_block": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:     schema.TypeString,
																Required: true,
															},
															"vlan_id": {
																Type:     schema.TypeInt,
																Optional: true,
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
						"public_network_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_networks": {
										Type:       schema.TypeList,
										Optional:   true,
										Computed:   true,
										ConfigMode: schema.SchemaConfigModeAttr,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_public_network": {
													Type:       schema.TypeList,
													Required:   true,
													MaxItems:   1,
													ConfigMode: schema.SchemaConfigModeAttr,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:     schema.TypeString,
																Required: true,
															},
															"ips": {
																Type:     schema.TypeSet,
																Required: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
															"status_description": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"compute_slaac_ip": {
																Type:     schema.TypeBool,
																Optional: true,
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
								},
							},
						},
					},
				},
			},
			"storage_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"root_partition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"raid": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "NO_RAID",
									},
									"size": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  -1,
									},
								},
							},
						},
					},
				},
			},
			"gpu_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"long_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
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
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			isNullPrivIps, isNullDhcp := markNotSetPrivateNetworkFields(diff)
			isNullComputeSlaacIp := markNotSetPublicNetworkComputeSlaacIpForCustomize(diff)
			oldRawNCint, rawNCint := diff.GetChange("network_configuration")
			var oldIds []string
			var oldPubIds []string
			if oldRawNCint != nil && len(oldRawNCint.([]interface{})) > 0 {
				networkConfiguration := oldRawNCint.([]interface{})[0]
				networkConfigurationItem := networkConfiguration.(map[string]interface{})
				if networkConfigurationItem["private_network_configuration"] != nil && len(networkConfigurationItem["private_network_configuration"].([]interface{})) > 0 {
					privateNetworkConfiguration := networkConfigurationItem["private_network_configuration"].([]interface{})[0]
					privateNetworkConfigurationItem := privateNetworkConfiguration.(map[string]interface{})
					oldRawPN := privateNetworkConfigurationItem["private_networks"].([]interface{})
					oldIds = make([]string, len(oldRawPN))
					for i := range oldRawPN {
						if oldRawPN[i] != nil && len(oldRawPN[i].(map[string]interface{})["server_private_network"].([]interface{})) > 0 {
							ispnItem := oldRawPN[i].(map[string]interface{})["server_private_network"].([]interface{})[0].(map[string]interface{})
							id := ispnItem["id"].(string)
							oldIds[i] = id
						}
					}
				}
				if networkConfigurationItem["public_network_configuration"] != nil && len(networkConfigurationItem["public_network_configuration"].([]interface{})) > 0 {
					publicNetworkConfiguration := networkConfigurationItem["public_network_configuration"].([]interface{})[0]
					publicNetworkConfigurationItem := publicNetworkConfiguration.(map[string]interface{})
					oldRawPN := publicNetworkConfigurationItem["public_networks"].([]interface{})
					oldPubIds = make([]string, len(oldRawPN))
					for i := range oldRawPN {
						if oldRawPN[i] != nil && len(oldRawPN[i].(map[string]interface{})["server_public_network"].([]interface{})) > 0 {
							ispnItem := oldRawPN[i].(map[string]interface{})["server_public_network"].([]interface{})[0].(map[string]interface{})
							id := ispnItem["id"].(string)
							oldPubIds[i] = id
						}
					}
				}
			}
			if rawNCint != nil && len(rawNCint.([]interface{})) > 0 {
				networkConfiguration := rawNCint.([]interface{})[0]
				networkConfigurationItem := networkConfiguration.(map[string]interface{})
				if networkConfigurationItem["private_network_configuration"] != nil && len(networkConfigurationItem["private_network_configuration"].([]interface{})) > 0 {
					privateNetworkConfiguration := networkConfigurationItem["private_network_configuration"].([]interface{})[0]
					privateNetworkConfigurationItem := privateNetworkConfiguration.(map[string]interface{})
					rawPN := privateNetworkConfigurationItem["private_networks"].([]interface{})
					for i := range rawPN {
						if rawPN[i] != nil && len(rawPN[i].(map[string]interface{})["server_private_network"].([]interface{})) > 0 {
							ispnItem := rawPN[i].(map[string]interface{})["server_private_network"].([]interface{})[0].(map[string]interface{})
							id := ispnItem["id"].(string)
							ipsInput := ispnItem["ips"].(*schema.Set).List()
							// Customizing value [""] for empty array of ips to be replaced with value [] because of block to attribute change in provider
							if len(ipsInput) == 1 && ipsInput[0] == "" {
								ispnItem["ips"] = make([]interface{}, 0)
							}
							// Correcting Terraform's index-based prior state fallback
							if len(rawPN) == len(isNullPrivIps) {
								if i < len(oldIds) && oldIds[i] == id {
									// Do nothing
								} else if isNullPrivIps[i] {
									ispnItem["ips"] = make([]interface{}, 0)
								}
							}
							if len(rawPN) == len(isNullDhcp) {
								if i < len(oldIds) && oldIds[i] == id {
									// Do nothing
								} else if isNullDhcp[i] {
									ispnItem["dhcp"] = false
								}
							}
						}
					}
				}
				if networkConfigurationItem["public_network_configuration"] != nil && len(networkConfigurationItem["public_network_configuration"].([]interface{})) > 0 {
					publicNetworkConfiguration := networkConfigurationItem["public_network_configuration"].([]interface{})[0]
					publicNetworkConfigurationItem := publicNetworkConfiguration.(map[string]interface{})
					rawPN := publicNetworkConfigurationItem["public_networks"].([]interface{})
					for i := range rawPN {
						if rawPN[i] != nil && len(rawPN[i].(map[string]interface{})["server_public_network"].([]interface{})) > 0 {
							ispnItem := rawPN[i].(map[string]interface{})["server_public_network"].([]interface{})[0].(map[string]interface{})
							id := ispnItem["id"].(string)
							ipsInput := ispnItem["ips"].(*schema.Set).List()
							// Customizing value [""] for empty array of ips to be replaced with value [] because of block to attribute change in provider
							if len(ipsInput) == 1 && ipsInput[0] == "" {
								ispnItem["ips"] = make([]interface{}, 0)
							}
							// Correcting Terraform's index-based prior state fallback
							if len(rawPN) == len(isNullComputeSlaacIp) {
								if i < len(oldPubIds) && oldPubIds[i] == id {
									// Do nothing
								} else if isNullComputeSlaacIp[i] {
									ispnItem["compute_slaac_ip"] = false
								}
							}
						}
					}
				}
				if err := diff.SetNew("network_configuration", rawNCint); err != nil {
					return fmt.Errorf("failed to override name: %w", err)
				}
			}
			return nil
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(receiver.BMCSDK)

	request := &bmcapiclient.ServerCreate{}
	request.Hostname = d.Get("hostname").(string)
	var desc = d.Get("description").(string)
	if len(desc) > 0 {
		request.Description = &desc
	}
	request.Os = d.Get("os").(string)
	request.Type = d.Get("type").(string)
	request.Location = d.Get("location").(string)
	var networkType = d.Get("network_type").(string)

	if len(networkType) > 0 {
		request.NetworkType = &networkType
	}

	var resId = d.Get("reservation_id").(string)
	if len(resId) > 0 {
		request.ReservationId = &resId
	}

	var prModel = d.Get("pricing_model").(string)
	if len(prModel) > 0 {
		request.PricingModel = &prModel
	}

	var installDefault = d.Get("install_default_ssh_keys").(bool)
	request.InstallDefaultSshKeys = &installDefault
	temp := d.Get("ssh_keys").(*schema.Set).List()
	keys := make([]string, len(temp))
	for i, v := range temp {
		keys[i] = fmt.Sprint(v)
	}
	//todo
	request.SshKeys = keys

	temp1 := d.Get("ssh_key_ids").(*schema.Set).List()
	keyIds := make([]string, len(temp1))
	for i, v := range temp1 {
		keyIds[i] = fmt.Sprint(v)
	}
	//todo
	request.SshKeyIds = keyIds

	temp2 := d.Get("rdp_allowed_ips").(*schema.Set).List()
	allowedIps := make([]string, len(temp2))
	for i, v := range temp2 {
		allowedIps[i] = fmt.Sprint(v)
	}
	bringLicense := d.Get("bring_your_own_license").(bool)

	temp3 := d.Get("management_access_allowed_ips").(*schema.Set).List()
	managementAccessAllowedIps := make([]string, len(temp3))
	for i, v := range temp3 {
		managementAccessAllowedIps[i] = fmt.Sprint(v)
	}
	installOsToRam := d.Get("install_os_to_ram").(bool)

	var datastoreName string
	if d.Get("esxi") != nil && len(d.Get("esxi").([]interface{})) > 0 {
		esxi := d.Get("esxi").([]interface{})[0]
		esxiItem := esxi.(map[string]interface{})
		if esxiItem["datastore_configuration"] != nil && len(esxiItem["datastore_configuration"].([]interface{})) > 0 {
			datastoreConfiguration := esxiItem["datastore_configuration"].([]interface{})[0]
			datastoreConfigurationItem := datastoreConfiguration.(map[string]interface{})
			datastoreName = datastoreConfigurationItem["datastore_name"].(string)
		}
	}

	var userData string
	if d.Get("cloud_init") != nil && len(d.Get("cloud_init").([]interface{})) > 0 {
		cloudInit := d.Get("cloud_init").([]interface{})[0]
		cloudInitItem := cloudInit.(map[string]interface{})
		userData = cloudInitItem["user_data"].(string)
	}

	var bootUrl, staticDhcpAddressV4 string
	var nativeVlanId int32
	if d.Get("ipxe") != nil && len(d.Get("ipxe").([]interface{})) > 0 {
		iPXE := d.Get("ipxe").([]interface{})[0]
		iPXEItem := iPXE.(map[string]interface{})
		if len(iPXEItem["url"].(string)) > 0 {
			bootUrl = iPXEItem["url"].(string)
		}
		if iPXEItem["native_vlan_configuration"] != nil && len(iPXEItem["native_vlan_configuration"].([]interface{})) > 0 {
			nativeVlanConf := iPXEItem["native_vlan_configuration"].([]interface{})[0]
			nativeVlanConfItem := nativeVlanConf.(map[string]interface{})
			nativeVlanId = int32(nativeVlanConfItem["vlan_id"].(int))
			staticDhcpAddressV4 = nativeVlanConfItem["static_dhcp_address_v4"].(string)
		}
	}

	var controllerAddress, controllerVersion, controllerAuthKey string
	var netris bool
	if d.Get("netris_softgate") != nil && len(d.Get("netris_softgate").([]interface{})) > 0 {
		netrisSoftgate := d.Get("netris_softgate").([]interface{})[0]
		netrisSoftgateItem := netrisSoftgate.(map[string]interface{})
		controllerAddress = netrisSoftgateItem["controller_address"].(string)
		controllerVersion = netrisSoftgateItem["controller_version"].(string)
		controllerAuthKey = netrisSoftgateItem["controller_auth_key"].(string)
		netris = true
	}

	if len(temp2) > 0 || bringLicense || len(temp3) > 0 || installOsToRam || len(datastoreName) > 0 || len(userData) > 0 || len(bootUrl) > 0 || netris {
		dtoOsConfiguration := bmcapiclient.OsConfiguration{}

		if len(temp2) > 0 || bringLicense {
			dtoWindows := bmcapiclient.OsConfigurationWindows{}
			if len(temp2) > 0 {
				dtoWindows.RdpAllowedIps = allowedIps
			}
			dtoWindows.BringYourOwnLicense = &bringLicense
			dtoOsConfiguration.Windows = &dtoWindows
		}
		if len(temp3) > 0 {
			dtoOsConfiguration.ManagementAccessAllowedIps = managementAccessAllowedIps
		}
		if installOsToRam {
			dtoOsConfiguration.InstallOsToRam = &installOsToRam
		}
		if len(datastoreName) > 0 {
			esxiObject := bmcapiclient.EsxiOsConfiguration{}
			datastoreObject := bmcapiclient.EsxiDatastoreConfiguration{}
			datastoreObject.DatastoreName = datastoreName
			esxiObject.DatastoreConfiguration = &datastoreObject
			dtoOsConfiguration.Esxi = &esxiObject
		}
		if len(userData) > 0 {
			cloudInitObject := bmcapiclient.OsConfigurationCloudInit{}
			cloudInitObject.UserData = &userData
			dtoOsConfiguration.CloudInit = &cloudInitObject
		}
		if len(bootUrl) > 0 {
			iPXEObject := bmcapiclient.OsConfigurationIPXE{}
			nativeVlanConfObject := bmcapiclient.OsConfigurationIPXENativeVlanConfiguration{}
			iPXEObject.Url = bootUrl
			if nativeVlanId > 0 {
				nativeVlanConfObject.VlanId = &nativeVlanId
			}
			if len(staticDhcpAddressV4) > 0 {
				nativeVlanConfObject.StaticDhcpAddressV4 = &staticDhcpAddressV4
			}
			iPXEObject.NativeVlanConfiguration = &nativeVlanConfObject
			dtoOsConfiguration.IPXE = &iPXEObject
		}
		if netris {
			netrisSoftgateObject := bmcapiclient.OsConfigurationNetrisSoftgate{}
			netrisSoftgateObject.ControllerAddress = &controllerAddress
			netrisSoftgateObject.ControllerVersion = &controllerVersion
			netrisSoftgateObject.ControllerAuthKey = &controllerAuthKey
			dtoOsConfiguration.NetrisSoftgate = &netrisSoftgateObject
		}
		request.OsConfiguration = &dtoOsConfiguration
	}

	tags := d.Get("tags").([]interface{})
	if len(tags) > 0 {
		tagsObject := make([]bmcapiclient.TagAssignmentRequest, len(tags))
		for i, j := range tags {
			tarObject := bmcapiclient.TagAssignmentRequest{}
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

	query := &dto.Query{}
	var force = d.Get("force").(bool)
	query.Force = force

	// network block
	if d.Get("network_configuration") != nil && len(d.Get("network_configuration").([]interface{})) > 0 {

		networkConfiguration := d.Get("network_configuration").([]interface{})[0]
		networkConfigurationItem := networkConfiguration.(map[string]interface{})

		networkConfigurationObject := bmcapiclient.NetworkConfiguration{}
		gatewayAddress := networkConfigurationItem["gateway_address"].(string)
		if len(gatewayAddress) > 0 {
			networkConfigurationObject.GatewayAddress = &gatewayAddress
		}
		if networkConfigurationItem["private_network_configuration"] != nil && len(networkConfigurationItem["private_network_configuration"].([]interface{})) > 0 {
			privateNetworkConfiguration := networkConfigurationItem["private_network_configuration"].([]interface{})[0]
			privateNetworkConfigurationItem := privateNetworkConfiguration.(map[string]interface{})

			gatewayAddress := privateNetworkConfigurationItem["gateway_address"].(string)
			configurationType := privateNetworkConfigurationItem["configuration_type"].(string)
			privateNetworks := privateNetworkConfigurationItem["private_networks"].([]interface{})

			if len(gatewayAddress) > 0 || len(configurationType) > 0 || len(privateNetworks) > 0 {
				privateNetworkConfigurationObject := bmcapiclient.PrivateNetworkConfiguration{}
				if len(gatewayAddress) > 0 {
					privateNetworkConfigurationObject.GatewayAddress = &gatewayAddress
				}

				if len(configurationType) > 0 {
					privateNetworkConfigurationObject.ConfigurationType = &configurationType
				}

				networkConfigurationObject.PrivateNetworkConfiguration = &privateNetworkConfigurationObject
				if len(privateNetworks) > 0 {
					isNullPrivIps := markNotSetPrivateNetworkIps(d)
					serPrivateNets := make([]bmcapiclient.ServerPrivateNetwork, len(privateNetworks))

					for k, j := range privateNetworks {
						serverPrivateNetworkObject := bmcapiclient.ServerPrivateNetwork{}

						privateNetworkItem := j.(map[string]interface{})

						serverPrivateNetwork := privateNetworkItem["server_private_network"].([]interface{})[0]
						serverPrivateNetworkItem := serverPrivateNetwork.(map[string]interface{})

						id := serverPrivateNetworkItem["id"].(string)
						tempIps := serverPrivateNetworkItem["ips"].(*schema.Set).List()

						netIps := make([]string, len(tempIps))
						for i, v := range tempIps {
							netIps[i] = fmt.Sprint(v)
						}
						dhcp := serverPrivateNetworkItem["dhcp"].(bool)

						if (len(id)) > 0 {
							serverPrivateNetworkObject.Id = id
						}
						if len(privateNetworks) == len(isNullPrivIps) && isNullPrivIps[k] {
							// Do nothing because ips is an empty field
						} else {
							serverPrivateNetworkObject.Ips = netIps
						}

						serverPrivateNetworkObject.Dhcp = &dhcp

						serPrivateNets[k] = serverPrivateNetworkObject
					}
					privateNetworkConfigurationObject.PrivateNetworks = serPrivateNets
				}
			}
		}
		if networkConfigurationItem["ip_blocks_configuration"] != nil && len(networkConfigurationItem["ip_blocks_configuration"].([]interface{})) > 0 {
			ipBlocksConfiguration := networkConfigurationItem["ip_blocks_configuration"].([]interface{})[0]
			ipBlocksConfigurationItem := ipBlocksConfiguration.(map[string]interface{})

			confType := ipBlocksConfigurationItem["configuration_type"].(string)
			ipBlocks := ipBlocksConfigurationItem["ip_blocks"].([]interface{})

			if len(confType) > 0 || len(ipBlocks) > 0 {
				ipBlocksConfigurationObject := bmcapiclient.IpBlocksConfiguration{}
				if len(confType) > 0 {
					ipBlocksConfigurationObject.ConfigurationType = &confType
				}

				networkConfigurationObject.IpBlocksConfiguration = &ipBlocksConfigurationObject
				if len(ipBlocks) > 0 {

					serIpBlocks := make([]bmcapiclient.ServerIpBlock, len(ipBlocks))

					for k, j := range ipBlocks {
						serverIpBlockObject := bmcapiclient.ServerIpBlock{}

						ipBlockItem := j.(map[string]interface{})

						serverIpBlock := ipBlockItem["server_ip_block"].([]interface{})[0]
						serverIpBlockItem := serverIpBlock.(map[string]interface{})

						id := serverIpBlockItem["id"].(string)
						vlanId := int32(serverIpBlockItem["vlan_id"].(int))

						if (len(id)) > 0 {
							serverIpBlockObject.Id = id
						}
						serverIpBlockObject.VlanId = &vlanId
						serIpBlocks[k] = serverIpBlockObject
					}
					ipBlocksConfigurationObject.IpBlocks = serIpBlocks
				}
			}
		}
		if networkConfigurationItem["public_network_configuration"] != nil && len(networkConfigurationItem["public_network_configuration"].([]interface{})) > 0 {
			publicNetworkConfiguration := networkConfigurationItem["public_network_configuration"].([]interface{})[0]
			publicNetworkConfigurationItem := publicNetworkConfiguration.(map[string]interface{})
			publicNetworks := publicNetworkConfigurationItem["public_networks"].([]interface{})

			if len(publicNetworks) > 0 {
				publicNetworkConfigurationObject := bmcapiclient.PublicNetworkConfiguration{}
				networkConfigurationObject.PublicNetworkConfiguration = &publicNetworkConfigurationObject
				serPublicNets := make([]bmcapiclient.ServerPublicNetwork, len(publicNetworks))

				for k, j := range publicNetworks {
					serverPublicNetworkObject := bmcapiclient.ServerPublicNetwork{}

					publicNetworkItem := j.(map[string]interface{})

					serverPublicNetwork := publicNetworkItem["server_public_network"].([]interface{})[0]
					serverPublicNetworkItem := serverPublicNetwork.(map[string]interface{})

					id := serverPublicNetworkItem["id"].(string)
					tempIps := serverPublicNetworkItem["ips"].(*schema.Set).List()

					netIps := make([]string, len(tempIps))
					for i, v := range tempIps {
						netIps[i] = fmt.Sprint(v)
					}
					computeSlaacIp := serverPublicNetworkItem["compute_slaac_ip"].(bool)

					if (len(id)) > 0 {
						serverPublicNetworkObject.Id = id
					}
					serverPublicNetworkObject.Ips = netIps
					serverPublicNetworkObject.ComputeSlaacIp = &computeSlaacIp

					serPublicNets[k] = serverPublicNetworkObject
				}
				publicNetworkConfigurationObject.PublicNetworks = serPublicNets
			}
		}
		request.NetworkConfiguration = &networkConfigurationObject
		// b, _ := json.MarshalIndent(request, "", "  ")
		// log.Printf("request object is" + string(b))
	}
	// end of network block

	// storage block
	if d.Get("storage_configuration") != nil && len(d.Get("storage_configuration").([]interface{})) > 0 {
		storageConfiguration := d.Get("storage_configuration").([]interface{})[0]
		storageConfigurationItem := storageConfiguration.(map[string]interface{})

		storageConfigurationObject := bmcapiclient.StorageConfiguration{}

		if storageConfigurationItem["root_partition"] != nil && len(storageConfigurationItem["root_partition"].([]interface{})) > 0 {
			rootPartition := storageConfigurationItem["root_partition"].([]interface{})[0]
			rootPartitionItem := rootPartition.(map[string]interface{})

			rootPartitionObject := bmcapiclient.StorageConfigurationRootPartition{}

			raid := rootPartitionItem["raid"].(string)
			if len(raid) > 0 {
				rootPartitionObject.Raid = &raid
			}
			size := rootPartitionItem["size"].(int)
			size32 := int32(size)
			rootPartitionObject.Size = &size32

			storageConfigurationObject.RootPartition = &rootPartitionObject
		}
		request.StorageConfiguration = &storageConfigurationObject
	}
	// end of storage block

	requestCommand := server.NewCreateServerCommandWithQuery(client, *request, query)

	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	} else {

		d.SetId(resp.Id)
		d.Set("password", resp.Password)
		if resp.OsConfiguration != nil {
			d.Set("root_password", resp.OsConfiguration.RootPassword)
			d.Set("management_ui_url", resp.OsConfiguration.ManagementUiUrl)
			netrisController := make([]interface{}, 1)
			netrisControllerItem := make(map[string]interface{})
			if resp.OsConfiguration.NetrisController != nil {
				if resp.OsConfiguration.NetrisController.HostOs != nil {
					netrisControllerItem["host_os"] = *resp.OsConfiguration.NetrisController.HostOs
				}
				if resp.OsConfiguration.NetrisController.NetrisWebConsoleUrl != nil {
					netrisControllerItem["netris_web_console_url"] = *resp.OsConfiguration.NetrisController.NetrisWebConsoleUrl
				}
				if resp.OsConfiguration.NetrisController.NetrisUserPassword != nil {
					netrisControllerItem["netris_user_password"] = *resp.OsConfiguration.NetrisController.NetrisUserPassword
				}
			}
			netrisController[0] = netrisControllerItem
			d.Set("netris_controller", netrisController)
		}

		waitResultError := resourceWaitForCreate(resp.Id, &client)
		if waitResultError != nil {
			return waitResultError
		}
	}

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	serverID := d.Id()
	requestCommand := server.NewGetServerCommand(client, serverID)
	resp, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	d.Set("status", resp.Status)
	d.Set("hostname", resp.Hostname)
	d.Set("description", resp.Description)
	d.Set("os", resp.Os)
	d.Set("type", resp.Type)
	d.Set("location", resp.Location)
	d.Set("cpu", resp.Cpu)
	d.Set("cpu_count", resp.CpuCount)
	d.Set("cores_per_cpu", resp.CoresPerCpu)
	d.Set("cpu_frequency_in_ghz", resp.CpuFrequency)
	d.Set("ram", resp.Ram)
	d.Set("storage", resp.Storage)
	d.Set("network_type", resp.NetworkType)
	d.Set("action", "")
	var privateIPs []interface{}
	for _, v := range resp.PrivateIpAddresses {
		privateIPs = append(privateIPs, v)
	}
	d.Set("private_ip_addresses", privateIPs)
	var publicIPs []interface{}
	for _, k := range resp.PublicIpAddresses {
		publicIPs = append(publicIPs, k)
	}
	d.Set("public_ip_addresses", publicIPs)
	d.Set("reservation_id", resp.ReservationId)
	d.Set("pricing_model", resp.PricingModel)

	d.Set("cluster_id", resp.ClusterId)
	if resp.OsConfiguration != nil && resp.OsConfiguration.ManagementAccessAllowedIps != nil {
		var mgmntAccessAllowedIps []interface{}
		for _, k := range resp.OsConfiguration.ManagementAccessAllowedIps {
			mgmntAccessAllowedIps = append(mgmntAccessAllowedIps, k)
		}
		d.Set("management_access_allowed_ips", mgmntAccessAllowedIps)
	}

	if resp.OsConfiguration != nil && resp.OsConfiguration.Windows != nil {
		if resp.OsConfiguration.Windows.RdpAllowedIps != nil {
			var rdpAllowedIps []interface{}
			for _, k := range resp.OsConfiguration.Windows.RdpAllowedIps {
				rdpAllowedIps = append(rdpAllowedIps, k)
			}
			d.Set("rdp_allowed_ips", rdpAllowedIps)
		}
		d.Set("bring_your_own_license", resp.OsConfiguration.Windows.BringYourOwnLicense)
	}

	if resp.OsConfiguration != nil {
		d.Set("install_os_to_ram", resp.OsConfiguration.InstallOsToRam)
		if resp.OsConfiguration.Esxi != nil && resp.OsConfiguration.Esxi.DatastoreConfiguration != nil {
			esxi := make([]interface{}, 1)
			esxiItem := make(map[string]interface{})
			datastoreConfiguration := make([]interface{}, 1)
			datastoreConfigurationItem := make(map[string]interface{})
			datastoreConfigurationItem["datastore_name"] = resp.OsConfiguration.Esxi.DatastoreConfiguration.DatastoreName
			datastoreConfiguration[0] = datastoreConfigurationItem
			esxiItem["datastore_configuration"] = datastoreConfiguration
			esxi[0] = esxiItem
			d.Set("esxi", esxi)
		}
		if resp.OsConfiguration.CloudInit != nil && resp.OsConfiguration.CloudInit.UserData != nil {
			cloudInit := make([]interface{}, 1)
			cloudInitItem := make(map[string]interface{})
			cloudInitItem["user_data"] = *resp.OsConfiguration.CloudInit.UserData
			cloudInit[0] = cloudInitItem
			d.Set("cloud_init", cloudInit)
		}
		if resp.OsConfiguration.IPXE != nil {
			iPXE := make([]interface{}, 1)
			iPXEItem := make(map[string]interface{})
			iPXEItem["url"] = resp.OsConfiguration.IPXE.Url
			nativeVlanConfResp := resp.OsConfiguration.IPXE.NativeVlanConfiguration
			if nativeVlanConfResp != nil {
				nativeVlanConf := make([]interface{}, 1)
				nativeVlanConfItem := make(map[string]interface{})
				if nativeVlanConfResp.VlanId != nil {
					nativeVlanConfItem["vlan_id"] = int(*nativeVlanConfResp.VlanId)
				}
				if nativeVlanConfResp.StaticDhcpAddressV4 != nil {
					nativeVlanConfItem["static_dhcp_address_v4"] = *nativeVlanConfResp.StaticDhcpAddressV4
				}
				if nativeVlanConfResp.Status != nil {
					nativeVlanConfItem["status"] = *nativeVlanConfResp.Status
				}
				nativeVlanConf[0] = nativeVlanConfItem
				iPXEItem["native_vlan_configuration"] = nativeVlanConf
			}
			iPXE[0] = iPXEItem
			d.Set("ipxe", iPXE)
		}
		if resp.OsConfiguration.NetrisSoftgate != nil {
			netrisSoftgate := make([]interface{}, 1)
			netrisSoftgateItem := make(map[string]interface{})
			if resp.OsConfiguration.NetrisSoftgate.HostOs != nil {
				netrisSoftgateItem["host_os"] = *resp.OsConfiguration.NetrisSoftgate.HostOs
			}
			if d.Get("netris_softgate") != nil && len(d.Get("netris_softgate").([]interface{})) > 0 {
				netrisSoftgateInput := d.Get("netris_softgate").([]interface{})[0]
				netrisSoftgateInputItem := netrisSoftgateInput.(map[string]interface{})
				netrisSoftgateItem["controller_address"] = netrisSoftgateInputItem["controller_address"].(string)
				netrisSoftgateItem["controller_version"] = netrisSoftgateInputItem["controller_version"].(string)
				netrisSoftgateItem["controller_auth_key"] = netrisSoftgateInputItem["controller_auth_key"].(string)
			}
			netrisSoftgate[0] = netrisSoftgateItem
			d.Set("netris_softgate", netrisSoftgate)
		}
	}

	if resp.ProvisionedOn != nil {
		d.Set("provisioned_on", resp.ProvisionedOn.String())
	}

	if resp.Tags != nil && len(resp.Tags) > 0 {
		var tagsInput = d.Get("tags").([]interface{})
		tags := flattenServerTags(resp.Tags, tagsInput)
		if err := d.Set("tags", tags); err != nil {
			return err
		}
	}

	var ncInput = d.Get("network_configuration").([]interface{})
	networkConfiguration, err := flattenNetworkConfiguration(&resp.NetworkConfiguration, ncInput)
	if err != nil {
		return err
	}

	if err := d.Set("network_configuration", networkConfiguration); err != nil {
		return err
	}

	var gpuConf bmcapiclient.GpuConfiguration
	if resp.GpuConfiguration != nil {
		gpuConf = *resp.GpuConfiguration
	}
	gpuConfiguration := flattenGpuConfiguration(gpuConf)
	d.Set("gpu_configuration", gpuConfiguration)

	d.Set("superseded_by", resp.SupersededBy)
	d.Set("supersedes", resp.Supersedes)

	return nil
}
func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("action") {
		client := m.(receiver.BMCSDK)
		//var requestCommand helpercommand.Executor
		newStatus := d.Get("action").(string)

		switch newStatus {
		case "powered-on":
			//do power-on request
			serverID := d.Id()
			requestCommand := server.NewPowerOnServerCommand(client, serverID)
			_, err := requestCommand.Execute()
			if err != nil {
				return err
			}
			waitResultError := resourceWaitForPowerON(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}
		case "powered-off":
			//power off request

			serverID := d.Id()

			requestCommand := server.NewPowerOffServerCommand(client, serverID)
			_, err := requestCommand.Execute()
			if err != nil {
				return err
			}
			waitResultError := resourceWaitForPowerOff(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}
		case "reboot":
			//reboot

			serverID := d.Id()
			isIPXE := strings.Contains(d.Get("os").(string), "ipxe")
			rebootRequest := &bmcapiclient.RebootRequest{}
			bootType := "STANDARD"
			if isIPXE {
				bootType = "IPXE"
				if d.Get("ipxe") != nil && len(d.Get("ipxe").([]interface{})) > 0 {
					iPXE := d.Get("ipxe").([]interface{})[0]
					iPXEItem := iPXE.(map[string]interface{})
					if len(iPXEItem["url"].(string)) > 0 {
						url1 := iPXEItem["url"].(string)
						ipxeUrl := bmcapiclient.NullableString{}
						ipxeUrl.Set(&url1)
						rebootRequest.IpxeUrl = ipxeUrl
					}
				}
			}
			rebootRequest.BootType = &bootType

			requestCommand := server.NewRebootServerCommand(client, serverID, *rebootRequest)
			_, err := requestCommand.Execute()
			if err != nil {
				return err
			}
			waitResultError := resourceWaitForCreate(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}
		case "reset": //Deprecated
			//reset
			request := &bmcapiclient.ServerReset{}
			temp := d.Get("ssh_keys").(*schema.Set).List()
			keys := make([]string, len(temp))
			for i, v := range temp {
				keys[i] = fmt.Sprint(v)
			}
			request.SshKeys = keys
			var installDefault = d.Get("install_default_ssh_keys").(bool)
			request.InstallDefaultSshKeys = &installDefault

			temp1 := d.Get("ssh_key_ids").(*schema.Set).List()
			keyIds := make([]string, len(temp1))
			for i, v := range temp1 {
				keyIds[i] = fmt.Sprint(v)
			}
			request.SshKeyIds = keyIds

			dtoOsConfiguration := bmcapiclient.OsConfigurationMap{}
			isWindows := strings.Contains(d.Get("os").(string), "windows")
			isEsxi := strings.Contains(d.Get("os").(string), "esxi")

			if isWindows {
				//log.Printf("Waiting for server windows to be reseted...")
				dtoWindows := bmcapiclient.OsConfigurationWindows{}
				temp2 := d.Get("rdp_allowed_ips").(*schema.Set).List()
				allowedIps := make([]string, len(temp2))
				for i, v := range temp2 {
					allowedIps[i] = fmt.Sprint(v)
				}

				dtoWindows.RdpAllowedIps = allowedIps
				dtoOsConfiguration.Windows = &dtoWindows
				dtoOsConfiguration.Esxi = nil
				request.OsConfiguration = &dtoOsConfiguration
			}

			if isEsxi {
				//log.Printf("Waiting for server esxi to be reseted...")
				dtoEsxi := bmcapiclient.OsConfigurationMapEsxi{}
				temp3 := d.Get("management_access_allowed_ips").(*schema.Set).List()
				managementAccessAllowedIps := make([]string, len(temp3))
				for i, v := range temp3 {
					managementAccessAllowedIps[i] = fmt.Sprint(v)
				}
				dtoEsxi.ManagementAccessAllowedIps = managementAccessAllowedIps
				dtoOsConfiguration.Esxi = &dtoEsxi
				dtoOsConfiguration.Windows = nil
				request.OsConfiguration = &dtoOsConfiguration

			}
			requestCommand := server.NewResetServerCommand(client, d.Id(), *request)
			resp, err := requestCommand.Execute()
			if err != nil {
				return err
			}
			d.Set("password", resp.Password)

			if resp.OsConfiguration != nil && resp.OsConfiguration.Esxi != nil {
				d.Set("root_password", resp.OsConfiguration.Esxi.RootPassword)
				d.Set("management_ui_url", resp.OsConfiguration.Esxi.ManagementUiUrl)
			}

			waitResultError := resourceWaitForCreate(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}

		case "shutdown":

			serverID := d.Id()

			requestCommand := server.NewShutDownServerCommand(client, serverID)
			_, err := requestCommand.Execute()
			if err != nil {
				return err
			}
			waitResultError := resourceWaitForPowerOff(d.Id(), &client)
			if waitResultError != nil {
				return waitResultError
			}

		case "default":
			return fmt.Errorf("unsupported action")
		}

	} else if d.HasChange("pricing_model") {
		client := m.(receiver.BMCSDK)
		//var requestCommand command.Executor
		//reserve action
		request := &bmcapiclient.ServerReserve{}
		//request.Id = d.Id()
		request.PricingModel = d.Get("pricing_model").(string)

		requestCommand := server.NewReserveServerCommand(client, d.Id(), *request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else if d.HasChange("transfer_reservation_to") {
		client := m.(receiver.BMCSDK)
		request := &bmcapiclient.ReservationTransferDetails{}
		serverID := d.Id()
		request.TargetServerId = d.Get("transfer_reservation_to").(string)

		requestCommand := server.NewTransferServerReservationCommand(client, serverID, *request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else if d.HasChange("ipxe") {
		client := m.(receiver.BMCSDK)
		serverID := d.Id()
		request := &bmcapiclient.OsConfigurationIPXE{}
		nativeVlanConfObject := bmcapiclient.OsConfigurationIPXENativeVlanConfiguration{}
		if d.Get("ipxe") != nil && len(d.Get("ipxe").([]interface{})) > 0 {
			iPXE := d.Get("ipxe").([]interface{})[0]
			iPXEItem := iPXE.(map[string]interface{})
			if len(iPXEItem["url"].(string)) > 0 {
				request.Url = iPXEItem["url"].(string)
			}
			if iPXEItem["native_vlan_configuration"] != nil && len(iPXEItem["native_vlan_configuration"].([]interface{})) > 0 {
				nativeVlanConf := iPXEItem["native_vlan_configuration"].([]interface{})[0]
				nativeVlanConfItem := nativeVlanConf.(map[string]interface{})
				nativeVlanId := int32(nativeVlanConfItem["vlan_id"].(int))
				if nativeVlanId > 0 {
					nativeVlanConfObject.VlanId = &nativeVlanId
				}
				staticDhcpAddressV4 := nativeVlanConfItem["static_dhcp_address_v4"].(string)
				if len(staticDhcpAddressV4) > 0 {
					nativeVlanConfObject.StaticDhcpAddressV4 = &staticDhcpAddressV4
				}
			}
			request.NativeVlanConfiguration = &nativeVlanConfObject
		}

		requestCommand := server.NewUpdateServerIPXECommand(client, serverID, *request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else if d.HasChange("tags") {
		tags := d.Get("tags").([]interface{})
		client := m.(receiver.BMCSDK)
		serverID := d.Id()

		var request []bmcapiclient.TagAssignmentRequest

		if len(tags) > 0 {
			request = make([]bmcapiclient.TagAssignmentRequest, len(tags))

			for i, j := range tags {
				tarObject := bmcapiclient.TagAssignmentRequest{}
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
		requestCommand := server.NewSetServerTagsCommand(client, serverID, request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else if d.HasChange("network_configuration") {
		client := m.(receiver.BMCSDK)
		serverID := d.Id()
		query := &dto.Query{}
		var force = d.Get("force").(bool)
		query.Force = force
		oldInterface, newInterface := d.GetChange("network_configuration")
		old := oldInterface.([]interface{})
		new := newInterface.([]interface{})

		if len(new) != 1 || len(old) != 1 {
			return fmt.Errorf("unsupported action")
		}
		ncOldMap := old[0].(map[string]interface{})
		ncNewMap := new[0].(map[string]interface{})
		if d.HasChange("network_configuration.0.gateway_address") {
			return fmt.Errorf("unsupported action, gateway_address has changed")
		} else if d.HasChange("network_configuration.0.ip_blocks_configuration") {
			return fmt.Errorf("unsupported action, ip_blocks_configuration has changed")
		}
		var pncNew, pncOld, pnNew, pnOld []interface{}
		if (ncNewMap["private_network_configuration"]) != nil && len(ncNewMap["private_network_configuration"].([]interface{})) > 0 {
			pncNew = ncNewMap["private_network_configuration"].([]interface{})
		}
		if (ncOldMap["private_network_configuration"]) != nil && len(ncOldMap["private_network_configuration"].([]interface{})) > 0 {
			pncOld = ncOldMap["private_network_configuration"].([]interface{})
		}
		var pncNewMap, pncOldMap map[string]interface{}
		if len(pncNew) > 0 && pncNew[0] != nil {
			pncNewMap = pncNew[0].(map[string]interface{})
		}
		if len(pncOld) > 0 && pncOld[0] != nil {
			pncOldMap = pncOld[0].(map[string]interface{})
		}
		if pncNewMap["gateway_address"] != pncOldMap["gateway_address"] {
			return fmt.Errorf("unsupported action, gateway_address (deprecated) has changed")
		}
		if pncNewMap["configuration_type"] != pncOldMap["configuration_type"] {
			return fmt.Errorf("unsupported action, configuration_type has changed")
		}
		if pncNewMap["private_networks"] != nil {
			pnNew = pncNewMap["private_networks"].([]interface{})
		}
		if pncOldMap["private_networks"] != nil {
			pnOld = pncOldMap["private_networks"].([]interface{})
		}
		var newIds []string
		var newIpss [][]string
		var newDhcps []bool
		isNullPrivIps := markNotSetPrivateNetworkIps(d)
		if len(pnNew) > 0 {
			for i, j := range pnNew {
				pnNewMap := j.(map[string]interface{})
				spnNew := pnNewMap["server_private_network"].([]interface{})[0]
				spnNewMap := spnNew.(map[string]interface{})
				newId := spnNewMap["id"].(string)
				tempIps := spnNewMap["ips"].(*schema.Set).List()
				newIps := make([]string, len(tempIps))
				for i, v := range tempIps {
					newIps[i] = fmt.Sprint(v)
				}
				// Mark an empty array of IPs
				if len(newIps) == 0 && len(isNullPrivIps) == len(pnNew) && !isNullPrivIps[i] {
					newIps = make([]string, 1)
				}
				newDhcp := spnNewMap["dhcp"].(bool)
				newIds = append(newIds, newId)
				newIpss = append(newIpss, newIps)
				newDhcps = append(newDhcps, newDhcp)
			}
		}
		var oldIds []string
		var oldIpss [][]string
		var oldDhcps []bool
		if len(pnOld) > 0 {
			for _, j := range pnOld {
				pnOldMap := j.(map[string]interface{})
				spnOld := pnOldMap["server_private_network"].([]interface{})[0]
				spnOldMap := spnOld.(map[string]interface{})
				oldId := spnOldMap["id"].(string)
				tempIps := spnOldMap["ips"].(*schema.Set).List()
				oldIps := make([]string, len(tempIps))
				for i, v := range tempIps {
					oldIps[i] = fmt.Sprint(v)
				}
				oldDhcp := spnOldMap["dhcp"].(bool)
				oldIds = append(oldIds, oldId)
				oldIpss = append(oldIpss, oldIps)
				oldDhcps = append(oldDhcps, oldDhcp)
			}
		}
		for i, j := range newIds {
			for k := range oldIds {
				if oldIds[k] == j {
					isNullPrivIps := markNotSetPrivateNetworkIps(d)
					checkIps, err := compareInputIps(isNullPrivIps[i], oldIpss[k], newIpss[i])
					if err != nil {
						return err
					} else if !checkIps {
						return fmt.Errorf("unsupported action, ips has changed on private network [id=%s]", j)
					} else if oldDhcps[k] != newDhcps[i] {
						return fmt.Errorf("unsupported action, dhcp has changed on private network [id=%s]", j)
					}
				}
			}
		}
		var sameIds []string
		var idExists bool
		for _, l := range newIds {
			idExists = false
			for _, n := range oldIds {
				if n == l {
					idExists = true
				}
			}
			if idExists {
				sameIds = append(sameIds, l)
			}
		}
		for o, p := range newIds {
			idExists = false
			for _, r := range sameIds {
				if p == r {
					idExists = true
				}
			}
			if !idExists {
				request := &bmcapiclient.ServerPrivateNetwork{}
				request.Id = p
				newIps := newIpss[o]
				if len(newIds) == len(isNullPrivIps) && isNullPrivIps[o] {
					// Do nothing because ips is an empty field
				} else if len(newIps) == 1 && newIps[0] == "" {
					request.Ips = make([]string, 0)
				} else {
					request.Ips = newIps
				}
				request.Dhcp = &newDhcps[o]

				requestCommand := server.NewAddServer2PrivateNetworkCommandWithQuery(client, serverID, *request, query)
				_, err := requestCommand.Execute()
				if err != nil {
					return err
				}
				waitResultError := resourceWaitForPrivateNetworkAssign(d.Id(), p, &client)
				if waitResultError != nil {
					return waitResultError
				}
			}
		}
		for _, t := range oldIds {
			idExists = false
			for _, v := range sameIds {
				if t == v {
					idExists = true
				}
			}
			if !idExists {
				requestCommand := server.NewDeleteServerPrivateNetworkCommand(client, serverID, t)
				_, err := requestCommand.Execute()
				if err != nil {
					return err
				}
				waitResultError := resourceWaitForPrivateNetworkUnassign(d.Id(), t, &client)
				if waitResultError != nil {
					return waitResultError
				}
			}
		}
		var pbncNew, pbncOld, pbnNew, pbnOld []interface{}
		if (ncNewMap["public_network_configuration"]) != nil && len(ncNewMap["public_network_configuration"].([]interface{})) > 0 {
			pbncNew = ncNewMap["public_network_configuration"].([]interface{})
		}
		if (ncOldMap["public_network_configuration"]) != nil && len(ncOldMap["public_network_configuration"].([]interface{})) > 0 {
			pbncOld = ncOldMap["public_network_configuration"].([]interface{})
		}
		var pbncNewMap, pbncOldMap map[string]interface{}
		if len(pbncNew) > 0 && pbncNew[0] != nil {
			pbncNewMap = pbncNew[0].(map[string]interface{})
		}
		if len(pbncOld) > 0 && pbncOld[0] != nil {
			pbncOldMap = pbncOld[0].(map[string]interface{})
		}
		if pbncNewMap["public_networks"] != nil {
			pbnNew = pbncNewMap["public_networks"].([]interface{})
		}
		if pbncOldMap["public_networks"] != nil {
			pbnOld = pbncOldMap["public_networks"].([]interface{})
		}
		var newPubIds []string
		var newPubIpss [][]string
		var newSlaacs []bool
		if len(pbnNew) > 0 {
			for _, j := range pbnNew {
				pbnNewMap := j.(map[string]interface{})
				spbnNew := pbnNewMap["server_public_network"].([]interface{})[0]
				spbnNewMap := spbnNew.(map[string]interface{})
				newPubId := spbnNewMap["id"].(string)
				tempPubIps := spbnNewMap["ips"].(*schema.Set).List()
				newPubIps := make([]string, len(tempPubIps))
				for i, v := range tempPubIps {
					newPubIps[i] = fmt.Sprint(v)
				}
				newSlaac := spbnNewMap["compute_slaac_ip"].(bool)
				newPubIds = append(newPubIds, newPubId)
				newPubIpss = append(newPubIpss, newPubIps)
				newSlaacs = append(newSlaacs, newSlaac)
			}
		}
		var oldPubIds []string
		var oldPubIpss [][]string
		var oldSlaacs []bool
		if len(pbnOld) > 0 {
			for _, j := range pbnOld {
				pbnOldMap := j.(map[string]interface{})
				spbnOld := pbnOldMap["server_public_network"].([]interface{})[0]
				spbnOldMap := spbnOld.(map[string]interface{})
				oldPubId := spbnOldMap["id"].(string)
				tempPubIps := spbnOldMap["ips"].(*schema.Set).List()
				oldPubIps := make([]string, len(tempPubIps))
				for i, v := range tempPubIps {
					oldPubIps[i] = fmt.Sprint(v)
				}
				oldSlaac := spbnOldMap["compute_slaac_ip"].(bool)
				oldPubIds = append(oldPubIds, oldPubId)
				oldPubIpss = append(oldPubIpss, oldPubIps)
				oldSlaacs = append(oldSlaacs, oldSlaac)
			}
		}
		for i, j := range newPubIds {
			for k := range oldPubIds {
				if oldPubIds[k] == j {
					checkIps, err := compareInputIps(false, oldPubIpss[k], newPubIpss[i])
					if err != nil {
						return err
					} else if !checkIps {
						return fmt.Errorf("unsupported action, ips has changed on public network [id=%s]", j)
					} else if oldSlaacs[k] != newSlaacs[i] {
						return fmt.Errorf("unsupported action, compute_slaac_ip has changed on public network [id=%s]", j)
					}
				}
			}
		}
		var samePubIds []string
		var pubIdExists bool
		for _, l := range newPubIds {
			pubIdExists = false
			for _, n := range oldPubIds {
				if n == l {
					pubIdExists = true
				}
			}
			if pubIdExists {
				samePubIds = append(samePubIds, l)
			}
		}
		for o, p := range newPubIds {
			pubIdExists = false
			for _, r := range samePubIds {
				if p == r {
					pubIdExists = true
				}
			}
			if !pubIdExists {
				request := &bmcapiclient.ServerPublicNetwork{}
				request.Id = p
				request.Ips = newPubIpss[o]
				isNullComputeSlaacIp := markNotSetPublicNetworkComputeSlaacIpForUpdate(d)
				var ipv6 bool
				for _, l := range newPubIpss[o] {
					if strings.Contains(l, ":") {
						ipv6 = true
					}
				}
				if ipv6 && len(isNullComputeSlaacIp) == len(newPubIds) {
					if isNullComputeSlaacIp[o] {
						//Do nothing
					} else {
						request.ComputeSlaacIp = &newSlaacs[o]
					}
				}
				requestCommand := server.NewAddServer2PublicNetworkCommandWithQuery(client, serverID, *request, query)
				_, err := requestCommand.Execute()
				if err != nil {
					return err
				}
				waitResultError := resourceWaitForPublicNetworkAssign(d.Id(), p, &client)
				if waitResultError != nil {
					return waitResultError
				}
			}
		}
		for _, t := range oldPubIds {
			pubIdExists = false
			for _, v := range samePubIds {
				if t == v {
					pubIdExists = true
				}
			}
			if !pubIdExists {
				requestCommand := server.NewDeleteServerPublicNetworkCommand(client, serverID, t)
				_, err := requestCommand.Execute()
				if err != nil {
					return err
				}
				waitResultError := resourceWaitForPublicNetworkUnassign(d.Id(), t, &client)
				if waitResultError != nil {
					return waitResultError
				}
			}
		}
	} else if d.HasChange("hostname") || d.HasChange("description") {
		client := m.(receiver.BMCSDK)
		serverID := d.Id()
		request := &bmcapiclient.ServerPatch{}
		var hostname = d.Get("hostname").(string)
		request.Hostname = &hostname
		var desc = d.Get("description").(string)
		request.Description = &desc
		requestCommand := server.NewPatchServerCommand(client, serverID, *request)
		_, err := requestCommand.Execute()
		if err != nil {
			return err
		}
	} else if d.HasChange("delete_ip_blocks") || d.HasChange("force") {
		return resourceServerRead(d, m)
	} else {
		return fmt.Errorf("unsupported action")
	}
	return resourceServerRead(d, m)

}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(receiver.BMCSDK)
	serverID := d.Id()

	var deleteIpBlocks = d.Get("delete_ip_blocks").(bool)

	relinquishIpBlock := bmcapiclient.RelinquishIpBlock{}
	relinquishIpBlock.DeleteIpBlocks = &deleteIpBlocks

	requestCommand := server.NewDeprovisionServerCommand(client, serverID, relinquishIpBlock)

	_, err := requestCommand.Execute()
	if err != nil {
		return err
	}

	return nil
}

func resourceWaitForCreate(id string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for server %s to be created...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "resetting", "rebooting"},
		Target:     []string{"powered-on", "powered-off"},
		Refresh:    refreshForCreate(client, id),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to switch to target state: %v", id, err)
	}

	return nil
}

func resourceWaitForPowerON(id string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for server %s to power on...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"powered-off"},
		Target:     []string{"powered-on"},
		Refresh:    refreshForCreate(client, id),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to power on: %v", id, err)
	}

	return nil
}

func resourceWaitForPowerOff(id string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for server %s to power off...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"powered-on"},
		Target:     []string{"powered-off"},
		Refresh:    refreshForCreate(client, id),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to power off: %v", id, err)
	}

	return nil
}

func refreshForCreate(client *receiver.BMCSDK, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		requestCommand := server.NewGetServerCommand(*client, id)

		resp, err := requestCommand.Execute()
		if err != nil {
			return 0, "", err
		} else {
			return 0, resp.Status, nil
		}
	}
}

func resourceWaitForPrivateNetworkAssign(id string, netId string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for server %s to change private network configuration...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"unassigned", "assigning", "in-progress"},
		Target:     []string{"assigned"},
		Refresh:    refreshForPrivateNetworksChange(client, id, netId),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to change private network configuration: %v", id, err)
	}
	return nil
}

func resourceWaitForPrivateNetworkUnassign(id string, netId string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for server %s to change private network configuration...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"assigned", "unassigning", "in-progress"},
		Target:     []string{"unassigned"},
		Refresh:    refreshForPrivateNetworksChange(client, id, netId),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to change private network configuration: %v", id, err)
	}
	return nil
}

func refreshForPrivateNetworksChange(client *receiver.BMCSDK, id string, netId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var status string
		requestCommand := server.NewGetServerCommand(*client, id)

		resp, err := requestCommand.Execute()
		if err != nil {
			return 0, "", err
		} else if resp.NetworkConfiguration.PrivateNetworkConfiguration == nil {
			return 0, "unassigned", nil
		} else if nets := resp.NetworkConfiguration.PrivateNetworkConfiguration.PrivateNetworks; len(nets) > 0 {
			var netFound bool
			for _, j := range nets {
				if j.Id == netId {
					netFound = true
					if j.StatusDescription != nil {
						status = *j.StatusDescription
						break
					}
				}
			}
			if !netFound {
				status = "unassigned"
			}
			return 0, status, nil
		} else {
			return 0, "unassigned", nil
		}
	}
}

func resourceWaitForPublicNetworkAssign(id string, netId string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for server %s to change public network configuration...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"unassigned", "assigning", "in-progress"},
		Target:     []string{"assigned"},
		Refresh:    refreshForPublicNetworksChange(client, id, netId),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to change public network configuration: %v", id, err)
	}
	return nil
}

func resourceWaitForPublicNetworkUnassign(id string, netId string, client *receiver.BMCSDK) error {
	log.Printf("Waiting for server %s to change public network configuration...", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"assigned", "unassigning", "in-progress"},
		Target:     []string{"unassigned"},
		Refresh:    refreshForPublicNetworksChange(client, id, netId),
		Timeout:    pnapRetryTimeout,
		Delay:      pnapRetryDelay,
		MinTimeout: pnapRetryMinTimeout,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for server (%s) to change public network configuration: %v", id, err)
	}
	return nil
}

func refreshForPublicNetworksChange(client *receiver.BMCSDK, id string, netId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var status string
		requestCommand := server.NewGetServerCommand(*client, id)

		resp, err := requestCommand.Execute()
		if err != nil {
			return 0, "", err
		} else if resp.NetworkConfiguration.PublicNetworkConfiguration == nil {
			return 0, "unassigned", nil
		} else if nets := resp.NetworkConfiguration.PublicNetworkConfiguration.PublicNetworks; len(nets) > 0 {
			var netFound bool
			for _, j := range nets {
				if j.Id == netId {
					netFound = true
					if j.StatusDescription != nil {
						status = *j.StatusDescription
						break
					}
				}
			}
			if !netFound {
				status = "unassigned"
			}
			return 0, status, nil
		} else {
			return 0, "unassigned", nil
		}
	}
}

func flattenNetworkConfiguration(netConf *bmcapiclient.NetworkConfiguration, ncInput []interface{}) ([]interface{}, error) {
	if len(ncInput) == 0 {
		ncInput = make([]interface{}, 1)
		n := make(map[string]interface{})
		ncInput[0] = n
	}
	nci := ncInput[0]
	nciMap := nci.(map[string]interface{})
	var err error

	if netConf != nil {
		if netConf.GatewayAddress != nil {
			nciMap["gateway_address"] = *netConf.GatewayAddress
		}
		if netConf.PrivateNetworkConfiguration != nil {
			prNetConf := *netConf.PrivateNetworkConfiguration
			pnc := make([]interface{}, 1)
			if (nciMap["private_network_configuration"]) != nil && len(nciMap["private_network_configuration"].([]interface{})) > 0 {
				pnc = nciMap["private_network_configuration"].([]interface{})
			}
			pncItem := make(map[string]interface{})
			if len(pnc) > 0 && pnc[0] != nil {
				pncItem = pnc[0].(map[string]interface{})
			}
			if prNetConf.GatewayAddress != nil {
				pncItem["gateway_address"] = *prNetConf.GatewayAddress
			}
			if prNetConf.PrivateNetworks != nil {
				privateNetworks := prNetConf.PrivateNetworks
				pncItem, err = readServerPrivateNetworks(pncItem, privateNetworks)
				if err != nil {
					return nil, err
				}
			}
			pnc[0] = pncItem
			nciMap["private_network_configuration"] = pnc
		}
		if netConf.IpBlocksConfiguration != nil {
			ipBlocksConf := *netConf.IpBlocksConfiguration
			ibc := make([]interface{}, 1)
			if (nciMap["ip_blocks_configuration"]) != nil && len(nciMap["ip_blocks_configuration"].([]interface{})) > 0 {
				ibc = nciMap["ip_blocks_configuration"].([]interface{})
			}
			ibcItem := make(map[string]interface{})
			if len(ibc) > 0 && ibc[0] != nil {
				ibcItem = ibc[0].(map[string]interface{})
			}
			if ipBlocksConf.IpBlocks != nil {
				ipBlocks := ipBlocksConf.IpBlocks
				ibcItem = readServerIpBlocks(ibcItem, ipBlocks)
			}
			ibc[0] = ibcItem
			nciMap["ip_blocks_configuration"] = ibc
		}
		if netConf.PublicNetworkConfiguration != nil {
			pubNetConf := *netConf.PublicNetworkConfiguration
			pnc := make([]interface{}, 1)
			if (nciMap["public_network_configuration"]) != nil && len(nciMap["public_network_configuration"].([]interface{})) > 0 {
				pnc = nciMap["public_network_configuration"].([]interface{})
			}
			pncItem := make(map[string]interface{})
			if len(pnc) > 0 && pnc[0] != nil {
				pncItem = pnc[0].(map[string]interface{})
			}
			if pubNetConf.PublicNetworks != nil {
				pubNet := pubNetConf.PublicNetworks
				pncItem, err = readServerPublicNetworks(pncItem, pubNet)
				if err != nil {
					return nil, err
				}
			}
			pnc[0] = pncItem
			nciMap["public_network_configuration"] = pnc
		}
		return ncInput, nil
	} else {
		return nil, nil
	}
}

func flattenServerTags(tagsRead []bmcapiclient.TagAssignment, tagsInput []interface{}) []interface{} {
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

func supressUserDefinedNetworkType(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if len(oldValue) > 0 {
		if newValue == "USER_DEFINED" || newValue == "PRIVATE_ONLY" || newValue == "PUBLIC_ONLY" || newValue == "PUBLIC_AND_PRIVATE" || newValue == "NONE" {
			return true
		}
	}
	return false
}

// readServerPrivateNetworks reads server private networks from API and sorts them in the same order as in configuration
func readServerPrivateNetworks(pncItem map[string]interface{}, prNet []bmcapiclient.ServerPrivateNetwork) (map[string]interface{}, error) {
	pn1 := make([]interface{}, 0)
	pn2 := make([]interface{}, 0)
	if pncItem["private_networks"] != nil && len(pncItem["private_networks"].([]interface{})) > 0 {
		pni := pncItem["private_networks"].([]interface{})
		// writing networks that are in the configuration, in the same order as in the configuration
		for k := range pni {
			for _, j := range prNet {
				if pni[k] != nil && len(pni[k].(map[string]interface{})["server_private_network"].([]interface{})) > 0 {
					ispnItem := pni[k].(map[string]interface{})["server_private_network"].([]interface{})[0].(map[string]interface{})
					if ispnItem["id"] == j.Id {
						pnItem := make(map[string]interface{})
						spn := make([]interface{}, 1)
						spnItem := make(map[string]interface{})
						spnItem["id"] = j.Id
						ipsInput := ispnItem["ips"].(*schema.Set).List()
						ipsi, err := resolveIps(ipsInput, j.Ips)
						if err != nil {
							return nil, err
						}
						spnItem["ips"] = ipsi
						if j.Dhcp != nil {
							spnItem["dhcp"] = *j.Dhcp
						}
						if j.StatusDescription != nil {
							spnItem["status_description"] = *j.StatusDescription
						}
						if j.VlanId != nil {
							spnItem["vlan_id"] = *j.VlanId
						}
						spn[0] = spnItem
						pnItem["server_private_network"] = spn
						pn1 = append(pn1, pnItem)
					}
				}
			}
		}
		// writing networks that are not in the configuration, if any
		for _, j := range prNet {
			sameId := false
			for k := range pni {
				if pni[k] != nil && len(pni[k].(map[string]interface{})["server_private_network"].([]interface{})) > 0 {
					ispnItem := pni[k].(map[string]interface{})["server_private_network"].([]interface{})[0].(map[string]interface{})
					if ispnItem["id"] == j.Id {
						sameId = true
					}
				}
			}
			if !sameId {
				pnItem := make(map[string]interface{})
				spn := make([]interface{}, 1)
				spnItem := make(map[string]interface{})
				spnItem["id"] = j.Id
				if j.Ips != nil {
					ips := make([]interface{}, len(j.Ips))
					for k, l := range j.Ips {
						ips[k] = l
					}
					spnItem["ips"] = ips
				}
				if j.Dhcp != nil {
					spnItem["dhcp"] = *j.Dhcp
				}
				if j.StatusDescription != nil {
					spnItem["status_description"] = *j.StatusDescription
				}
				if j.VlanId != nil {
					spnItem["vlan_id"] = *j.VlanId
				}
				spn[0] = spnItem
				pnItem["server_private_network"] = spn
				pn2 = append(pn2, pnItem)
			}
		}
		if len(pn2) > 0 {
			for i := range pn2 {
				pn1 = append(pn1, pn2[i])
			}
		}
	} else {
		pn1 = make([]interface{}, len(prNet))
		// writing networks in the same order as in the api response
		for i, j := range prNet {
			pnItem := make(map[string]interface{})
			spn := make([]interface{}, 1)
			spnItem := make(map[string]interface{})
			spnItem["id"] = j.Id
			if j.Ips != nil {
				ips := make([]interface{}, len(j.Ips))
				for k, l := range j.Ips {
					ips[k] = l
				}
				spnItem["ips"] = ips
			}
			if j.Dhcp != nil {
				spnItem["dhcp"] = *j.Dhcp
			}
			if j.StatusDescription != nil {
				spnItem["status_description"] = *j.StatusDescription
			}
			if j.VlanId != nil {
				spnItem["vlan_id"] = *j.VlanId
			}
			spn[0] = spnItem
			pnItem["server_private_network"] = spn
			pn1[i] = pnItem
		}
	}
	pncItem["private_networks"] = pn1
	return pncItem, nil
}

// resolveIps returns configuration value of IPs if it is the same as API value (only written in different format)
// In other cases it returns the API response value
func resolveIps(ipsInput []interface{}, ipsApi []string) ([]interface{}, error) {
	if ipsApi != nil {
		ipsApiMono, err := divideIpsRange(ipsApi)
		if err != nil {
			return nil, err
		}

		ipsInputS := make([]string, len(ipsInput))
		for m, n := range ipsInput {
			ipsInputS[m] = n.(string)
		}
		ipsInputMono, err := divideIpsRange(ipsInputS)
		if err != nil {
			return nil, err
		}

		ipsInputMonoPurged := removeDuplicateIps(ipsInputMono)

		if compareIps(ipsApiMono, ipsInputMonoPurged) {
			return ipsInput, nil
		} else {
			ips := make([]interface{}, len(ipsApi))
			for o, p := range ipsApi {
				ips[o] = p
			}
			return ips, nil
		}
	} else {
		return ipsInput, nil
	}
}

// readServerIpBlocks reads server ip blocks from API and sorts them in the same order as in configuration
func readServerIpBlocks(ibcItem map[string]interface{}, ipBlocks []bmcapiclient.ServerIpBlock) map[string]interface{} {
	ib1 := make([]interface{}, 0)
	ib2 := make([]interface{}, 0)
	if ibcItem["ip_blocks"] != nil && len(ibcItem["ip_blocks"].([]interface{})) > 0 {
		ibi := ibcItem["ip_blocks"].([]interface{})
		// writing ip blocks that are in the configuration, in the same order as in the configuration
		for k := range ibi {
			for _, j := range ipBlocks {
				if ibi[k] != nil && len(ibi[k].(map[string]interface{})["server_ip_block"].([]interface{})) > 0 {
					isibItem := ibi[k].(map[string]interface{})["server_ip_block"].([]interface{})[0].(map[string]interface{})
					if isibItem["id"] == j.Id {
						ibItem := make(map[string]interface{})
						sib := make([]interface{}, 1)
						sibItem := make(map[string]interface{})
						sibItem["id"] = j.Id
						if j.VlanId != nil {
							sibItem["vlan_id"] = *j.VlanId
						}
						sib[0] = sibItem
						ibItem["server_ip_block"] = sib
						ib1 = append(ib1, ibItem)
					}
				}
			}
		}
		// writing ip blocks that are not in the configuration, if any
		for _, j := range ipBlocks {
			sameId := false
			for k := range ibi {
				if ibi[k] != nil && len(ibi[k].(map[string]interface{})["server_ip_block"].([]interface{})) > 0 {
					isibItem := ibi[k].(map[string]interface{})["server_ip_block"].([]interface{})[0].(map[string]interface{})
					if isibItem["id"] == j.Id {
						sameId = true
					}
				}
			}
			if !sameId {
				ibItem := make(map[string]interface{})
				sib := make([]interface{}, 1)
				sibItem := make(map[string]interface{})
				sibItem["id"] = j.Id
				if j.VlanId != nil {
					sibItem["vlan_id"] = *j.VlanId
				}
				sib[0] = sibItem
				ibItem["server_ip_block"] = sib
				ib2 = append(ib2, ibItem)
			}
		}
		if len(ib2) > 0 {
			for i := range ib2 {
				ib1 = append(ib1, ib2[i])
			}
		}
	} else {
		ib1 = make([]interface{}, len(ipBlocks))
		// writing ip blocks in the same order as in the api response
		for i, j := range ipBlocks {
			ibItem := make(map[string]interface{})
			sib := make([]interface{}, 1)
			sibItem := make(map[string]interface{})
			sibItem["id"] = j.Id
			if j.VlanId != nil {
				sibItem["vlan_id"] = *j.VlanId
			}
			sib[0] = sibItem
			ibItem["server_ip_block"] = sib
			ib1[i] = ibItem
		}
	}
	ibcItem["ip_blocks"] = ib1
	return ibcItem
}

// readServerPublicNetworks reads server public networks from API and sorts them in the same order as in configuration
func readServerPublicNetworks(pncItem map[string]interface{}, pubNet []bmcapiclient.ServerPublicNetwork) (map[string]interface{}, error) {
	pn1 := make([]interface{}, 0)
	pn2 := make([]interface{}, 0)
	if pncItem["public_networks"] != nil && len(pncItem["public_networks"].([]interface{})) > 0 {
		pni := pncItem["public_networks"].([]interface{})
		// writing networks that are in the configuration, in the same order as in the configuration
		for k := range pni {
			for _, j := range pubNet {
				if pni[k] != nil && len(pni[k].(map[string]interface{})["server_public_network"].([]interface{})) > 0 {
					ispnItem := pni[k].(map[string]interface{})["server_public_network"].([]interface{})[0].(map[string]interface{})
					if ispnItem["id"] == j.Id {
						pnItem := make(map[string]interface{})
						spn := make([]interface{}, 1)
						spnItem := make(map[string]interface{})
						spnItem["id"] = j.Id
						ipsInput := ispnItem["ips"].(*schema.Set).List()
						var ipv6 bool
						for _, l := range j.Ips {
							if strings.Contains(l, ":") {
								ipv6 = true
							}
						}
						if ipv6 {
							spnItem["ips"] = ipsInput
						} else {
							ipsi, err := resolveIps(ipsInput, j.Ips)
							if err != nil {
								return nil, err
							}
							spnItem["ips"] = ipsi
						}
						if j.StatusDescription != nil {
							spnItem["status_description"] = *j.StatusDescription
						}
						spnItem["compute_slaac_ip"] = ispnItem["compute_slaac_ip"].(bool)
						if j.VlanId != nil {
							spnItem["vlan_id"] = *j.VlanId
						}
						spn[0] = spnItem
						pnItem["server_public_network"] = spn
						pn1 = append(pn1, pnItem)
					}
				}
			}
		}
		// writing networks that are not in the configuration, if any
		for _, j := range pubNet {
			sameId := false
			for k := range pni {
				if pni[k] != nil && len(pni[k].(map[string]interface{})["server_public_network"].([]interface{})) > 0 {
					ispnItem := pni[k].(map[string]interface{})["server_public_network"].([]interface{})[0].(map[string]interface{})
					if ispnItem["id"] == j.Id {
						sameId = true
					}
				}
			}
			if !sameId {
				pnItem := make(map[string]interface{})
				spn := make([]interface{}, 1)
				spnItem := make(map[string]interface{})
				spnItem["id"] = j.Id
				if j.Ips != nil {
					ips := make([]interface{}, len(j.Ips))
					for k, l := range j.Ips {
						ips[k] = l
					}
					spnItem["ips"] = ips
				}
				if j.StatusDescription != nil {
					spnItem["status_description"] = *j.StatusDescription
				}
				if j.VlanId != nil {
					spnItem["vlan_id"] = *j.VlanId
				}
				spn[0] = spnItem
				pnItem["server_public_network"] = spn
				pn2 = append(pn2, pnItem)
			}
		}
		if len(pn2) > 0 {
			for i := range pn2 {
				pn1 = append(pn1, pn2[i])
			}
		}
	} else {
		pn1 = make([]interface{}, len(pubNet))
		// writing networks in the same order as in the api response
		for i, j := range pubNet {
			pnItem := make(map[string]interface{})
			spn := make([]interface{}, 1)
			spnItem := make(map[string]interface{})
			spnItem["id"] = j.Id
			if j.Ips != nil {
				ips := make([]interface{}, len(j.Ips))
				for k, l := range j.Ips {
					ips[k] = l
				}
				spnItem["ips"] = ips
			}
			if j.StatusDescription != nil {
				spnItem["status_description"] = *j.StatusDescription
			}
			if j.VlanId != nil {
				spnItem["vlan_id"] = *j.VlanId
			}
			spn[0] = spnItem
			pnItem["server_public_network"] = spn
			pn1[i] = pnItem
		}
	}
	pncItem["public_networks"] = pn1
	return pncItem, nil
}

// divideIpsRange transforms a slice of IP addresses in range format to a slice of individual IP addresses.
func divideIpsRange(ipsRanged []string) ([]string, error) {
	var ipsMono []string
	for _, j := range ipsRanged {
		if strings.Contains(j, "-") {
			firstLast := strings.Split(j, "-")
			if len(firstLast) > 1 {
				first := strings.TrimSpace(firstLast[0])
				last := strings.TrimSpace(firstLast[1])
				firstAddr, err := netip.ParseAddr(first)
				if err != nil {
					return nil, err
				}
				lastAddr, err := netip.ParseAddr(last)
				if err != nil {
					return nil, err
				}
				nextAddr := firstAddr.Next()

				num1 := binary.BigEndian.Uint32(firstAddr.AsSlice())
				num2 := binary.BigEndian.Uint32(lastAddr.AsSlice())
				for i := num1; i <= num2; i++ {
					if nextAddr == lastAddr {
						ipsMono = append(ipsMono, first, last)
						break
					} else {
						ipsMono = append(ipsMono, first)
						first = nextAddr.String()
						nextAddr = nextAddr.Next()
					}
				}
			}
		} else {
			j = strings.TrimSpace(j)
			ipsMono = append(ipsMono, j)
		}
	}
	return ipsMono, nil
}

// compareIps compares slices of individual IP addresses and returns true if they are equal or false if they are not equal.
func compareIps(ips1 []string, ips2 []string) bool {
	var num int
	if len(ips1) == len(ips2) {
		for _, j := range ips1 {
			ip1, _ := netip.ParseAddr(j)
			var exists bool
			for _, l := range ips2 {
				ip2, _ := netip.ParseAddr(l)
				if ip1 == ip2 {
					num += 1
					exists = true
				}
			}
			if !exists {
				return false
			}
		}
		if len(ips1) == num {
			return true
		}
	}
	return false
}

// removeDuplicateIps removes duplicates of IP addresses
func removeDuplicateIps(ips []string) []string {
	presentIps := make(map[string]bool)
	var ipsPurged []string
	for _, k := range ips {
		if _, l := presentIps[k]; !l {
			presentIps[k] = true
			ipsPurged = append(ipsPurged, k)
		}
	}
	return ipsPurged
}

// compareInputIps compares slices of Ips from configuration. If they are the same returns true, if they are different returns false
func compareInputIps(isNullIps bool, oldIps []string, newIps []string) (bool, error) {
	if len(oldIps) > 0 && len(newIps) > 0 {
		oldIpsMono, err := divideIpsRange(oldIps)
		if err != nil {
			return false, err
		}
		oldIpsMonoPurged := removeDuplicateIps(oldIpsMono)
		newIpsMono, err := divideIpsRange(newIps)
		if err != nil {
			return false, err
		}
		newIpsMonoPurged := removeDuplicateIps(newIpsMono)

		if len(oldIpsMonoPurged) == 1 && oldIpsMonoPurged[0] == "" && len(newIpsMonoPurged) == 0 {
			return true, nil
		} else if compareIps(oldIpsMonoPurged, newIpsMonoPurged) {
			return true, nil
		} else {
			return false, nil
		}
	} else if len(oldIps) == 0 && len(newIps) == 0 {
		return true, nil
	} else if len(oldIps) == 0 && len(newIps) == 1 && newIps[0] == "" {
		return true, nil
	} else if isNullIps && len(newIps) == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// markNotSetPrivateNetworkFields returns which of the private network's fields ips and dhcp, are not set in configuration
func markNotSetPrivateNetworkFields(diff *schema.ResourceDiff) ([]bool, []bool) {
	var isNullPrivIps []bool
	var isNullDhcp []bool
	rawConfig := diff.GetRawConfig()
	objType := rawConfig.Type()
	if ok := objType.HasAttribute("network_configuration"); ok {
		netConf := rawConfig.GetAttr("network_configuration")
		if !netConf.IsNull() && netConf.IsKnown() && netConf.LengthInt() > 0 {
			netConfItem := netConf.Index(cty.NumberIntVal(0))
			pNetConf := netConfItem.GetAttr("private_network_configuration")
			if !pNetConf.IsNull() && pNetConf.IsKnown() && pNetConf.LengthInt() > 0 {
				pNetConfItem := pNetConf.Index(cty.NumberIntVal(0))
				pNet := pNetConfItem.GetAttr("private_networks")
				if !pNet.IsNull() && pNet.IsKnown() && pNet.LengthInt() > 0 {
					for i := 0; i < pNet.LengthInt(); i++ {
						isNullPrivIps = append(isNullPrivIps, false)
						isNullDhcp = append(isNullDhcp, false)
						pNetItem := pNet.Index(cty.NumberIntVal(int64(i)))
						spn := pNetItem.GetAttr("server_private_network")
						if !spn.IsNull() && spn.IsKnown() && spn.LengthInt() > 0 {
							spnItem := spn.Index(cty.NumberIntVal(0))
							ips := spnItem.GetAttr("ips")
							if ips.IsNull() {
								isNullPrivIps[i] = true
							}
							dhcp := spnItem.GetAttr("dhcp")
							if dhcp.IsNull() {
								isNullDhcp[i] = true
							}
						}
					}
				}
			}
		}
	}
	return isNullPrivIps, isNullDhcp
}

// markNotSetPublicNetworkComputeSlaacIp returns which of the public network's compute_slaac_ip are not set in configuration
func markNotSetPublicNetworkComputeSlaacIpForCustomize(diff *schema.ResourceDiff) []bool {
	var isNullComputeSlaacIp []bool
	rawConfig := diff.GetRawConfig()
	objType := rawConfig.Type()
	if ok := objType.HasAttribute("network_configuration"); ok {
		netConf := rawConfig.GetAttr("network_configuration")
		if !netConf.IsNull() && netConf.IsKnown() && netConf.LengthInt() > 0 {
			netConfItem := netConf.Index(cty.NumberIntVal(0))
			pNetConf := netConfItem.GetAttr("public_network_configuration")
			if !pNetConf.IsNull() && pNetConf.IsKnown() && pNetConf.LengthInt() > 0 {
				pNetConfItem := pNetConf.Index(cty.NumberIntVal(0))
				pNet := pNetConfItem.GetAttr("public_networks")
				if !pNet.IsNull() && pNet.IsKnown() && pNet.LengthInt() > 0 {
					for i := 0; i < pNet.LengthInt(); i++ {
						isNullComputeSlaacIp = append(isNullComputeSlaacIp, false)
						pNetItem := pNet.Index(cty.NumberIntVal(int64(i)))
						spn := pNetItem.GetAttr("server_public_network")
						if !spn.IsNull() && spn.IsKnown() && spn.LengthInt() > 0 {
							spnItem := spn.Index(cty.NumberIntVal(0))
							slaac := spnItem.GetAttr("compute_slaac_ip")
							if slaac.IsNull() {
								isNullComputeSlaacIp[i] = true
							}
						}
					}
				}
			}
		}
	}
	return isNullComputeSlaacIp
}

// markNotSetPublicNetworkComputeSlaacIp returns which of the public network's compute_slaac_ip are not set in configuration
func markNotSetPublicNetworkComputeSlaacIpForUpdate(d *schema.ResourceData) []bool {
	var isNullComputeSlaacIp []bool
	rawConfig := d.GetRawConfig()
	objType := rawConfig.Type()
	if ok := objType.HasAttribute("network_configuration"); ok {
		netConf := rawConfig.GetAttr("network_configuration")
		if !netConf.IsNull() && netConf.IsKnown() && netConf.LengthInt() > 0 {
			netConfItem := netConf.Index(cty.NumberIntVal(0))
			pNetConf := netConfItem.GetAttr("public_network_configuration")
			if !pNetConf.IsNull() && pNetConf.IsKnown() && pNetConf.LengthInt() > 0 {
				pNetConfItem := pNetConf.Index(cty.NumberIntVal(0))
				pNet := pNetConfItem.GetAttr("public_networks")
				if !pNet.IsNull() && pNet.IsKnown() && pNet.LengthInt() > 0 {
					for i := 0; i < pNet.LengthInt(); i++ {
						isNullComputeSlaacIp = append(isNullComputeSlaacIp, false)
						pNetItem := pNet.Index(cty.NumberIntVal(int64(i)))
						spn := pNetItem.GetAttr("server_public_network")
						if !spn.IsNull() && spn.IsKnown() && spn.LengthInt() > 0 {
							spnItem := spn.Index(cty.NumberIntVal(0))
							slaac := spnItem.GetAttr("compute_slaac_ip")
							if slaac.IsNull() {
								isNullComputeSlaacIp[i] = true
							}
						}
					}
				}
			}
		}
	}
	return isNullComputeSlaacIp
}

// markNotSetPrivateNetworkIps returns which of the private network's ips are not set in configuration
func markNotSetPrivateNetworkIps(d *schema.ResourceData) []bool {
	var isNullPrivIps []bool
	rawConfig := d.GetRawConfig()
	objType := rawConfig.Type()
	if ok := objType.HasAttribute("network_configuration"); ok {
		netConf := rawConfig.GetAttr("network_configuration")
		if !netConf.IsNull() && netConf.IsKnown() && netConf.LengthInt() > 0 {
			netConfItem := netConf.Index(cty.NumberIntVal(0))
			pNetConf := netConfItem.GetAttr("private_network_configuration")
			if !pNetConf.IsNull() && pNetConf.IsKnown() && pNetConf.LengthInt() > 0 {
				pNetConfItem := pNetConf.Index(cty.NumberIntVal(0))
				pNet := pNetConfItem.GetAttr("private_networks")
				if !pNet.IsNull() && pNet.IsKnown() && pNet.LengthInt() > 0 {
					for i := 0; i < pNet.LengthInt(); i++ {
						isNullPrivIps = append(isNullPrivIps, false)
						pNetItem := pNet.Index(cty.NumberIntVal(int64(i)))
						spn := pNetItem.GetAttr("server_private_network")
						if !spn.IsNull() && spn.IsKnown() && spn.LengthInt() > 0 {
							spnItem := spn.Index(cty.NumberIntVal(0))
							ips := spnItem.GetAttr("ips")
							if ips.IsNull() {
								isNullPrivIps[i] = true
							}
						}
					}
				}
			}
		}
	}
	return isNullPrivIps
}

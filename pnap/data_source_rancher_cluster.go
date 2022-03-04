package pnap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PNAP/go-sdk-helper-bmc/command/ranchersolutionapi/cluster"
	"github.com/PNAP/go-sdk-helper-bmc/receiver"
)

func dataSourceRancherCluster() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceRancherClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initial_cluster_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"server_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRancherClusterRead(d *schema.ResourceData, m interface{}) error {
	if len(d.Get("name").(string)) > 0 {
		client := m.(receiver.BMCSDK)

		requestCommand := cluster.NewGetClustersCommand(client)
		resp, err := requestCommand.Execute()
		if err != nil {
			return err
		}

		if len(d.Get("id").(string)) > 0 {
			numOfClusters := 0
			for _, instance := range resp {
				if instance.Name != nil && instance.Id != nil {
					name := *instance.Name
					id := *instance.Id
					if name == d.Get("name").(string) && id == d.Get("id").(string) {
						numOfClusters++
						d.SetId(id)
						d.Set("id", id)
						d.Set("name", name)
						if instance.Description != nil {
							d.Set("description", *instance.Description)
						}
						d.Set("location", instance.Location)
						if instance.InitialClusterVersion != nil {
							d.Set("initial_cluster_version", *instance.InitialClusterVersion)
						}
						if instance.NodePools != nil {
							np := make([]interface{}, 0)
							nodePools := flattenNodePools(*instance.NodePools, np)
							if err := d.Set("node_pools", nodePools); err != nil {
								return err
							}
						}
						if instance.Metadata != nil {
							metadata := make([]interface{}, 1)
							serverMetadata := make(map[string]interface{})
							if instance.Metadata.Url != nil {
								serverMetadata["url"] = *instance.Metadata.Url
							}
							metadata[0] = serverMetadata
							d.Set("metadata", metadata)
						}
						if instance.StatusDescription != nil {
							d.Set("status_description", *instance.StatusDescription)
						}
					}
				}
			}
			if numOfClusters > 1 {
				return fmt.Errorf("too many clusters with id %s and name %s (found %d, expected 1)", d.Get("id").(string), d.Get("name").(string), numOfClusters)
			}
		} else {
			numOfClusters := 0
			for _, instance := range resp {
				if instance.Name != nil && instance.Id != nil {
					name := *instance.Name
					id := *instance.Id
					if name == d.Get("name").(string) {
						numOfClusters++
						d.SetId(id)
						d.Set("id", id)
						d.Set("name", name)
						if instance.Description != nil {
							d.Set("description", *instance.Description)
						}
						d.Set("location", instance.Location)
						if instance.InitialClusterVersion != nil {
							d.Set("initial_cluster_version", *instance.InitialClusterVersion)
						}
						if instance.NodePools != nil {
							np := make([]interface{}, 0)
							nodePools := flattenNodePools(*instance.NodePools, np)
							if err := d.Set("node_pools", nodePools); err != nil {
								return err
							}
						}
						if instance.Metadata != nil {
							metadata := make([]interface{}, 1)
							serverMetadata := make(map[string]interface{})
							if instance.Metadata.Url != nil {
								serverMetadata["url"] = *instance.Metadata.Url
							}
							metadata[0] = serverMetadata
							d.Set("metadata", metadata)
						}
						if instance.StatusDescription != nil {
							d.Set("status_description", *instance.StatusDescription)
						}
					}
				}
			}
			if numOfClusters > 1 {
				return fmt.Errorf("too many clusters with name %s (found %d, expected 1)", d.Get("name").(string), numOfClusters)
			}
		}

	} else if len(d.Get("id").(string)) > 0 {
		client := m.(receiver.BMCSDK)
		clusterID := d.Get("id").(string)
		requestCommand := cluster.NewGetClusterCommand(client, clusterID)
		resp, err := requestCommand.Execute()
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("unknown cluster identifier")
		}
		d.SetId(*resp.Id)
		d.Set("id", *resp.Id)
		if resp.Name != nil {
			d.Set("name", *resp.Name)
		}
		if resp.Description != nil {
			d.Set("description", *resp.Description)
		}
		d.Set("location", resp.Location)
		if resp.InitialClusterVersion != nil {
			d.Set("initial_cluster_version", *resp.InitialClusterVersion)
		}
		if resp.NodePools != nil {
			np := make([]interface{}, 0)
			nodePools := flattenNodePools(*resp.NodePools, np)
			if err := d.Set("node_pools", nodePools); err != nil {
				return err
			}
		}
		if resp.Metadata != nil {
			metadata := make([]interface{}, 1)
			serverMetadata := make(map[string]interface{})
			if resp.Metadata.Url != nil {
				serverMetadata["url"] = *resp.Metadata.Url
			}
			metadata[0] = serverMetadata
			d.Set("metadata", metadata)
		}
		if resp.StatusDescription != nil {
			d.Set("status_description", *resp.StatusDescription)
		}
	}
	return nil
}

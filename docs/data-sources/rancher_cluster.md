---
layout: "pnap"
page_title: "phoenixNAP: pnap_rancher_cluster"
sidebar_current: "docs-pnap-datasource-rancher-cluster"
description: |-
  Provides a phoenixNAP Rancher Cluster datasource. This can be used to read Rancher Server deployment details.
---

# pnap_rancher_cluster Datasource

Provides a phoenixNAP Rancher Cluster datasource. This can be used to read Rancher Server deployment details.



## Example Usage

Fetch a Rancher Cluster by ID or name and show it's details in alphabetical order. 

```hcl
# Fetch a Rancher Cluster
data "pnap_rancher_cluster" "test" {
  id = "123"
  name = "Rancher-Deployment-1"
}

# Show the Rancher Cluster details
output "rancher-cluster" {
  value = data.pnap_rancher_cluster.test
}
```

## Argument Reference

The following arguments are supported:

* `id` - The cluster (Rancher Cluster) identifier.
* `name` - Cluster name.


## Attributes Reference

The following attributes are exported:

* `id` - The cluster identifier.
* `name` - Cluster name.
* `description` - Cluster description.
* `location` - Deployment location.
* `initial_cluster_version` - The Rancher version that was installed on the cluster during the first creation process.
* `node_pools` - The node pools associated with the cluster.
    * `name` - The name of the node pool.
    * `node_count` - Number of configured nodes.
    * `server_type` - Node server type.
    * `nodes` - The nodes associated with this node pool.        
        * `server_id` - The server identifier.
* `metadata` - Connection parameters to use to connect to the Rancher Server Administrative GUI.
    * `url` - The Rancher Server URL.    
* `status_description` - The cluster status.

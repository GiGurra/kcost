# kcost

`kcost` is a cli tool (& quick hack) to calculate the cost of a Google Kubernetes Engine (GKE)
Autopilot cluster. It iterates over all pods in the cluster, checks their resource consumption against a price table,
and calculates the total cost.

### How it works

`kcost` uses `kubectl` to interact with the Kubernetes cluster. It retrieves information about the nodes and pods in the
cluster, except those in excluded namespaces, and calculates the cost based on the CPU and memory requested by each pod.

### Prerequisites

* Go 1.21 or higher
* kubectl installed and configured to interact with your GKE cluster

### Installation

```
go install github.com/gigurra/kcost@latest
```

### Configuration

`kcost` uses a configuration file named config.yaml. This file should be located in the same directory from where you run
the `kcost` command. The configuration file contains the price details for GKE Autopilot and the namespaces to exclude
from the cost calculation.

Here is an example of a config.yaml file:

```yaml
prices:
  gke:
    autopilot:
      spot:
        cpu: 11
        ram: 1
      regular:
        cpu: 36
        ram: 4
namespaces:
  excluded:
    - 'gke-gmp-system'
    - 'gke-managed-filestorecsi'
    - 'gmp-public'
    - 'kube-node-lease'
    - 'kube-public'
    - 'kube-system'
```

Prices are in currency units per GB (ram) and CPU cores (cpu).
The example prices match the best Autopilot prices in Euros currently (2023-09-25).

### Example

```shell
> kcost

2023/09/25 23:14:53 INFO 
2023/09/25 23:14:53 INFO -----------PRICE FOR NAMESPACE default------------
2023/09/25 23:14:53 INFO  + pod ubuntu-deployment-b8c49ddcf-2qb5b: spot=true, cpu=0.500000, memory=0.500000 => price=6.000000
2023/09/25 23:14:53 INFO  = 6.000000
2023/09/25 23:14:53 INFO 
2023/09/25 23:14:53 INFO -->> TOTAL PRICE: 6.000000
```

The cpu and memory units are in GB and cores respectively.

### License

`kcost` is released under the MIT License.

### Future

* Maybe add some cli args
* Maybe include a default config
* Maybe go eat some ice cream!

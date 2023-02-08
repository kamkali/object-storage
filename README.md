# object-storage
This is a distributed object storage implementation with discovery and load balancing.

## How to run
You can find `Makefile` with useful targets.

As an entry point, you can just do:
```sh
make run
```
Which is pretty much just `docker compose up --build`. App run this way should work out of the box.

### Other targets 
<details open>
<summary>Tests</summary>

- Unit test:
```sh
make test
```
- Integration test (can take a while to init):
```sh
make itest
```
</details>

<details open>
<summary>Tools</summary>

- Tools:
```sh
make tools
```
- Code generation:
```sh
make generate
```
- Format:
```sh
make format
```
- Lint:
```sh
make lint
```
</details>


## High level architecture

<details open>
<summary>Containers</summary>

```mermaid
    C4Container
    title Container diagram for Amazing Object Storage
    Person(customer, Customer, "User of the Gateway", $tags="v1.0")

    Rel(customer, gw, "Calls", "HTTP")

    Container_Boundary(c1, "Object Storage") {
        Container(gw, "Object Gateway", "Golang", "Routes API calls to storage nodes")
        ContainerDb_Ext(node1, "MinIO Node", "Object Storage", "Stores objects")
        ContainerDb_Ext(node2, "MinIO Node", "Object Storage", "Stores objects")
        ContainerDb_Ext(node3, "MinIO Node", "Object Storage", "Stores objects")
    }

    Rel(gw, node1, "Uses", "HTTP")
    
    Rel(gw, node2, "Uses", "HTTP")
    Rel(gw, node3, "Uses", "HTTP")
```

</details>

<details open>
<summary>Components</summary>

```mermaid
    C4Container
    title Container diagram for Amazing Object Storage
    Person(customer, Customer, "User of the Gateway", $tags="v1.0")

    Rel(customer, gwserver, "Calls", "HTTP")

    Container_Boundary(c3, "Object Gateway") {
        Component(gwserver, "Gateway server", "", "Handles user requests")
        Component(service, "Storage service", "", "Validates requests, gets nodes from manager and passes objects")
        Component(manager, "Manager", "", "Manages discovery and load balancing")
        Component(loadbalancer, "Load Balancer", "", "Keeps track of nodes and assigns server to request")

        
        Component(docker, "DockerDiscoverer", "", "Docker Client — implements Discoverer")    
        Component(discovery, "Discoverer", "", "Discovers new nodes")
        Component(node, "Node", "", "Interface for storage node")
        Component(ring, "RingLoadBalancer", "", "Implements Load Balancer")
        
        Component(minio, "Minio Node", "", "MinIO Client — implements Load Balancer")
    }

    ContainerDb_Ext(node1, "MinIO Node", "Object Storage", "Stores objects")
    Rel(gwserver, service, "Uses", "")
    Rel(service, manager, "Uses", "")
    Rel(manager, loadbalancer, "Uses", "")
    Rel(manager, discovery, "Uses", "")
    Rel(loadbalancer, node, "Has many", "")
    Rel(manager, discovery, "Uses", "")
    Rel(ring, loadbalancer, "Implements", "")
    Rel(docker, discovery, "Implements", "")
    Rel(minio, node, "Implements", "")
    Rel(minio, node1, "Uses", "")
```

</details>


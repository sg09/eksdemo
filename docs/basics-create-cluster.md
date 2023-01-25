# Create an Amazon EKS Cluster with Bottlerocket Nodes

`eksdemo` can manage applications in any EKS cluster and the cluster doesn’t have to be created by `eksdemo`. You can use `eksctl` to create the cluster and then manage application using `eksdemo`. However, there are a number of benefits to using `eksdemo` to create your cluster:
* Cluster logging is enabled by default
* OIDC is enabled by default so IAM Roles for Service Accounts (IRSA) works out of the box
* The Managed Node Group ASG max is set to 10 so cluster autoscaling can work out of the box
* Private networking for nodes is set by default
* VPC CNI is configured as a Managed Add-on and configured with IRSA by default
* t3.large instances used by default instead of m5.large for cost savings, but can be easily changed with the `--instance` flag or the shorthand `-i`
* To use containerd as the CRI on Amazon EKS optimized Amazon Linux AMIs is as easy as using the `--containerd` boolean flag
* To create a Fargate profile that selects workloads in the “fargate” namespace, use the `--fargate` boolean flag
* Choose a supported EKS version with the `--version` flag or the shorthand `-v` like `-v 1.21`
* Using a different OS like Bottlerocket or Ubuntu is as easy as `--os bottlerocket` or `--os ubuntu`
* To use IPv6 networking, set the `--ipv6` boolean flag
* If you need to further customize the config, add the `--dry-run` flag and it will output the eksctl YAML config file and you can copy/paste it into a file, make your edits and run `eksctl create cluster -f cluster.yaml`

In this section we will walk through the process of creating an Amazon EKS cluster using `eksdemo` that highlights some of the benefits from the list above. First, review the usage and options of the **`eksdemo create cluster`** command using the help flag `--help` or the shorthand `-h`.

```
» eksdemo create cluster -h
Create EKS Cluster

Usage:
  eksdemo create cluster NAME [flags]

Aliases:
  cluster, clusters

Flags:
      --containerd        use containerd runtime
      --dry-run           don't create, just print out all creation steps
      --fargate           create a Fargate profile
  -h, --help              help for cluster
  -i, --instance string   instance type (default "t3.large")
      --ipv6              use IPv6 networking
      --max int           max nodes (default 10)
      --min int           min nodes
      --no-roles          don't create IAM roles
  -N, --nodes int         desired number of nodes (default 2)
      --os string         Operating System (default "AmazonLinux2")
  -v, --version string    Kubernetes version (default "1.23")

Global Flags:
      --profile string   use the specific profile from your credential file
      --region string    the region to use, overrides config/env settings
  ```

In this example, we would like the following customizations:
* Name our cluster “test”
* Use Bottlerocket nodes instead of Amazon Linux 2
* Use `t3.xlarge` instances instead t3.large
* Create a Managed Node Group with 3 nodes instances instead of 2

The command for this is **`eksdemo create cluster test --os bottlerocket -i t3.xlarge -N 3`**

Before you run the command, let’s dive a bit deeper and understand exactly how `eksdemo` will use and configure `eksctl` to create the cluster. We can do that with the `--dry-run` flag.

```
» eksdemo create cluster test --os bottlerocket -i t3.xlarge -N 3 --dry-run

Eksctl Resource Manager Dry Run:
eksctl create cluster -f -
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: test
  region: us-west-2
  version: "1.22"

addons:
- name: vpc-cni

cloudWatch:
  clusterLogging:
    enableTypes: ["*"]

iam:
  withOIDC: true
  serviceAccounts:
  - metadata:
    name: aws-load-balancer-controller
    <snip>
  - metadata:
    name: cluster-autoscaler
    <snip>
  - metadata:
    name: external-dns
    <snip>
  - metadata:
    name: karpenter
    <snip>

managedNodeGroups:
- name: main
  amiFamily: Bottlerocket
  iam:
    attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy
    - arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly
    - arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore
  instanceType: t3.xlarge
  minSize: 0
  desiredCapacity: 3
  maxSize: 10
  privateNetworking: true
  spot: false
```

You’ll notice that `eksdemo` automatically creates the IAM Roles used for IRSA for the most commonly deployed applications: AWS Load Balancer Controller, Cluster Autoscaler, ExternalDNS and Karpenter. This speeds up installation of the applications later as you don’t have to wait for CloudFormation to create the IAM Roles. To opt out, you can use the `--no-roles` flag.

After reviewing the output above, go ahead and create your cluster.

```
» eksdemo create cluster test --os bottlerocket -i t3.xlarge -N 3
2022-07-11 23:06:06 [ℹ]  eksctl version 0.105.0
2022-07-11 23:06:06 [ℹ]  using region us-west-2
2022-07-11 23:06:06 [ℹ]  setting availability zones to [us-west-2c us-west-2d us-west-2b]
2022-07-11 23:06:06 [ℹ]  subnets for us-west-2c - public:192.168.0.0/19 private:192.168.96.0/19
2022-07-11 23:06:06 [ℹ]  subnets for us-west-2d - public:192.168.32.0/19 private:192.168.128.0/19
2022-07-11 23:06:06 [ℹ]  subnets for us-west-2b - public:192.168.64.0/19 private:192.168.160.0/19
2022-07-11 23:06:06 [ℹ]  nodegroup "main" will use "" [Bottlerocket/1.22]
2022-07-11 23:06:06 [ℹ]  using Kubernetes version 1.22
2022-07-11 23:06:06 [ℹ]  creating EKS cluster "test" in "us-west-2" region with managed nodes
2022-07-11 23:06:06 [ℹ]  1 nodegroup (main) was included (based on the include/exclude rules)
<snip>
2022-07-11 23:25:22 [ℹ]  waiting for CloudFormation stack "eksctl-test-nodegroup-main"
2022-07-11 23:25:22 [ℹ]  waiting for the control plane availability...
2022-07-11 23:25:22 [✔]  saved kubeconfig as "/Users/awsuser/.kube/config"
2022-07-11 23:25:22 [ℹ]  no tasks
2022-07-11 23:25:22 [✔]  all EKS cluster resources for "test" have been created
2022-07-11 23:25:23 [ℹ]  kubectl command should work with "/Users/awsuser/.kube/config", try 'kubectl get nodes'
2022-07-11 23:25:23 [✔]  EKS cluster "test" in "us-west-2" region is ready
```

To view the status and info about your cluster you can run the **`eksdemo get cluster`** command.

```
» eksdemo get cluster
+----------+--------+---------+---------+----------+----------+---------+
|   Age    | Status | Cluster | Version | Platform | Endpoint | Logging |
+----------+--------+---------+---------+----------+----------+---------+
| 17 hours | ACTIVE | blue    |    1.22 | eks.4    | Public   | true    |
| 10 hours | ACTIVE | *test   |    1.22 | eks.4    | Public   | true    |
+----------+--------+---------+---------+----------+----------+---------+
* Indicates current context in local kubeconfig
```

To view detail on the node group, use the **`eksdemo get nodegroup`** command. For this get command and many others, there is a required `--cluster <cluster-name>` flag.

```
» eksdemo get nodegroup --cluster test
+----------+--------+------+-------+-----+-----+----------------+-----------+-------------+
|   Age    | Status | Name | Nodes | Min | Max |    Version     |   Type    | Instance(s) |
+----------+--------+------+-------+-----+-----+----------------+-----------+-------------+
| 10 hours | ACTIVE | main |     3 |   0 |  10 | 1.8.0-a6233c22 | ON_DEMAND | t3.xlarge   |
+----------+--------+------+-------+-----+-----+----------------+-----------+-------------+
```

To view detail on the nodes, use the **`eksdemo get nodes`** command. Here we’ll use the `-c` flag which is the shorthand version of the `--cluster` flag.

```
» eksdemo get nodes -c test
+----------+--------------------+---------------------+------------+-----------+-----------+
|   Age    |        Name        |     Instance Id     |    Zone    | Nodegroup |   Type    |
+----------+--------------------+---------------------+------------+-----------+-----------+
| 10 hours | ip-192-168-110-154 | i-0091caca84b3f86fe | us-west-2c | main      | t3.xlarge |
| 10 hours | ip-192-168-142-85  | i-0e5f24935921f0256 | us-west-2d | main      | t3.xlarge |
| 10 hours | ip-192-168-172-58  | i-0676551dbd61b7bc7 | us-west-2b | main      | t3.xlarge |
+----------+--------------------+---------------------+------------+-----------+-----------+
* Names end with ".us-west-2.compute.internal"
```

Congratulations, your EKS cluster with 3 Bottlerocket `t3.xlarge` nodes is now ready!  In the future if you want to see more detail from a get command you can use `-o yaml` or `-o json` and you will see the raw AWS API response in full. For example, you can try running **`eksdemo get cluster test -o yaml`**. You can also run **`eksdemo get`** to see all the options available.

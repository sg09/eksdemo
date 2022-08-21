# `eksdemo` - kubectl-like CLI for Amazon EKS
The easy button for testing, learning and demoing Amazon EKS:
* Install complex applications and dependencies with a single command
* Extensive application catalog (51 CNCF and open source projects)
* Easily customize application installs with simple command line flags
* Query and search AWS resources with kubectl-like get commands

## Why `eksdemo`?
While creating an EKS cluster is fairly easy thanks to [`eksctl`](https://eksctl.io/), installing and configuring applications on EKS is complex, time consuming and error-prone. One of the most powerful feature of `eksdemo` is its extensive application catalog that can installed (including dependencies) with a single command.

For example, the command: `eksdemo install autoscaling-karpenter -c <cluster-name>`
will:
* Create the Karpenter Node IAM Role
* Create the Karpenter Controller IAM Role (IRSA)
* Add an entry to the `aws-auth` ConfigMap for the Karpenter Node IAM Role
* Install Karpenter
* Create a default Karpenter Provisioner

## No Magic
Application installs are:
* Transparent
    * The `--dry-run` flag prints out all the steps `eksdemo` will take to create dependencies and install the application
* Customizable
    * Each application has optional flags for common configuration options
    * The `--set` flag is available to override any settings in a Helm chart's values file 
* Managed by Helm
    * `eksdemo` embeds Helm as a library and it's used to install all applications, even those that don't have a Helm chart

## `eksdemo` Application Catalog

`eksdemo` comes with an extensive application catalog. Each application can be installed with a single command:
`eksdemo install <application> -c <cluster-name> [flags]`

To install applications under a group, you can use either a space or a hyphen. For example, each of the below are valid:
`eksdemo install ingress nginx`
`eksdemo install ingress-nginx`

The application catalog includes:

* `ack` — AWS Operators for Kubernetes (ACK)
    * `apigatewayv2-controller` — ACK API Gateway v2 Controller
    * `ec2-controller` — ACK EC2 Controller
    * `ecr-controller` — ACK ECR Controller
    * `eks-controller` — ACK EKS Controller
    * `s3-controller` — ACK S3 Controller
* `adot-operator` — AWS Distro for OpenTelemetry Operator
* `appmesh-controller` — AWS App Mesh Controller
* `argo` — Get stuff done with Kubernetes!
    * `cd` — Declarative continuous deployment for Kubernetes
    * `workflows` — Workflow engine for Kubernetes
* `autoscaling` — Kubernetes Autoscaling Applications
    * `cluster-autoscaler` — Kubernetes Cluster Autoscaler
    * `goldilocks` — Get your resource requests "Just Right"
    * `inflate` — Example App to Demonstrate Autoscaling
    * `karpenter` — Karpenter Node Autoscaling
    * `keda` — Kubernetes-based Event Driven Autoscaling
    * `vpa` — Vertical Pod Autoscaler
* `aws-fluent-bit` — AWS Fluent Bit
* `aws-lb-controller` — AWS Load Balancer Controller
* `cert-manager` — Cloud Native Certificate Management
* `cilium` — eBPF-based Networking, Observability, Security
* `container-insights` — CloudWatch Container Insights
    * `adot-collector` — Container Insights ADOT Metrics
    * `cloudwatch-agent` — Container Insights CloudWatch Agent Metrics
    * `fluent-bit` — Container Insights Fluent Bit Logs
    * `prometheus` — CloudWatch Container Insights monitoring for Prometheus
* `crossplane` — Cloud Native Control Plane
* `example` — Example Applications
    * `eks-workshop` — EKS Workshop Example Microservices
    * `game-2048` — Example Game 2048
    * `kube-ops-view` — Kubernetes Operational View
    * `wordpress` — WordPress Blog
* `external-dns` — ExternalDNS
* `falco` — Cloud Native Runtime Security
* `flux` — GitOps family of projects
    * `controllers` — Flux Controllers
    * `sync` — Flux GitRepository to sync with
* `grafana-amp` — Grafana with Amazon Managed Prometheus (AMP)
* `harbor` — Cloud Native Registry
* `ingress` — Ingress Controllers
    * `contour` — Ingress Controller using Envoy proxy
    * `emissary` — Open Source API Gateway from Ambassador
    * `nginx` — NGINX Ingress Controller
* `istio` — Istio Service Mesh
    * `base` — Istio Base (includes CRDs)
    * `istiod` — Istio Control Plane
* `keycloak-amg` — Keycloak SAML iDP for Amazon Managed Grafana
* `kube-prometheus` — Kube Prometheus Stack
* `kubecost` — Monitor & Reduce Kubernetes Spend
* `metrics-server` — Kubernetes Metric Server
* `policy` — Kubernetes Policy Controllers
    * `kyverno` — Kubernetes Native Policy Management
    * `opa-gatekeeper` — Policy Controller for Kubernetes
* `prometheus-amp` — Prometheus with Amazon Managed Prometheus (AMP)
* `storage` — Kubernetes Storage Solutions
    * `ebs-csi` — Amazon EBS CSI driver
    * `efs-csi` — Amazon EFS CSI driver
    * `fsx-csi` — Amazon FSx for Lustre CSI Driver
    * `openebs` — Kubernetes storage simplified
* `velero` — Backup and migrate Kubernetes applications

## `eksdemo` vs EKS Blueprints

Both `eksdemo` and [EKS Blueprints](https://aws.amazon.com/blogs/containers/bootstrapping-clusters-with-eks-blueprints/) automate the creation of EKS clusters and install commonly used applications. Why would you use `eksdemo` for testing, learning and demoing EKS?

| `eksdemo` | EKS Blueprints |
------------|-----------------
Use cases: testing, learning, demoing EKS | Use cases: customers deploying to prod and non-prod environments
Kubectl-like CLI installs apps with single command | Infrastructure as Code (IaC) built on Terraform or CDK
Imperative tooling is great for iterative testing | Declarative IaC tooling is not designed for iterative testing
Used to get up and running quickly | Used to drive drive standards and communicate vetted architecture patterns  for utilizing EKS within customer organizations

## Prerequisites

1. AWS Account with Administrator access
2. Route53 Public Hosted Zone (Optional but strongly recommended)
    1. You can update the domain registration of your existing domain (using any domain registrar) to [change the name servers for the domain to use the four Route 53 name servers](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/migrate-dns-domain-inactive.html#migrate-dns-update-domain-inactive). 
    2. You can still use `eksdemo` if you don’t have a Route53 Hosted Zone. Most applications that explose an Ingress resource default to deploying a Service of type LoadBalancer if you don't use the `--ingress-host` flag and your connection to the application will be unencrypted.

## Install `eksdemo`

`eksdemo` is a Golang binary and releases include support for Mac, Linux and Windows running on x86 or arm64. The easiest way to install is using Homebrew.

### Install using Homebrew

[Homebrew](https://brew.sh/) installation method is supported for Mac and Linux. Using the Terminal, enter the following commands:

```
brew tap aaroniscode/tap
brew install eksdemo
```

Note: Depending on how you originally installed `eksctl`, you may receive the error: `eksctl is already installed from homebrew/core!`  This is because `eksdemo` uses the official Weaveworks tap `weaveworks/tap` as a dependency. _ONLY IF you receive that error_, run the following commands:

```
brew uninstall eksctl
brew install eksdemo
```

### Install Manually

Open https://github.com/aaroniscode/eksdemo/releases/latest in your browser, navigate to Assets, and locate the binary that matches your operation system and platform. Download the file, uncompress and copy to a location of your choice that is in your path. A common location on Mac and Linux is `/usr/local/bin`. Note that `eksctl` is required and [must be installed](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html) as well.

### Set the AWS Region

`eksdemo` requires that a default AWS region is configured. There are 2 ways to configure this:

1. Using the [AWS CLI Configuration file](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html). On Linux and MacOS this file is located in ~/.aws/config. You can use the aws configure command to set the region.
2. Set the `AWS_REGION` environment variable to the desired default region. Unless you set the environment variable in your shell profile, you will need to set this every time you open a new terminal.

### Validate Install

To validate installation you can run the `eksdemo version` command and confirm you are running the latest version. The output will be similar to below:

```
» eksdemo version
eksdemo version info: cmd.Version{Version:"0.1.0", Date:"2022-08-21T17:57:03Z", Commit:"9441aa8"}
```

To validate the AWS region is set, you can run `eksdemo get cluster` which will list running EKS clusters in the default region. If you don’t have any EKS clusters in the region, you will get the response: `No resources found.`

```
» eksdemo get cluster
+-------+--------+---------+---------+----------+----------+---------+
|  Age  | Status | Cluster | Version | Platform | Endpoint | Logging |
+-------+--------+---------+---------+----------+----------+---------+
| 1 day | ACTIVE | *blue   |    1.23 | eks.1    | Public   | true    |
+-------+--------+---------+---------+----------+----------+---------+
* Indicates current context in local kubeconfig
```

# Tutorials

## Create an Amazon EKS Cluster with Bottlerocket Nodes

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
* If you need to further customize the config, add the `--dry-run` flag and it will output the eksctl YAML config file and you can copy/paste it into a file, make your edits and run `eksctl create cluster -f cluster.yaml` 

In this section we will walk through the process of creating an Amazon EKS cluster using `eksdemo` that highlights some of the benefits from the list above. First, review the usage and options of the `eksdemo create cluster` command using the help flag `--help` or the shorthand `-h`.

```
» eksdemo create cluster
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
      --max int           max nodes (default 10)
      --min int           min nodes
      --no-roles          don't create IAM roles
  -N, --nodes int         desired number of nodes (default 2)
      --os string         Operating System (default "AmazonLinux2")
  -v, --version string    Kubernetes version (default "1.23")
  ```

In this example, we would like the following customizations:
* Name our cluster “test”
* Use Bottlerocket nodes instead of Amazon Linux 2
* Use `t3.xlarge` instances instead t3.large
* Create a Managed Node Group with 3 nodes instances instead of 2

The command for this would be: 
`eksdemo create cluster test --os bottlerocket -i t3.xlarge -N 3`

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

To view the status and info about your cluster you can run the  `eksdemo get cluster` command.

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

To view detail on the node group, use the `eksdemo get nodegroup` command. For this get command and many others, there is a required `--cluster <cluster-name>` flag. 

```
» eksdemo get nodegroup --cluster test
+----------+--------+------+-------+-----+-----+----------------+-----------+-------------+
|   Age    | Status | Name | Nodes | Min | Max |    Version     |   Type    | Instance(s) |
+----------+--------+------+-------+-----+-----+----------------+-----------+-------------+
| 10 hours | ACTIVE | main |     3 |   0 |  10 | 1.8.0-a6233c22 | ON_DEMAND | t3.xlarge   |
+----------+--------+------+-------+-----+-----+----------------+-----------+-------------+
```

To view detail on the nodes, use the `eksdemo get nodes` command. Here we’ll use the `-c` flag which is the shorthand version of the `--cluster` flag.

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

Congratulations, your EKS cluster with 3 Bottlerocket `t3.xlarge` nodes is now ready!  In the future if you want to see more detail from a get command you can use `-o yaml` or `-o json` and you will see the raw AWS API response in full. For example, you can try running `eksdemo get cluster test -o yaml`. You can also run `eksdemo get` to see all the options available.
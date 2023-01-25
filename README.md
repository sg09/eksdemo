# `eksdemo` - kubectl-like CLI for Amazon EKS
The easy button for testing, learning and demoing Amazon EKS:
* Install complex applications and dependencies with a single command
* Extensive application catalog (over 50 CNCF, open source and related projects)
* Customize application installs easily with simple command line flags
* Query and search AWS resources with kubectl-like get commands

## Table of Contents
* [Why `eksdemo`?](#why-eksdemo)
* [No Magic](#no-magic)
* [Application Catalog](#application-catalog)
* [`eksdemo` vs EKS Blueprints](#eksdemo-vs-eks-blueprints)
* [Prerequisites](#prerequisites)
* [Install eksdemo](#install-eksdemo)
* [Tutorials](#tutorials)
  * [Basics](#basics)
    * [Create an Amazon EKS Cluster with Bottlerocket Nodes](/docs/basics-create-cluster.md)
    * [Install Game 2048 Example Application with TLS using an ACM Certificate](/docs/basics-install-game-2048.md)
  * [Advanced](#advanced)
    * [Install Karpenter autoscaler and test node provisioning and consolidation](/docs/install-karpenter.md)

## Why `eksdemo`?
While creating an EKS cluster is fairly easy thanks to [`eksctl`](https://eksctl.io/), manually installing and configuring applications on EKS is complex, time consuming and error-prone. One of the most powerful feature of `eksdemo` is its extensive application catalog that can installed (including dependencies) with a single command.

For example, the command: **`eksdemo install autoscaling-karpenter -c <cluster-name>`**
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

## Application Catalog

`eksdemo` comes with an extensive application catalog. Each application can be installed with a single command:
**`eksdemo install <application> -c <cluster-name> [flags]`**

To install applications under a group, you can use either a space or a hyphen. For example, each of the following are valid:
**`eksdemo install ingress nginx`** or **`eksdemo install ingress-nginx`**

The application catalog includes:

* `ack` — AWS Operators for Kubernetes (ACK)
    * `apigatewayv2-controller` — ACK API Gateway v2 Controller
    * `ec2-controller` — ACK EC2 Controller
    * `ecr-controller` — ACK ECR Controller
    * `eks-controller` — ACK EKS Controller
    * `prometheusservice-controller` -- ACK Amazon Managed Prometheus Controller
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
* `crossplane` — Cloud Native Control Planes
* `example` — Example Applications
    * `eks-workshop` — EKS Workshop Example Microservices
    * `game-2048` — Example Game 2048
    * `kube-ops-view` — Kubernetes Operational View
    * `podinfo` — Go app w/microservices best practices
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
* `kubecost` — Visibility Into Kubernetes Spend
    * `eks` — EKS optimized bundle of Kubecost
    * `eks-amp` — EKS optimized Kubecost using Amazon Managed Prometheus
    * `vendor` — Vendor distribution of Kubecost
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

To validate installation you can run the **`eksdemo version`** command and confirm you are running the latest version. The output will be similar to below:

```
» eksdemo version
eksdemo version info: cmd.Version{Version:"0.2.0", Date:"2022-08-31T05:20:29Z", Commit:"b9d662f"}
```

To validate the AWS region is set, you can run **`eksdemo get cluster`** which will list running EKS clusters in the default region. If you don’t have any EKS clusters in the region, you will get the response: `No resources found.`

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

The Basics tutorials provide detailed knowledge on how `eksdemo` works. It's recommended you review the Basics tutorials before proceeding to Advanced tutorial as they assume this knowlege.

## Basics
* [Create an Amazon EKS Cluster with Bottlerocket Nodes](/docs/basics-create-cluster.md)
* [Install Game 2048 Example Application with TLS using an ACM Certificate](/docs/basics-install-game-2048.md)

## Advanced
* [Install Karpenter autoscaler and test node provisioning and consolidation](/docs/install-karpenter.md)
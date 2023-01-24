# Install and Test Karpenter Autoscaling

With `eksdemo`, it’s very easy to setup, install and test Karpenter autoscaling. Karpenter can be installed with a single command, and will automate the following pre-install and post-install step: 

1. Create the EC2 Spot Service Linked Role (if it doesn't already exist)
2. Create the Karpenter Controller IAM Role (IRSA)
3. Create the Karpenter Node IAM Role
4. Create an SQS Queue and EventBridge rules for native Spot Termination Handling
5. Add an entry to the `aws-auth` ConfigMap for the Karpenter Node IAM Role
6. Install the Karpenter Helm Chart
7. Create default Karpenter `Provisioner` and `AWSNodeTemplate` Custom Resources

This tutorial walks through the installation of the Karpenter Autoscaler and the example Inflate application to trigger an autoscaling event. It also tests Node Consolidation.

1. [Install Karpenter Autoscaler](#install-karpenter-autoscaler)
2. [Test Automatic Node Provisioning](#test-automatic-node-provisioning)
3. [Test Node Consolidation](#test-node-consolidation)
4. [(Optional) Inspect Karpenter IAM Roles](#optional-inspect-karpenter-iam-roles)

## Install Karpenter Autoscaler

In this section we walk through the process of installing the Karpenter Autoscaler. The command for performing the installation is: `eksdemo install autoscaling-karpenter -c <cluster-name>`

Let’s expore the command and it’s options by using the -h help shorthand flag.

```
» eksdemo install autoscaling-karpenter -h
Install autoscaling-karpenter

Usage:
  eksdemo install autoscaling-karpenter [flags]

Flags:
  -A, --ami-family string        provisioner ami family (default "AL2")
      --chart-version string     chart version (default "v0.22.1")
  -c, --cluster string           cluster to install application (required)
      --dry-run                  don't install, just print out all installation steps
  -h, --help                     help for karpenter
  -n, --namespace string         namespace to install (default "karpenter")
      --replicas int             number of replicas for the controller deployment (default 1)
      --service-account string   service account name (default "karpenter")
      --set strings              set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)
  -T, --ttl-after-empty int      provisioner ttl seconds after empty (disables consolidation)
      --use-previous             use previous working chart/app versions ("v0.21.1"/"v0.21.1")
  -v, --version string           application version (default "v0.22.1")
```

The Karpenter specific flags are:
* `--ami-family` -- This sets the AMI Family on the default `AWSNodeTemplate`. Options include AL2, Bottlerocket and Ubuntu.
* `--replicas` -- `eksdemo` defaults to only 1 replica for easier log viewing in a demo environment. You can use this flag to increase to the default Karpenter Helm chart value of 2 replicas for high availability.
* `--ttl-after-emty` -- `eksdemo` enables [Consolidation](https://karpenter.sh/docs/concepts/#consolidation) on the default `Provisioner` and this option disables consolidation and sets a Time To Live (TTL) instead.

The `eksdemo` install of Karpenter is identical to the [Getting Started with eksctl](https://karpenter.sh/docs/getting-started/getting-started-with-eksctl/) instructions with the following small differences:
* The Karpenter instructions create a cluster with `karpenter.sh/discovery: <cluster-name>` tags. This is not required and any `eksctl` or `eksdemo` created EKS cluster will work.
* The SQS Queue and EventBridge rules to support native Spot Termination are in a separate CloudFormation stack.
* The SQS Queue name is `karpenter-<cluster-name>` instead of `<cluster-name>`.
* The four EventBridge rule names follow the pattern `karpenter-<cluster-name>-<rule-name>` instead of the automated CloudFormation generated name.
* The Karpenter controller deployment defaults to 1 replica instead of 2.
* There are a few changes to the  default Karpenter `Provisioner` and `AWSNodeTemplate` Custom Resources:
    * Support for both `on-demand` and `spot` nodes is configured instead of only `spot`. It will still default to using Spot  nodes but you can request `on-demand` using the `karpenter.sh/capacity-type` label selector.
    * Consolidated is enabled, but you can disable it using the `--ttl-after-empty` install flag.
    * The `subnetSelector` uses `eksctl` default tags to select only Private subnets.
    * The `securityGroupSelector` uses the `aws:eks:cluster-name: <cluster-name>` tag selector to use the EKS cluster security group.

Optionally, if you want to see details about all the prequisite items that are created when you install Karpenter, you can run the Karpenter install command with the `--dry-run` flag to first inspect all the actions that will be performed. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install autoscaling-karpenter -c <cluster-name> --dry-run
Creating 5 dependencies for autoscaling-karpenter
<snip>
Helm Installer Dry Run:
<snip>
Creating 1 post-install resources for autoscaling-karpenter
Creating post-install resource: karpenter-default-provisioner
<snip>
```

Now, install Karpenter.
```
» eksdemo install autoscaling-karpenter -c <cluster-name>
Creating 5 dependencies for autoscaling-karpenter
<snip>
Downloading Chart: oci://public.ecr.aws/karpenter/karpenter:v0.22.1
Helm installing...
2023/01/23 15:24:01 creating 1 resource(s)
2023/01/23 15:24:01 CRD awsnodetemplates.karpenter.k8s.aws is already present. Skipping.
2023/01/23 15:24:01 creating 1 resource(s)
2023/01/23 15:24:02 CRD provisioners.karpenter.sh is already present. Skipping.
2023/01/23 15:24:08 creating 1 resource(s)
2023/01/23 15:24:09 creating 21 resource(s)
2023/01/23 15:24:10 beginning wait for 21 resources with timeout of 5m0s
Using chart version "v0.22.1", installed "autoscaling-karpenter" version "v0.22.1" in namespace "karpenter"
Creating 1 post-install resources for karpenter
Creating post-install resource: karpenter-default-provisioner
Creating Provisioner "default"
Creating AWSNodeTemplate "default"
```

## Test Automatic Node Provisioning

The test of automatic node provisioning is identical to the [Automatic Node Provisioning](https://karpenter.sh/docs/getting-started/getting-started-with-eksctl/#automatic-node-provisioning) Getting Started instructions. `eksdemo` uses the same example Deployment manifest but uses Helm to install it so it can be managed and uninstalled. Additionally there are a few options. Let's review the help.

```
» eksdemo install autoscaling-inflate -h
Install autoscaling-inflate

Usage:
  eksdemo install autoscaling-inflate [flags]

Flags:
  -c, --cluster string     cluster to install application (required)
      --dry-run            don't install, just print out all installation steps
  -h, --help               help for autoscaling-inflate
  -n, --namespace string   namespace to install (default "inflate")
      --on-demand          request on-demand instances using karpenter node selector
      --replicas int       number of replicas for the deployment (default 0)
      --spread             use topology spread constraints to spread across zones
```
The Inflate specific flags are:
* `--on-demand` -- Requests On-Demand nodes by adding a Node Selector for the `karpenter.sh/capacity-type: on-demand` label 
* `--replicas` -- Sets the number of replicas at install time instead of scaling after install
* `--spread` -- Spreads the replicas across availability zones using Topology Spread Constraint and `topologyKey: topology.kubernetes.io/zone`

Let's test autoscaling with 10 replicas and spread across AZs. We will first review the manifest with the `--dry-run` flag. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install autoscaling-inflate -c <cluster-name> --replicas 10 --spread --dry-run

Manifest Installer Dry Run:
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inflate
  namespace: inflate
spec:
  replicas: 10
  selector:
    matchLabels:
      app: inflate
  template:
    metadata:
      labels:
        app: inflate
    spec:
      terminationGracePeriodSeconds: 0
      containers:
        - name: inflate
          image: public.ecr.aws/eks-distro/kubernetes/pause:3.2
          resources:
            requests:
              cpu: 1
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: topology.kubernetes.io/zone
        whenUnsatisfiable: ScheduleAnyway
        labelSelector:
          matchLabels:
            app: inflate
```

Now, remove the `--dry-run` flag and instal the Inflate app to trigger an autoscaling event.

```
» eksdemo install autoscaling-inflate -c <cluster-name> --replicas 10 --spread
Helm installing...
2023/01/23 16:31:26 creating 1 resource(s)
2023/01/23 16:31:26 creating 1 resource(s)
Using chart version "n/a", installed "autoscaling-inflate" version "n/a" in namespace "inflate"
```

Wait a few moments and then list the EC2 instances in your EKS cluster's VPC using the `eksdemo get ec2-instances` command. The `-c` shorthard cluster flag is optional and filters the instance list to your EKS Cluster VPC.


```
» eksdemo get ec2-instances -c <cluster-name>
+-----------+---------+---------------------+--------------------------------+--------------+------------+
|    Age    |  State  |         Id          |              Name              |     Type     |    Zone    |
+-----------+---------+---------------------+--------------------------------+--------------+------------+
| 5 hours   | running | i-0700781fbdfc33254 | blue-main-Node                 | t3.large     | us-west-2d |
| 1 minute  | running | i-0b9d97609d89c4e87 | karpenter.sh/provisioner-na... | *c5ad.xlarge | us-west-2d |
| 5 hours   | running | i-0f4aef0997638970d | blue-main-Node                 | t3.large     | us-west-2b |
| 1 minute  | running | i-01833ebd5a38c541c | karpenter.sh/provisioner-na... | *m3.xlarge   | us-west-2b |
| 1 minute  | running | i-0fdc7907aab2654ea | karpenter.sh/provisioner-na... | *m3.xlarge   | us-west-2c |
+-----------+---------+---------------------+--------------------------------+--------------+------------+
* Indicates Spot Instance
```

In the example above, you can see that Karpenter has created 3 Spot instances all in separate availabiilty zones: `us-west-2b`, `us-west-2c` and `us-west-2d`. The 10 pods are spread across all 3 AZ's due to the topology spread constraints.

## Test Node Consolidation

To test consolidation, let's reduce the number of replicas for the Inflate deployment from 10 to 5.

```
» kubectl -n inflate scale deploy/inflate --replicas 5
deployment.apps/inflate scaled
```

You may need to wait a few minutes for Karpenter's consolidation logic to make a decision to replace or terminate nodes. Again, use the `eksdemo get ec2` command to view the EC2 instances. This time we will use the shorthand alias `ec2`.

```
» eksdemo get ec2 -c <cluster-name>
+------------+------------+---------------------+--------------------------------+--------------+------------+
|    Age     |   State    |         Id          |              Name              |     Type     |    Zone    |
+------------+------------+---------------------+--------------------------------+--------------+------------+
| 6 hours    | running    | i-0700781fbdfc33254 | blue-main-Node                 | t3.large     | us-west-2d |
| 11 minutes | running    | i-0b9d97609d89c4e87 | karpenter.sh/provisioner-na... | *c5ad.xlarge | us-west-2d |
| 6 hours    | running    | i-0f4aef0997638970d | blue-main-Node                 | t3.large     | us-west-2b |
| 11 minutes | terminated | i-01833ebd5a38c541c | karpenter.sh/provisioner-na... | *m3.xlarge   | us-west-2b |
| 11 minutes | running    | i-0fdc7907aab2654ea | karpenter.sh/provisioner-na... | *m3.xlarge   | us-west-2c |
+------------+------------+---------------------+--------------------------------+--------------+------------+
* Indicates Spot Instance
```

In the example above, Karpenter decided to terminated one of the Spot instances in `us-west-2b` because there is already a Managed Node Group node running in `us-west-2b` so the deployment is still spread across AZ's.

## (Optional) Inspect Karpenter IAM Roles 

The Karpenter install process creates 2 IAM Roles:
* Karpenter Controller IRSA (IAM Role for Service Accounts)
* Karpenter Node IAM Role

Use the `eksdemo get iam-role` command and the `--search` flag to find the roles.

```
» eksdemo get iam-role --search karpenter
+---------+----------------------------------+
|   Age   |              Role                |
+---------+----------------------------------+
| 7 hours | eksdemo.blue.karpenter.karpenter |
| 6 hours | KarpenterNodeRole-blue           |
+---------+----------------------------------+
```

`eksdemo` uses a specific naming convention for IRSA roles: `eksdemo.<cluster-name>.<namespace>.<serviceaccount-name>`. To view the permissions assigned to the role, use the `eksdemo get iam-policy` command and the `--role` command which lists only the policies assigned to the role. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get iam-policy --role eksdemo.<cluster-name>.karpenter.karpenter
+-----------------------------------------------------------------+--------+-------------+
|                               Name                              |  Type  | Description |
+-----------------------------------------------------------------+--------+-------------+
| eksctl-blue-addon-iamserviceaccount-karpenter-karpenter-Policy1 | Inline |             |
+-----------------------------------------------------------------+--------+-------------+
```

To view the details of the policy, including the policy document, you can use the `--output` flag or `-o` shorthand flag to output the raw AWS API responses in either JSON or YAML. In the example we'll use YAML.

```
» eksdemo get iam-policy --role  eksdemo.<cluster-name>.karpenter.karpenter -o yaml
InlinePolicies:
- Name: eksctl-blue-addon-iamserviceaccount-karpenter-karpenter-Policy1
  PolicyDocument: |-
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Action": [
                    "ec2:CreateLaunchTemplate",
                    "ec2:CreateFleet",
                    "ec2:RunInstances",
                    "ec2:CreateTags",
                    "ec2:TerminateInstances",
                    "ec2:DeleteLaunchTemplate",
                    "ec2:DescribeLaunchTemplates",
                    "ec2:DescribeInstances",
                    "ec2:DescribeSecurityGroups",
                    "ec2:DescribeSubnets",
                    "ec2:DescribeImages",
                    "ec2:DescribeInstanceTypes",
                    "ec2:DescribeInstanceTypeOfferings",
                    "ec2:DescribeAvailabilityZones",
                    "ec2:DescribeSpotPriceHistory",
                    "ssm:GetParameter",
                    "pricing:GetProducts"
                ],
                "Resource": "*",
                "Effect": "Allow"
            },
            {
                "Action": [
                    "sqs:DeleteMessage",
                    "sqs:GetQueueUrl",
                    "sqs:GetQueueAttributes",
                    "sqs:ReceiveMessage"
                ],
                "Resource": "arn:aws:sqs:us-west-2:123456789012:karpenter-blue",
                "Effect": "Allow"
            },
            {
                "Action": [
                    "iam:PassRole"
                ],
                "Resource": "arn:aws:iam::123456789012:role/KarpenterNodeRole-blue",
                "Effect": "Allow"
            }
        ]
    }
ManagedPolicies: []
```

Next, let's inspect the Karpenter Node IAM Role. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get iam-policy --role KarpenterNodeRole-<cluster-name>
+------------------------------------+---------+-------------------------------------------+
|                Name                |  Type   |                Description                |
+------------------------------------+---------+-------------------------------------------+
| AmazonSSMManagedInstanceCore       | AWS Mgd | The policy for Amazon EC2 Role            |
|                                    |         | to enable AWS Systems Manager             |
|                                    |         | service core functionality.               |
+------------------------------------+---------+-------------------------------------------+
| AmazonEKS_CNI_Policy               | AWS Mgd | This policy provides the Amazon VPC       |
|                                    |         | CNI Plugin (amazon-vpc-cni-k8s) the       |
|                                    |         | permissions it requires to modify         |
|                                    |         | the IP address configuration on your      |
|                                    |         | EKS worker nodes. This permission set     |
|                                    |         | allows the CNI to list, describe, and     |
|                                    |         | modify Elastic Network Interfaces on      |
|                                    |         | your behalf. More information on the      |
|                                    |         | AWS VPC CNI Plugin is available here:     |
|                                    |         | https://github.com/aws/amazon-vpc-cni-k8s |
+------------------------------------+---------+-------------------------------------------+
| AmazonEC2ContainerRegistryReadOnly | AWS Mgd | Provides read-only access to              |
|                                    |         | Amazon EC2 Container Registry             |
|                                    |         | repositories.                             |
+------------------------------------+---------+-------------------------------------------+
| AmazonEKSWorkerNodePolicy          | AWS Mgd | This policy allows Amazon EKS             |
|                                    |         | worker nodes to connect to                |
|                                    |         | Amazon EKS Clusters.                      |
+------------------------------------+---------+-------------------------------------------+
```

If you want to view the policy document details you can run the above command again adding `-o yaml`.
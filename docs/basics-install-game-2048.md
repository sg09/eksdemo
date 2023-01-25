# Install Game 2048 Example Application with TLS using an ACM Certificate

`eksdemo` makes it extremely easy to install applications from it’s extensive application catalog in your EKS clusters. In this section we will walk through the installation of the Game 2048 example application with TLS using an AWS Certificate Manager (ACM) certificate. There are a few prerequisite steps and applications to install.

1. [Create ACM Certificate](#create-acm-certificate) — This certificate will be attached to the ALB to provide a secure HTTPS connection
2. [Install AWS Load Balancer Controller](#install-aws-load-balancer-controller) — Will serve as an Ingress controller provisioning and configuring the ALB
3. [Install ExternalDNS](#install-externaldns) — Will add a Route 53 record to your Hosted Zone for the Game 2048 application
4. [Install Game 2048 Example Application](#install-game-2048-example-application) — Will use optional configuration flags to specify an Ingress with TLS
5. [(Optional) Game 2048 Installation Configurations](#Optional-game-2048-installation-configurations) — How to deploy without a Hosted Zone using CLB or NLB

## Create ACM Certificate

A publicly trusted certificate is required to access the Game 2048 example application securely over TLS easily with your web browser of choice. While it’s possible to use a self-signed or custom CA signed certificate, many modern browsers try to prevent access to such sites as they are commonly used for fraud, malware and phishing. If you don’t have a Hosted Zone configured in Route 53, you can skip this section and deploy the application insecurely over HTTP.

As mentioned in the [Prerequisites](#prerequisites) above, you will need a domain that you own configured as a Hosted Zone in Route 53. An alternative is creating a free <alias>.people.aws.dev domain in SuperNova. The instructions below and in the following sections will refer to the domain as `example.com`. Please replace all instances of `example.com` with your Route 53 Hosted Zone.

Use the **`eksdemo create acm-certificate`** command to create an ACM wildcard certificate that can be used not only for this application, but other applications in the tutorial as well. **Replace `example.com` with your Route 53 Hosted Zone.**

```
» eksdemo create acm-certificate "*.example.com"
Creating ACM Certificate request for: *.example.com...done
Created ACM Certificate Id: 835e5e51-9720-43f1-8193-5be657e3d649
Validating domain "*.example.com" using hosted zone "example.com"
Waiting for certificate to be issued...done
```

`eksdemo` will not only create the ACM certificate request, but also add the necessary TXT record entries to your Hosted Zone required to validate the certificate. Confirm the certificate details using the **`eksdemo get acm-certificate`** command.

```
» eksdemo get acm-certificate
+-----------+---------------+--------+--------+
|    Age    |     Name      | Status | In Use |
+-----------+---------------+--------+--------+
| 44 weeks  | *.eks.dev     | ISSUED | Yes    |
| 7 minutes | *.example.com | ISSUED | No     |
+-----------+---------------+--------+--------+
```

## Install AWS Load Balancer Controller

The AWS Load Balancer Controller manages AWS Elastic Load Balancers for a Kubernetes clusters. The controller provision a Network Load Balancer (NLB) when you create a Kubernetes service of type `LoadBalancer` and provisions an Application Load Balancer (ALB) when you create a Kubernetes `Ingress`. The install of the Game 2048 example application includes an `Ingress` resources that instructs the AWS Load Balancer Controller to provision an ALB that will enable access to the application over the Internet.

In this section we will walk through the process of installing the AWS Load Balancer Controller application. The command for performing the installation is: **`eksdemo install aws-lb-controller -c <cluster-name>`**

Let’s learn a bit more about the command and it’s options before we continue by using the `-h` help shorthand flag.

```
» eksdemo install aws-lb-controller -h
Install aws-lb-controller

Usage:
  eksdemo install aws-lb-controller [flags]

Aliases:
  aws-lb-controller, aws-lb, awslb

Flags:
      --chart-version string     chart version (default "1.4.5")
  -c, --cluster string           cluster to install application (required)
      --default                  set as the default IngressClass for the cluster
      --dry-run                  don't install, just print out all installation steps
  -h, --help                     help for aws-lb-controller
  -n, --namespace string         namespace to install (default "awslb")
      --service-account string   service account name (default "aws-load-balancer-controller")
      --set strings              set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --use-previous             use previous working chart/app versions ("1.4.4"/"v2.4.3")
  -v, --version string           application version (default "v2.4.4")

Global Flags:
      --profile string   use the specific profile from your credential file
      --region string    the region to use, overrides config/env settings
```

The help content provides a lot of valuable information at a glance. The default chart and application versions, namespace and service account names are included along with optional flags to modify the defaults if desired.

A very powerful optional flag is the `--dry-run` flag. This will print out details about any dependencies and exactly how the application install will take place so there is no mystery about the steps `eksdemo` is taking to install your application. Let’s use the the `--dry-run` flag to understand how the AWS Load Balancer Controller will be installed.

```
Creating 1 dependencies for aws-lb-controller
Creating dependency: aws-lb-controller-irsa

Eksctl Resource Manager Dry Run:
eksctl create iamserviceaccount -f - --approve
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: blue
  region: us-west-2

iam:
  withOIDC: true
  serviceAccounts:
  - metadata:
      name: aws-load-balancer-controller
      namespace: awslb
    roleName: eksdemo.blue.awslb.aws-load-balancer-controller
    roleOnly: true
    attachPolicy:
      <snip>

Helm Installer Dry Run:
+---------------------+----------------------------------+
| Application Version | v2.4.4                           |
| Chart Version       | 1.4.5                            |
| Chart Repository    | https://aws.github.io/eks-charts |
| Chart Name          | aws-load-balancer-controller     |
| Release Name        | aws-lb-controller                |
| Namespace           | awslb                            |
| Wait                | false                            |
+---------------------+----------------------------------+
Set Values: []
Values File:
---
replicaCount: 1
image:
  tag: v2.4.4
fullnameOverride: aws-load-balancer-controller
clusterName: blue
serviceAccount:
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.awslb.aws-load-balancer-controller
  name: aws-load-balancer-controller
region: us-west-2
vpcId: vpc-08a68dc8b440fec75
```

From the `--dry-run` output above, you can see that there is one dependency —  an IAM Role. This role is associated with the AWS Load Balancer Controller’s service account. This is security best practices feature for EKS called [IAM Roles for Service Accounts (IRSA)](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html). `eksdemo` uses `eksctl` to create the IAM Role and the dry run output includes the exact configuration that will be used.

Additionally, the output includes details on how the application will be installed. Most applications, including the AWS Load Balancer Controller, are installed using a Helm chart. The dry run information for Helm installs includes 3 sections:

1. A table with the Helm chart repository URL and name, the chart and application versions and the release name and namespace where the application will be installed.
2. Any `--set` flag variables to override the chart’s `values.yaml` defaults or the values file configuration in the next section. See the Helm documentation for more details on the [format and limitations of the --set flag](https://helm.sh/docs/intro/using_helm/#the-format-and-limitations-of---set).
3. The opinionated values file settings built into the `eksdemo` application catalog. Some of these settings can be change with optional flags. If a flag is not available for the value you wish to change, the `--set` flag can be used to override any value.

With this application and with many others, a number of values file settings are automatically populated. In the example above, the `region`, `vpcID` and AWS Account ID in the IRSA annotation are dynamically updated each time `eksdemo` runs.

With a thorough understanding of how the application install process works, let’s install the AWS Load Balancer controller.

```
» eksdemo install aws-lb-controller -c blue
Creating 1 dependencies for aws-lb-controller
Creating dependency: aws-lb-controller-irsa
2022-11-13 08:27:36 [ℹ]  4 existing iamserviceaccount(s) (awslb/aws-load-balancer-controller,external-dns/external-dns,karpenter/karpenter,kube-system/cluster-autoscaler) will be excluded
2022-11-13 08:27:36 [ℹ]  1 iamserviceaccount (awslb/aws-load-balancer-controller) was excluded (based on the include/exclude rules)
2022-11-13 08:27:36 [!]  serviceaccounts that exist in Kubernetes will be excluded, use --override-existing-serviceaccounts to override
2022-11-13 08:27:36 [ℹ]  no tasks
Downloading Chart: https://aws.github.io/eks-charts/aws-load-balancer-controller-1.4.5.tgz
Helm installing...
2022/11/13 08:27:40 creating 2 resource(s)
2022/11/13 08:27:40 Clearing discovery cache
2022/11/13 08:27:40 beginning wait for 2 resources with timeout of 1m0s
2022/11/13 08:27:48 creating 1 resource(s)
2022/11/13 08:27:48 creating 12 resource(s)
Using chart version "1.4.5", installed "aws-lb-controller" version "v2.4.4" in namespace "awslb"
NOTES:
AWS Load Balancer controller installed!
```

## Install ExternalDNS

ExternalDNS is a [Kubernetes SIGs](https://github.com/kubernetes-sigs) project that synchronizes exposed Kubernetes Services and Ingresses with DNS providers. It [watches the Kubernetes API](https://kubernetes.io/docs/reference/using-api/api-concepts/) for new `Service` and `Ingress` resources to determine which DNS records to configure. The install of the Game 2048 example application includes an `Ingress` resource with a hostname that ExternalDNS will automatically configure in Route 53.

In this section we will walk through the process of installing ExternalDNS. The command for performing the installation is:
**`eksdemo install external-dns -c <cluster-name>`**

Before you continue with the installation, you are encouraged to explore the help for external-dns with the `-h` flag and the dry run output with the `--dry-run` flag. The exact syntax for the commands are:
**`eksdemo install external-dns -h`** and **`eksdemo install external-dns -c blue --dry-run`**

When you are ready to continue, proceed with installing ExternalDNS.

```
» eksdemo install external-dns -c blue
Creating 1 dependencies for external-dns
Creating dependency: external-dns-irsa
2022-11-13 12:56:23 [ℹ]  4 existing iamserviceaccount(s) (awslb/aws-load-balancer-controller,external-dns/external-dns,karpenter/karpenter,kube-system/cluster-autoscaler) will be excluded
2022-11-13 12:56:23 [ℹ]  1 iamserviceaccount (external-dns/external-dns) was excluded (based on the include/exclude rules)
2022-11-13 12:56:23 [!]  serviceaccounts that exist in Kubernetes will be excluded, use --override-existing-serviceaccounts to override
2022-11-13 12:56:23 [ℹ]  no tasks
Downloading Chart: https://github.com/kubernetes-sigs/external-dns/releases/download/external-dns-helm-chart-1.11.0/external-dns-1.11.0.tgz
Helm installing...
2022/11/13 12:56:30 creating 1 resource(s)
2022/11/13 12:56:30 creating 5 resource(s)
Using chart version "1.11.0", installed "external-dns" version "v0.12.2" in namespace "external-dns"
NOTES:
***********************************************************************
* External DNS                                                        *
***********************************************************************
  Chart version: 1.11.0
  App version:   v0.12.2
  Image tag:     k8s.gcr.io/external-dns/external-dns:v0.12.2
***********************************************************************
```

Let’s verify that the applications were installed properly with the **`eksdemo get application`** command. Since this command is specific to a given EKS cluster, the `-c <cluster-name>` flag is required.

```
» eksdemo get application -c blue
+-------------------+--------------+---------+----------+--------+
|       Name        |  Namespace   | Version |  Status  | Chart  |
+-------------------+--------------+---------+----------+--------+
| aws-lb-controller | awslb        | v2.4.4  | deployed | 1.4.5  |
| external-dns      | external-dns | v0.12.2 | deployed | 1.11.0 |
+-------------------+--------------+---------+----------+--------+
```

From the output above we can see that both applications are successfully deployed in the EKS cluster named `blue`. `eksdemo` is using Helm as a Golang client library and the output above is very similar to running a `helm list --all-namespaces` command. Because Helm is bundled as a part of `eksdemo`, that you don’t need to have Helm installed on your system to install or manage any of the applications in the `eksdemo` application catalog.

When ExternalDNS is deployed on AWS, it will query Route 53 for a list of Hosted Zones. IAM Roles for Service Accounts (IRSA) is used to give permissions to access Route 53. You can quickly see all the IAM Roles configured for IRSA by using the **`eksdemo get iam-role -c <cluster-name>`** command. Include the `--last-used` or `-L` shorthand flag to see when the role was last used.

```
» eksdemo get iam-role -c blue -L
oidc.eks.us-west-2.amazonaws.com%2Fid%2F84D3CFB801297D007D945709D8F1C0F6
+----------+-------------------------------------------------+------------+
|   Age    |                      Role                       | Last Used  |
+----------+-------------------------------------------------+------------+
| 14 hours | eksctl-blue-addon-vpc-cni-Role1-1PXCY1L5F2C05   | 1 hour     |
| 14 hours | eksdemo.blue.awslb.aws-load-balancer-controller | -          |
| 14 hours | eksdemo.blue.external-dns.external-dns          | 29 minutes |
| 14 hours | eksdemo.blue.karpenter.karpenter                | -          |
| 14 hours | eksdemo.blue.kube-system.cluster-autoscaler     | -          |
+----------+-------------------------------------------------+------------+
```

Notice that IAM Roles have been created for Cluster Autoscaler and Karpenter even though they haven’t been installed. See [Create an Amazon EKS Cluster with Bottlerocket Nodes](#create-an-amazon-eks-cluster-with-bottlerocket-nodes) for more detail on this and how to disable this feature.

### Install Game 2048 Example Application

The [Game 2048](https://play2048.co/) example application is included as [part of the EKS documentation](https://docs.aws.amazon.com/eks/latest/userguide/alb-ingress.html#application-load-balancer-sample-application) to test and validate the successful deployment of the AWS Load Balancer Controller. The install of the Game 2048 example application includes an `Ingress` resources that instructs the AWS Load Balancer Controller to provision an ALB that will enable access to the application over the Internet.

In this section we will walk through the process of installing the Game 2048 example application. The command for performing the installation is **`eksdemo install example-game-2048 -c <cluster-name>`**

Let’s learn a bit more about the command and it’s options before we continue by using the `-h` help shorthand flag.

```
» eksdemo install example-game-2048 -h
Install example-game-2048

Usage:
  eksdemo install example-game-2048 [flags]

Aliases:
  example-game-2048, example-game2048, example-2048

Flags:
  -c, --cluster string         cluster to install application (required)
      --dry-run                don't install, just print out all installation steps
  -h, --help                   help for example-game-2048
      --ingress-class string   name of IngressClass (default "alb")
  -I, --ingress-host string    hostname for Ingress with TLS (default is Service of type LoadBalancer)
  -n, --namespace string       namespace to install (default "game-2048")
  -X, --nginx-pass string      basic auth password for admin user (only valid with --ingress-class=nginx)
      --nlb                    use NLB instead of CLB (when not using Ingress)
      --replicas int           number of replicas for the deployment (default 1)
      --target-type string     target type when deploying NLB or ALB Ingress (default "ip")
      --use-previous           use previous working chart/app versions (""/"latest")
  -v, --version string         application version (default "latest")

Global Flags:
      --profile string   use the specific profile from your credential file
      --region string    the region to use, overrides config/env settings
```

You’ll notice above there is an optional `--ingress-host` flag with a `-I` shorthand version of the flag. For this application and others that have external access, `eksdemo` defaults to using a Service of type `LoadBalancer` without any encryption (HTTPS). If you have a Hosted Zone configured in Route 53, you will include the Ingress Host flag with the fully qualified domain name for the application, like `-I game2048.example.com`.

Since Game 2048 is included in the EKS documentation as a manifest file, let’s use the the `--dry-run` flag to understand how the application will be installed. **Replace `example.com` with your Hosted Zone.**

```
» eksdemo install example-game-2048 -c blue -I game2048.example.com --dry-run

Manifest Installer Dry Run:
---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: game-2048
  name: deployment-2048
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: app-2048
  replicas: 1
  template:
    metadata:
      labels:
        app.kubernetes.io/name: app-2048
    spec:
      containers:
      - image: public.ecr.aws/l6m2t8p7/docker-2048:latest
        imagePullPolicy: Always
        name: app-2048
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  namespace: game-2048
  name: service-2048
  annotations:
    {}
spec:
  ports:
    - port: 80
      targetPort: 80
      protocol: TCP
  type: ClusterIP
  selector:
    app.kubernetes.io/name: app-2048
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: game-2048
  name: ingress-2048
  annotations:
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS":443}]'
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/ssl-redirect: '443'
    alb.ingress.kubernetes.io/target-type: ip
spec:
  ingressClassName: alb
  rules:
    - host: game2048.example.com
      http:
        paths:
        - path: /
          pathType: Prefix
          backend:
            service:
              name: service-2048
              port:
                number: 80
  tls:
  - hosts:
    - game2048.example.com
```

The dry run output is different from the previous example and simply displays the manifest to be installed. When a Helm chart is not available for an application, the manifest is included in the EKS binary. The manifest is stored as a template and is rendered dynamically each time `eksdemo` is run and can change significantly depending on the flags used. You can run the command again without the `-I` flag to see how the Service object type is changed to `LoadBalancer` and the Ingress resource is removed.

One of the benefits of using a Helm chart is that applications can be easily managed and uninstalled. A powerful feature of `eksdemo` is that all applications are installed as Helm charts even if the application is only available as a manifest. Since `eksdemo` bundles Helm as a library, it dynamically generate a Helm chart in memory from the manifest files before deploying the application.

Now that you know how a manifest install works, let’s install the Game 2048 example application. **Replace `example.com` with your Hosted Zone.**

```
» eksdemo install example-game-2048 -c blue -I game2048.example.com
Helm installing...
2022/11/15 19:45:20 creating 1 resource(s)
2022/11/15 19:45:20 creating 3 resource(s)
Using chart version "n/a", installed "example-game-2048" version "latest" in namespace "game-2048"
```

Let’s check the status of all three of our installed applications, understanding that they are all installed as Helm charts.

```
» eksdemo get application -c blue
+-------------------+--------------+---------+----------+--------+
|       Name        |  Namespace   | Version |  Status  | Chart  |
+-------------------+--------------+---------+----------+--------+
| aws-lb-controller | awslb        | v2.4.4  | deployed | 1.4.5  |
| example-game-2048 | game-2048    | latest  | deployed | n/a    |
| external-dns      | external-dns | v0.12.2 | deployed | 1.11.0 |
+-------------------+--------------+---------+----------+--------+
```

The Ingress resource that is created as part the Game 2048 example application install will trigger the AWS Load Balancer Controller to create an ALB. This will take a few minutes to provision. You can check on the status of the ALB by using the **`eksdemo get load-balancer`** command. For this command, the `-c <cluster-name>` flag is optional, and if used it will filter the query to ELB’s in the same VPC as the `blue` EKS cluster.

```
» eksdemo get load-balancer -c blue
+-----------+--------+----------------------------------+------+-------+-----+-----+
|    Age    | State  |               Name               | Type | Stack | AZs | SGs |
+-----------+--------+----------------------------------+------+-------+-----+-----+
| 3 minutes | active | k8s-game2048-ingress2-0d50dcef8e | ALB  | ipv4  |   3 |   2 |
+-----------+--------+----------------------------------+------+-------+-----+-----+
* Indicates internal load balancer
```

If the state shows as `provisioning`, wait a moment and run the command again.

Next let’s confirm that ExternalDNS has created a Route 53 resource record for our application. The command to query Route 53 records is **`eksdemo get dns-records --zone <zone-name>`.** `eksdemo` has a lot of shorthand aliases for commands and flags and you can discover these by using the `--help` flag on any command. For the `get dns-records` command we’ll use the command alias `dns` and for the `--zone` flag, we’ll use the shorthand `-z`.

**Replace `example.com` with your Hosted Zone.**

```
» eksdemo get dns -z example.com
+----------------------------+------+---------------------------------------------------------------------+
|          Name              | Type |                                Value                                |
+----------------------------+------+---------------------------------------------------------------------+
| example.com                | NS   | ns-1855.awsdns-39.co.uk.                                            |
|                            |      | ns-1452.awsdns-53.org.                                              |
|                            |      | ns-921.awsdns-51.net.                                               |
|                            |      | ns-35.awsdns-04.com.                                                |
| example.com                | SOA  | ns-1855.awsdns-39.co.uk.                                            |
|                            |      | awsdns-hostmaster.amazon.com.                                       |
|                            |      | 1 7200 900 1209600 86400                                            |
| cname-game2048.example.com | TXT  | "heritage=external-dns,external-dns/owner=blue,external-dns/reso... |
| game2048.example.com       | A    | k8s-game2048-ingress2-0d50dcef8e-334176506.us-west-2.elb.amazona... |
| game2048.example.com       | TXT  | "heritage=external-dns,external-dns/owner=blue,external-dns/reso... |
+----------------------------+------+---------------------------------------------------------------------+
```

We can see that an `A` record has been created for `game2048.example.com` that points to the DNS name of the ALB. Next open your web browser and enter `https://game2048.example.com` (**replace `example.com` with your Hosted Zone**) to load your Game 2048 example application!

![Game 2048 Screenshot](/docs/images/game-2048-screenshot.jpg "Game 2048 Screenshot")

Congratulations! You’ve successfully deployed the Game 2048 example application over HTTPS with a publicly trusted certificate!

NOTE: It’s possible you may have to wait for DNS to propagate. The time depends on your local ISP and operating system. If you get a DNS resolution error, you can wait and try again later. Or if you’d like to troubleshoot a bit further, A2 Hosting has a Knowledge base article [How to test DNS with dig and nslookup](https://www.a2hosting.com/kb/getting-started-guide/internet-and-networking/troubleshooting-dns-with-dig-and-nslookup).

Tips:

* Wait a minute or two after the Route 53 A record is created before querying on your computer. I’ve found that if I perform a lookup too fast before DNS has propagated, the operating system can cache the response for some time.
* On my Mac I’ve found that `dig` will directly query the local name servers and will have the latest data and `nslookup` will use the host cache that can have stale data.
* If you believe your DNS cache is to blame, consider this article [How to Flush DNS Cache: Windows and Mac](https://constellix.com/news/how-to-flush-dns-cache-windows-mac).

## (Optional) Game 2048 Installation Configurations

If you don’t have a Hosted Zone or want to deploy the Game 2048 example application unencrypted over HTTP, you can run the command without the `--ingress-host` flag or `-I` shorthand flag: **`eksdemo install example-game-2048 -c blue`**

By default, the application will deployed with a Service of type `LoadBalancer`, which will deploy a Classic Load Balancer (CLB). There are a number of flags that allow you to choose more options:

```
Flags:
      --ingress-class string   name of IngressClass (default "alb")
  -I, --ingress-host string    hostname for Ingress with TLS (default is Service of type LoadBalancer)
  -X, --nginx-pass string      basic auth password for admin user (only valid with --ingress-class=nginx)
      --nlb                    use NLB instead of CLB (when not using Ingress)
      --target-type string     target type when deploying NLB or ALB Ingress (default "ip")
```

To expose the application unencrypted as a Service using an NLB in Instance mode, the command is:
**`eksdemo install example-game-2048 -c blue --nlb --target-type instance`**

To expose the application encrypted as an Ingress using Nginx Ingress, the command is:
**`eksdemo install example-game-2048 -c blue -I game2048.example.com --ingress-class nginx`**

If exposing using a Service and NLB, you will need to have AWS Load Balancer Controller installed. If exposing using an Ingress, you will need to have the Ingress Controller and ExternalDNS installed. Also, if using an IngressClass other than `alb`, you will need to have cert-manager installed.

# AWS

As student you are restricted in certain actions. You can't create a code pipeline due to missing rights.

## Integration with different VCS with CodeBuild

![AWS vcs source](img/aws/cloudbuild_2.png)

Supported in console:

- Amazon S3
- AWS CodeCommit
- Github and Github Enterprise
- Bitbucket

Similar with gcloud, gitlab can also be used with to start builds.
Two possible ways are:
- aws console and ci script that pushes changes
- aws fargate with coordinators that push into an aws S3 bucket which triggers a vloud build
- aws lambda function that pushes into a bucket

## Github actions + aws


## Github 

### Overview:


### Pricing

<https://calculator.aws/#/>
<https://aws.amazon.com/amazon-linux-2/>
<https://aws.amazon.com/rds/postgresql/pricing/>

The AWS CodeBuild free tier includes 100 build minutes of build.general1.small per month. The CodeBuild free tier does not expire automatically at the end of your 12-month AWS Free Tier term. It is available to new and existing AWS customers.

#### CodeBuild

| Compute instance type 	| Memory 	| vCPU 	| Linux price per build minute 	| Windows price per build minute 	|
|:-:	|:-:	|:-:	|:-:	|:-:	|
| general1.small 	| 3 GB 	| 2 	| $0.005 	| N/A 	|
| general1.medium 	| 7 GB 	| 4 	| $0.01 	| N/A 	|
| arm1.large 	| 16 GiB 	| 8 	| $0.0175 	| N/A 	|
| general1.large 	| 15 GB 	| 8 	| $0.02 	| N/A 	|
| general1.2xlarge 	| 144 GiB 	| 72 	| $0.25 	| N/A 	|
| gpu1.large 	| 244 GiB 	| 32 	| $0.80 	| N/A 	|


You may incur additional charges if your builds transfer data or use other AWS services. For example, you may incur charges from Amazon CloudWatch Logs for build log streams, Amazon S3 for build artifact storage, and AWS Key Management Service for encryption. You may also incur additional charges if you use AWS CodeBuild with AWS CodePipeline.

If you want to use the secret manager:

Pricing
$0.40 per secret per month
$0.05 per 10,000 API calls

#### Container Registry

<https://console.aws.amazon.com/ecr/home?region=us-east-1#>
<https://aws.amazon.com/ecr/pricing/>

AWS Free Tier
As part of the AWS Free Tier, new Amazon ECR customers get 500 MB-month of storage for one year for your private repositories.

As a new or existing customer, Amazon ECR offers you 50 GB-month of always-free storage for your public repositories. You can transfer 500 GB of data to the internet for free from a public repository each month anonymously (without using an AWS account.) If you sign up for an AWS account, or authenticate to ECR with an existing AWS Account, you can transfer 5 TB of data to the internet for free from a public repository each month, and you get unlimited bandwidth for free when transferring data from a public repository to AWS compute resources in any AWS Region.

 Storage:
    Storage is $0.10 per GB-month for data stored in private or public repositories.

Data transferred from private repositories:

|  	| Pricing 	|
|:-:	|:-:	|
| Data Transfer IN 	|  	|
| All data transfer in 	| $0.00 per GB 	|
| Data Transfer OUT ** 	|  	|
| Up to 1 GB / Month 	| $0.00 per GB 	|
| Next 9.999 TB / Month 	| $0.09 per GB 	|
| Next 40 TB / Month 	| $0.085 per GB 	|
| Next 100 TB / Month 	| $0.07 per GB 	|
| Greater than 150 TB / Month 	| $0.05 per GB 	|

Prices for public repositories:
|  	| Pricing 	|
|:-:	|:-:	|
| Greater than 5 TB / Month to non AWS Regions  	| $0.09 per GB |

Data transfer "in" and “out” refers to transfer into and out of Amazon Elastic Container Registry. Data transferred between Amazon Elastic Container Registry and Amazon EC2 within a single region is free of charge (i.e., $0.00 per GB). Data transferred between Amazon Elastic Container Registry and Amazon EC2 in different regions will be charged at Internet Data Transfer rates on both sides of the transfer.

### Steps

#### Create Trigger

Before you create a codebuild project, you need to review the IAM policies.  
You need CloudWatch Logs , S3, Systems Manager, CodeCommit, CodeBuild and Elastic Container Registry  rights otherwise the build will fail.

1. Go to AWS CodeBuild
2. Create a new build project
3. Select a project name
4. Select your VCS Repository (Github)
   1. If Github is not allowed access yet, you can select repositories that AWS has access to
5. Add a webhook event
   1. Add condition for a new tag push -> HEAD_REF : `^refs/tags/v\d+\.\d+\.\d+`  
   2. Add another filter group for pull requests
6. Manage the environment, [amazon linux 2](https://aws.amazon.com/amazon-linux-2/) is the recommended system as it provides packages and configurations with many aws tools
   1. Set the privileged flag as we want to build docker images
7. Additional specifications can be set for the build vm
   1. Timeouts, resources, certificates, file systems ..
8. Add environment variables
   1. You have the option to use it as plaintext, parameter or with the secrets manager
9. Artifacts are not produced by our build, or rather they are pushed to the docker registry and as such we will not have any artifacts that should be saved
10. Select CloudWatch logs, as they provide additional insight into how many builds failed and more important stats.

In AWS ECR, you can have only have one repository per image but each repository can have multiple versions of a single image. ECR only allows pushes if a repository has already been created.
In our case we have to create 3 repositories.

#### Building with cloudbuild.yaml

<https://docs.aws.amazon.com/codebuild/latest/userguide/build-spec-ref.html>

Few issues you might have:
- Incorrect policies set for cloudbuild
  - <https://docs.aws.amazon.com/codebuild/latest/userguide/setting-up.html#setting-up-service-role>
  - <https://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-access.html>
  - <https://docs.aws.amazon.com/AmazonECR/latest/userguide/repository-policies.html>
- Retry build fails
  - Changes to environment variables are not updated for a build that is restarted even if different values for parameters are shown in the build job (likely a bug)

### Amazon RDS

To access a RDS instance from outside you need to select a VPC with a security Group that has incoming connections allowed. Setting source to `0.0.0.0/0` for inbound connections allows every IPv4 with every port.

#### Pricing

<https://aws.amazon.com/rds/postgresql/pricing/>


Multi AZ:
Select Create Replica in Different Zone to have Amazon RDS maintain a synchronous standby replica in a different Availability Zone than the DB instance. Amazon RDS will automatically fail over to the standby in the case of a planned or unplanned outage of the primary. 

On-Demand DB Instances:

On-Demand DB Instances let you pay for compute capacity by the hour your DB Instance runs with no long-term commitments. This frees you from the costs and complexities of planning, purchasing, and maintaining hardware and transforms what are commonly large fixed costs into much smaller variable costs.


| Standard Instances - Current Generation 	| Price Per Hour 	|
|:-:	|:-:	|
| db.t3.micro 	| $0.021 	|
| db.t3.small 	| $0.042 	|
| db.t3.medium 	| $0.084 	|
| ..... 	| ..... 	|
| db.m5.16xlarge 	| $6.784 	|
| db.m5.24xlarge 	| $10.176 	|


Reserved instances:

Amazon RDS Reserved Instances give you the option to reserve a DB instance for a one or three year term and in turn receive a significant discount compared to the On-Demand Instance pricing for the DB instance. Amazon RDS provides three RI payment options -- No Upfront, Partial Upfront, All Upfront -- that enable you to balance the amount you pay upfront with your effective hourly price.

Amazon RDS Reserved Instances provide size flexibility for the PostgreSQL database engine. With size flexibility, your RI’s discounted rate will automatically apply to usage of any size in the same instance family (M5, T3, R5, etc.)

Please note that Reserved Instance prices don't cover storage or I/O costs. To learn more about features, payment options and rules, please visit our Reserved Instances page.

Region:
1 Year  
db.t3.micro  
Frankfurt  

| Payment Option 	| Upfront 	| Monthly* 	| Effective Hourly** 	| Savings over On-Demand 	| On-Demand Hourly 	|
|:-:	|:-:	|:-:	|:-:	|:-:	|:-:	|
|                        No Upfront                      	|                        $0                      	|                        $10.439                      	|                                                 $0.014                                             	|                        32%                      	|                          $0.0210                        	|
|                        Partial Upfront                      	|                        $60                      	|                        $4.964                      	|                                                 $0.014                                             	|                        35%                      	|                          $0.0210                        	|
|                        All Upfront                      	|                        $117                      	|                        $0.000                      	|                                                 $0.013                                             	|                        36%                      	|                          $0.0210                        	|

1 Year db.r5.24xlarge  
Frankfurt  

| Payment Option 	| Upfront 	| Monthly* 	| Effective Hourly** 	| Savings over On-Demand 	| On-Demand Hourly 	|
|:-:	|:-:	|:-:	|:-:	|:-:	|:-:	|
|                        No Upfront                      	|                        $0                      	|                        $7,069.612                      	|                                                 $9.684                                             	|                        34%                      	|                          $14.6400                        	|
|                        Partial Upfront                      	|                        $40,398                      	|                        $3,366.468                      	|                                                 $9.223                                             	|                        37%                      	|                          $14.6400                        	|
|                        All Upfront                      	|                        $79,179                      	|                        $0.000                      	|                                                 $9.039                                             	|                        38%                      	|                          $14.6400                        	|


### Elastic Beanstalk

Load balancer only with:  
High availability
High availability (using Spot and On-Demand instances)


Needs IAM permissions to load from ecr, error warnings give no info whatsoever if permissions are missing:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowEbAuth",
            "Effect": "Allow",
            "Action": [
                "ecr:GetAuthorizationToken"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Sid": "AllowPull",
            "Effect": "Allow",
            "Resource": [
                "arn:aws:ecr:*:YOUR_ID:repository/*"
            ],
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability"
            ]
        }
    ]
}
```

In addition the security group needs to be set and this:

1

Fixed. Answering myself, for future reference.

There were two basic problems:

- EB really (really!) wants to get connections on port 80 from the public internet. At deployment time it creates a new security group for the environment, which opens port 80. Editing that group to open more ports does not help because of problem #2.
    The default configuration for a single-instance docker environment will deploy with nginx as a reverse proxy, mapping port 80 of the instance to whatever port was configured as the HostPort in Dockerrun.aws.json (or to the ContainerPort if no HostPort is defined). This is a problem for a Minecraft client/server connection, because nginx is at bottom a web server, and the client sends packets that are not valid HTTP requests.

So, the solution is to:

- Make the client connect to port 80, specifying it as IP_ADDRESS:80
    Remove nginx from the configuration. The easiest way to do so is through the Web UI: after the EB environment launches, click the Configuration link, then the Modify button in the Software section; select None from the Proxy Server pulldown at the top, then click the Apply Configuration button. The environment will be re-deployed, but with port 80 mapped straight to the docker container through iptables, without a reverse proxy in between.


#### Deployment policies

<https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/using-features.rolling-version-deploy.html?icmpid=docs_elasticbeanstalk_console>

AWS Elastic Beanstalk provides several options for how deployments are processed, including deployment policies (All at once, Rolling, Rolling with additional batch, Immutable, and Traffic splitting) and options that let you configure batch size and health check behavior during deployments. By default, your environment uses all-at-once deployments. If you created the environment with the EB CLI and it's a scalable environment (you didn't specify the --single option), it uses rolling deployments.

With rolling deployments, Elastic Beanstalk splits the environment's Amazon EC2 instances into batches and deploys the new version of the application to one batch at a time. It leaves the rest of the instances in the environment running the old version of the application. During a rolling deployment, some instances serve requests with the old version of the application, while instances in completed batches serve other requests with the new version. For details, see How rolling deployments work.

To maintain full capacity during deployments, you can configure your environment to launch a new batch of instances before taking any instances out of service. This option is known as a rolling deployment with an additional batch. When the deployment completes, Elastic Beanstalk terminates the additional batch of instances.

Immutable deployments perform an immutable update to launch a full set of new instances running the new version of the application in a separate Auto Scaling group, alongside the instances running the old version. Immutable deployments can prevent issues caused by partially completed rolling deployments. If the new instances don't pass health checks, Elastic Beanstalk terminates them, leaving the original instances untouched.

Traffic-splitting deployments let you perform canary testing as part of your application deployment. In a traffic-splitting deployment, Elastic Beanstalk launches a full set of new instances just like during an immutable deployment. It then forwards a specified percentage of incoming client traffic to the new application version for a specified evaluation period. If the new instances stay healthy, Elastic Beanstalk forwards all traffic to them and terminates the old ones. If the new instances don't pass health checks, or if you choose to abort the deployment, Elastic Beanstalk moves traffic back to the old instances and terminates the new ones. There's never any service interruption. For details, see How traffic-splitting deployments work. 

The Application deployments section of the Rolling updates and deployments page has the following options for application deployments:
Deployment policy

- All at once – Deploy the new version to all instances simultaneously. All instances in your environment are out of service for a short time while the deployment occurs.
- Rolling – Deploy the new version in batches. Each batch is taken out of service during the deployment phase, reducing your environment's capacity by the number of instances in a batch.
- Rolling with additional batch – Deploy the new version in batches, but first launch a new batch of instances to ensure full capacity during the deployment process.
- Immutable – Deploy the new version to a fresh group of instances by performing an immutable update.
- Traffic splitting – Deploy the new version to a fresh group of instances and temporarily split incoming client traffic between the existing application version and the new one.

For the Rolling and Rolling with additional batch deployment policies you can configure:

- Batch size – The size of the set of instances to deploy in each batch.  
  Choose Percentage to configure a percentage of the total number of EC2 instances in the Auto Scaling group (up to 100 percent), or choose Fixed to configure a fixed number of instances (up to the maximum instance count in your environment's Auto Scaling configuration).

For the Traffic splitting deployment policy you can configure the following:

- Traffic split – The initial percentage of incoming client traffic that Elastic Beanstalk shifts to environment instances running the new application version you're deploying.  
  Traffic splitting evaluation time – The time period, in minutes, that Elastic Beanstalk waits after an initial healthy deployment before proceeding to shift all incoming client traffic to the new application version that you're deploying.

### Elastic Container Service

<https://docs.aws.amazon.com/AmazonECS/latest/developerguide/Welcome.html>
Amazon Elastic Container Service (Amazon ECS) is a highly scalable, fast container management service that makes it easy to run, stop, and manage containers on a cluster. Your containers are defined in a task definition that you use to run individual tasks or tasks within a service. In this context, a service is a configuration that enables you to run and maintain a specified number of tasks simultaneously in a cluster. You can run your tasks and services on a serverless infrastructure that is managed by AWS Fargate. Alternatively, for more control over your infrastructure, you can run your tasks and services on a cluster of Amazon EC2 instances that you manage.

Amazon ECS enables you to launch and stop your container-based applications by using simple API calls. You can also retrieve the state of your cluster from a centralized service and have access to many familiar Amazon EC2 features.

You can schedule the placement of your containers across your cluster based on your resource needs, isolation policies, and availability requirements. With Amazon ECS, you don't have to operate your own cluster management and configuration management systems or worry about scaling your management infrastructure.

Amazon ECS can be used to create a consistent build and deployment experience, to manage and scale batch and Extract-Transform-Load (ETL) workloads, and to build sophisticated application architectures on a microservices model. For more information about Amazon ECS use cases and scenarios, see Container Use Cases.

The AWS container services team maintains a public roadmap on GitHub. The roadmap contains information about what the teams are working on and enables AWS customers to provide direct feedback. For more information, see AWS Containers Roadmap.
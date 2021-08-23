# AWS

1. Follow instructions in [Pulumi Get Started Docs on AWS](https://www.pulumi.com/docs/get-started/aws/begin/)
2. [Create a key pair using Amazon EC2](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html#having-ec2-create-your-key-pair)
3. Run `pulumi stack init tidb-cluster-aws` to create a pulumi stack corresponding to `Pulumi.tidb-cluster-aws.yaml` config file. You can skip this if stack already exists.
4. Change `Pulumi.tidb-cluster-aws.yaml` to your desired configuration. Note that `ec2KeyName` should be the same as the name you create in step 2, and `ec2KeyPemFilename` should point to the private key file you downloaded.
5. Run `pulumi up -y` and `./post_script.sh`, and the ec2 instances are ready to go!
    1. The ec2 sshd startup time is unpredictable. If the `./post_script.sh` script failed, please try again later.
    2. You can ssh to each other within the ec2 instances with just their public/private IPs, like `ssh 172.31.1.2`.
    3. You can access the ports from everywhere described in [TiDB Network Requirements](https://docs.pingcap.com/tidb/stable/hardware-and-software-requirements#network-requirements).
6. Get `tiup` on one ec2 instance, spin up the TiDB cluster, run benchmark, and do whatever you like!
    1. There's a `tiup` topology template file named `tidb-cluster.tiup-topology.yaml` in this project, you can start there.
7. Run `pulumi destroy -y` to destroy all the resources

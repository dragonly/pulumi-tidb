package main

import (
	"fmt"
	"io/ioutil"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type ClusterConfig struct {
	Nodes []NodeConfig
}
type NodeConfig struct {
	InstanceType string
	Count        int
}

func main() {
	// publicIps := make([]string, 0)
	// var ec2KeyPemFilename string
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		ec2KeyName := cfg.Require("ec2KeyName")
		// ec2KeyPemFilename = cfg.Require("ec2KeyPemFilename")
		resourceTag := cfg.Require("resourceTag")
		var clusterCfg ClusterConfig
		cfg.RequireObject("cluster", &clusterCfg)

		initScriptBytes, err := ioutil.ReadFile("init_script.sh")
		if err != nil {
			return err
		}
		// fmt.Print(string(initScriptBytes))

		// Create a new security group for port 80.
		ingressArray := make([]ec2.SecurityGroupIngressInput, 0)
		ingressPorts := []int{
			22,          // ssh
			4000, 10080, // tidb
			20160, 20180, // tikv
			2379, 2380, // pd
			9000, 8123, 3930, 20170, 20292, 8234, // tiflash
			8250,       // pump
			8249,       // drainer
			8300,       // ticdc
			9090,       // prometheus
			9100,       // node_exporter
			9115,       // blackbox_exporter
			3000,       // grafana
			9093, 9094, // alertmanager
		}
		for _, port := range ingressPorts {
			ingressArray = append(ingressArray, makeSecGroupIngressArgs(port))
		}
		group, err := ec2.NewSecurityGroup(ctx, "resourceTag", &ec2.SecurityGroupArgs{
			Description: pulumi.StringPtr(resourceTag),
			Ingress:     ec2.SecurityGroupIngressArray(ingressArray),
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(65535),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		})
		if err != nil {
			return err
		}

		// Get the ID for the latest Ubuntu Linux AMI.
		mostRecent := true
		ami, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
			Filters: []ec2.GetAmiFilter{
				// ubuntu 20.04
				{
					Name:   "image-id",
					Values: []string{"ami-03d5c68bab01f3496"},
				},
			},
			Owners:     []string{"099720109477"},
			MostRecent: &mostRecent,
		})
		if err != nil {
			return err
		}

		// instancesFuture := make([]*ec2.Instance, 0)
		for _, node := range clusterCfg.Nodes {
			for count := 0; count < node.Count; count++ {
				// Create a simple web server using the startup script for the instance.
				instanceURN := fmt.Sprintf("%s-%d", node.InstanceType, count)
				instance, err := ec2.NewInstance(ctx, instanceURN, &ec2.InstanceArgs{
					Tags:                pulumi.StringMap{"Name": pulumi.String(resourceTag)},
					InstanceType:        pulumi.String(node.InstanceType),
					VpcSecurityGroupIds: pulumi.StringArray{group.ID()},
					Ami:                 pulumi.String(ami.Id),
					KeyName:             pulumi.String(ec2KeyName),
					// NOTE: UserData seems to be not working
					// 		 Currently just setup the ec2 and networks, manually install tiup and start tidb cluster is ok for now
					// TODO: try AWS OpsWorks / AWS CloudFormation
					UserData: pulumi.String(string(initScriptBytes)),
				})
				if err != nil {
					return err
				}

				// instancesFuture = append(instancesFuture, instance)
				// Export the resulting server's IP address and DNS name.
				ctx.Export(fmt.Sprintf("%s-publicIp", instanceURN), instance.PublicIp)
				ctx.Export(fmt.Sprintf("%s-privateIp", instanceURN), instance.PrivateIp)
			}
		}

		// for _, instance := range instancesFuture {
		// 	output := instance.PublicIp.ApplyT(func(publicIp string) string {
		// 		publicIps = append(publicIps, publicIp)
		// 		return publicIp
		// 	})
		// 	pulumi.Printf("publicIps: %s\n", output)
		// }

		return nil
	})

	// for _, publicIp := range publicIps {
	// 	text := fmt.Sprintf("scp -i %s %s ubuntu@%s:/home/ubuntu/.ssh/id_rsa", ec2KeyPemFilename, ec2KeyPemFilename, publicIp)
	// 	cmd := exec.Command(text)
	// 	cmd.Stdout = os.Stdout
	// 	fmt.Printf("cmd.Run(): %v\n", cmd.Run())
	// }
}

func makeSecGroupIngressArgs(port int) ec2.SecurityGroupIngressArgs {
	return ec2.SecurityGroupIngressArgs{
		Protocol:   pulumi.String("tcp"),
		FromPort:   pulumi.Int(port),
		ToPort:     pulumi.Int(port),
		CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
	}
}

# TiDB ❤️ Pulumi

Deploy TiDB with Pulumi effortlessly.

It should be easy to spin up some virtual machines, and deploy a TiDB cluster there for development & testing.

This project is a work in progress. Contributions are appreciated!

# Instructions

- [AWS](./aws/README.md)

# Cheat Sheet

- Run `pulumi stack output` to check the public / private IPs of ec2 instances.
- Run `ssh -i /path/to/key/pem ubuntu@ip-address` to access each ec2 instance.
- Run `scp -i /path/to/key/pem` to copy files between your local machine and ec2 instances.
- Instances could be inspected at [AWS Console](https://us-west-2.console.aws.amazon.com/ec2/v2/home?region=us-west-2#Instances:).

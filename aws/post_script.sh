#!/bin/bash

pemFilename=$(pulumi config get ec2KeyPemFilename)
nodes=$(pulumi config get --path 'cluster.nodes')
node_num=$(echo $nodes | jq 'length')

for (( n=0; n<$node_num; n++ ))
do
    node=$(echo $nodes | jq ".[$n]")
    count=$(echo $node | jq '.count')
    instanceType=$(echo $node | jq -r '.instanceType')
    for (( i=0; i<$count; i++ ))
    do
        ip=$(pulumi stack output "$instanceType-$i-publicIp")
        scp -o "StrictHostKeyChecking no" -i $pemFilename $pemFilename ubuntu@$ip:/home/ubuntu/.ssh/id_rsa
    done
done

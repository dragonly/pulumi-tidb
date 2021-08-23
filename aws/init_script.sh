#cloud-boothook
#!/bin/bash
sudo apt update -y
sudo apt install curl -y
curl --proto '=https' --tlsv1.2 -sSf https://tiup-mirrors.pingcap.com/install.sh | sh
nohup $HOME/.tiup/bin/tiup playground --host 0.0.0.0 --tiflash 0 2>&1 > tiup.log &

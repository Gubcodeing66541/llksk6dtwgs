# 安装docker
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun

# 源码目录
cd /home/code/server/docker  [code是代码目录]
docker compose up -d

docker ps

## https
yum install -y epel-release certbot

cd /home/server/docker/nginx


sudo certbot certonly --standalone --preferred-challenges http -d fpn.agorns.com \
--config-dir ./ssl/fpn.agorns.com --cert-name fpn.agorns.com	
1.adambuckland43@gmail.com
2.y
3.y

# email:
adambuckland43@gmail.com

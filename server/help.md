#### 安装前
开放6379  3306端口

sudo yum update -y

#### 安装docker
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun

#### 安装docker
sudo yum-config-manager \
--add-repo \
https://download.docker.com/linux/centos/docker-ce.repo


systemctl start docker

#### 安装mysql
docker pull mysql:latest
docker run -itd --name mysql-server -p 3306:3306 -e MYSQL_ROOT_PASSWORD=laimi151813 mysql

#### 安装redis
docker run -p 6379:6379 --name redis -d redis:latest --requirepass "laimi151813"

#### docker 进入mysql 创建数据库`
docker exec -it mysql-server /bin/bash
mysql -h 127.0.0.1 -u root -p
laimi151813   
create database chat;
exit;
按curl + q + p  退出容器



wget https://s3.amazonaws.com/bitly-downloads/nsq/nsq-1.2.0.linux-amd64.go1.12.9.tar.gz
tar -xf nsq-1.2.0.linux-amd64.go1.12.9.tar.gz
cd nsq-1.2.0.linux-amd64.go1.12.9/bin/


# 启动nsqlookupd  nsqd nsqadmin
nohup ./nsqlookupd ./nsqlookupd > nsqlookupd.log 2>&1 &
nohup ./nsqd --lookupd-tcp-address=127.0.0.1:4160  > nsqd.log 2>&1 &
nohup ./nsqadmin --lookupd-http-address=127.0.0.1:4161  > nsqadmin.log 2>&1 &
 

#### 修改config配置文件的IP
vi config.json

#### 启动服务
chmod -R 777 server & nohub ./server &

#### 日志命令行安装
go get github.com/go-swagger/go-swagger
cd Gopath/pkg/mod/github.com/go-swagger
go install github.com/swaggo/swag/cmd/swag

go get -u github.com/swaggo/swag/cmd/swag

swagger version


#### 注释日志生成
swag init -g ./server.go

#### 打包提交代码
set goos=linux
go build server.go
git add .
git commit -m 'u'
git push

#### 只做重启或者启动
cd /home/laimi_group
chmod -R 777 server
chmod -R 777 tools
ps -ef | grep ./server | grep -v grep| awk '{print $2}' | xargs kill -9
ps -ef | grep ./tools | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./server > server.log 2>&1 &
nohup ./tools > tools.log 2>&1 &

#### 重启或者更新代码
cd /home/laimi_group
rm -rf server
rm -rf tools
git reset --hard HEAD && git pull
chmod -R 777 server
chmod -R 777 tools

#### 服务器重启
cd /home/laimi_group/server
ps -ef | grep ./server | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./server > server.log 2>&1 &

#### 定时任务重启 注意只能主服务器启动 其他服务区不能启动
ps -ef | grep ./tools | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./tools > tools.log 2>&1 &

#### 查看请求日志
tail -fn50 server.log


#### 整个服务器重启处理
systemctl start docker
docker start redis
ps -ef | grep ./tools | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./tools > tools.log 2>&1 &
## 以上四个命令只能在主服务器执行（前两句是重启docker和redis,后两句是重启定时任务）
cd /home
ps -ef | grep ./server | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./server > server.log 2>&1 &
## 以上三个命令是重启服务器（主服务器和副服务器都要执行）


#### nginx配置和https
yum install nginx
systemctl start nginx
server {
    #开启gzip
    gzip  on;  
    #低于1kb的资源不压缩
    gzip_min_length 1k;
    #压缩级别1-9，越大压缩率越高，同时消耗cpu资源也越多，建议设置在5左右。
    gzip_comp_level 9;
    #需要压缩哪些响应类型的资源，多个空格隔开。不建议压缩图片.
    gzip_types text/plain application/javascript application/x-javascript text/javascript text/xml text/css;  
    #配置禁用gzip条件，支持正则。此处表示ie6及以下不启用gzip（因为ie低版本不支持）
    gzip_disable "MSIE [1-6]\.";  
    #是否添加“Vary: Accept-Encoding”响应头
    gzip_vary on;
    client_max_body_size 20m;
    client_body_buffer_size 512k;
    client_header_buffer_size 2k;

   location / {
        proxy_pass http://localhost:81;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
   }
}



yum install -y epel-release
yum install -y certbot
yum install -y python3-certbot-nginx
yum install certbot python-certbot-nginx
certbot --nginx -d voice.iutye.com
server {
    listen 80;
    server_name api.rtasuygd.xyz;
    location / {
        proxy_pass http://localhost:81;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}


#### 仅仅更新 不需要再次申请证书
certbot --nginx  certonly

# live
mkdir ./live
cd /live

docker run --rm -it -v$PWD:/output livekit/generate
sudo ./init_script.sh
docker compose up -d

端口开通 7880-7880 80 443








# 安装Certbot和Nginx（如果尚未安装）
sudo apt-get update
sudo apt-get install software-properties-common
sudo add-apt-repository universe
sudo add-apt-repository ppa:certbot/certbot
sudo apt-get update
sudo apt-get install certbot python3-certbot-nginx nginx

# 运行Certbot以自动获取证书
sudo certbot --nginx

# 如果你想手动指定域名，可以使用以下命令
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com


set GOOS=linux
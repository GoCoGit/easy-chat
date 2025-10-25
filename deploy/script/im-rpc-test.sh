#!/bin/bash
reso_addr='crpi-fj5aof5qiib3zd9p.cn-beijing.personal.cr.aliyuncs.com/easy-chat-goco/im-rpc-dev'
tag='latest'

pod_ip="192.168.18.48"

container_name="easy-chat-im-rpc-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-im -v /easy-im/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 10002:10002 -e POD_IP=${pod_ip} --name=${container_name} -d ${reso_addr}:${tag}
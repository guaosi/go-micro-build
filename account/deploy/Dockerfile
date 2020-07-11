FROM centos:7

ADD bin/account /

RUN chmod 777 /account

# 通过先设置一个环境变量，然后在容器运行时传入环境变量具体的值，达到外部指定参数运行的效果

ENV PARAMS=""

ENTRYPOINT ["sh","-c","/account $PARAMS"]
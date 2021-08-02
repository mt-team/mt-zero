FROM mt-zero:v1

ARG appname
ENV TZ Asia/Shanghai
ENV APP_NAME default
ENV ENV_TAG local

ADD ./bin/$appname /app/
ADD ./src/$appname/etc/* /app/
ADD ./run.sh /app/

ENTRYPOINT ./run.sh $APP_NAME $ENV_TAG

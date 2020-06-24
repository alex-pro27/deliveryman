STATIC_ROOT=$(cat $DELIVERY_API_CONF | awk '/STATIC_ROOT:(.*)/{print $2}')
SCHEMA_PATH=$STATIC_ROOT/schema
[ ! -d $SCHEMA_PATH ] && mkdir -p $SCHEMA_PATH
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--swagger_out=logtostderr=true:$SCHEMA_PATH \
--grpc-gateway_out=logtostderr=true:$GOPATH/src \
--go_out=plugins=grpc:$GOPATH/src \
--proto_path=./proto $1.proto
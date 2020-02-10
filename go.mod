module github.com/videocoin/cloud-users

go 1.12

require (
	cloud.google.com/go v0.37.4 // indirect
	github.com/dchest/authcookie v0.0.0-20120917135355-fbdef6e99866
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/locales v0.12.1
	github.com/go-playground/universal-translator v0.16.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/jinzhu/gorm v1.9.12
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/videocoin/cloud-api v0.3.0
	github.com/videocoin/cloud-pkg v0.0.6
	github.com/videocoin/videocoinapis-admin v0.1.1 // indirect
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd
	google.golang.org/grpc v1.21.1
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.0
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace github.com/videocoin/cloud-api => ../cloud-api

replace github.com/videocoin/cloud-pkg => ../cloud-pkg

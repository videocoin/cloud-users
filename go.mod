module github.com/videocoin/cloud-users

go 1.12

require (
	github.com/dchest/authcookie v0.0.0-20120917135355-fbdef6e99866
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/locales v0.12.1
	github.com/go-playground/universal-translator v0.16.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/jinzhu/gorm v1.9.9
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/videocoin/cloud-api v0.1.178
	github.com/videocoin/cloud-pkg v0.0.5
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	google.golang.org/grpc v1.21.1
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.0
)

// replace github.com/videocoin/cloud-api => ../cloud-api
// replace github.com/videocoin/cloud-pkg => ../cloud-pkg

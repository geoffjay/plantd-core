module github.com/geoffjay/plantd/proxy

go 1.14

require (
	github.com/geoffjay/plantd/core v0.0.0-20220321025006-41732bbf807b
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/validator/v10 v10.10.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/ugorji/go v1.2.7 // indirect
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/geoffjay/plantd/core => ../core

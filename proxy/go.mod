module github.com/geoffjay/plantd/proxy

go 1.14

require (
	github.com/geoffjay/plantd/core v0.0.0-20220321025006-41732bbf807b
	github.com/gin-gonic/gin v1.9.0
	github.com/sirupsen/logrus v1.8.1
)

replace github.com/geoffjay/plantd/core => ../core

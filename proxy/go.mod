module github.com/geoffjay/plantd/proxy

go 1.14

require (
	github.com/geoffjay/plantd/core v0.0.0-20220206064530-0e8772b4ac8e
	github.com/gin-gonic/gin v1.7.7
	github.com/sirupsen/logrus v1.8.1
)

replace (
	github.com/geoffjay/plantd/core => ../core
)

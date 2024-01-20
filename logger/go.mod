module github.com/geoffjay/plantd/logger

go 1.14

require (
	github.com/geoffjay/plantd/core v0.0.0-20220326162816-f51114c82ae5
	github.com/lib/pq v1.10.4
	github.com/sirupsen/logrus v1.8.1
)

replace github.com/geoffjay/plantd/core => ../core

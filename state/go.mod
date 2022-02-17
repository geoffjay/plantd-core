module github.com/geoffjay/plantd-core/state

go 1.14

require (
	github.com/geoffjay/plantd-core/core v0.0.0-20220208051425-4545e2bd9a55
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.2.2
	go.etcd.io/bbolt v1.3.6
)

replace github.com/geoffjay/plantd-core/core => ../core

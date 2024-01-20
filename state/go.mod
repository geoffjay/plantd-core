module github.com/geoffjay/plantd/state

go 1.14

require (
	github.com/geoffjay/plantd/core v0.0.0-20220321025006-41732bbf807b
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	go.etcd.io/bbolt v1.3.6
)

replace github.com/geoffjay/plantd/core => ../core

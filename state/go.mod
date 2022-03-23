module github.com/geoffjay/plantd/state

go 1.14

require (
	github.com/geoffjay/plantd/core v0.0.0-20220321025006-41732bbf807b
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.8.2 // indirect
	github.com/stretchr/testify v1.7.1
	go.etcd.io/bbolt v1.3.6
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
)

replace github.com/geoffjay/plantd/core => ../core

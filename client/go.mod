module github.com/geoffjay/plantd/client

require (
	github.com/geoffjay/plantd/core v0.0.0-20220208051425-4545e2bd9a55
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
)

replace github.com/geoffjay/plantd/core => ../core

go 1.14

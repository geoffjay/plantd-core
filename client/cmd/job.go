package cmd

// import (
// 	"fmt"
// 	"log"
// 	"strings"

// 	plantd "github.com/geoffjay/plantd/core/service"

// 	"github.com/spf13/cobra"
// )

// var (
// 	jobName       string
// 	jobValue      string
// 	jobProperties string

// 	jobCmd = &cobra.Command{
// 		Use:   "job",
// 		Short: "Perform job functions with plantd",
// 		Args:  cobra.ExactArgs(1),
// 		Run:   job,
// 	}
// )

// func init() {
// 	jobCmd.PersistentFlags().StringVarP(&jobName, "name", "n", "", "job name")
// 	jobCmd.PersistentFlags().StringVarP(&jobValue, "value", "", "", "job value")
// 	jobCmd.PersistentFlags().StringVarP(&jobProperties, "properties", "p", "", "job properties")
// }

// func job(cmd *cobra.Command, args []string) {
// 	client, err := plantd.NewClient(endpoint)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	properties := make(map[string]string)
// 	if jobProperties != "" {
// 		list := strings.Split(jobProperties, ",")
// 		for _, property := range list {
// 			kv := strings.Split(property, "=")
// 			properties[kv[0]] = kv[1]
// 		}
// 	}

// 	id := args[0]
// 	response, err := client.SubmitJob(id, jobName, jobValue, properties)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(response.String())
// }

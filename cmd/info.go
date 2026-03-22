package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/scttfrdmn/ood-sagemaker-training-adapter/internal/sagemakertraining"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <job-name>",
	Short: "Print full SageMaker training job details as JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := sagemakertraining.New(ctx, region)
		if err != nil {
			return err
		}
		detail, err := client.DescribeTrainingJob(ctx, args[0])
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(detail)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

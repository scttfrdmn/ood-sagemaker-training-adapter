package cmd

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/ood-sagemaker-training-adapter/internal/sagemakertraining"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <job-name>",
	Short: "Stop a SageMaker training job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := sagemakertraining.New(ctx, region)
		if err != nil {
			return err
		}
		if err := client.StopTrainingJob(ctx, args[0]); err != nil {
			return err
		}
		fmt.Printf("Training job %s stopped\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

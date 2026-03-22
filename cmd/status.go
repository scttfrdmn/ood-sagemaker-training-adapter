package cmd

import (
	"context"
	"encoding/json"
	"os"

	smtypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	internalood "github.com/scttfrdmn/ood-sagemaker-training-adapter/internal/ood"
	"github.com/scttfrdmn/ood-sagemaker-training-adapter/internal/sagemakertraining"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <job-name>",
	Short: "Get the status of a SageMaker training job",
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

		js := internalood.JobStatus{
			ID:     args[0],
			Status: smStateToOod(detail.TrainingJobStatus),
		}
		if detail.FailureReason != nil {
			js.Message = *detail.FailureReason
		}

		return json.NewEncoder(os.Stdout).Encode(js)
	},
}

func smStateToOod(s smtypes.TrainingJobStatus) string {
	switch s {
	case smtypes.TrainingJobStatusInProgress, smtypes.TrainingJobStatusStopping:
		return internalood.StatusRunning
	case smtypes.TrainingJobStatusCompleted:
		return internalood.StatusCompleted
	case smtypes.TrainingJobStatusFailed:
		return internalood.StatusFailed
	case smtypes.TrainingJobStatusStopped:
		return internalood.StatusCancelled
	default:
		return internalood.StatusUnknown
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

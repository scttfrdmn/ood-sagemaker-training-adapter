package cmd

import (
	"github.com/spf13/cobra"
)

var (
	region   string
	roleArn  string
	outputS3 string
)

var rootCmd = &cobra.Command{
	Use:   "ood-sagemaker-training-adapter",
	Short: "OOD compute adapter for AWS SageMaker Training Jobs",
	Long:  "Translates Open OnDemand job submissions to AWS SageMaker Training Jobs API calls.",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region")
	rootCmd.PersistentFlags().StringVar(&roleArn, "role-arn", "", "IAM role ARN for the training job execution")
	rootCmd.PersistentFlags().StringVar(&outputS3, "output-s3", "", "S3 URI for training job output artifacts")
}

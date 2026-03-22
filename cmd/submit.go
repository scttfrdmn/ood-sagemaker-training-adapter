package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scttfrdmn/ood-sagemaker-training-adapter/internal/sagemakertraining"
	"github.com/spf13/cobra"
)

// JobSpec is the SageMaker training-specific job submission payload.
type JobSpec struct {
	AlgorithmImage  string            `json:"algorithm_image"`
	RoleArn         string            `json:"role_arn,omitempty"`
	InputS3         string            `json:"input_s3,omitempty"`
	OutputS3        string            `json:"output_s3,omitempty"`
	InstanceType    string            `json:"instance_type"`
	InstanceCount   int               `json:"instance_count,omitempty"`
	VolumeSizeGB    int               `json:"volume_size_gb,omitempty"`
	Walltime        string            `json:"walltime,omitempty"`
	Hyperparameters map[string]string `json:"hyperparameters,omitempty"`
	Env             map[string]string `json:"env,omitempty"`
	JobName         string            `json:"job_name"`
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit an OOD job to SageMaker Training",
	Long:  "Reads a JSON job spec from stdin and creates a SageMaker training job.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var spec JobSpec
		if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
			return fmt.Errorf("decode job spec: %w", err)
		}

		if spec.JobName == "" {
			return fmt.Errorf("job spec must include job_name")
		}
		if spec.AlgorithmImage == "" {
			return fmt.Errorf("job spec must include algorithm_image")
		}
		if spec.InstanceType == "" {
			return fmt.Errorf("job spec must include instance_type")
		}

		effectiveRole := roleArn
		if effectiveRole == "" {
			effectiveRole = spec.RoleArn
		}
		if effectiveRole == "" {
			return fmt.Errorf("--role-arn is required (or set role_arn in job spec)")
		}

		effectiveOutput := outputS3
		if effectiveOutput == "" {
			effectiveOutput = spec.OutputS3
		}
		if effectiveOutput == "" {
			return fmt.Errorf("--output-s3 is required (or set output_s3 in job spec)")
		}

		ctx := context.Background()
		client, err := sagemakertraining.New(ctx, region)
		if err != nil {
			return err
		}

		jobName, err := client.CreateTrainingJob(ctx, sagemakertraining.TrainingJobSpec{
			AlgorithmImage:  spec.AlgorithmImage,
			RoleArn:         effectiveRole,
			InputS3:         spec.InputS3,
			OutputS3:        effectiveOutput,
			InstanceType:    spec.InstanceType,
			InstanceCount:   int32(spec.InstanceCount),
			VolumeSizeGB:    int32(spec.VolumeSizeGB),
			Walltime:        spec.Walltime,
			Hyperparameters: spec.Hyperparameters,
			Env:             spec.Env,
			JobName:         spec.JobName,
		})
		if err != nil {
			return err
		}

		fmt.Println(jobName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}

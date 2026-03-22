// Package sagemakertraining wraps the SageMaker Training Jobs API for the OOD adapter.
package sagemakertraining

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

// Client wraps the AWS SageMaker client.
type Client struct {
	svc    *sagemaker.Client
	region string
}

// New creates a SageMaker client using the default AWS credential chain.
func New(ctx context.Context, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &Client{svc: sagemaker.NewFromConfig(cfg), region: region}, nil
}

// TrainingJobSpec holds the parameters for a SageMaker training job.
type TrainingJobSpec struct {
	AlgorithmImage  string
	RoleArn         string
	InputS3         string
	OutputS3        string
	InstanceType    string
	InstanceCount   int32
	VolumeSizeGB    int32
	Walltime        string // HH:MM:SS
	Hyperparameters map[string]string
	Env             map[string]string
	JobName         string
}

// walltimeToSeconds converts HH:MM:SS to seconds. Returns 86400 on parse failure.
func walltimeToSeconds(walltime string) int32 {
	parts := strings.Split(walltime, ":")
	if len(parts) != 3 {
		return 86400
	}
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	s, _ := strconv.Atoi(parts[2])
	total := h*3600 + m*60 + s
	if total <= 0 {
		return 86400
	}
	return int32(total)
}

// CreateTrainingJob submits a SageMaker training job and returns the job name (which is the ID).
func (c *Client) CreateTrainingJob(ctx context.Context, spec TrainingJobSpec) (string, error) {
	instanceCount := spec.InstanceCount
	if instanceCount <= 0 {
		instanceCount = 1
	}
	volumeSize := spec.VolumeSizeGB
	if volumeSize <= 0 {
		volumeSize = 30
	}

	input := &sagemaker.CreateTrainingJobInput{
		TrainingJobName: aws.String(spec.JobName),
		RoleArn:         aws.String(spec.RoleArn),
		AlgorithmSpecification: &types.AlgorithmSpecification{
			TrainingImage:     aws.String(spec.AlgorithmImage),
			TrainingInputMode: types.TrainingInputModeFile,
		},
		ResourceConfig: &types.ResourceConfig{
			InstanceType:   types.TrainingInstanceType(spec.InstanceType),
			InstanceCount:  aws.Int32(instanceCount),
			VolumeSizeInGB: aws.Int32(volumeSize),
		},
		OutputDataConfig: &types.OutputDataConfig{
			S3OutputPath: aws.String(spec.OutputS3),
		},
		StoppingCondition: &types.StoppingCondition{
			MaxRuntimeInSeconds: aws.Int32(walltimeToSeconds(spec.Walltime)),
		},
	}

	if spec.InputS3 != "" {
		input.InputDataConfig = []types.Channel{
			{
				ChannelName: aws.String("training"),
				DataSource: &types.DataSource{
					S3DataSource: &types.S3DataSource{
						S3Uri:                  aws.String(spec.InputS3),
						S3DataType:             types.S3DataTypeS3Prefix,
						S3DataDistributionType: types.S3DataDistributionFullyReplicated,
					},
				},
			},
		}
	}

	if len(spec.Hyperparameters) > 0 {
		input.HyperParameters = spec.Hyperparameters
	}
	if len(spec.Env) > 0 {
		input.Environment = spec.Env
	}

	_, err := c.svc.CreateTrainingJob(ctx, input)
	if err != nil {
		return "", fmt.Errorf("sagemaker CreateTrainingJob: %w", err)
	}
	return spec.JobName, nil
}

// DescribeTrainingJob returns the current detail of a SageMaker training job.
func (c *Client) DescribeTrainingJob(ctx context.Context, jobName string) (*sagemaker.DescribeTrainingJobOutput, error) {
	out, err := c.svc.DescribeTrainingJob(ctx, &sagemaker.DescribeTrainingJobInput{
		TrainingJobName: aws.String(jobName),
	})
	if err != nil {
		return nil, fmt.Errorf("sagemaker DescribeTrainingJob: %w", err)
	}
	return out, nil
}

// StopTrainingJob stops a SageMaker training job.
func (c *Client) StopTrainingJob(ctx context.Context, jobName string) error {
	_, err := c.svc.StopTrainingJob(ctx, &sagemaker.StopTrainingJobInput{
		TrainingJobName: aws.String(jobName),
	})
	if err != nil {
		return fmt.Errorf("sagemaker StopTrainingJob: %w", err)
	}
	return nil
}

//go:build integration

package sagemakertraining_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	substrate "github.com/scttfrdmn/substrate"

	. "github.com/scttfrdmn/ood-sagemaker-training-adapter/internal/sagemakertraining"
)

// TestCreateDescribeStopTrainingJob_Substrate exercises the SageMaker training job
// lifecycle (CreateTrainingJob → DescribeTrainingJob → StopTrainingJob) against
// the substrate emulator.
func TestCreateDescribeStopTrainingJob_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	jobName := fmt.Sprintf("ood-training-%d", time.Now().UnixNano())

	spec := TrainingJobSpec{
		JobName:        jobName,
		AlgorithmImage: "123456789012.dkr.ecr.us-east-1.amazonaws.com/my-training:latest",
		RoleArn:        "arn:aws:iam::123456789012:role/SageMakerTrainingRole",
		InputS3:        "s3://my-bucket/training-data/",
		OutputS3:       "s3://my-bucket/model-output/",
		InstanceType:   "ml.m5.xlarge",
		InstanceCount:  1,
		VolumeSizeGB:   30,
		Walltime:       "02:00:00",
		Hyperparameters: map[string]string{
			"epochs":    "10",
			"lr":        "0.001",
		},
	}

	// CreateTrainingJob — returns the job name as ID
	returnedName, err := client.CreateTrainingJob(ctx, spec)
	if err != nil {
		t.Fatalf("CreateTrainingJob: %v", err)
	}
	if returnedName != jobName {
		t.Errorf("CreateTrainingJob: got name %q, want %q", returnedName, jobName)
	}
	t.Logf("created training job: %s", returnedName)

	// DescribeTrainingJob
	detail, err := client.DescribeTrainingJob(ctx, jobName)
	if err != nil {
		t.Fatalf("DescribeTrainingJob: %v", err)
	}
	if detail == nil {
		t.Fatal("DescribeTrainingJob: got nil output")
	}
	t.Logf("training job status: %s", detail.TrainingJobStatus)

	// StopTrainingJob
	err = client.StopTrainingJob(ctx, jobName)
	if err != nil {
		t.Fatalf("StopTrainingJob: %v", err)
	}
	t.Log("training job stopped successfully")
}

// TestDescribeTrainingJob_NotFound_Substrate verifies that DescribeTrainingJob
// returns an error for a job that was never created.
func TestDescribeTrainingJob_NotFound_Substrate(t *testing.T) {
	ts := substrate.StartTestServer(t)
	t.Setenv("AWS_ENDPOINT_URL", ts.URL)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()
	client, err := New(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = client.DescribeTrainingJob(ctx, "no-such-job-ever")
	if err == nil {
		t.Fatal("expected error for non-existent training job, got nil")
	}
	if !strings.Contains(err.Error(), "sagemaker") {
		t.Logf("error (acceptable): %v", err)
	}
}

package cmd

import (
	"testing"

	smtypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

func TestSmStateToOod(t *testing.T) {
	tests := []struct {
		state    smtypes.TrainingJobStatus
		expected string
	}{
		{smtypes.TrainingJobStatusInProgress, "running"},
		{smtypes.TrainingJobStatusStopping, "running"},
		{smtypes.TrainingJobStatusCompleted, "completed"},
		{smtypes.TrainingJobStatusFailed, "failed"},
		{smtypes.TrainingJobStatusStopped, "cancelled"},
		{smtypes.TrainingJobStatus("UNKNOWN_STATE"), "undetermined"},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			got := smStateToOod(tt.state)
			if got != tt.expected {
				t.Errorf("smStateToOod(%q) = %q, want %q", tt.state, got, tt.expected)
			}
		})
	}
}

# ood-sagemaker-training-adapter

An [Open OnDemand](https://openondemand.org/) compute adapter that translates OOD job submissions into [AWS SageMaker Training Jobs](https://docs.aws.amazon.com/sagemaker/latest/dg/how-it-works-training.html).

> **Distinct from `ood-sagemaker-adapter`** — that adapter targets SageMaker Studio apps (interactive sessions). This adapter targets batch training workloads via the SageMaker Training Jobs API.

## Commands

| Command | Description |
|---------|-------------|
| `submit` | Read a JSON job spec from stdin and create a SageMaker training job |
| `status <job-name>` | Return OOD-normalized status JSON for a training job |
| `delete <job-name>` | Stop a running training job |
| `info <job-name>` | Print full SageMaker job detail as JSON |

## Job spec (stdin for `submit`)

```json
{
  "job_name": "my-training-job-001",
  "algorithm_image": "763104351884.dkr.ecr.us-east-1.amazonaws.com/pytorch-training:2.1.0-gpu-py310",
  "instance_type": "ml.p3.2xlarge",
  "instance_count": 1,
  "volume_size_gb": 50,
  "walltime": "04:00:00",
  "input_s3": "s3://my-bucket/input/",
  "output_s3": "s3://my-bucket/output/",
  "role_arn": "arn:aws:iam::123456789012:role/SageMakerRole",
  "hyperparameters": {
    "epochs": "10",
    "lr": "0.001"
  },
  "env": {
    "MY_VAR": "value"
  }
}
```

`role_arn` and `output_s3` may also be supplied as CLI flags (`--role-arn`, `--output-s3`); flags take precedence over the spec.

## Global flags

```
--region      AWS region (default: us-east-1)
--role-arn    IAM execution role ARN
--output-s3   S3 URI for output artifacts
```

## Status mapping

| SageMaker status | OOD status |
|------------------|------------|
| InProgress, Stopping | running |
| Completed | completed |
| Failed | failed |
| Stopped | cancelled |
| (other) | undetermined |

## Building

```bash
go build -o ood-sagemaker-training-adapter .
```

## Authentication

Uses the standard AWS credential chain (environment variables, `~/.aws/credentials`, EC2/ECS instance metadata, etc.).

## License

MIT — see [LICENSE](LICENSE).

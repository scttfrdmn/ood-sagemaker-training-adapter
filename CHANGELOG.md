# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial scaffold — OOD compute adapter for AWS SageMaker Training Jobs, translating Open OnDemand job submissions into SageMaker training jobs.
- CLI commands: `submit` (JSON job spec from stdin → SageMaker training job), `status <job-name>` (OOD-normalized status), `delete <job-name>` (stop a training job), and `info <job-name>` (full SageMaker job detail as JSON).
- Unit tests for status state mapping.
- Substrate integration tests for the SageMaker training job lifecycle.
- CI workflow with pinned action SHAs.

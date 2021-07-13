# actions-plan-preview

An action that comments PipeCD's PlanPreview result on GitHub pull request. This action can be used for all application kinds: Kubernetes, Terraform, CloudRun, Lambda, Amazon ECS.

See https://pipecd.dev/docs/user-guide/plan-preview/ for more details about this feature.

## Screenshots

- Plan-preview comment shows the diff of an Kubernetes application

![](/assets/plan-preview-comment-kubernetes.png)

- Plan-preview comment shows the diff of an Terraform application

![](/assets/plan-preview-comment-terraform.png)

## Usage

By adding the following workflow to your `.github` directory (eg: `.github/workflows/plan-preview.yaml`) you will have:

- The result of plan-preview will be automatically commented on pull request when it is opened or updated 
- You can leave a `/pipecd plan-preview` comment on pull request to trigger a plan-preview manually


``` yaml
on:
  pull_request:
    branches:
      - main
    types: [opened, synchronize, reopened]
  issue_comment:
    types: [created]

jobs:
  plan-preview:
    runs-on: ubuntu-latest
    if: "github.event_name == 'pull_request'"
    steps:
      - uses: pipe-cd/actions-plan-preview@v1.0.0
        with:
          address: ${{ secrets.PIPECD_ADDRESS }}
          api-key: ${{ secrets.PIPECD_PLAN_PREVIEW_API_KEY }}
          token: ${{ secrets.GITHUB_TOKEN }}

  plan-preview-on-comment:
    runs-on: ubuntu-latest
    if: "github.event_name == 'issue_comment' && github.event.issue.pull_request && startsWith(github.event.comment.body, '/pipecd plan-preview')"
    steps:
      - uses: pipe-cd/actions-plan-preview@v1.0.0
        with:
          address: ${{ secrets.PIPECD_ADDRESS }}
          api-key: ${{ secrets.PIPECD_PLAN_PREVIEW_API_KEY }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Name                            | Description                                                                                       | Required | Default Value |
|---------------------------------|---------------------------------------------------------------------------------------------------|:--------:|:-------------:|
| address                         | The address of PipeCD's control-plane.                                                            |    yes   |               |
| api-key                         | The API key with READ_WRITE role used by pipectl while communicating with PipeCD's control-plane. |    yes   |               |
| token                           | The GITHUB_TOKEN secret.                                                                          |    yes   |               |

## Contributing

Source code for this action is placing at [pipe-cd/pipe](https://github.com/pipe-cd/pipe/tree/master/dockers/actions-plan-preview) repository.
Please send pull request to that repository to update.

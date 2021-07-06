# actions-plan-preview

An action that comments PipeCD's PlanPreview result on GitHub pull request.

![](/assets/plan-preview-comment.png)

## Usage

Adding a new workflow (eg: `.github/workflows/plan-preview.yaml`) with the content as below:

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

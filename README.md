# actions-plan-preview

An action that comments PipeCD's PlanPreview result on GitHub pull request. This action can be used for all application kinds: Kubernetes, Terraform, CloudRun, Lambda, Amazon ECS.

See https://pipecd.dev/docs/user-guide/plan-preview/ for more details about this feature.

**NOTE**: The source code of this GitHub Action is placing under the tool directory of of [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd/tree/master/tool) repository. If you want to make a pull request or raise an issue, please send it to [pipe-cd/pipecd](https://github.com/pipe-cd/pipecd) repository.

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
name: PipeCD

on:
  pull_request:
    branches:
      - main
    types: [opened, synchronize, reopened]
  issue_comment:
    types: [created]

jobs:
  plan-preview:
    name: Plan Preview
    runs-on: ubuntu-latest
    if: "github.event_name == 'pull_request'"
    steps:
      - uses: pipe-cd/actions-plan-preview@v1.7.2
        with:
          address: ${{ secrets.PIPECD_API_ADDRESS }}
          api-key: ${{ secrets.PIPECD_PLAN_PREVIEW_API_KEY }}
          token: ${{ secrets.GITHUB_TOKEN }}

  plan-preview-on-comment:
    name: Plan Preview
    runs-on: ubuntu-latest
    if: "github.event_name == 'issue_comment' && github.event.issue.pull_request && startsWith(github.event.comment.body, '/pipecd plan-preview')"
    steps:
      - uses: pipe-cd/actions-plan-preview@v1.7.2
        with:
          address: ${{ secrets.PIPECD_API_ADDRESS }}
          api-key: ${{ secrets.PIPECD_PLAN_PREVIEW_API_KEY }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

### Push events

To run actions-plan-preview after automatically creating PRs on push events using [GITHUB_TOKEN](https://docs.github.com/en/actions/using-workflows/triggering-a-workflow#triggering-a-workflow-from-a-workflow), it goes as follows.

```yaml
name: PipeCD

on:
  push:
    branches: pr-target-branch
jobs:
  create-pr:
    runs-on: ubuntu-latest
    if: "github.event.created"
    env:
      GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    steps:
      - uses: actions/checkout@v2
      - id: create-pr
        run: |
          PR=$(gh pr create --title "Update" --body "" | awk -F "/" '{print $NF}')
          echo "::set-output name=pull-request-number::$PR"
      - uses: pipe-cd/actions-plan-preview@v1.8.0
        with:
          address: ${{ secrets.PIPECD_API_ADDRESS }}
          api-key: ${{ secrets.PIPECD_PLAN_PREVIEW_API_KEY }}
          token: ${{ secrets.GITHUB_TOKEN }}
          pull-request-number: ${{ steps.create-pr.outputs.pull-request-number }}
```

## Inputs

| Name                            | Description                                                                                       | Required | Default Value |
|---------------------------------|---------------------------------------------------------------------------------------------------|:--------:|:-------------:|
| address                         | The API address of PipeCD's control-plane.                                                        |    yes   |               |
| api-key                         | The API key with READ_WRITE role used by pipectl while communicating with PipeCD's control-plane. |    yes   |               |
| token                           | The GITHUB_TOKEN secret.                                                                          |    yes   |               |
| pull-request-number             | PR Number needed for push event.                                                                  |   false  |               |

# terraform-concurrent-workspace-poc

POC of combining multiple runs of terraform (represented by multiple `terraform-exec` clients) with Terraform CLI's concept of a workspace.

Each `terraform-exec` client is associated with a workspace by setting the workspace on the client via `SetEnv()` with the environment variable `TF_WORKSPACE`.

This allows terraform command lines to be run within different workspaces in parallel.

To run POC:
1. Download terraform binary and put in path
2. `make run`
3. Observe that 'resources' are created in the `instance-east` and `instance-west` directories

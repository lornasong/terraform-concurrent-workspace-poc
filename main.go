package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type TerraformDriver struct {
	workspace string
	client    *tfexec.Terraform
}

func NewTerraformDriver(workingDir string, workspace string) (*TerraformDriver, error) {
	tf, err := tfexec.NewTerraform(workingDir, "")
	if err != nil {
		return nil, err
	}
	env := make(map[string]string)
	env["TF_WORKSPACE"] = workspace
	tf.SetEnv(env)

	return &TerraformDriver{
		workspace: workspace,
		client:    tf,
	}, nil
}

func (td *TerraformDriver) init(ctx context.Context) error {
	fmt.Println("init", td.workspace)
	if err := td.client.Init(ctx); err != nil {
		return err
	}
	return nil
}

func (td *TerraformDriver) apply(ctx context.Context, errchan chan error) {
	fmt.Println("apply", td.workspace)
	if err := td.client.Apply(ctx); err != nil {
		errchan <- err
	}
}

func main() {

	workingDir := "dir"

	// creating a terraform driver for each workspace (east and west)
	east, err := NewTerraformDriver(workingDir, "east")
	if err != nil {
		panic(err)
	}
	west, err := NewTerraformDriver(workingDir, "west")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err = east.init(ctx); err != nil {
		panic(err)
	}
	if err = west.init(ctx); err != nil {
		panic(err)
	}

	errchan := make(chan error)
	go east.apply(ctx, errchan)
	go west.apply(ctx, errchan)

	for err := range errchan {
		fmt.Println(err)
	}
}

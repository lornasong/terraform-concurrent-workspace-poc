package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func main() {

	workingDir := "dir"

	// creating a client for each workspace (east and west)
	east, err := tfexec.NewTerraform(workingDir, "")
	if err != nil {
		panic(err)
	}
	eastEnv := make(map[string]string)
	eastEnv["TF_WORKSPACE"] = "east"
	east.SetEnv(eastEnv)

	west, err := tfexec.NewTerraform(workingDir, "")
	if err != nil {
		panic(err)
	}
	westEnv := make(map[string]string)
	westEnv["TF_WORKSPACE"] = "west"
	west.SetEnv(westEnv)

	// workspaces := []*tfexec.Terraform()
	ctx := context.Background()
	fmt.Println("initializing east client")
	if err = east.Init(ctx); err != nil {
		panic(err)
	}
	fmt.Println("initializing west client")
	if err = west.Init(ctx); err != nil {
		panic(err)
	}

	errchan := make(chan error)
	go apply(ctx, east, errchan)
	go apply(ctx, west, errchan)

	for err := range errchan {
		fmt.Println(err)
	}

	// if err = east.Apply(ctx); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("applying west client")
	// wg.Add(1)

	// if err = west.Apply(ctx); err != nil {
	// 	panic(err)
	// }
}

func apply(ctx context.Context, tf *tfexec.Terraform, errchan chan error) {
	fmt.Println("apply")
	if err := tf.Apply(ctx); err != nil {
		errchan <- err
	}
}

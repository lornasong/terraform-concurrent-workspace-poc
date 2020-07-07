package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func main() {

	// stuff from configuration
	workingDir := "dir"
	instances := []string{"east", "west"}

	// creating terraform drivers for each instance (east and west)
	drivers := make([]*TerraformDriver, len(instances))
	for ix, instance := range instances {
		driver, err := NewTerraformDriver(workingDir, instance)
		if err != nil {
			panic(err)
		}
		drivers[ix] = driver
	}

	ctx := context.Background()

	// initalize the drivers
	for _, driver := range drivers {
		if err := driver.init(ctx); err != nil {
			panic(err)
		}
	}

	// run applies concurrently
	errchan := make(chan error)
	for _, driver := range drivers {
		go driver.apply(ctx, errchan)
	}

	// watch for errors
	for err := range errchan {
		fmt.Println(err)
	}
}

// TerraformDriver is a struct to hold terraform-exec client and instance
// related information.
type TerraformDriver struct {
	instance string
	client   *tfexec.Terraform
}

// NewTerraformDriver creates a driver per instance. This instance is translated into
// terraform CLI's workspace concept
func NewTerraformDriver(workingDir string, instance string) (*TerraformDriver, error) {
	tf, err := tfexec.NewTerraform(workingDir, "")
	if err != nil {
		return nil, err
	}
	env := make(map[string]string)
	env["TF_WORKSPACE"] = instance
	tf.SetEnv(env)

	return &TerraformDriver{
		instance: instance,
		client:   tf,
	}, nil
}

func (td *TerraformDriver) init(ctx context.Context) error {
	fmt.Println("init", td.instance)
	if err := td.client.Init(ctx); err != nil {
		return err
	}
	return nil
}

func (td *TerraformDriver) apply(ctx context.Context, errchan chan error) {
	fmt.Println("apply", td.instance)
	if err := td.client.Apply(ctx); err != nil {
		errchan <- err
	}
}

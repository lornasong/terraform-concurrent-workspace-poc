package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// TerraformWorkspaceEnvVar is set for each terraform driver where
// the value is the instance name. Use workspace to represent an
// instance
const TerraformWorkspaceEnvVar = "TF_WORKSPACE"

func main() {
	// stuff from configuration
	workingDir := "dir"
	instances := []string{"east", "west"}

	// creating terraform drivers for each instance
	drivers, err := NewDrivers(workingDir, instances)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// initalize all the drivers (note: could make this concurrent)
	if err = drivers.InitAll(ctx); err != nil {
		panic(err)
	}

	// run applies concurrently
	drivers.ApplyAll(ctx)

	// watch for errors
	for err := range drivers.ErrChan() {
		fmt.Println(err)
	}
}

// Drivers is a struct to hold all the drivers for all instances
// Has functions to to call action across all drivers.
type Drivers struct {
	drivers []*InstanceDriver
	errchan chan error
}

// NewDrivers creates a Driver struct and creates InstanceDriver for all instances
func NewDrivers(workingDir string, instances []string) (*Drivers, error) {
	drivers := make([]*InstanceDriver, len(instances))
	for ix, instance := range instances {
		driver, err := NewInstanceDriver(workingDir, instance)
		if err != nil {
			return nil, err
		}
		drivers[ix] = driver
	}

	errchan := make(chan error)

	return &Drivers{
		drivers: drivers,
		errchan: errchan,
	}, nil
}

// InitAll calls initialize across all instance drivers
// Could probably make this concurrent too
func (d *Drivers) InitAll(ctx context.Context) error {
	for _, driver := range d.drivers {
		if err := driver.init(ctx); err != nil {
			return err
		}
	}
	return nil
}

// ApplyAll calls apply _concurrently_ across all instance drivers
// Errors are in driver.ErrChan()
func (d *Drivers) ApplyAll(ctx context.Context) {
	for _, driver := range d.drivers {
		go driver.apply(ctx, d.errchan)
	}
}

func (d *Drivers) ErrChan() chan error {
	return d.errchan
}

// InstanceDriver is a struct to hold terraform-exec client for a single instance
type InstanceDriver struct {
	instance string
	client   *tfexec.Terraform
}

// NewInstanceDriver creates a driver per instance. This instance is translated into
// terraform CLI's workspace concept
func NewInstanceDriver(workingDir string, instance string) (*InstanceDriver, error) {
	tf, err := tfexec.NewTerraform(workingDir, "")
	if err != nil {
		return nil, err
	}
	env := make(map[string]string)
	env[TerraformWorkspaceEnvVar] = instance
	tf.SetEnv(env)

	return &InstanceDriver{
		instance: instance,
		client:   tf,
	}, nil
}

func (td *InstanceDriver) init(ctx context.Context) error {
	fmt.Println("init", td.instance)
	if err := td.client.Init(ctx); err != nil {
		return err
	}
	return nil
}

func (td *InstanceDriver) apply(ctx context.Context, errchan chan error) {
	fmt.Println("apply", td.instance)
	if err := td.client.Apply(ctx); err != nil {
		errchan <- err
	}
}

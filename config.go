package main

import (
	"fmt"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"strconv"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	communicator.Config `mapstructure:",squash"`

	Url              string `mapstructure:"url"`
	Username         string `mapstructure:"username"`
	Password         string `mapstructure:"password"`

	Template   string `mapstructure:"template"`
	VMName     string `mapstructure:"vm_name"`
	FolderName string `mapstructure:"folder_name"`
	DCName     string `mapstructure:"dc_name"`

	Cpus            string `mapstructure:"cpus"`
	ShutdownCommand string `mapstructure:"shutdown_command"`
	Ram             string `mapstructure:"RAM"`
	//TODO: add more options

	ctx      interpolate.Context
}

func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := new(Config)
	err := config.Decode(c, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &c.ctx,
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	// Accumulate any errors
	errs := new(packer.MultiError)

	// Prepare config(s)
	errs = packer.MultiErrorAppend(errs, c.Config.Prepare(&c.ctx)...)

	// Check the required params
	if c.Url == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("URL required"))
	}
	if c.Username == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Username required"))
	}
	if c.Password == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Password required"))
	}
	if c.Template == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Template VM name required"))
	}
	if c.VMName == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Target VM name required"))
	}

	// Verify numeric parameters if present
	if c.Cpus != "" {
		if _, err = strconv.Atoi(c.Cpus); err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("Invalid number of cpu sockets"))
		}
	}
	if c.Ram != "" {
		if _, err = strconv.Atoi(c.Ram); err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("Invalid number for Ram"))
		}
	}

	// Warnings
	var warnings []string

	if len(errs.Errors) > 0 {
		return nil, warnings, errs
	}

	return c, warnings, nil
}
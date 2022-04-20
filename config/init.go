package config

import "github.com/vrischmann/envconfig"

const prefix = "VISIAUTH"

type initFunc func(conf interface{}) error

var initFuncs = []initFunc{
	func(conf interface{}) error {
		return envconfig.InitWithPrefix(conf, prefix)
	},
	func(conf interface{}) error {
		return envconfig.Init(conf)
	},
}

func Init(conf interface{}) error {
	var errs []error
	for _, fn := range initFuncs {
		if err := fn(conf); err != nil {
			errs = append(errs, err)
		} else {
			return nil
		}
	}

	return errs[0]
}

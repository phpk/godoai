package cmd

import (
	"godoai/deps"
	"godoai/libs"
	"godoai/progress"
)

func InitSystem() error {

	err := libs.LoadConfig()
	if err != nil {
		return err
	}
	err = libs.InitConfig()
	if err != nil {
		return err
	}
	err = deps.InitDir()
	if err != nil {
		return err
	}
	//	libs.InitLinese()
	// err = libs.InitLinese()
	// if err != nil {
	// 	return err
	// }
	err = progress.StartCmd("ollama")
	if err != nil {
		return err
	}
	return nil
}

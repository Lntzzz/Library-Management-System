package conf

import (
	"path/filepath"

	"github.com/spf13/viper"
)

func newViperWithConfPath(confPath string) *viper.Viper {
	dir, file := filepath.Split(confPath)
	ext := filepath.Ext(file)
	filename := file[:len(file)-len(ext)]

	// Set up Viper
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType(ext[1:]) // Remove the dot from the extension
	if dir == "" {
		v.AddConfigPath(".")
	} else {
		v.AddConfigPath(dir)
	}
	v.SetConfigType("yaml")
	return v
}

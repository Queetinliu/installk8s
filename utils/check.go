package utils
import (
	//"os"
	//"log"
	//"path/filepath"
	//"fmt"
	//"strings"
	//"github.com/spf13/viper"
)

/*
func checkcfgfileexists(file string) (bool,error) {
	curpath, err := os.Getwd()
	if err != nil {
		log.Fatal("can't get current dir")
		return false,err
	}
	fileinfo, err := os.Stat(filepath.Join(curpath, string(file)))
	if err != nil {
		if os.IsNotExist(err){
			err := writesample(file)
			if err != nil {
				log.Fatal("write sample error")
				return false,err

		}
		fmt.Println(fileinfo)
		log.Fatal(file + " not exist.\nwe have create the sampl config.yaml,please edit it and run install again")
		return false,err
		}
	}
	return true,nil
}

func checkpackagesexists(dir string) (bool,error){
	curpath, err := os.Getwd()
	if err != nil {
		log.Fatal("can't get current dir")
		return false,err
	}
	fileinfo, err := os.Stat(filepath.Join(curpath, string(dir)))
	if err != nil {
		log.Fatal(err)
		return false,err
	}
	if !fileinfo.IsDir() {
		log.Fatal("the installed dir not exist")
		return false,nil
	}
    return true,nil
}


func writesample(filename string) error {
	configname := strings.Split(filename, ".")[0]
	configtype := strings.Split(filename, ".")[1]
	viper.SetConfigName(configname)
	viper.SetConfigType(configtype)
	viper.AddConfigPath(".")
	viper.Set(Samplefiles,nil)
	err := viper.SafeWriteConfig()
	if err != nil {
		log.Fatal("write config failed: ", err)
		return err
	}
	return nil
}
*/
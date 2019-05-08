package shared

import (
	"compress/flate"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/mholt/archiver"
	"github.com/sibeshkar/jiminy-env/utils"
	"github.com/spf13/viper"
)

type PluginConfig struct {
	Repository  string
	EnvName     string
	Tasks       []string
	Tag         string
	Link        string
	Directory   string
	BinaryFile  string
	IncludeDirs []string
}

func CreatePluginConfig(link string) PluginConfig {

	var linkString []string = strings.Split(link, "/")
	var envName []string = strings.Split(linkString[1], "-")

	homedir := UserHomeDir()
	config := PluginConfig{
		Repository: linkString[0],
		EnvName:    envName[0],
		Tag:        envName[1],
		Link:       link,
		Directory:  homedir + "/" + ".jiminy/plugins/" + linkString[0] + "/" + linkString[1] + "/",
		BinaryFile: homedir + "/" + ".jiminy/plugins/" + linkString[0] + "/" + linkString[1] + "/" + linkString[1],
	}

	return config
}

func UserHomeDir() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("USERPROFILE")
	case "plan9":
		return os.Getenv("home")
	case "nacl", "android":
		return "/"
	case "darwin":
		if runtime.GOARCH == "arm" || runtime.GOARCH == "arm64" {
			return "/"
		}
		fallthrough
	default:
		return os.Getenv("HOME")
	}
}

//Install plugin from local folder (make sure config.json is present)
func Install(filepath string) {
	var t PluginConfig

	fileDir := utils.AbsPathify(filepath)

	err := os.Chdir(fileDir)
	if err != nil {
		panic(err)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(fileDir)
	viper.SetConfigType("json")

	fmt.Println(fileDir)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&t)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        true,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}

	os.Remove(t.BinaryFile + ".zip")

	// List of Files to Zip
	err = z.Archive(append(t.IncludeDirs, t.BinaryFile, "config.json"), t.BinaryFile+".zip")
	if err != nil {
		fmt.Println(err)
	}

	// m := CreatePluginConfig(t.Link)
	// m.Tasks = t.Tasks
	// m.IncludeDirs = t.IncludeDirs

	// if exist, _ := utils.Exists(m.Directory); exist != true {
	// 	os.MkdirAll(m.Directory, os.ModePerm)
	// }

	// file, _ := json.MarshalIndent(m, "", " ")

	// _ = ioutil.WriteFile(m.Directory+"config.json", file, 0644)

	// err = archiver.Unarchive(t.BinaryFile+".zip", m.Directory)
	//os.Remove(t.BinaryFile)
	InstallFromArchive(t.BinaryFile + ".zip")

}

func InstallFromArchive(zipfile string) {

	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        true,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}

	stringName := strings.TrimRight(zipfile, ".zip")

	var t PluginConfig

	filePath := stringName + "-tmp"
	fmt.Println(filePath)

	err := z.Unarchive(zipfile, filePath)
	if err != nil {
		fmt.Println("Error unarchiving from temp dir")
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(filePath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&t)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	m := CreatePluginConfig(t.Link)
	m.Tasks = t.Tasks
	m.IncludeDirs = t.IncludeDirs

	if exist, _ := utils.Exists(m.Directory); exist != true {
		os.MkdirAll(m.Directory, os.ModePerm)
	}

	os.Remove(m.BinaryFile)

	err = z.Unarchive(zipfile, m.Directory)

	file, _ := json.MarshalIndent(m, "", " ")

	_ = ioutil.WriteFile(m.Directory+"config.json", file, 0644)

	//utils.CopyDir(filePath, m.Directory)
	os.RemoveAll(filePath)

}

func InstallFromLink(link string) {
	return
}

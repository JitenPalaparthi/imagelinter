package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/JitenPalaparthi/imagelinter/pkg/cmdhelper"
	imagewrapper "github.com/JitenPalaparthi/imagelinter/pkg/imagewrapper"
	imglint "github.com/JitenPalaparthi/imagelinter/pkg/lint"
	"github.com/google/uuid"
)

var (
	//go:embed config/imagelintconfig.yaml
	data string
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var pathFlag = flag.String("path", wd, "path to be provided")                                    // default is current working directory
	var configPathFlag = flag.String("config", "", "path for the configuration file to be provided") // default config is the config.json file that is there in the imagelint path
	var showSumary = flag.Bool("summary", false, "to get summary pass summary=true;to off either dont pass or summary=false")
	var detailedSummary = flag.String("details", "Fail", "detailed summary can be Fail,Pass,Not-Identified or Pull-Failed")
	flag.Parse()
	var imc *imglint.ImageLintConfig
	if *configPathFlag == "" {
		imc, err = imglint.NewFromContent([]byte(data))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		imc, err = imglint.New(*configPathFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = imc.Init(*pathFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User given Path", *pathFlag)
	chelper := &cmdhelper.CmdHelper{Writer: nil}
	fmt.Println("Total number of images to process:", len(imc.ImageMap))
	c := 0
	isFatal := false
	for key := range imc.ImageMap {
		c++
		fmt.Println("Currently processing", c, " out of ", len(imc.ImageMap), "image(s)")
		containerName := uuid.New().String()
		wrapper, err := imagewrapper.New(strings.Trim(key, " "), containerName, chelper)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = wrapper.PullImage()
		if err != nil {
			// isFatal = true
			imc.OnPullFail("Pulling Image Failed:"+err.Error(), key)
			continue
		}
		skip, err := wrapper.Validate(imc.SuccessValidators)
		if skip && err == nil {
			// Success
			imc.OnPass("According to the image history, this is not an Alpine Image", key)
			continue
		}
		wrapper.CreateContainer()
		if wrapper.IsContainerExists() {
			_, err := wrapper.ContainerCP("/etc/os-release", "./")
			if err == nil {
				osdata, err := ioutil.ReadFile("os-release")
				if err == nil {
					err = os.Remove("os-release")
					if err != nil {
						log.Println(err)
					}
					if strings.Contains(string(osdata), "Alpine") {
						imc.OnFail("Alpine Image", key)
						isFatal = true
						wrapper.DeleteContainer()
						continue
					} else {
						imc.OnPass("Not an Alpine image", key)
						wrapper.DeleteContainer()
						continue
					}
				} else {
					imc.OnNotIdentifed("error in reading the file:"+err.Error()+"data:"+string(osdata), key)
					// isFatal = true
					// TODO
				}
			}
			_, err = wrapper.ContainerCP("/usr/lib/os-release", "./")
			if err == nil {
				osdata, err := ioutil.ReadFile("os-release")
				if err == nil {
					err = os.Remove("os-release")
					if err != nil {
						log.Println(err)
					}
					if strings.Contains(string(osdata), "Alpine") {
						imc.OnFail("Alpine Image", key)
						isFatal = true
						wrapper.DeleteContainer()
						continue
					} else {
						imc.OnPass("Not an Alpine image", key)
						wrapper.DeleteContainer()
						continue
					}
				} else {
					imc.OnNotIdentifed("error in reading the file:"+err.Error()+"data:"+string(osdata), key)
					// isFatal = true
					// TODO
				}
			}
			_, err = wrapper.ContainerCP("/licenses/LICENSE", "./")
			if err == nil {
				data, err := ioutil.ReadFile("LICENSE")
				if err == nil && len(data) > 0 {
					err = os.Remove("LICENSE")
					if err != nil {
						log.Println(err)
					}
					if strings.Contains(string(data), "Apache License") {
						imc.OnPass("Valid license file found", key)
						wrapper.DeleteContainer()
						continue
					}
				}
			}
			wrapper.DeleteContainer()
			_, err = wrapper.RunCommand("run", `--entrypoint=/bin/busybox`, "--name", containerName+"a", key)
			if err == nil {
				result, err := wrapper.RunCommand("logs", containerName+"a")
				wrapper.RunCommand("rm", "-f", containerName+"a")
				if strings.Contains(strings.ToUpper(result), "BUSYBOX") && err == nil {
					imc.OnPass("According to the containr bins, this is not an Alpine Image;Its busybox", key)
					continue
				}
			}
		}
		imc.OnNotIdentifed("Cound not find Container OS", key)
		// isFatal = true
		wrapper.DeleteContainer()
	}
	switch *detailedSummary {
	case "Fail", "fail", "FAIL":
		imc.ShowFailSummary()
	case "Pass", "pass", "PASS":
		imc.ShowPassSummary()
	case "Not Identified", "not identified", "NOT IDENTIFIED", "Not identified", "not Identified":
		imc.ShowNotIdentifiedSummary()
	case "Pull Failed", "", "PULL FAILED", "Pull failed", "pull Failed":
		imc.ShowPullFailedSummary()
	case "ALL", "all":
		imc.ShowPassSummary()
		imc.ShowPullFailedSummary()
		imc.ShowNotIdentifiedSummary()
		imc.ShowFailSummary()
	}
	if *showSumary {
		imc.ShowSummary()
		fmt.Println()
	}
	if isFatal {
		os.Exit(1)
	}

}

package lint

import "fmt"

const (
	colorReset   = "\033[0m" // Reset
	colorBlue    = "\033[34m"
	colorRed     = "\033[31m" // Fail
	colorPurple  = "\033[35m" // Not Identified
	colorGreen   = "\033[32m"
	colorMagenta = "\033[35m" // Pull Failed
)

func (imc *ImageLintConfig) ShowSummary() {
	pass := 0
	fail := 0
	notIdentified := 0
	pullFailed := 0
	for _, imgs := range imc.ImageMap {
		if len(imgs) > 0 {
			img := imgs[0]
			switch img.Status {
			case "Pass":
				pass = pass + len(imgs)
			case "Fail":
				fail = fail + len(imgs)
			case "Not Identified":
				notIdentified = notIdentified + len(imgs)
			case "Pull Failed":
				pullFailed = pullFailed + len(imgs)
			}
		}
	}
	fmt.Println("-------------------------------", string(colorGreen), "Summary", string(colorReset), "---------------------------------------------")
	fmt.Println("Total  Images                :", len(imc.ImageMap))
	fmt.Println("Total  Occurrences           :", pass+fail+notIdentified+pullFailed)
	fmt.Println("Total", string(colorBlue), "Pass", string(colorReset), "                :", pass)
	fmt.Println("Total", string(colorRed), "Fail", string(colorReset), "                :", fail)
	fmt.Println("Total", string(colorMagenta), "Pull Failed", string(colorReset), "         :", pullFailed)
	fmt.Println("Total", string(colorPurple), "Not Identified", string(colorReset), "      :", notIdentified)
	fmt.Println("[Note:All totals for Pass|Fail|Not Identified|Pull Failed are based on number of occurrences]")
	fmt.Println("--------------------------------------------------------------------------------------------")
}

func (imc *ImageLintConfig) ShowFailSummary() {
	fmt.Println("-------------------------------", string(colorGreen), "Fail Summary", string(colorReset), "---------------------------------------------")
	fail := 0
	for image, imgs := range imc.ImageMap {
		if len(imgs) > 0 {
			img := imgs[0]
			if img.Status == "Fail" {
				fail = fail + len(imgs)
				fmt.Println("Image:    ", image)
				fmt.Println("File Path:", img.Path)
				fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
				fmt.Println("")
			}
		}
	}
	fmt.Println("Total", string(colorRed), "Fail", string(colorReset), "          :", fail)
}

func (imc *ImageLintConfig) ShowPassSummary() {
	fmt.Println("-------------------------------", string(colorGreen), "Pass Summary", string(colorReset), "---------------------------------------------")
	pass := 0
	for image, imgs := range imc.ImageMap {
		if len(imgs) > 0 {
			img := imgs[0]
			if img.Status == "Pass" {
				pass = pass + len(imgs)
				fmt.Println("Image:    ", image)
				fmt.Println("File Path:", img.Path)
				fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
				fmt.Println("")
			}
		}
	}
	fmt.Println("Total", string(colorBlue), "Pass", string(colorReset), "          :", pass)
}

func (imc *ImageLintConfig) ShowNotIdentifiedSummary() {
	fmt.Println("-------------------------------", string(colorGreen), "Not-Identified Summary", string(colorReset), "---------------------------------------------")
	notIdentified := 0
	for image, imgs := range imc.ImageMap {
		if len(imgs) > 0 {
			img := imgs[0]
			if img.Status == "Not Identified" {
				notIdentified = notIdentified + len(imgs)
				fmt.Println("Image:    ", image)
				fmt.Println("File Path:", img.Path)
				fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
				fmt.Println("")
			}
		}
	}
	fmt.Println("Total", string(colorBlue), "Not Identified", string(colorPurple), "          :", notIdentified)
}

func (imc *ImageLintConfig) ShowPullFailedSummary() {
	fmt.Println("-------------------------------", string(colorGreen), "Pull-Failed Summary", string(colorReset), "---------------------------------------------")
	pullFailed := 0
	for image, imgs := range imc.ImageMap {
		if len(imgs) > 0 {
			img := imgs[0]
			if img.Status == "Pull Failed" {
				pullFailed = pullFailed + len(imgs)
				fmt.Println("Image:    ", image)
				fmt.Println("File Path:", img.Path)
				fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
				fmt.Println("")
			}
		}
	}
	fmt.Println("Total", string(colorBlue), "Pull Failed", string(colorMagenta), "          :", pullFailed)
}

func (imc *ImageLintConfig) OnPass(message, image string) {
	fmt.Println("Status: ", string(colorBlue), "Pass", string(colorReset))
	fmt.Println("Image:   ", image)
	fmt.Println("Message: ", message)
	fmt.Println("Total ", len(imc.ImageMap[image]), " file(s) contain(s) this image")
	for i, img := range imc.ImageMap[image] {
		fmt.Println("File Path:", img.Path)
		fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
		imc.ImageMap[image][i].Status = "Pass"
	}
	fmt.Println()
}
func (imc *ImageLintConfig) OnFail(message, image string) {
	fmt.Println("Status: ", string(colorRed), "Fail", string(colorReset))
	fmt.Println("Image:   ", image)
	fmt.Println("Error:   ", message)
	fmt.Println("Total ", len(imc.ImageMap[image]), " file(s) contain(s) this image")
	for i, img := range imc.ImageMap[image] {
		fmt.Println("File Path:", img.Path)
		fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
		imc.ImageMap[image][i].Status = "Fail"
	}
	fmt.Println()
}
func (imc *ImageLintConfig) OnNotIdentifed(message, image string) {
	fmt.Println("Status: ", string(colorPurple), "Not Identified", string(colorReset))
	fmt.Println("Image:   ", image)
	fmt.Println("Error:   ", message)
	fmt.Println("Total ", len(imc.ImageMap[image]), " file(s) contain(s) this image")
	for i, img := range imc.ImageMap[image] {
		fmt.Println("File Path:", img.Path)
		fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
		imc.ImageMap[image][i].Status = "Not Identified"
	}
	fmt.Println()
}
func (imc *ImageLintConfig) OnPullFail(message, image string) {
	fmt.Println("Status: ", string(colorMagenta), "Pull Failed", string(colorReset))
	fmt.Println("Image:   ", image)
	fmt.Println("Error:   ", message)
	fmt.Println("Total ", len(imc.ImageMap[image]), " file(s) contain(s) this image")
	for i, img := range imc.ImageMap[image] {
		fmt.Println("File Path:", img.Path)
		fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
		imc.ImageMap[image][i].Status = "Pull Failed"
	}
	fmt.Println()
}

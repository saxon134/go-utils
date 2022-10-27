package saVideo

import (
	"path/filepath"
)

func (V Video) Screenshot(timeInSeconds float64, outputFilename string) (string, error) {
	compileKwArgs := V.prepareKwArgs([]string{"c:a", "c:v", "af"})
	compileKwArgs["ss"] = timeInSeconds
	compileKwArgs["vframes"] = "1"

	abs, absError := filepath.Abs(outputFilename)
	if absError != nil {
		return "", absError
	}

	err := V.stream.Output(abs, compileKwArgs).OverWriteOutput().Run()
	if err != nil {
		return "", err
	}

	return abs, nil
}

package utils

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const splitSize = 32
const u2ModelDownloadURL = "https://github.com/danielgatis/rembg/releases/download/v0.0.0/u2net.onnx"

var destFolder string
var modelFilePath string
var mu sync.Mutex

func init() {
	userHomeFolder, err := os.UserHomeDir()

	if err != nil {
		panic("cannot get user home directory")
	}

	destFolder = filepath.Join(userHomeFolder, ".u2net")
	modelFilePath = filepath.Join(destFolder, "u2net.onnx")
}

func ChkU2ModelDownload() bool {
	info, err := os.Stat(modelFilePath)

	return err == nil && !info.IsDir()
}

func U2Download() error {

	headResponse, err := http.Head(u2ModelDownloadURL)

	if err != nil {
		return fmt.Errorf("cannot construct a HEAD request, reason: %s", err.Error())
	}

	if headResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("http HEAD returned a failure status code: %d", headResponse.StatusCode)
	}

	bodySize, err := strconv.Atoi(headResponse.Header.Get("Content-Length"))

	if err != nil {
		return fmt.Errorf("cannot read the numeric value of content length header")
	}

	wg := &sync.WaitGroup{}
	wg.Add(splitSize)

	bufSplitted := make([][]byte, splitSize)

	perChunkSize := int(math.Floor(float64(bodySize) / float64(splitSize)))

	chunkBroken := false

	for i := 0; i < splitSize; i++ {
		go func(chunkIndex int) {
			defer wg.Done()

			if chunkBroken {
				return
			}

			startByte := chunkIndex * perChunkSize
			endByte := (chunkIndex+1)*perChunkSize - 1

			if chunkIndex == splitSize-1 {
				endByte = bodySize
			}

			req, err := http.NewRequest("GET", u2ModelDownloadURL, nil)
			if err != nil {
				bufSplitted[chunkIndex] = nil

				mu.Lock()
				chunkBroken = true
				mu.Unlock()

				return
			}

			req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", startByte, endByte))

			resp, err := http.DefaultClient.Do(req)

			if err != nil {
				bufSplitted[chunkIndex] = nil

				mu.Lock()
				chunkBroken = true
				mu.Unlock()

				return
			}
			defer resp.Body.Close()

			fileContents, err := io.ReadAll(resp.Body)
			if err != nil {
				bufSplitted[chunkIndex] = nil

				mu.Lock()
				chunkBroken = true
				mu.Unlock()

				return
			}

			bufSplitted[chunkIndex] = fileContents
		}(i)
	}

	wg.Wait()

	if chunkBroken {
		return fmt.Errorf("one of download chunk is broken")
	}

	file, err := os.OpenFile(modelFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

	if err != nil {
		return fmt.Errorf("cannot open result file u2net.onnx, reason: %s", err.Error())
	}
	defer file.Close()

	var currentFileTop int64 = 0
	for _, chunk := range bufSplitted {
		file.Seek(currentFileTop, 0)
		file.Write(chunk)
		currentFileTop += int64(len(chunk))
	}

	return nil
}

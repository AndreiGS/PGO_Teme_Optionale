package pkg

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	"log"
	"runtime"

	ort "github.com/yalue/onnxruntime_go"
)

var (
	mean = []float32{23.06966, -0.00452, 1.5650475, 3.119159}
	std  = []float32{8.339648, 14.216207, 0.92338675, 1.8202546}
)

type PredictResponse struct {
	Real float32 `json:"real"`
	Imag float32 `json:"imag"`
}

func InitializeDLL() {
	onnxruntimeLibPath := getDefaultSharedLibPath()
	ort.SetSharedLibraryPath(onnxruntimeLibPath)
}

func Generate(batch [][]float32) ([]PredictResponse, error) {
	prediction, e := predict(batch)
	if e != nil {
		return nil, e
	}
	return prediction, nil
}

func scaleMatrixInPlace(matrix [][]float32, mean []float32, std []float32) {
	cols := len(matrix[0])

	// Check if mean and std have same length as columns
	if len(mean) != cols || len(std) != cols {
		panic("mean and std must have the same length as matrix columns")
	}

	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = (matrix[i][j] - mean[j]) / std[j]
		}
	}
}

func getDefaultSharedLibPath() string {
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			return "./third_party/onnxruntime.dll"
		}
	}
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			return "./third_party/onnxruntime_arm64.dylib"
		}
	}
	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "arm64" {
			return "./third_party/onnxruntime_arm64.so"
		}
		return "./third_party/onnxruntime.so"
	}
	log.Fatalf("Unable to determine a path to the onnxruntime shared library"+
		" for OS \"%s\" and architecture \"%s\".\n", runtime.GOOS,
		runtime.GOARCH)
	return ""
}

func getInput(batch [][]float32) []float32 {
	input := make([]float32, len(batch)*4)
	for i := 0; i < len(batch); i++ {
		for j := 0; j < len(batch[i]); j++ {
			input[i*4+j] = batch[i][j]
		}
	}
	return input
}

func getOutput(output []float32) []PredictResponse {
	response := make([]PredictResponse, len(output)/2)
	for i := 0; i < len(output); i += 2 {
		response[i/2] = PredictResponse{Real: output[i], Imag: output[i+1]}
	}
	return response
}

func predict(batch [][]float32) ([]PredictResponse, error) {
	scaleMatrixInPlace(batch, mean, std)

	batchSize := int64(len(batch))

	e := ort.InitializeEnvironment()
	if e != nil {
		return nil, fmt.Errorf("error initializing the onnxruntime library: %w", e)
	}
	defer ort.DestroyEnvironment()

	inputShape := ort.NewShape(batchSize, 4)
	inputData := getInput(batch)
	input, e := ort.NewTensor(inputShape, inputData)
	if e != nil {
		return nil, fmt.Errorf("error creating input tensor: %w", e)
	}
	defer input.Destroy()

	output, e := ort.NewEmptyTensor[float32](ort.NewShape(batchSize, 2))
	if e != nil {
		return nil, fmt.Errorf("error creating output tensor: %w", e)
	}
	defer output.Destroy()

	session, e := ort.NewAdvancedSession("./third_party/autoencoder2_freezed.onnx",
		[]string{"input"}, []string{"output"},
		[]ort.ArbitraryTensor{input}, []ort.ArbitraryTensor{output}, nil)
	if e != nil {
		return nil, fmt.Errorf("error creating network session: %w", e)
	}
	defer session.Destroy()

	e = session.Run()
	if e != nil {
		return nil, fmt.Errorf("error running the network: %w", e)
	}

	outputData := output.GetData()
	return getOutput(outputData), nil
}

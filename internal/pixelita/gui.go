package pixelita

import (
	"app/internal/config"
	"app/internal/filesystem"
	"app/internal/imageencode"
	"app/internal/imagetype"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
)

type List struct {
	SourcePath string
	Path       string
	Format     string
	Size       string
	Status     string
	NewSize    string
}

var (
	showDemoWindow bool

	backend imgui.Backend[imgui.GLFWWindowFlags]

	posterization int32
	speed         int32
	minQuality    int32
	maxQuality    int32
	quality       [2]int32 = [2]int32{minQuality, maxQuality}
	jpgQuality    int32

	list      []List
	plainList []string
	conf      *config.Config

	fs      filesystem.FileSystem
	imgtype imagetype.ImageType
	pngEnc  imageencode.PNGEncoder
	jpgEnc  imageencode.JPEGEncoder

	workersNum int32

	status string
)

func (px *Pixelita) StartGUI() {
	px.backendInit()
}

func callback(data imgui.InputTextCallbackData) int {
	fmt.Println("got call back")
	return 0
}

func mainWIndow() {
	// Main viewport position
	posX := imgui.MainViewport().Pos().X
	posY := imgui.MainViewport().Pos().Y
	imgui.SetNextWindowPos(imgui.Vec2{X: posX, Y: posY})

	// Main viewport size
	width := imgui.MainViewport().Size().X
	height := imgui.MainViewport().Size().Y
	imgui.SetNextWindowSize(imgui.Vec2{X: width, Y: height})

	if showDemoWindow {
		imgui.ShowDemoWindowV(&showDemoWindow)
	}

	var is_opened *bool
	imgui.BeginV("Main Window", is_opened, imgui.WindowFlagsNoTitleBar|imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoCollapse)
	var flags imgui.TableFlags = imgui.TableFlagsResizable /* + imgui.TableFlagsBorders*/
	if imgui.BeginTableV("container", 1, flags, imgui.Vec2{0.0, 0.0}, 0) {
		imgui.TableNextColumn()
		if imgui.BeginTable("mainLayout", 2) {

			imgui.TableSetupColumnV("settings", imgui.TableColumnFlagsWidthFixed, 380.0, 600.0)
			imgui.TableSetupColumnV("files", imgui.TableColumnFlagsWidthStretch, 600.0, 600.0)

			imgui.TableNextColumn()

			// imgui.BeginGroup()
			imgui.SeparatorText("General")
			imgui.SetNextItemWidth(200.0)
			imgui.SliderInt("Concurrent threads", &workersNum, 1, 16)
			imgui.SameLine()
			imgui.Text("(?)")
			if imgui.IsItemHoveredV(imgui.HoveredFlagsAllowWhenDisabled) {
				imgui.SetTooltip("Specifies the number of simultaneous operations that can be performed concurrently.\nHigher values may improve processing speed but consume more system resources.\nAdjust based on your system's capabilities and performance requirements.")
			}

			imgui.SeparatorText("PNG Settings (Image Quantization)")

			imgui.SetNextItemWidth(200.0)
			imgui.SliderInt("Posterizing level", &posterization, 0, 4)
			imgui.SameLine()
			imgui.Text("(?)")
			if imgui.IsItemHoveredV(imgui.HoveredFlagsAllowWhenDisabled) {
				imgui.SetTooltip("Ignores given number of least significant bits in all channels, posterizing image to 2^bits levels.\n0 gives full quality. Use 2 for VGA or 16-bit RGB565 displays.\n4 if image is going to be output on a RGB444/RGBA4444 display (e.g. low-quality textures on Android).")
			}

			imgui.SetNextItemWidth(200.0)
			imgui.SliderInt2("Min max quality", &quality, 0, 100)
			imgui.SameLine()
			imgui.Text("(?)")
			if imgui.IsItemHoveredV(imgui.HoveredFlagsAllowWhenDisabled) {
				imgui.SetTooltip("Quality is in range 0 (worst) to 100 (best) and values are analoguous to JPEG quality (i.e. 80 is usually good enough).\nQuantization will attempt to use the lowest number of colors needed to achieve maximum quality.\nMaximum value of 100 is the default and means conversion as good as possible.\nIf it's not possible to convert the image with at least minimum quality (i.e. 256 colors is not enough to meet the minimum quality),\nthen Image.Quantize() will fail. The default minimum is 0 (proceeds regardless of quality)")
			}

			imgui.SetNextItemWidth(200.0)
			imgui.SliderInt("Speed", &speed, 1, 10)
			imgui.SameLine()
			imgui.Text("(?)")
			if imgui.IsItemHoveredV(imgui.HoveredFlagsAllowWhenDisabled) {
				imgui.SetTooltip("Higher speed levels disable expensive algorithms and reduce quantization precision.\nThe default speed is 3. Speed 1 gives marginally better quality at significant CPU cost.\nSpeed 10 has usually 5%% lower quality, but is 8 times faster than the default.\nHigh speeds combined with Quality parameter will use more colors than necessary\nand will be less likely to meet minimum required quality.")
			}

			imgui.SeparatorText("JPG/JPEG Settings")
			imgui.SetNextItemWidth(200.0)
			imgui.SliderInt("Min max quality", &jpgQuality, 0, 100)
			imgui.SameLine()

			// imgui.Checkbox("Show demo window", &showDemoWindow)

			// imgui.EndGroup()
			imgui.TableNextColumn()

			if imgui.BeginChildStrV("child", imgui.Vec2{0, height - 140}, false, imgui.WindowFlagsAlwaysAutoResize) {
				imgui.Text(" ")
				imgui.SameLine()

				var flags imgui.TableFlags = imgui.TableFlagsResizable + imgui.TableFlagsBorders
				if imgui.BeginTableV("files", 5, flags, imgui.Vec2{0.0, 0.0}, 0) {
					imgui.TableSetupColumnV("Path", imgui.TableColumnFlagsWidthStretch, 400.0, 0)
					imgui.TableSetupColumnV("Format", imgui.TableColumnFlagsWidthFixed, 60.0, 0)
					imgui.TableSetupColumnV("Size", imgui.TableColumnFlagsWidthFixed, 60.0, 0)
					imgui.TableSetupColumnV("New size", imgui.TableColumnFlagsWidthFixed, 60.0, 0)
					imgui.TableSetupColumnV("Status", imgui.TableColumnFlagsWidthFixed, 60.0, 0)

					imgui.TableHeadersRow()
					if len(list) > 0 {
						for _, elem := range list {
							imgui.TableNextColumn()
							imgui.Text(elem.Path)

							imgui.TableNextColumn()
							extension := strings.ReplaceAll(filepath.Ext(elem.Format), ".", "")
							imgui.Text(extension)

							imgui.TableNextColumn()
							imgui.Text(elem.Size)

							imgui.TableNextColumn()
							imgui.Text(elem.NewSize)

							imgui.TableNextColumn()
							imgui.Text(elem.Status)

						}
					}

					imgui.EndTable()
				}
				imgui.EndChild()
			}

			imgui.EndTable()
		}

		imgui.TableNextRow()
		imgui.TableNextColumn()

		imgui.SeparatorText("Actions")

		if imgui.ButtonV("SCAN UPLOADS", imgui.Vec2{X: 0, Y: 48}) {
			list = []List{}

			// Check if conf.UploadDir exists, create it if it doesn't
			_, err := os.Stat(conf.UploadDir)
			if os.IsNotExist(err) {
				fmt.Printf("Upload directory '%s' does not exist. Creating it...\n", conf.UploadDir)
				err = os.MkdirAll(conf.UploadDir, 0755) // Create the directory and parent directories if needed
				if err != nil {
					fmt.Printf("Error creating upload directory: %v\n", err)
					return
				}
			}

			go func() {
				l, err := fs.GetImageFiles(conf.UploadDir)
				if err != nil {
					panic("error getting image files from uploads dir")
				}

				for _, file := range l {
					relativePath := strings.TrimPrefix(file, conf.UploadDir)
					if strings.HasPrefix(relativePath, "\\") {
						relativePath = relativePath[1:] // Remove leading '/'
					}

					var fileSize int64
					fileInfo, err := os.Stat(file)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}

					fileSize = fileInfo.Size()
					fileSizeString := formatFileSize(fileSize)
					addToList(file, relativePath, filepath.Ext(file), fileSizeString)
				}
			}()

			status = "COMPRESS FILES"
		}

		imgui.SameLine()

		if imgui.ButtonV(status, imgui.Vec2{X: 0, Y: 48}) {
			if len(list) > 0 {
				if status == "COMPRESS FILES" {
					go func() {
						status = "WORKING..."
						numWorkers := workersNum // Adjust the number of workers as needed
						jobs := make(chan int, len(list))
						results := make(chan CompressionResult, len(list))

						// Start the worker pool
						for w := 1; w <= int(numWorkers); w++ {
							go worker(w, jobs, results)
						}

						// Add jobs to the job queue
						for i := range list {
							jobs <- i
						}
						close(jobs)

						// Collect the results and update the UI
						for range list {
							result := <-results
							list[result.index].Status = result.status
							list[result.index].NewSize = result.newSize
						}
						// for i := range list {
						// 	var processedData []byte
						// 	// Read uploaded file content
						// 	uploadedData, err := fs.ReadFile(list[i].SourcePath)
						// 	if err != nil {
						// 		log.Printf("Error reading image data: %v", err)
						// 		list[i].Status = "UKNOWN"
						// 		return
						// 	}

						// 	// Check the file format
						// 	fileFormat, err := imagetype.New().GetFormatByExtension(list[i].Format)
						// 	if err != nil {
						// 		log.Printf("Error detecting image format: %v", err)
						// 		continue
						// 	}

						// 	switch fileFormat {
						// 	case imagetype.JPEG:
						// 		// Compress JPEG
						// 		processedData, err = jpgEnc.Encode(uploadedData, int(jpgQuality))
						// 		if err != nil {
						// 			log.Printf("Error compressing image: %v", err)
						// 			list[i].Status = "ERROR"
						// 			continue
						// 		}
						// 	case imagetype.PNG:
						// 		// Compress PNG
						// 		processedData, err = pngEnc.Encode(uploadedData, int(posterization), int(minQuality), int(maxQuality), int(speed))
						// 		if err != nil {
						// 			log.Printf("Error compressing image: %v", err)
						// 			list[i].Status = "ERROR"
						// 			continue
						// 		}
						// 	default:
						// 		log.Printf("Unsupported image format: %v", fileFormat)
						// 	}

						// 	// Save processed content to file
						// 	destFile := filepath.Join(conf.ProcessedDir, list[i].Path)
						// 	err = fs.SaveFile(destFile, processedData)
						// 	if err != nil {
						// 		log.Printf("Error saving converted image to file: %v", err)
						// 	}
						// 	list[i].Status = "DONE"

						// 	var fileSize int64
						// 	fileInfo, err := os.Stat(destFile)
						// 	if err != nil {
						// 		fmt.Println("Error:", err)
						// 		return
						// 	}

						// 	fileSize = fileInfo.Size()
						// 	fileSizeString := formatFileSize(fileSize)
						// 	list[i].NewSize = fileSizeString
						// }
						status = "CLEAR LIST"
					}()

				}
			}

			if status == "CLEAR LIST" {
				list = []List{}
				status = "COMPRESS FILES"
			}
		}

		imgui.EndTable()
	}

	imgui.End()
}

func afterCreateContext() {

}

func loop() {
	mainWIndow()

}

func beforeDestroyContext() {
	imgui.PlotDestroyContext()
}

func (px *Pixelita) backendInit() {

	posterization = 0
	speed = 3
	minQuality = int32(0)
	maxQuality = int32(100)
	quality = [2]int32{minQuality, maxQuality}
	jpgQuality = 80
	workersNum = 4

	conf = config.NewConfig()
	fs = filesystem.NewFileSystemImpl(conf)
	imgtype = imagetype.New()
	pngEnc = *imageencode.NewPNGEncoder()
	jpgEnc = *imageencode.NewJPEGEncoder()

	status = "COMPRESS FILES"

	backend = imgui.CreateBackend(imgui.NewGLFWBackend())
	backend.SetAfterCreateContextHook(afterCreateContext)
	backend.SetBeforeDestroyContextHook(beforeDestroyContext)

	backend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))
	backend.CreateWindow("Pixelita - JPG and PNG compressor", 1024, 580)

	backend.SetDropCallback(func(p []string) {
		fmt.Printf("drop triggered: %v", p)
	})

	backend.SetCloseCallback(func(b imgui.Backend[imgui.GLFWWindowFlags]) {
		fmt.Println("window is closing")
	})

	// backend.SetIcons(img)
	backend.SetTargetFPS(120)

	io := imgui.CurrentIO()
	io.Fonts().AddFontFromFileTTF("font", float32(20))

	style := imgui.NewStyle()
	style.SetWindowPadding(imgui.Vec2{X: 15, Y: 15})
	style.SetFramePadding(imgui.Vec2{X: 10, Y: 2})
	style.SetCellPadding(imgui.Vec2{X: 6, Y: 4})
	style.SetItemSpacing(imgui.Vec2{X: 6, Y: 10})
	style.SetItemInnerSpacing(imgui.Vec2{X: 6, Y: 6})
	style.SetIndentSpacing(float32(20))
	style.SetScrollbarSize(float32(20))
	style.SetGrabMinSize(float32(20))
	style.SetGrabRounding(float32(3))
	style.SetFrameRounding(float32(3))
	style.SetScrollbarRounding(float32(3))
	style.SetTabRounding(float32(3))

	io.Ctx().SetStyle(*style)

	// imgui.StyleColorsLight()

	backend.Run(loop)
}

// Create a new List instance for each image file and add it to the list slice
func addToList(sourcePath, path string, format string, size string) {
	image := List{
		SourcePath: sourcePath,
		Path:       path,
		Format:     format,
		Size:       size,
	}

	list = append(list, image)

	// fmt.Println(list)
}

// Function to format file size with appropriate unit (KB, MB, GB, etc.)
func formatFileSize(fileSize int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
		TB = 1 << 40
		PB = 1 << 50
		EB = 1 << 60
	)

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	value := float64(fileSize)

	var i int
	for i = 0; i < len(units); i++ {
		if value < 1024 {
			break
		}
		value /= 1024
	}

	return fmt.Sprintf("%.2f %s", value, units[i])
}

func worker(id int, jobs <-chan int, results chan<- CompressionResult) {
	for i := range jobs {
		// Process the image
		var processedData []byte
		// Read uploaded file content
		uploadedData, err := fs.ReadFile(list[i].SourcePath)
		if err != nil {
			log.Printf("Error reading image data: %v", err)
			list[i].Status = "UKNOWN"
			return
		}

		// Check the file format
		fileFormat, err := imagetype.New().GetFormatByExtension(list[i].Format)
		if err != nil {
			log.Printf("Error detecting image format: %v", err)
			continue
		}

		switch fileFormat {
		case imagetype.JPEG:
			// Compress JPEG
			processedData, err = jpgEnc.Encode(uploadedData, int(jpgQuality))
			if err != nil {
				log.Printf("Error compressing image: %v", err)
				list[i].Status = "ERROR"
				continue
			}
		case imagetype.PNG:
			// Compress PNG
			processedData, err = pngEnc.Encode(uploadedData, int(posterization), int(minQuality), int(maxQuality), int(speed))
			if err != nil {
				log.Printf("Error compressing image: %v", err)
				list[i].Status = "ERROR"
				continue
			}
		default:
			log.Printf("Unsupported image format: %v", fileFormat)
		}

		// Save processed content to file
		destPath := filepath.Join(conf.ProcessedDir, list[i].Path)
		err = fs.SaveFile(destPath, processedData)
		if err != nil {
			log.Printf("Error saving converted image to file: %v", err)
		}
		list[i].Status = "DONE"

		var fileSize int64
		fileInfo, err := os.Stat(destPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fileSize = fileInfo.Size()
		fileSizeString := formatFileSize(fileSize)
		list[i].NewSize = fileSizeString

		// Simulate some work (replace this with your image processing)
		time.Sleep(time.Millisecond * 500)

		// Send the result back to the main goroutine
		results <- CompressionResult{
			index:   i,
			status:  "DONE",
			newSize: fileSizeString,
		}
	}
}

type CompressionResult struct {
	index   int
	status  string
	newSize string
}

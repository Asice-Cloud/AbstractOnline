package verification

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 请求结构体
type VerifyRequest struct {
	SliderPosition int `json:"slider_position"`
}

// 响应结构体
type VerifyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Worker Pool的大小
const workerPoolSize = 20

var (
	gapPosition   int
	allowedError  = 20
	backgroundImg image.Image
	sliderImg     image.Image
)

func init() {
	rand.Seed(time.Now().UnixNano())
	loadImages()
}

func loadImages() {
	file, err := os.Open("./static/img/look.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	backgroundImg, err = jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	gapPosition = rand.Intn(backgroundImg.Bounds().Dx() - 100)
	fmt.Printf("Gap position: %d\n", gapPosition) // 添加日志
	sliderImg = createSliderImage(backgroundImg, gapPosition)
	backgroundImg = createBackgroundWithGap(backgroundImg, gapPosition)
}

func createSliderImage(img image.Image, x int) image.Image {
	rect := image.Rect(0, 0, 100, 100)
	slider := image.NewRGBA(rect)
	draw.Draw(slider, rect, img, image.Point{X: x, Y: 100}, draw.Src)
	return slider
}

func createBackgroundWithGap(img image.Image, x int) image.Image {
	rect := image.Rect(x, 100, x+100, 200)
	backgroundWithGap := image.NewRGBA(img.Bounds())
	draw.Draw(backgroundWithGap, img.Bounds(), img, image.Point{}, draw.Src)

	// 用透明色填充缺口部分
	transparent := image.NewUniform(color.RGBA{0, 0, 0, 0})
	draw.Draw(backgroundWithGap, rect, transparent, image.Point{}, draw.Src)
	return backgroundWithGap
}

func verifySlider(req VerifyRequest) VerifyResponse {
	fmt.Printf("Received slider position: %d\n", req.SliderPosition) // 添加日志
	if abs(req.SliderPosition-gapPosition) <= allowedError {
		return VerifyResponse{Success: true, Message: "验证通过"}
	}
	return VerifyResponse{Success: false, Message: "验证失败"}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func worker(id int, jobs <-chan VerifyRequest, results chan<- VerifyResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range jobs {
		results <- verifySlider(req)
	}
}

func GetBackground(c *gin.Context) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, backgroundImg, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成背景图像"})
		return
	}
	c.Data(http.StatusOK, "image/jpeg", buf.Bytes())
}

func Slider(c *gin.Context) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, sliderImg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成滑块图像"})
		return
	}
	c.Data(http.StatusOK, "image/png", buf.Bytes())
}

func Verify(c *gin.Context) {
	jobs := make(chan VerifyRequest, workerPoolSize)
	results := make(chan VerifyResponse, workerPoolSize)

	var wg sync.WaitGroup
	for w := 1; w <= workerPoolSize; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "缺少滑块位置"})
		return
	}

	jobs <- req
	close(jobs) // Close the jobs channel to signal no more jobs

	go func() {
		wg.Wait()
		close(results) // Close the results channel after all workers are done
	}()

	res := <-results
	fmt.Printf("Verification result: %v\n", res) // Add logging for the verification result
	c.JSON(http.StatusOK, res)
}

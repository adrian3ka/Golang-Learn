package main
import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"fmt"
	"runtime"
	"image"
	"github.com/disintegration/imaging"

)
type Image struct{
	ID string
	UserID string
	Name string
	Location string
	Size int64
	CreatedAt time.Time
	Description string
}

const imageIDLength = 10

func NewImage(user *User) *Image{
	return &Image{
		ID: GenerateID("img", imageIDLength),
		UserID: user.ID,
		CreatedAt: time.Now(),
	}
}

// A map of accepted mime types and their file extension
var mimeExtensions = map[string]string{
	"image/png": ".png",
	"image/jpeg": ".jpg",
	"image/gif": ".gif",
}

// A map of accepted mime types and their file extension
var fileExtensions = map[string]string{
	".png": ".png",
	".jpg": ".jpg",
	".gif": ".gif",
}

func (image *Image) CreateFromURL(imageURL string) error {
	// Get the response from the URL
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	// Make sure we have a response
	if response.StatusCode != http.StatusOK {
		return errImageURLInvalid
	}
	defer response.Body.Close()
	
	// Ascertain the type of file we downloaded
	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	fmt.Println(mimeType)
	if err != nil {
		return errInvalidImageType
	}
	// Get an extension for the file
	ext, valid := mimeExtensions[mimeType]
	if !valid {
		return errInvalidImageType
	}
	// Get a name from the URL
	image.Name = filepath.Base(imageURL)
	image.Location = image.ID + ext
	// Open a file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()
	// Copy the entire response to the output flie
	size, err := io.Copy(savedFile, response.Body)
	if err != nil {
		return err
	}
	// The returned value from io.Copy is the number of bytes copied
	image.Size = size
	
	// Create the various resizes of the images
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}
	
	// Save our image to the store
	return globalImageStore.Save(image)
}

func (image *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {

	// Move our file to an appropriate place, with an appropriate name
	image.Name = headers.Filename
	image.Location = image.ID + filepath.Ext(image.Name)
	ext := filepath.Ext(image.Name)
	_, valid := fileExtensions[ext]
	if !valid {
		return errInvalidImageType
	}
	
	// Open a file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}

	defer savedFile.Close()

	// Copy the uploaded file to the target location
	size, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}
	image.Size = size

	// Create the various resizes of the images
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}
	
	// Save the image to the database
	return globalImageStore.Save(image)
}
//Kenapa bisa pake im doang?, lihat di main.go, ingat2 lagi
/*
router.ServeFiles(
	"/im/*filepath",
	http.Dir("data/images/"),
)
*/
func (image *Image) StaticRoute() string {
	return "/im/" + image.Location
}
func (image *Image) ShowRoute() string {
	return "/image/" + image.ID
}
func (image *Image) StaticThumbnailRoute() string {
	return "/im/thumbnail/" + image.Location
}
func (image *Image) StaticPreviewRoute() string {
	return "/im/preview/" + image.Location
}

var widthThumbnail = 400
func (image *Image) resizeThumbnail(errorChan chan error, srcImage image.Image) {
	dstImage := imaging.Thumbnail(srcImage, widthThumbnail, widthThumbnail, imaging.Lanczos)

	destination := "./data/images/thumbnail/" + image.Location
	errorChan <- imaging.Save(dstImage, destination)
}

var widthPreview = 800
func (image *Image) resizePreview(errorChan chan error, srcImage image.Image) {
	size := srcImage.Bounds().Size()
	ratio := float64(size.Y) / float64(size.X)
	targetHeight := int(float64(widthPreview) * ratio)
	dstImage := imaging.Resize(srcImage, widthPreview, targetHeight,imaging.Lanczos)
	destination := "./data/images/preview/" + image.Location
	errorChan <- imaging.Save(dstImage, destination)
}

func (image *Image) CreateResizedImages() error {
	// Generate an image from file
	srcImage, err := imaging.Open("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	// Create a channel to receive errors on
	errorChan := make(chan error)
	// Process each size
	go image.resizePreview(errorChan, srcImage)
	go image.resizeThumbnail(errorChan, srcImage)
	// Wait for images to finish resizing
	for i := 0; i < 2; i++ {
		err := <-errorChan
		fmt.Println("err img resize>> ",err)
		if err == nil {
			return err
		}
	}
	return err
}

func init() {
	// Ensure our goroutines run across all cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
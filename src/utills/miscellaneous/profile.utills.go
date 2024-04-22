package utills

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"strings"
)

func RandomProfileImage() ([]byte, error) {
	var profileImages = []string{
		"/assets/profile/1.jpg",
		"/assets/profile/2.png",
		"/assets/profile/3.png"}

	// generate a random index
	var random_idx = rand.Intn(len(profileImages))
	current_wd, _ := os.Getwd()
	current_wd = strings.Split(current_wd, "src")[0]
	// check if the current working directory ends with a slash
	if current_wd[len(current_wd)-1] != '/' {
		current_wd = current_wd + "/"
	}
	current_wd = current_wd + "src"
	var target_file = fmt.Sprintf("%s%s", current_wd, profileImages[random_idx])
	fmt.Println(target_file)

	return GetImage(target_file)
}

func RandomServiceImage() ([]byte, error) {
	var serviceImages = []string{
		"/assets/service/1.jpg",
		"/assets/service/2.jpg"}

	// generate a random index
	var random_idx = rand.Intn(len(serviceImages))
	current_wd, _ := os.Getwd()
	current_wd = strings.Split(current_wd, "src")[0] + "src"
	var target_file = fmt.Sprintf("%s%s", current_wd, serviceImages[random_idx])
	fmt.Println(target_file)

	return GetImage(target_file)
}

func GetImage(target_file string) ([]byte, error) {
	file, err := os.Open(target_file)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)

	return buf.Bytes(), err
}

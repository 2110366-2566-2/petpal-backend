package utills

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
)

func RandomProfileImage() ([]byte, error) {
	var profileImages = []string{
		"/assets/profile/1.jpg",
		"/assets/profile/2.png",
		"/assets/profile/3.png"}

	// generate a random index
	var random_idx = rand.Intn(len(profileImages))
	current_wd, _ := os.Getwd()
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

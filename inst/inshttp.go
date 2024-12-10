package inst

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func File(postURL string) (string, error) {

	url := fmt.Sprintf("http://150.241.113.11:3000/api/video?url={%s}", postURL)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cant read resp body")
	}

	return string(body), nil

	// // Парсим HTML-страницу
	// doc, err := goquery.NewDocumentFromReader(resp.Body)
	// if err != nil {
	// 	return "", err, 0
	// }

	// metaTag := doc.Find("meta[property='og:type' content='video']")

	// if metaTag.Length() == 0 {
	// 	fmt.Printf("image not found")
	// 	// Ищем теги <meta> с изображением
	// 	metaTag := doc.Find("meta[property='og:image']")
	// 	if metaTag.Length() == 0 {
	// 		return "", fmt.Errorf("contnet not found"), 0
	// 	} else {
	// 		file, err := Img(postURL)
	// 		if err != nil {
	// 			return "", fmt.Errorf("can't save img :%v", err), 0
	// 		}
	// 		return file, nil, 1
	// 	}
	// } else {
	// 	file, err := video(postURL)
	// 	if err != nil {
	// 		return "", fmt.Errorf("cant save video: %v", err), 0
	// 	}
	// 	return file, nil, 2
	// }

}

func video(postURL string) (string, error) {
	// Send an HTTP GET request to the Instagram post URL
	resp, err := http.Get(postURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch the Instagram post page")
	}

	// Parse the HTML of the page using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Find the <meta> tag with property="og:video"
	metaTag := doc.Find(`meta[property="og:video"]`)
	if metaTag.Length() == 0 {
		return "", errors.New("video URL not found in the Instagram post")
	}

	// Extract the content attribute, which contains the video URL
	videoURL, exists := metaTag.Attr("content")
	if !exists {
		return "", errors.New("failed to extract video URL")
	}

	return videoURL, nil
}

func Img(postURL string) (string, error) {

	// Получаем ссылку на изображение
	imageURL, err := fetchInstagramImage(postURL)
	if err != nil {
		return "", fmt.Errorf("Failed to fetch Instagram image: %v", err)
	}

	fmt.Printf("Image URL: %s\n", imageURL)

	// Скачиваем изображение
	filename := "downloaded_image.jpg"
	fpath, err := downloadImage(imageURL, filename)
	if err != nil {
		return "", fmt.Errorf("Failed to download image: %v", err)
	}

	fmt.Printf("Image downloaded successfully: %s\n", filename)

	return fpath, nil
}

func downloadImage(url, filename string) (string, error) {
	// Отправляем запрос для загрузки изображения
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Создаем файл
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Копируем данные изображения в файл
	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return file.Name(), err
}

func fetchInstagramImage(postURL string) (string, error) {
	// Отправляем запрос к Instagram
	resp, err := http.Get(postURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Парсим HTML-страницу
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Ищем теги <meta> с изображением
	metaTag := doc.Find("meta[property='og:image']")
	if metaTag.Length() == 0 {
		return "", fmt.Errorf("image not found")
	}

	// Извлекаем URL изображения
	imageURL, exists := metaTag.Attr("content")
	if !exists {
		return "", fmt.Errorf("image URL not found")
	}

	return imageURL, nil
}

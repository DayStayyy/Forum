package fileUpload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// UploadFile : Place file to fileServer/userIMgs and return it's path or default image if no image selected
func UploadFile(r *http.Request, userName string) string {
	file, handler, err := r.FormFile("profilPicture")
	if err != nil {
		fmt.Println(file, handler, err)
		return "../../../fileServer/userImgs/default.png"
	}

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(f, file)
	if err != nil {
		log.Fatal(err)
		return "../../../fileServer/userImgs/default.png"
	}

	filePath := "fileServer/userImgs/" + strings.ReplaceAll(userName, " ", "_") + handler.Filename[strings.LastIndex(handler.Filename, "."):]
	formats := [7]string{".bmp", ".jpg", ".jpeg", ".gif", ".png", ".tif", ".tiff"}

	for i := 0; i < len(formats); i++ {
		if fileExists(filePath[0:strings.LastIndex(filePath, ".")] + formats[i]) {
			os.Remove(filePath[0:strings.LastIndex(filePath, ".")] + formats[i])
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
		return "../../../fileServer/userImgs/default.png"
	}
	
	err = os.Rename(handler.Filename, filePath)
	if err != nil {
		log.Fatal(err)
		return "../../../fileServer/userImgs/default.png"
	}

	

	return "../../../fileServer/userImgs/" + strings.ReplaceAll(userName, " ", "_") + handler.Filename[strings.LastIndex(handler.Filename, "."):]
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

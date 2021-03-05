package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

var  (
	upload_url string = "https://upload.imagekit.io/api/v1/files/upload"
	file_name string = "favicon-516140983.ico"
	api_secret_key string = "PRIVATE_KEY"
)
func UploadMultipartFile(client *http.Client, uri, key, path string) (*http.Response, error) {

	body, writer := io.Pipe()

	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	mwriter := multipart.NewWriter(writer)
	req.Header.Add("Content-Type", mwriter.FormDataContentType())
	req.SetBasicAuth(api_secret_key, "")
	
	errchan := make(chan error)
		

	go func() {
		
		defer close(errchan)

		defer writer.Close()

		defer mwriter.Close()

		
		file, err := os.Open(path)
		if err != nil {
			errchan <- err
	 		return
	 	}
		
		defer file.Close()
	

		err := mwriter.WriteField("fileName", file_name); err != nil {
			errchan <- err	
			return 
		}
	
		w, err := mwriter.CreateFormFile("file", path)
		if err != nil {
			errchan <- err
			return
		}
		
		
		if written, err := io.Copy(w, file); err != nil {

			errchan <- fmt.Errorf("error copying %s (%d bytes written): %v", path, written, err)
			return
		}
	///	log.Println("FIKA")
		if err := mwriter.Close(); err != nil {
			errchan <- err
			return
		}
	}()

	resp, err := client.Do(req)
	//log.Println(err)
	merr := <-errchan

	if err != nil || merr != nil {
		return resp, fmt.Errorf("http error: %v, multipart error: %v", err, merr)
	}

	return resp, nil
}

func main() {
	path, _ := os.Getwd()

	path += "/"+file_name

	client := &http.Client{}

	resp, err := UploadMultipartFile(client, upload_url, "file", path)

	if err != nil {

		log.Println(err)
	} else {

		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)

		_, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
	}
}

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	uploader := NewClientUploader(client)

	if err := uploader.List(ctx, "estudos-bucket"); err != nil {
		panic(err)
	}

	for {
		// files, err := dir.ReadDir(1)
		// if err != nil {
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	log.Printf("Error reading directory: %s\n", err)
		// 	continue
		// }

		// f, err := os.Open(files[0].Name())
		// if err != nil {
		// 	log.Printf("Error reading directory: %s\n", err)
		// 	continue
		// }
		// defer f.Close()

		content, _ := os.ReadFile("/home/jailton/Git/go-storage/tmp/CPF.pdf")

		if err := uploader.Upload(ctx, "estudos-bucket", "documents", "CPF.pdf", "application/pdf", content); err != nil {
			log.Println(err)
			continue
		}
	}

}

type clientUploader struct {
	Client *storage.Client
}

func NewClientUploader(client *storage.Client) *clientUploader {
	return &clientUploader{Client: client}
}

func (c *clientUploader) Upload(ctx context.Context, bucketName, folderPath, fileName, contentType string, content []byte) error {
	buffer := bytes.NewBuffer(content)

	bucket := c.Client.Bucket(bucketName).Object(folderPath + "/" + fileName).NewWriter(ctx)
	bucket.ContentType = contentType
	defer bucket.Close()

	if _, err := io.Copy(bucket, buffer); err != nil {
		return err
	}

	return nil
}

func (c *clientUploader) List(ctx context.Context, bucketName string) error {
	bucket := c.Client.Bucket(bucketName)
	it := bucket.Objects(ctx, &storage.Query{
		Prefix: "documents/CPF",
	})

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("bucket(%v).Objects: %v", bucket, err)
		}
		fmt.Println(attrs.Name)
	}
	return nil
}

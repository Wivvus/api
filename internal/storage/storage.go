package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func newClient() *s3.Client {
	return s3.New(s3.Options{
		Region:       os.Getenv("DO_SPACES_REGION"),
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(os.Getenv("DO_SPACES_KEY"), os.Getenv("DO_SPACES_SECRET"), "")),
		BaseEndpoint: aws.String(os.Getenv("DO_SPACES_ENDPOINT")),
		UsePathStyle: true,
	})
}

func publicURL(objectKey string) string {
	return fmt.Sprintf("https://%s.%s.digitaloceanspaces.com/%s",
		os.Getenv("DO_SPACES_BUCKET"),
		os.Getenv("DO_SPACES_REGION"),
		objectKey,
	)
}

func putObject(body io.Reader, objectKey, contentType string) (string, error) {
	_, err := newClient().PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("DO_SPACES_BUCKET")),
		Key:         aws.String(objectKey),
		Body:        body,
		ContentType: aws.String(contentType),
		ACL:         "public-read",
	})
	if err != nil {
		return "", err
	}
	return publicURL(objectKey), nil
}

// UploadAvatar uploads a multipart file to Spaces and returns the public URL.
func UploadAvatar(file multipart.File, filename, contentType string) (string, error) {
	return putObject(file, fmt.Sprintf("avatars/%s", filename), contentType)
}

// CopyFromURL downloads an image from srcURL and uploads it to Spaces under filename.
func CopyFromURL(srcURL, filename string) (string, error) {
	resp, err := http.Get(srcURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "image/jpeg"
	}

	return putObject(bytes.NewReader(data), fmt.Sprintf("avatars/%s", filename), ct)
}

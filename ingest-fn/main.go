package main

import (
	"context"
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
)

type templateValues struct {
	Domain    string
	Id        string
	ImageFile string
	Width     int
	Height    int
	SiteName  string
	Type      string
}

func panicErr(err error) {
	if err != nil {
		pc, _, lineNum, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc)
		fmt.Printf("[PANIC] %s at %d\n", funcName.Name(), lineNum)
		panic(err)
	}
}

func warnErr(err error) {
	if err != nil {
		pc, _, lineNum, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc)
		fmt.Printf("[WARN] %s at %d\n", funcName.Name(), lineNum)
		fmt.Println(err.Error())
	}
}

func handler(ctx context.Context, s3Event events.S3Event) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	panicErr(err)
	s3Client := s3.NewFromConfig(cfg)
	snsClient := sns.NewFromConfig(cfg)

	siteName := os.Getenv("SITE_NAME")
	domain := os.Getenv("DOMAIN")
	targetBucket := os.Getenv("HOSTING_BUCKET")
	topicArn := os.Getenv("SNS_TOPIC_ARN")

	for _, record := range s3Event.Records {
		extension := strings.ToLower(path.Ext(record.S3.Object.Key))
		if extension != ".jpg" && extension != ".jpeg" && extension != ".png" && extension != ".gif" {
			// Unsupported file - skip
			fmt.Printf("[WARN] File %s skipped - only jpg/jpeg/png/gif file extensions supported\n", record.S3.Object.Key)
			continue
		}
		fmt.Printf("[INFO] Processing %s from bucket %s\n", record.S3.Object.Key, record.S3.Bucket.Name)
		fmt.Printf("[INFO] Decoded %s to %s\n", record.S3.Object.Key, record.S3.Object.URLDecodedKey)

		// Make getObject request for input image
		getObjectInput := &s3.GetObjectInput{
			Bucket: aws.String(record.S3.Bucket.Name),
			Key:    aws.String(record.S3.Object.URLDecodedKey),
		}
		getObjectOutput, err := s3Client.GetObject(ctx, getObjectInput)
		warnErr(err)
		if err != nil {
			// Retry with converted `+` -> ` ` characters
			getObjectInput := &s3.GetObjectInput{
				Bucket: aws.String(record.S3.Bucket.Name),
				Key:    aws.String(strings.Replace(record.S3.Object.URLDecodedKey, "+", " ", -1)),
			}
			getObjectOutput, err = s3Client.GetObject(ctx, getObjectInput)
		}
		panicErr(err)

		// Generate UUID for upload filename/ID
		filename := uuid.New().String()

		// Download input image for dimensions checking later
		outFile, err := os.Create(fmt.Sprintf("/tmp/%s%s", filename, extension))
		panicErr(err)
		byteCount, err := io.Copy(outFile, getObjectOutput.Body)
		fmt.Printf("[INFO] wrote %d bytes\n", byteCount)
		panicErr(err)
		warnErr(getObjectOutput.Body.Close())
		panicErr(outFile.Close())

		// Read downloaded file and assign image attributes
		imageFile, err := os.Open(fmt.Sprintf("/tmp/%s%s", filename, extension))
		panicErr(err)
		imageReader, _, err := image.DecodeConfig(imageFile)
		panicErr(err)
		width := imageReader.Width
		height := imageReader.Height
		imageContentType := fmt.Sprintf("image/%s", extension[1:])
		if extension == ".jpg" {
			imageContentType = "image/jpeg"
		}
		panicErr(imageFile.Close())

		// Generate HTML with image based on template
		html, err := template.ParseFiles("template.gohtml")
		panicErr(err)
		htmlFile, err := os.Create(fmt.Sprintf("/tmp/%s.html", filename))
		panicErr(err)
		templateParams := templateValues{
			Domain:    domain,
			Id:        filename,
			ImageFile: filename + extension,
			Width:     width,
			Height:    height,
			SiteName:  siteName,
			Type:      imageContentType,
		}
		err = html.Execute(htmlFile, templateParams)
		panicErr(err)
		panicErr(htmlFile.Close())

		// Copy input image to hosting bucket
		copyInput := &s3.CopyObjectInput{
			Bucket:               aws.String(targetBucket),
			CopySource:           aws.String(fmt.Sprintf("%s/%s", record.S3.Bucket.Name, record.S3.Object.URLDecodedKey)),
			Key:                  aws.String(filename + "/" + filename + extension),
			ContentType:          aws.String(imageContentType),
			ServerSideEncryption: types.ServerSideEncryptionAes256,
		}
		_, err = s3Client.CopyObject(ctx, copyInput)
		panicErr(err)

		// Upload generated HTML to hosting bucket
		htmlFile, err = os.Open(fmt.Sprintf("/tmp/%s.html", filename))
		panicErr(err)
		putObjectInput := &s3.PutObjectInput{
			Bucket:               aws.String(targetBucket),
			Key:                  aws.String(filename + "/index.html"),
			Body:                 htmlFile,
			ContentType:          aws.String("text/html"),
			ServerSideEncryption: types.ServerSideEncryptionAes256,
		}
		_, err = s3Client.PutObject(ctx, putObjectInput)
		panicErr(err)
		warnErr(htmlFile.Close())

		// Cleanup temporary files
		err = os.Remove(fmt.Sprintf("/tmp/%s%s", filename, extension))
		warnErr(err)
		err = os.Remove(fmt.Sprintf("/tmp/%s.html", filename))
		warnErr(err)

		// Notify SNS subscribers of completion
		publishInput := &sns.PublishInput{
			Message:  aws.String(fmt.Sprintf("Image processing is now complete - available at:\nhttps://%s/%s/", domain, filename)),
			Subject:  aws.String("Image processing complete"),
			TopicArn: aws.String(topicArn),
		}
		_, err = snsClient.Publish(ctx, publishInput)
		panicErr(err)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}

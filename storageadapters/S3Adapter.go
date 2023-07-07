/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 * S3 Storage Adapter Implementation.
 */

package storageadapters

import (
	"bytes"
	"context"
	"io/ioutil"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "k8s.io/klog/v2"
)

var credCache = make(map[string]*aws.CredentialsCache)
var s3ClientCacheMutex sync.Once
var awsS3Client *s3.Client

func getCredCache(key, secret string) *aws.CredentialsCache {
	if credCache[key] == nil {
		credCache[key] = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(key, secret, ""))
	}
	return credCache[key]
}

func getCachedS3Client(ctx context.Context, key, secret, region string) *s3.Client {

	s3ClientCacheMutex.Do(func() {
		awsS3Client = s3.New(s3.Options{
			Region:      region,
			Credentials: getCredCache(key, secret),
		})
	})

	return awsS3Client
}

type IS3Client interface {
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	DeleteBucket(ctx context.Context, params *s3.DeleteBucketInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
}

type AwsS3Impl struct {
	*s3.Client
}

type Config struct {
	AccessKey string
	SecretKey string
	Bucket    string
	EndPoint  string
	Region    string
	URLPrefix string
	URLSuffix string
}

func getS3ClientFromConfig() *s3.Client {
	localS3Config := Config{
		AccessKey: "tW_KDPCzz12MZiokdYRULLXlAaVh9ZuL",
		SecretKey: "OUe2waPYFv2QrP9b2tbVSgc04WkwMrIp",
		Bucket:    "",
		EndPoint:  "http://10.45.48.219",
		Region:    "us-east-1",
		URLPrefix: "",
		URLSuffix: "",
	}

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "http://10.45.48.219",
			SigningRegion: "us-east-1",
			// For some s3-compatible object stores, converting the hostname is not required,
			// and not setting this option will result in not being able to access the corresponding object store address.
			HostnameImmutable: true,
		}, nil
	})

	awsConfig, err := s3config.LoadDefaultConfig(context.TODO(),
		s3config.WithEndpointResolverWithOptions(resolver),
		s3config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(localS3Config.AccessKey, localS3Config.SecretKey, "")),
	)
	if err != nil {
		log.Errorf("Error Occured during Aws Config %v", err)
	}
	s3Client := s3.NewFromConfig(awsConfig)
	return s3Client
}

func (impl *AwsS3Impl) CreateBucket(ctx context.Context, input *s3.CreateBucketInput, optFns ...func(*s3.Options)) (
	*s3.CreateBucketOutput, error) {

	s3client := getS3ClientFromConfig()

	output, err := s3client.CreateBucket(ctx, input)
	if err != nil {
		log.Info("Create bucket error")
		return output, err
	}
	// TODO check if bucket exists waiter is needed
	log.Info("Create bucket output", output)

	return output, nil

}

func (impl *AwsS3Impl) DeleteBucket(ctx context.Context, params *s3.DeleteBucketInput, optFns ...func(*s3.Options)) (
	*s3.DeleteBucketOutput, error) {
	s3client := getS3ClientFromConfig()
	output, err := s3client.DeleteBucket(ctx, params)
	if err != nil {
		log.Info("Delete bucket error")
		return output, err
	}
	// todo check if bucket exists waiter is needed
	log.Info("Delete bucket output", output)
	return output, nil
}

func (impl *AwsS3Impl) CreateWarehouseBucket(ctx context.Context, bucketName string) error {

	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	output, err := impl.CreateBucket(ctx, input)
	if err != nil {
		log.Errorf("Error Occured while creating the bucket %s", err)
		return err
	}
	log.Infof("Bucket got created successfully bucket %s", output.ResultMetadata)
	return nil
}

func (impl *AwsS3Impl) DeleteWarehouseBucket(ctx context.Context, bucketName string) error {

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}
	output, err := impl.DeleteBucket(ctx, input)
	if err != nil {
		log.Errorf("Error Occurred while Deleting the bucket %s", err)
		return err
	}
	log.Infof("Bucket got Deleted successfully bucket %s", output.ResultMetadata)
	return nil
}

func (impl *AwsS3Impl) UploadFileToBucket(ctx context.Context, bucketName string, pathToFile string, data []byte) error {
	objectKey := pathToFile
	s3Client := getS3ClientFromConfig()
	// Upload the file to S3
	output, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		log.Errorf("Error Occurred while uploading file to the bucket %s", err)
		return err
	}

	log.Infof("File uploaded to bucket successfully %s", output.ResultMetadata)
	return nil
}

func (impl *AwsS3Impl) DeleteFileFromBucket(ctx context.Context, bucketName string, pathToFile string) error {
	objectKey := pathToFile
	s3Client := getS3ClientFromConfig()
	// Delete the file from S3
	output, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Errorf("Error Occurred while deleting file from the bucket %s", err)
		return err
	}

	log.Infof("File successfully deleted from the bucket %s", output.ResultMetadata)
	return nil
}

func (impl *AwsS3Impl) DeleteAllFileFromBucket(ctx context.Context, bucketName string) error {
	s3Client := getS3ClientFromConfig()
	// List all objects in the bucket
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	for {
		listOutput, err := s3Client.ListObjectsV2(context.TODO(), listInput)
		if err != nil {
			log.Errorf("Error Occurred while listing files from the bucket %s", err)
			return err
		}

		// Delete each object
		for _, obj := range listOutput.Contents {
			_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    obj.Key,
			})
			if err != nil {
				log.Errorf("Error Occurred while deleting file from the bucket %s", err)
				return err
			}
			log.Infof("Deleted object:", *obj.Key)
		}

		// Check if there are more objects to retrieve
		if !listOutput.IsTruncated {
			break
		}

		// Set the continuation token to fetch the next page of objects
		listInput.ContinuationToken = listOutput.NextContinuationToken
	}

	log.Info("All objects deleted successfully.")
	return nil
}

func (impl *AwsS3Impl) ListFilesInBucket(ctx context.Context, bucketName string, key string) ([]string, error) {
	objectKey := key
	var fileList []string
	s3Client := getS3ClientFromConfig()
	// Delete the file from S3
	output, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(objectKey),
	})
	if err != nil {
		log.Errorf("Error Occurred while listing files from the bucket %s", err)
		return nil, err
	}

	log.Infof("Files successfully listed from the bucket %s", output.ResultMetadata)
	for _, item := range output.Contents {
		fileList = append(fileList, *item.Key)
	}
	return fileList, nil
}

func (impl *AwsS3Impl) GetFileInBucket(ctx context.Context, bucketName string, pathToFile string) ([]byte, error) {
	objectKey := pathToFile
	s3Client := getS3ClientFromConfig()
	// Delete the file from S3
	output, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Errorf("Error Occurred while getting file from the bucket %s", err)
		return nil, err
	}

	log.Infof("Files successfully downloaded from the bucket %s", output.ResultMetadata)
	body, err := ioutil.ReadAll(output.Body)
	defer output.Body.Close()

	if err != nil {
		log.Errorf("Error Occurred while parsing file contents %s", err)
		return nil, err
	}
	return body, nil
}

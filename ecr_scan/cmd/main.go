package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"time"

	//"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type ScanSpec struct {
	ImageID    imageID `locationName:"imageId" type:"structure" required:"true"`
	RegistryID string  `locationName:"registryId" type:"string"`
	Repository string  `locationName:"repositoryName" min:"2" type:"string" required:"true"`
}

type imageID struct {
	ImageDigest string `locationName:"imageDigest" type:"string"`
	ImageTag    string `locationName:"imageDigest" type:"string"`
}

func EcrNewSession() *ecr.ECR {
	svc := ecr.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("ap-southeast-1"))

	return svc
}

func AwsStartScan(svc *ecr.ECR, id string, it string, ri string, rn string) {
	specImage := imageID{
		ImageDigest: id,
		ImageTag:    it,
	}

	spec := &ScanSpec{
		ImageID:    specImage,
		RegistryID: ri,
		Repository: rn,
	}

	input := &ecr.StartImageScanInput{
		ImageId: &ecr.ImageIdentifier{
			ImageTag:    &spec.ImageID.ImageTag,
			ImageDigest: &spec.ImageID.ImageDigest,
		},
		RepositoryName: &spec.Repository,
		RegistryId:     &spec.RegistryID,
	}

	result, err := svc.StartImageScan(input)
	if err != nil {
		if aer, ok := err.(awserr.Error); ok {
			switch aer.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aer.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aer.Error())
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aer.Error())
			default:
				fmt.Println(aer.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}

func AwsListRepo(svc *ecr.ECR) []string {
	input := &ecr.DescribeRepositoriesInput{}

	result, err := svc.DescribeRepositories(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		//return
	}

	var repoName []string

	for _, r := range result.Repositories {
		repoName = append(repoName, *r.RepositoryName)
	}

	return repoName
}

func AWSDescribeImage(svc *ecr.ECR, rn string) (string, string, string, string) {
	repoName := rn
	input := &ecr.DescribeImagesInput{
		RepositoryName: &repoName,
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	var imagePushAt time.Time
	var imageTag string
	var imageDigest string
	var registryID string
	var repository string

	for _, r := range result.ImageDetails {
		if r.ImagePushedAt.After(imagePushAt) == true {
			for _, t := range r.ImageTags {
				imageTag = *t
			}
			imagePushAt = *r.ImagePushedAt
			imageDigest = *r.ImageDigest
			registryID = *r.RegistryId
			repository = *r.RepositoryName
		}
	}

	return imageDigest, imageTag, registryID, repository
}

func HandleLambda() {
	svc := EcrNewSession()
	repoArray := AwsListRepo(svc)

	for _, r := range repoArray {
		imageDigest, imageTag, registryID, repository := AWSDescribeImage(svc, r)
		AwsStartScan(svc, imageDigest, imageTag, registryID, repository)
	}
}

func main() {
	lambda.Start(HandleLambda)
}

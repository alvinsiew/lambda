package main

import (
	"fmt"
	"sort"

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

func HandlerStartScan(svc *ecr.ECR) {
	specImage := imageID{
		ImageDigest: "sha256:1e1a915f208cab016a981212218d465053572007b9122b5e149d16830e33751b",
		ImageTag:    "7-2.7.1-12.22.1",
	}

	spec := &ScanSpec{
		ImageID:    specImage,
		RegistryID: "355716222559",
		Repository: "centos-ruby-node",
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

func HandlerListRepo(svc *ecr.ECR) []string {
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

func HandleListImage(svc *ecr.ECR, rn string) *ecr.ListImagesOutput {
	//svc := ecr.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("ap-southeast-1"))
	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(rn),
	}

	result, err := svc.ListImages(input)
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

	//fmt.Println("result = ", reflect.TypeOf(result))
	//fmt.Println(result)

	return result
}

func main() {
	//lambda.Start(HandlerStartScan)
	svc := EcrNewSession()
	//HandleListImage(svc)
	abc := HandlerListRepo(svc)

	for _, r := range abc {
		fmt.Println(sort.Strings(string(HandleListImage(svc, r))))
	}

}

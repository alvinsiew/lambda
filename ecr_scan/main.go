package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"

	//"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

//type ScanSpec struct {
//	ImageID ImageID `json:"imageId"`
//	RegistryID string `json:"registry"`
//	Repository string `json:"repository"`
//}

type ScanSpec struct {
	//ImageID ImageIdentifier `locationName:"imageId" type:"structure" required:"true"`
	RegistryID string `locationName:"registryId" type:"string"`
	Repository string `locationName:"repositoryName" min:"2" type:"string" required:"true"`
}

type imageID struct {
	ImageDigest string `locationName:"imageDigest" type:"string"`
	ImageTag    string `locationName:"imageDigest" type:"string"`
}

//type ImageIdentifier struct {
//
//	// The sha256 digest of the image manifest.
//	ImageDigest string `locationName:"imageDigest" type:"string"`
//
//	// The tag used for the image.
//	ImageTag string `locationName:"imageTag" min:"1" type:"string"`
//	// contains filtered or unexported fields
//}

func ecrStartScan() {
	svc := ecr.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("ap-southeast-1"))

	//ImageDigest := "sha256:80d8b356e087631ac21bdf5aa51e4917bef73baa4e398a89a57c9e26f2fec342"
	//ImageTag := "latest"

	specImage := &imageID{
		ImageDigest: "sha256:1e1a915f208cab016a981212218d465053572007b9122b5e149d16830e33751b",
		ImageTag:    "latest",
	}

	id := aws.String(specImage.ImageDigest)
	it := aws.String(specImage.ImageTag)

	//fmt.Println(ImageDigest)

	spec := &ScanSpec{
		//ImageID: imageSpec,
		RegistryID: "355716222559",
		Repository: "centos-ruby-node",
	}

	//test := &ecr.ImageIdentifier{ImageTag: id, ImageDigest: it}
	//svc := ecr.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("ap-southeast-1"))
	input := &ecr.StartImageScanInput{
		ImageId: &ecr.ImageIdentifier{
			ImageTag:    it,
			ImageDigest: id,
		},
		RepositoryName: &spec.Repository,
		RegistryId:     &spec.RegistryID,
	}

	//input := &ecr.BatchDeleteImageInput{
	//	ImageIds: []*ecr.ImageIdentifier{
	//		{
	//			ImageTag: aws.String("precise"),
	//		},
	//	},
	//	RepositoryName: aws.String("ubuntu"),
	//}

	result, err := svc.StartImageScan(input)
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
		return
	}
	fmt.Println(result)
}

func LambdaHandler() {
	svc := ecr.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("ap-southeast-1"))
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
		return
	}

	fmt.Println(result)
}

func ecrImages() *ecr.ListImagesOutput {
	svc := ecr.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("ap-southeast-1"))
	input := &ecr.ListImagesInput{
		RepositoryName: aws.String("dbserver"),
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
	//ecrImages()

	//abc := ecrImages()
	lambda.Start(ecrStartScan)

	//ecrRepositories()
	//ecrStartScan()
}

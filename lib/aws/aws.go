package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go/aws"
)

func Init() {

	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	cognitoClient := cognitoidentityprovider.NewFromConfig(sdkConfig)

	fmt.Println("Let's list the user pools for your account.")
	var pools []types.UserPoolDescriptionType
	paginator := cognitoidentityprovider.NewListUserPoolsPaginator(
		cognitoClient,
		&cognitoidentityprovider.ListUserPoolsInput{
			MaxResults: aws.Int32(10),
		},
	)
	for paginator.HasMorePages() {

		output, err := paginator.NextPage(ctx)
		if err != nil {
			log.Printf("Couldn't get user pools. Here's why: %v\n", err)
		} else {
			pools = append(pools, output.UserPools...)
		}
	}

	if len(pools) == 0 {
		fmt.Println("You don't have any user pools!")
	} else {
		for _, pool := range pools {
			fmt.Printf("\t%v: %v\n", *pool.Name, *pool.Id)
		}
	}
}

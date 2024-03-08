package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type BucketmanagementStackProps struct {
	awscdk.StackProps
}


func NewBucketmanagementStack(scope constructs.Construct, id string, props *BucketmanagementStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	cloudEngineer := awsiam.NewRole(stack, jsii.String("CloudEngineerRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewArnPrincipal(jsii.String("arn:aws:iam::654654234409:root")),
	})

	dataEngineer := awsiam.NewRole(stack, jsii.String("DataEngineerRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewArnPrincipal(jsii.String("arn:aws:iam::654654234409:root")),
	})

	securityEngineer := awsiam.NewRole(stack, jsii.String("SecurityEngineerRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewArnPrincipal(jsii.String("arn:aws:iam::654654234409:root")),
	})

	// attach policies to USERS
	// should actually use json which is more convinent
	cloudEngineer.AttachInlinePolicy(awsiam.NewPolicy(stack, jsii.String("PolicyCreateBucket"), &awsiam.PolicyProps{
		Statements: &([]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(
				&awsiam.PolicyStatementProps{
					Actions: &([]*string{
						jsii.String("s3:CreateBucket"),
						jsii.String("s3:PutLifecycleConfiguration"),
						jsii.String("s3:GetLifecycleConfiguration"),
						jsii.String("s3:PutReplicationConfiguration"),
						jsii.String("s3:GetReplicationConfiguration"),
						jsii.String("s3:GetEncryptionConfiguration"),
						jsii.String("s3:PutEncryptionConfiguration"),
					}),
					Resources: &([]*string{jsii.String("arn:aws:s3:::*")}),
					Effect:    awsiam.Effect_ALLOW},
			),
		}),
	}))

	dataEngineer.AttachInlinePolicy(awsiam.NewPolicy(stack, jsii.String("PolicyReadData"), &awsiam.PolicyProps{
		Statements: &([]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{ // for data analytics
				Actions: &([]*string{
					jsii.String("s3:GetObject"), jsii.String("s3:List*"),
				}), Resources: &([]*string{jsii.String("arn:aws:s3:::*")}), Effect: awsiam.Effect_ALLOW}),
		}),
	}))

	securityEngineer.AttachInlinePolicy(
		awsiam.NewPolicy(stack, jsii.String("PolicyChangeKey"), &awsiam.PolicyProps{
			Statements: &([]awsiam.PolicyStatement{
				awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{ // for security engineer
					Actions: &([]*string{
						jsii.String("s3:PutEncryptionConfiguration"),
					}), Resources: &([]*string{jsii.String("arn:aws:s3:::*")}), Effect: awsiam.Effect_ALLOW}),
			}),
		}))

	var expireDuration float64 = 7
	lifeCycleRuleBucket := []*awss3.LifecycleRule{
		{
			Expiration: awscdk.Duration_Days(&expireDuration),
		},
	}
	// https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/v2/awss3#NewBucket
	awss3.NewBucket(stack, jsii.String("CustomerPIIBucket"), &awss3.BucketProps{
		LifecycleRules: &lifeCycleRuleBucket,       // only keep the data for 7 days
		Encryption:     awss3.BucketEncryption_KMS, // encrypting Data-at-Rest & Data-in-transit
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewBucketmanagementStack(app, "BucketmanagementStack", &BucketmanagementStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }

}

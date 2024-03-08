package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	// "github.com/aws/jsii-runtime-go"
)

// example tests. To run these tests, uncomment this file along with the
// example resource in bucketmanagement_test.go
func TestBucketmanagementStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewBucketmanagementStack(app, "MyStack", nil)

	// THEN
	template := assertions.Template_FromStack(stack, nil)

	roleType := "AWS::IAM::Role"
	rawPropsSecurity := `                    "assumeRolePolicyDocument": {
		"Statement": [
		  {
			"Action": "sts:AssumeRole",
			"Effect": "Allow",
			"Principal": {
			  "AWS": "arn:aws:iam::654654234409:root"
			}
		  }
		],
		"Version": "2012-10-17"
	  }`


	  var documentAsString  map[string]interface{}

	  json.Unmarshal([]byte(rawPropsSecurity), &documentAsString)

	// test the security role 
	template.HasResourceProperties( &roleType, documentAsString)

	// template.HasResourceProperties(jsii.String("AWS::SQS::Queue"), map[string]interface{}{
	// 	"VisibilityTimeout": 300,
	// })
}

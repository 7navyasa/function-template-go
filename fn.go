package main

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured"
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/composed"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"
	"github.com/7navyasa/function-template-go/input/v1beta1"
)

// InputParameters defines the input parameters for the Composition Function.
type InputParameters struct {
	TableName             string                 `json:"tableName"`
	BillingMode           string                 `json:"billingMode"`
	AttributeDefinitions  []AttributeDefinition  `json:"attributeDefinitions"`
	KeySchema             []KeySchema            `json:"keySchema"`
	ProvisionedThroughput *ProvisionedThroughput `json:"provisionedThroughput,omitempty"`
	EnableKinesisStream   bool                   `json:"enableKinesisStream"`
	KinesisStreamName     string                 `json:"kinesisStreamName,omitempty"`
}

// AttributeDefinition represents an attribute definition for the DynamoDB table.
type AttributeDefinition struct {
	AttributeName string `json:"attributeName"`
	AttributeType string `json:"attributeType"`
}

// KeySchema represents the key schema for the DynamoDB table.
type KeySchema struct {
	AttributeName string `json:"attributeName"`
	KeyType       string `json:"keyType"`
}

// ProvisionedThroughput represents the provisioned throughput settings for the DynamoDB table.
type ProvisionedThroughput struct {
	ReadCapacityUnits  int64 `json:"readCapacityUnits"`
	WriteCapacityUnits int64 `json:"writeCapacityUnits"`
}

// OutputResources defines the output resources for the Composition Function.
type OutputResources struct {
	DynamoDBTable resource.Composed  `json:"dynamoDBTable"`
	KinesisStream *resource.Composed `json:"kinesisStream,omitempty"`
}

// IntermediateResources defines any intermediate resources needed for the Composition Function.
type IntermediateResources struct {
	// Define intermediate resources if needed
}

// Compose is the main function that composes the resources based on the input parameters.
func Compose(ctx context.Context, in InputParameters) (*OutputResources, error) {
	// Create the DynamoDB table resource
	dynamoDBTable := composed.NewComposed(
		&unstructured.UnstructuredExtensionObject{
			Spec: unstructured.UnstructuredExtension{
				Object: &unstructured.UnstructuredObject{
					"apiVersion": "dynamodb.aws.crossplane.io/v1alpha1",
					"kind":       "Table",
					"metadata": map[string]interface{}{
						"name": in.TableName,
					},
					"spec": map[string]interface{}{
						"billingMode":           in.BillingMode,
						"attributeDefinitions":  in.AttributeDefinitions,
						"keySchema":             in.KeySchema,
						"provisionedThroughput": in.ProvisionedThroughput,
					},
				},
			},
		},
	)

	var kinesisStream *resource.Composed
	if in.EnableKinesisStream {
		// Create the Kinesis Stream resource
		kinesisStream = composed.NewComposed(
			&unstructured.UnstructuredExtensionObject{
				Spec: unstructured.UnstructuredExtension{
					Object: &unstructured.UnstructuredObject{
						"apiVersion": "streams.aws.crossplane.io/v1alpha1",
						"kind":       "Stream",
						"metadata": map[string]interface{}{
							"name": in.KinesisStreamName,
						},
						"spec": map[string]interface{}{
							"streamModeDetails": map[string]interface{}{
								"streamMode": "PROVISIONED",
							},
						},
					},
				},
			},
		)
	}

	return &OutputResources{
		DynamoDBTable: dynamoDBTable,
		KinesisStream: kinesisStream,
	}, nil
}

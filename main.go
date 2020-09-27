package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := dynamodb.NewTable(ctx, "sucursal_table", &dynamodb.TableArgs{
			Attributes: dynamodb.TableAttributeArray{
				&dynamodb.TableAttributeArgs{
					Name: pulumi.String("id"),
					Type: pulumi.String("S"),
				},
			},
			Name:          pulumi.String("sucursal_table"),
			BillingMode:   pulumi.String("PROVISIONED"),
			HashKey:       pulumi.String("id"),
			ReadCapacity:  pulumi.Int(5),
			WriteCapacity: pulumi.Int(5),
		})
		if err != nil {
			return err
		}
		return nil
	})
}

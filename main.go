package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := secretsmanager.New(sess)

	// Retrieve all secrets from Secrets Manager if no secretName is specified
	secretName := os.Getenv("secretName")
	var secrets []*secretsmanager.SecretListEntry
	// var err error
	if secretName == "" {
		secretsOutput, err := svc.ListSecrets(&secretsmanager.ListSecretsInput{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		secrets = secretsOutput.SecretList
	} else {
		// Retrieve the specified secret from Secrets Manager
		secret, err := svc.DescribeSecret(&secretsmanager.DescribeSecretInput{
			SecretId: aws.String(secretName),
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		secrets = append(secrets, &secretsmanager.SecretListEntry{
			ARN:  secret.ARN,
			Name: secret.Name,
		})
	}

	// Change these values to match your use case
	newKey := "test2"
	newValue := "test2"
	keyToDelete := "test1"

	// Loop through each secret and update as necessary
	for _, secret := range secrets {
		// Retrieve the current secret value
		currentSecret, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
			SecretId: secret.ARN,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Decode the secret value from JSON to a map[string]string
		secretMap := make(map[string]string)
		err = json.Unmarshal([]byte(*currentSecret.SecretString), &secretMap)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Add new key/value to the secret
		if newKey != "" && newValue != "" {
			secretMap[newKey] = newValue
		}

		// Modify existing key/value in the secret
		if keyToDelete != "" {
			delete(secretMap, keyToDelete)
		}

		// Encode the secret map back into JSON format
		newSecret, err := json.Marshal(secretMap)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Update the secret value in Secrets Manager
		_, err = svc.UpdateSecret(&secretsmanager.UpdateSecretInput{
			SecretId:     secret.ARN,
			SecretString: aws.String(string(newSecret)),
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Secret %s has been updated\n", *secret.Name)
	}
}

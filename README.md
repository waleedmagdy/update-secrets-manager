AWS Secrets Manager Update Script

This script updates the values of all or a specific secret in AWS Secrets Manager.

Usage:
  To update a specific secret:
    $ secretName=my-secret go run main.go

  To update all secrets:
    $ go run main.go

Behavior:
  - If the specified secretName environment variable is not set, the script will update all secrets in Secrets Manager.
  - If the specified secretName environment variable is set, the script will only update that specific secret.
  - The script can add new key/value to all secrets or to a specific secret.
  - The script can modify a key/value in all secrets or in a specific secret.
  - The script can delete a key/value in all secrets or in a specific secret.

Configuration:
  Modify the newKey, newValue, and keyToDelete variables in the script to match your use case.

Dependencies:
  - AWS SDK for Go (github.com/aws/aws-sdk-go)

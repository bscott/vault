package command

import (
	"fmt"
	"strings"
)

// RotateCommand is a Command that rotates the encryption key being used
type RotateCommand struct {
	Meta
}

func (c *RotateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("rotate", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	// Rotate the key
	err = client.Sys().Rotate()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error with key rotation: %s", err))
		return 2
	}

	// Print the key status
	status, err := client.Sys().KeyStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error reading audits: %s", err))
		return 2
	}

	c.Ui.Output(fmt.Sprintf("Key Term: %d", status.Term))
	c.Ui.Output(fmt.Sprintf("Installation Time: %v", status.InstallTime))
	return 0
}

func (c *RotateCommand) Synopsis() string {
	return "Rotates the backend encryption key used to persist data"
}

func (c *RotateCommand) Help() string {
	helpText := `
Usage: vault rotate [options]

  Rotates the backend encryption key which is used to secure data
  written to the storage backend. This is done by installing a new key
  which encrypts new data, while old keys are still used to decrypt
  secrets written previously. This is an online operation and is not
  disruptive.

General Options:

  -address=addr           The address of the Vault server.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-path is used.

  -tls-skip-verify        Do not verify TLS certificate. This is highly
                          not recommended. This is especially not recommended
                          for unsealing a vault.

`
	return strings.TrimSpace(helpText)
}

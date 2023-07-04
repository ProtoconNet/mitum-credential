package cmds

type CredentialCommand struct {
	CreateCredentialService CreateCredentialServiceCommand `cmd:"" name:"create-credential-service" help:"register credential service to contract account"`
	AddTemplate             AddTemplateCommand             `cmd:"" name:"add-template" help:"add template to credential service"`
	AssignCredentials       AssignCredentialsCommand       `cmd:"" name:"assign-credential" help:"assign credential"`
	RevokeCredentials       RevokeCredentialsCommand       `cmd:"" name:"revoke-credential" help:"revoke credential"`
}

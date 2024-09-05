package cmds

type CredentialCommand struct {
	RegisterModel RegisterModelCommand     `cmd:"" name:"register-model" help:"register credential service to contract account"`
	AddTemplate   AddTemplateCommand       `cmd:"" name:"add-template" help:"add template to credential service"`
	Issue         IssueCommand             `cmd:"" name:"issue" help:"issue credential"`
	Revoke        RevokeCredentialsCommand `cmd:"" name:"revoke" help:"revoke credential"`
}

package cmds

type CredentialCommand struct {
	CreateService CreateServiceCommand     `cmd:"" name:"create-service" help:"register credential service to contract account"`
	AddTemplate   AddTemplateCommand       `cmd:"" name:"add-template" help:"add template to credential service"`
	Assign        AssignCommand            `cmd:"" name:"assign" help:"assign credential"`
	Revoke        RevokeCredentialsCommand `cmd:"" name:"revoke" help:"revoke credential"`
}

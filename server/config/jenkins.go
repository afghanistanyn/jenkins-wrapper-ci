package config

type JenkinsConfig struct {
	Url 		string		`mapstructure:"url" json:"url" yaml:"url"`
	UserName	string		`mapstructure:"username" json:"username" yaml:"username"`
	Password    string		`mapstructure:"password" json:"password" yaml:"password"`
	Timeout		uint		`mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}

type GitCredential struct {
	GitServer			string	`mapstructure:"gitServer" json:"gitServer" yaml:"gitServer"`
	GitCredentialId    	string 	`mapstructure:"gitCredentialId" json:"gitCredentialId" yaml:"gitCredentialId"`
}
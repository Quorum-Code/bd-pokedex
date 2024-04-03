package config

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Clicfg, []string) error
}

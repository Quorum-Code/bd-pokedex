package config

type Clicfg struct {
	Cache    Cache
	Commands map[string]CliCommand
	MapLast  *string
	MapNext  *string
	MapPrev  *string
}

func NewClicfg() *Clicfg {
	return &Clicfg{}
}

package core

type TaskFile struct {
	Version string                  `yaml:"version"`
	Dotenv  []string                `yaml:"dotenv"`
	Tasks   map[string]TaskCommands `yaml:"tasks"`
}

type TaskCommands struct {
	Cmds []string `yaml:"cmds"`
}

package program

import (
	"fmt"
	"os"
)

type Topic struct {
	Name      string
	Output    string
	Languages []string
	Exit      bool
}

func (p *Topic) String() string {
	return fmt.Sprintf("Exit: %t, Name: %s, OutputPath: %s, Languages: %v", p.Exit, p.Name, p.Output, p.Languages)
}

func (p *Topic) ExitCLI() {
	if p.Exit {
		os.Exit(1)
	}
}

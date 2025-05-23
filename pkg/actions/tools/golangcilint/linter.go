package golangcilint

import (
	"encoding/json"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/style"
)

// ActionLinters completes linters
//
//	deadcode (Finds unused code)
//	dogsled (Checks assignments with too many blank identifiers)
func ActionLinters(config string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		args := []string{"linters", "--json"}
		if config != "" {
			args = append(args, "--config", config)
		}
		return carapace.ActionExecCommand("golangci-lint", args...)(func(output []byte) carapace.Action {
			var formatters struct {
				Enabled []struct {
					Name        string
					Description string
					Deprecated  bool
				}
				Disabled []struct {
					Name        string
					Description string
					Deprecated  bool
				}
			}
			if err := json.Unmarshal(output, &formatters); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			vals := make([]string, 0)
			for _, enabled := range formatters.Enabled {
				vals = append(vals, enabled.Name, enabled.Description, style.Carapace.KeywordPositive)
			}
			for _, disabled := range formatters.Disabled {
				vals = append(vals, disabled.Name, disabled.Description, style.Carapace.KeywordNegative)
			}
			return carapace.ActionStyledValuesDescribed(vals...)
		})
	}).Tag("linters")
}

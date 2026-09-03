// app/wizard.go
package app

type StepResult struct {
	Value string
	Err   error
}

type Step struct {
	Question string
	Validate func(input string) error // return non-nil to re-prompt
}

// RunWizard asks each step in order, retrying on validation failure,
// and returns collected answers keyed by step index or a caller-assigned name.
func (c *Controller) RunWizard(steps []Step) ([]string, bool) {
	answers := make([]string, len(steps))
	for i, step := range steps {
		for {
			input := c.ui.Prompt(step.Question)
			if input == "" {
				return nil, false // treat blank as cancel
			}
			if step.Validate != nil {
				if err := step.Validate(input); err != nil {
					c.ui.ShowError(err)
					continue
				}
			}
			answers[i] = input
			break
		}
	}
	return answers, true
}

//planned better system with a generic Wizard struct that can be used for more complex flows, but for now, this simple function is sufficient.

package attributes

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateAndReplaceForCommands validates the commands data for global attribute references and replaces them with the attribute value
func ValidateAndReplaceForCommands(attributes map[string]string, commands []v1alpha2.Command) error {

	for i := range commands {
		var err error

		// Validate various command types
		switch {
		case commands[i].Exec != nil:
			if err = validateAndReplaceForExecCommand(attributes, commands[i].Exec); err != nil {
				return err
			}
		case commands[i].Composite != nil:
			if err = validateAndReplaceForCompositeCommand(attributes, commands[i].Composite); err != nil {
				return err
			}
		case commands[i].Apply != nil:
			if err = validateAndReplaceForApplyCommand(attributes, commands[i].Apply); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForExecCommand validates the exec command data for global attribute references and replaces them with the attribute value
func validateAndReplaceForExecCommand(attributes map[string]string, exec *v1alpha2.ExecCommand) error {
	var err error

	if exec != nil {
		// Validate exec command line
		if exec.CommandLine, err = validateAndReplaceDataWithAttribute(exec.CommandLine, attributes); err != nil {
			return err
		}

		// Validate exec working dir
		if exec.WorkingDir, err = validateAndReplaceDataWithAttribute(exec.WorkingDir, attributes); err != nil {
			return err
		}

		// Validate exec label
		if exec.Label, err = validateAndReplaceDataWithAttribute(exec.Label, attributes); err != nil {
			return err
		}

		// Validate exec env
		if len(exec.Env) > 0 {
			if err = validateAndReplaceForEnv(attributes, exec.Env); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForCompositeCommand validates the composite command data for global attribute references and replaces them with the attribute value
func validateAndReplaceForCompositeCommand(attributes map[string]string, composite *v1alpha2.CompositeCommand) error {
	var err error

	if composite != nil {
		// Validate composite label
		if composite.Label, err = validateAndReplaceDataWithAttribute(composite.Label, attributes); err != nil {
			return err
		}
	}

	return nil
}

// validateAndReplaceForApplyCommand validates the apply command data for global attribute references and replaces them with the attribute value
func validateAndReplaceForApplyCommand(attributes map[string]string, apply *v1alpha2.ApplyCommand) error {
	var err error

	if apply != nil {
		// Validate apply label
		if apply.Label, err = validateAndReplaceDataWithAttribute(apply.Label, attributes); err != nil {
			return err
		}
	}

	return nil
}

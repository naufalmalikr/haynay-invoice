package bot

type outputFunc func(user *user) []string
type dynamicActionFunc func(step *step, user *user) *action
type editedStepFunc func(user *user) *step
type step struct {
	Name            string
	DynamicAction   dynamicActionFunc
	ImmediateAction *action
	ChoicesAction   map[string]*action
	FreeTextAction  *action
	Output          outputFunc

	editedStep editedStepFunc
}

func newStep(name string) *step {
	return &step{
		Name:          name,
		ChoicesAction: map[string]*action{},
	}
}

func (s *step) asDynamic(dynamicAction dynamicActionFunc) *step {
	s.DynamicAction = dynamicAction
	return s
}

func (s *step) immediateTo(
	nextStep string,
	process processFunc,
	rollback rollbackFunc,
) *step {
	s.ImmediateAction = &action{
		NextStep: nextStep,
		Process:  process,
		Rollback: rollback,
	}

	return s
}

func (s *step) allowFreeTextTo(
	nextStep string,
	process processFunc,
	rollback rollbackFunc,
) *step {
	s.FreeTextAction = &action{
		NextStep: nextStep,
		Process:  process,
		Rollback: rollback,
	}

	return s
}

func (s *step) rejectFreeText() *step {
	mustErrProcess := func(input string, user *user) error { return errInvalidChoice }
	return s.allowFreeTextTo("", mustErrProcess, noRollback)
}

func (s *step) addChoiceTo(
	nextStep string,
	process processFunc,
	rollback rollbackFunc,
	keywords ...string,
) *step {
	action := &action{
		NextStep: nextStep,
		Process:  process,
		Rollback: rollback,
	}

	for _, keyword := range keywords {
		s.ChoicesAction[keyword] = action
	}

	return s
}

func (s *step) cancelable() *step {
	return s.addChoiceTo("cancel", noProcess, noRollback, "batal", "cancel")
}

func (s *step) setEditedStep(editedStep editedStepFunc) *step {
	s.editedStep = editedStep
	return s.addChoiceTo("edit", noProcess, noRollback, "ubah", "edit")
}

func (s *step) editableFor(editedStep *step) *step {
	return s.setEditedStep(func(user *user) *step {
		return editedStep
	})
}

func (s *step) setOutput(output outputFunc) *step {
	s.Output = output
	return s
}

func (s *step) simpleOutput(message string) *step {
	return s.setOutput(func(user *user) []string {
		return []string{message}
	})
}

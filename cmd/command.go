package cmd

import "flag"

type FlagAppender func(set *flag.FlagSet, cf *CommandFlag) *interface{}
type FlagValidator func(cf *CommandFlag) error
type CommandAction func(c Command)

type CommandFlag struct {
	Name      string
	appender  FlagAppender
	validator FlagValidator
	valuePtr  *interface{}
}

func (receiver *CommandFlag) AppendMyself(targetFlagSet *flag.FlagSet) {
	receiver.valuePtr = receiver.appender(targetFlagSet, receiver)
}

func (receiver *CommandFlag) AddValidator(validator FlagValidator) {
	receiver.validator = validator
}

func (receiver *CommandFlag) ValidateMyself() error {
	if receiver.validator != nil {
		return receiver.validator(receiver)
	}
	return nil
}

func (receiver *CommandFlag) GetValuePtr() *interface{} {
	return receiver.valuePtr
}

func NewFlag(name string, appender FlagAppender) *CommandFlag {
	f := CommandFlag{Name: name, appender: appender}
	return &f
}

type Command struct {
	Name    string
	flagSet *flag.FlagSet
	Flags   []*CommandFlag
	action  CommandAction
}

func NewCommand(name string) *Command {
	f := flag.NewFlagSet(name, flag.ExitOnError)
	c := Command{Name: name, flagSet: f}
	return &c
}

func (c *Command) AddFlag(cf *CommandFlag) {
	c.Flags = append(c.Flags, cf)
	cf.AppendMyself(c.flagSet)
}

func (c Command) GetFlag(name string) *CommandFlag {
	for _, f := range c.Flags {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (c Command) Parse(args []string) (*flag.FlagSet, error) {
	err := c.flagSet.Parse(args)
	if err != nil {
		return nil, err
	} else {
		for _, f := range c.Flags {
			err := f.ValidateMyself()
			if err != nil {
				return nil, err
			}
		}
	}
	return c.flagSet, nil
}

func (c *Command) AddAction(action CommandAction) {
	c.action = action
}

func (c Command) Execute() {
	c.action(c)
}

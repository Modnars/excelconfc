package compiler

type Options struct {
	fileName      string
	sheetName     string
	groupFlag     uint8
	containerType uint8
}

type Option func(*Options)

func WithFileName(fileName string) Option {
	return func(o *Options) {
		o.fileName = fileName
	}
}

func WithSheetName(sheetName string) Option {
	return func(o *Options) {
		o.sheetName = sheetName
	}
}

func WithGroupFlag(groupFlag uint8) Option {
	return func(o *Options) {
		o.groupFlag = groupFlag
	}
}

func WithContainerType(containerType uint8) Option {
	return func(o *Options) {
		o.containerType = containerType
	}
}

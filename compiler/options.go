package compiler

type Options struct {
	filePath  string
	sheetName string
	enumSheet string
	outDir    string
	goPackage string
	addEnum   bool
}

type Option func(*Options)

func WithFilePath(filePath string) Option {
	return func(o *Options) {
		o.filePath = filePath
	}
}

func WithSheetName(sheetName string) Option {
	return func(o *Options) {
		o.sheetName = sheetName
	}
}

func WithEnumSheet(enumSheet string) Option {
	return func(o *Options) {
		o.enumSheet = enumSheet
	}
}

func WithOutDir(outDir string) Option {
	return func(o *Options) {
		o.outDir = outDir
	}
}

func WithGoPackage(goPackage string) Option {
	return func(o *Options) {
		o.goPackage = goPackage
	}
}

func WithAddEnum(addEnum bool) Option {
	return func(o *Options) {
		o.addEnum = addEnum
	}
}

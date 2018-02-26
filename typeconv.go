package fproto_gowrap

type TypeConverter interface {
	GetSources() []TypeConverterSource
	GetImports() []string
}

type TypeConverterSource struct {
	FilePath    string
	PackageName string
}

type ImportAliasGetter interface {
	GetAlias(importname string, defaultalias string) string
}

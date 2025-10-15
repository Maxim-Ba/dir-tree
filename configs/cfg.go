package configs

type FormatCfg struct {
	//Type of output (json, yml, xml)
	Type string
	//Output file
	OutputPath string
}

type Config struct {
	// Path to target dir
	Path string
	// Exclude types
	ExcludeTypes []string
	// Exclude paths (regexp)
	ExcludePaths []string
	Format       FormatCfg
	// Tree depth
	Depth int
}

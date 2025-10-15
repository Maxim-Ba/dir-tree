package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Maxim-Ba/dir-tree/configs"
	"github.com/Maxim-Ba/dir-tree/formatter"
	"github.com/Maxim-Ba/dir-tree/tree"
)

func main() {

	cfg, err := configs.ParseConfig()
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}


	root, err := tree.BuildTree(
		tree.BuildOptions{Path: cfg.Path,
			MaxDepth:     cfg.MaxDepth,
			ExcludePaths: cfg.ExcludePaths,
			ExcludeTypes: cfg.ExcludePaths,
			IncludeFiles: cfg.IncludeFiles,
			FollowLinks:  cfg.FollowLinks,
		})
	if err != nil {
		log.Fatalf("Error building tree: %v", err)
	}

	
	formattedOutput, err := formatter.Format(root, &cfg.Format)
	if err != nil {
		log.Fatalf("Error formatting tree: %v", err)
	}

	if err := saveOutput(formattedOutput, &cfg.Format); err != nil {
		log.Fatalf("Error saving output: %v", err)
	}

}
func saveOutput(data []byte, format *configs.FormatCfg) error {
	outputPath := format.OutputPath
	if outputPath == "" {
		// Вывод в stdout
		fmt.Println(string(data))
		return nil
	}

	ext := fmt.Sprintf(".%s", string(format.Type))
	if !hasExtension(outputPath, ext) {
		outputPath += ext
	}

	err := os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Tree successfully written to: %s\n", outputPath)
	return nil
}

func hasExtension(path, ext string) bool {
	return len(path) >= len(ext) && path[len(path)-len(ext):] == ext
}

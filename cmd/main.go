package main

import (
	"log"

	"github.com/keshon/pixelita/internal/config"
	"github.com/keshon/pixelita/internal/pixelita"

	"github.com/spf13/cobra"
)

func main() {
	config := config.NewConfig()
	px := pixelita.NewPixelita(config)

	cobra.MousetrapHelpText = ""

	var runCLI bool
	var jpgOnly bool
	var pngOnly bool
	var toWebp bool
	var quality int

	rootCmd := &cobra.Command{
		Use:   "pixelita",
		Short: "Pixelita is JPG and PNG compressor and converter to WebP",
		Run: func(cmd *cobra.Command, args []string) {
			if runCLI {
				px.StartCLI(jpgOnly, pngOnly, toWebp, quality)
			} else {
				px.StartGUI()
			}
		},
	}

	rootCmd.Flags().BoolVar(&runCLI, "cli", false, "Run in CLI mode")
	rootCmd.Flags().BoolVarP(&jpgOnly, "jpg", "j", false, "Optimize JPEG files only")
	rootCmd.Flags().BoolVarP(&pngOnly, "png", "p", false, "Optimize PNG files only")
	rootCmd.Flags().BoolVarP(&toWebp, "webp", "w", false, "Convert images to WebP format")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 80, "Compression ratio for JPEG or WebP: 1-100")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

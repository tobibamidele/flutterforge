package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flutterforge",
	Short: "FlutterForge - Flutter project scaffolder with feature-first architecture",
	Long: `FlutterForge helps you create Flutter projects with a standardized,
feature-first architecture and common setup tasks pre-configured.

Features:
- Feature-first directory structure
- go_router setup
- Firebase/Supabase integration
- Riverpod state management
- Custom theming with Material 3
- Google Fonts integration
- FCM setup
- Auth and Onboarding scaffolding`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

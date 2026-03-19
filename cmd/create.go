package cmd

import (
	"fmt"
	"os"

	"flutterforge/scaffold"
	"github.com/spf13/cobra"
)

var (
	projectName     string
	fontName        string
	backend         string
	backgroundColor string
	modalColor      string
	elevatedColor   string
	outlinedColor   string
	textColor       string
	primaryColor    string
	secondaryColor  string
	tertiaryColor   string
	orgIdentifier   string
	withOnboarding  bool
	withFcm         bool
)

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new Flutter project with feature-first architecture",
	Long: `Create a new Flutter project with your custom setup.

Examples:
  # Create with Firebase
  flutterforge create my_app --firebase --font Roboto
  
  # Create with Supabase and custom colors
  flutterforge create my_app --supabase --font Inter --background-color #121212 --text-color #FFFFFF
  
  # Full setup with onboarding and FCM
  flutterforge create my_app --firebase --onboarding --fcm --font Poppins`,
	Args: cobra.ExactArgs(1),
	Run:  runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&fontName, "font", "f", "Roboto", "Google Font family to use")
	createCmd.Flags().StringVar(&backend, "backend", "none", "Backend: firebase, supabase, or none")
	createCmd.Flags().StringVar(&backgroundColor, "background-color", "#121212", "App background color (hex)")
	createCmd.Flags().StringVar(&modalColor, "modal-color", "#1E1E1E", "Modal/dialog background color")
	createCmd.Flags().StringVar(&elevatedColor, "elevated-button-color", "#6750A4", "Elevated button color")
	createCmd.Flags().StringVar(&outlinedColor, "outlined-button-color", "#938F99", "Outlined button color")
	createCmd.Flags().StringVar(&textColor, "text-color", "#FFFFFF", "Primary text color")
	createCmd.Flags().StringVar(&primaryColor, "primary-color", "#6750A4", "Primary color")
	createCmd.Flags().StringVar(&secondaryColor, "secondary-color", "#CCC2DC", "Secondary color")
	createCmd.Flags().StringVar(&tertiaryColor, "tertiary-color", "#EFB8C8", "Tertiary color")
	createCmd.Flags().StringVar(&orgIdentifier, "org", "com.example", "Organization identifier")
	createCmd.Flags().BoolVar(&withOnboarding, "onboarding", true, "Include onboarding flow")
	createCmd.Flags().BoolVar(&withFcm, "fcm", false, "Include Firebase Cloud Messaging")
}

func runCreate(cmd *cobra.Command, args []string) {
	projectName = args[0]

	scaffold.ProjectName = projectName
	scaffold.FontName = fontName
	scaffold.Backend = backend
	scaffold.BackgroundColor = backgroundColor
	scaffold.ModalColor = modalColor
	scaffold.ElevatedColor = elevatedColor
	scaffold.OutlinedColor = outlinedColor
	scaffold.TextColor = textColor
	scaffold.PrimaryColor = primaryColor
	scaffold.SecondaryColor = secondaryColor
	scaffold.TertiaryColor = tertiaryColor
	scaffold.OrgIdentifier = orgIdentifier
	scaffold.WithOnboarding = withOnboarding
	scaffold.WithFcm = withFcm

	fmt.Printf("Creating Flutter project: %s\n", projectName)

	if err := scaffold.CreateFlutterProject(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Flutter project: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.CreateDirectoryStructure(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory structure: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.UpdatePubspec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating pubspec.yaml: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.CreateCoreFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating core files: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.CreateAuthFeature(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating auth feature: %v\n", err)
		os.Exit(1)
	}

	if withOnboarding {
		if err := scaffold.CreateOnboardingFeature(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating onboarding feature: %v\n", err)
			os.Exit(1)
		}
	}

	if err := scaffold.CreateHomeFeature(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating home feature: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.CreateSplashFeature(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating splash feature: %v\n", err)
		os.Exit(1)
	}

	if backend == "firebase" {
		if err := scaffold.SetupFirebase(); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting up Firebase: %v\n", err)
			os.Exit(1)
		}
	} else if backend == "supabase" {
		if err := scaffold.SetupSupabase(); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting up Supabase: %v\n", err)
			os.Exit(1)
		}
	}

	if withFcm && backend == "firebase" {
		if err := scaffold.SetupFcm(); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting up FCM: %v\n", err)
			os.Exit(1)
		}
	}

	if backend == "none" {
		if err := scaffold.UpdateMainDart(); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating main.dart: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\n✅ Project %s created successfully!\n", projectName)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  flutter pub get")
	fmt.Println("  flutter run")
}

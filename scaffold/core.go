package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ProjectName     string
	FontName        string
	Backend         string
	BackgroundColor string
	ModalColor      string
	ElevatedColor   string
	OutlinedColor   string
	TextColor       string
	PrimaryColor    string
	SecondaryColor  string
	TertiaryColor   string
	OrgIdentifier   string
	WithOnboarding  bool
	WithFcm         bool
)

func CreateFlutterProject() error {
	cmd := exec.Command("flutter", "create", ProjectName, "--org", OrgIdentifier)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("flutter create failed: %w", err)
	}
	return nil
}

func CreateDirectoryStructure() error {
	base := filepath.Join(ProjectName, "lib", "src")

	dirs := []string{
		// Core
		filepath.Join(base, "core", "constants"),
		filepath.Join(base, "core", "errors"),
		filepath.Join(base, "core", "router"),
		filepath.Join(base, "core", "theme"),
		filepath.Join(base, "core", "utils"),
		filepath.Join(base, "core", "widgets"),

		// Features
		filepath.Join(base, "features", "auth", "data", "datasources"),
		filepath.Join(base, "features", "auth", "data", "models"),
		filepath.Join(base, "features", "auth", "data", "repositories"),
		filepath.Join(base, "features", "auth", "domain", "entities"),
		filepath.Join(base, "features", "auth", "domain", "repositories"),
		filepath.Join(base, "features", "auth", "domain", "usecases"),
		filepath.Join(base, "features", "auth", "presentation", "screens"),
		filepath.Join(base, "features", "auth", "presentation", "widgets"),
		filepath.Join(base, "features", "auth", "presentation", "providers"),

		filepath.Join(base, "features", "home", "data", "datasources"),
		filepath.Join(base, "features", "home", "data", "models"),
		filepath.Join(base, "features", "home", "data", "repositories"),
		filepath.Join(base, "features", "home", "domain", "entities"),
		filepath.Join(base, "features", "home", "domain", "repositories"),
		filepath.Join(base, "features", "home", "domain", "usecases"),
		filepath.Join(base, "features", "home", "presentation", "screens"),
		filepath.Join(base, "features", "home", "presentation", "widgets"),
		filepath.Join(base, "features", "home", "presentation", "providers"),

		filepath.Join(base, "features", "onboarding", "data", "datasources"),
		filepath.Join(base, "features", "onboarding", "data", "models"),
		filepath.Join(base, "features", "onboarding", "data", "repositories"),
		filepath.Join(base, "features", "onboarding", "domain", "entities"),
		filepath.Join(base, "features", "onboarding", "domain", "repositories"),
		filepath.Join(base, "features", "onboarding", "domain", "usecases"),
		filepath.Join(base, "features", "onboarding", "presentation", "screens"),
		filepath.Join(base, "features", "onboarding", "presentation", "widgets"),
		filepath.Join(base, "features", "onboarding", "presentation", "providers"),

		filepath.Join(base, "features", "splash", "presentation", "screens"),
		filepath.Join(base, "features", "splash", "presentation", "widgets"),
		filepath.Join(base, "features", "splash", "presentation", "providers"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func UpdatePubspec() error {
	pubspecPath := filepath.Join(ProjectName, "pubspec.yaml")

	content := fmt.Sprintf(`name: %s
description: A new Flutter project.
publish_to: 'none'
version: 1.0.0+1

environment:
  sdk: ^3.7.0

dependencies:
  flutter:
    sdk: flutter

  # State Management
  flutter_riverpod: ^2.6.1
  riverpod_annotation: ^2.6.1

  # Routing
  go_router: ^14.8.1

  # Fonts
  google_fonts: ^6.2.1

  # UI Utilities
  flutter_animate: ^4.5.2
  gap: ^3.0.1
  flutter_svg: ^2.0.17

  # Utils
  equatable: ^2.0.7
  freezed_annotation: ^2.4.4
  json_annotation: ^4.9.0

`, ProjectName)

	if Backend == "firebase" {
		content += `  # Firebase
  firebase_core: ^3.13.0
  firebase_auth: ^5.5.0
  firebase_messaging: ^15.2.4
  firebase_analytics: ^11.6.0
`
	} else if Backend == "supabase" {
		content += `  # Supabase
  supabase_flutter: ^2.8.1
`
	}

	if WithOnboarding {
		content += `
  # Onboarding / Introduction
  smooth_page_indicator: ^1.2.0+3
`
	}

	content += `
dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^5.0.0
  build_runner: ^2.4.15
  freezed: ^2.5.8
  json_serializable: ^6.9.5
  riverpod_generator: ^2.6.5
`

	if err := os.WriteFile(pubspecPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write pubspec.yaml: %w", err)
	}

	return nil
}

func CreateCoreFiles() error {
	base := filepath.Join(ProjectName, "lib", "src", "core")

	// App Colors
	colorsContent := fmt.Sprintf(`import 'package:flutter/material.dart';

class AppColors {
  AppColors._();

  static Color get background => const Color(0xFF%s);
  static Color get modal => const Color(0xFF%s);
  static Color get elevatedButton => const Color(0xFF%s);
  static Color get outlinedButton => const Color(0xFF%s);
  static Color get text => const Color(0xFF%s);
  static Color get primary => const Color(0xFF%s);
  static Color get secondary => const Color(0xFF%s);
  static Color get tertiary => const Color(0xFF%s);
  
  // Additional semantic colors
  static Color get error => const Color(0xFFCF6679);
  static Color get success => const Color(0xFF4CAF50);
  static Color get warning => const Color(0xFFFFC107);
  static Color get info => const Color(0xFF2196F3);
  
  // Surface variants
  static Color get surface => const Color(0xFF1E1E1E);
  static Color get surfaceVariant => const Color(0xFF2D2D2D);
  static Color get outline => const Color(0xFF938F99);
  static Color get outlineVariant => const Color(0xFF49454F);
  
  // On colors
  static Color get onPrimary => const Color(0xFFFFFFFF);
  static Color get onSecondary => const Color(0xFF000000);
  static Color get onSurface => const Color(0xFFE6E1E5);
  static Color get onBackground => const Color(0xFFE6E1E5);
}
`,
		strings.TrimPrefix(BackgroundColor, "#"),
		strings.TrimPrefix(ModalColor, "#"),
		strings.TrimPrefix(ElevatedColor, "#"),
		strings.TrimPrefix(OutlinedColor, "#"),
		strings.TrimPrefix(TextColor, "#"),
		strings.TrimPrefix(PrimaryColor, "#"),
		strings.TrimPrefix(SecondaryColor, "#"),
		strings.TrimPrefix(TertiaryColor, "#"),
	)

	if err := os.WriteFile(filepath.Join(base, "constants", "app_colors.dart"), []byte(colorsContent), 0644); err != nil {
		return err
	}

	// App Strings
	stringsContent := `class AppStrings {
  AppStrings._();

  static const String appName = 'My App';
  
  // Auth
  static const String email = 'Email';
  static const String password = 'Password';
  static const String confirmPassword = 'Confirm Password';
  static const String login = 'Login';
  static const String signup = 'Sign Up';
  static const String logout = 'Logout';
  static const String forgotPassword = 'Forgot Password?';
  static const String noAccount = "Don't have an account?";
  static const String haveAccount = 'Already have an account?';
  static const String signUpHere = 'Sign up here';
  static const String loginHere = 'Login here';
  
  // Onboarding
  static const String skip = 'Skip';
  static const String next = 'Next';
  static const String getStarted = 'Get Started';
  static const String previous = 'Previous';
  
  // Errors
  static const String genericError = 'Something went wrong. Please try again.';
  static const String networkError = 'Please check your internet connection.';
  static const String invalidEmail = 'Please enter a valid email address.';
  static const String weakPassword = 'Password must be at least 8 characters.';
  static const String passwordMismatch = 'Passwords do not match.';
  
  // Success
  static const String loginSuccess = 'Login successful!';
  static const String signupSuccess = 'Account created successfully!';
  static const String logoutSuccess = 'Logged out successfully!';
}
`
	if err := os.WriteFile(filepath.Join(base, "constants", "app_strings.dart"), []byte(stringsContent), 0644); err != nil {
		return err
	}

	// Theme
	themeContent := fmt.Sprintf(`import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../constants/app_colors.dart';

class AppTheme {
  AppTheme._();

  static ThemeData get darkTheme {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      colorScheme: ColorScheme.dark(
        primary: AppColors.primary,
        onPrimary: AppColors.onPrimary,
        secondary: AppColors.secondary,
        onSecondary: AppColors.onSecondary,
        tertiary: AppColors.tertiary,
        surface: AppColors.surface,
        onSurface: AppColors.onSurface,
        error: AppColors.error,
      ),
      scaffoldBackgroundColor: AppColors.background,
      textTheme: _buildTextTheme(),
      appBarTheme: AppBarTheme(
        backgroundColor: AppColors.background,
        elevation: 0,
        centerTitle: true,
        titleTextStyle: GoogleFonts.%s(
          fontSize: 20,
          fontWeight: FontWeight.w600,
          color: AppColors.text,
        ),
        iconTheme: const IconThemeData(color: AppColors.text),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.elevatedButton,
          foregroundColor: AppColors.onPrimary,
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          textStyle: GoogleFonts.%s(
            fontSize: 16,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: AppColors.outlinedButton,
          side: const BorderSide(color: AppColors.outlinedButton),
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          textStyle: GoogleFonts.%s(
            fontSize: 16,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: AppColors.primary,
          textStyle: GoogleFonts.%s(
            fontSize: 14,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: AppColors.surfaceVariant,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide.none,
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide.none,
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: AppColors.primary, width: 2),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: AppColors.error),
        ),
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
        hintStyle: GoogleFonts.%s(
          color: AppColors.outline,
          fontSize: 16,
        ),
        labelStyle: GoogleFonts.%s(
          color: AppColors.onSurface,
          fontSize: 16,
        ),
      ),
      cardTheme: CardTheme(
        color: AppColors.surface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
      ),
      dialogTheme: DialogTheme(
        backgroundColor: AppColors.modal,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
      ),
      bottomSheetTheme: const BottomSheetThemeData(
        backgroundColor: AppColors.modal,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
        ),
      ),
      snackBarTheme: SnackBarThemeData(
        backgroundColor: AppColors.surface,
        contentTextStyle: GoogleFonts.%s(
          color: AppColors.onSurface,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        behavior: SnackBarBehavior.floating,
      ),
      dividerTheme: const DividerThemeData(
        color: AppColors.outlineVariant,
        thickness: 1,
      ),
      progressIndicatorTheme: const ProgressIndicatorThemeData(
        color: AppColors.primary,
      ),
    );
  }

  static TextTheme _buildTextTheme() {
    final fontFamily = GoogleFonts.%s();

    return TextTheme(
      displayLarge: fontFamily.copyWith(
        fontSize: 57,
        fontWeight: FontWeight.w400,
        letterSpacing: -0.25,
        color: AppColors.text,
      ),
      displayMedium: fontFamily.copyWith(
        fontSize: 45,
        fontWeight: FontWeight.w400,
        color: AppColors.text,
      ),
      displaySmall: fontFamily.copyWith(
        fontSize: 36,
        fontWeight: FontWeight.w400,
        color: AppColors.text,
      ),
      headlineLarge: fontFamily.copyWith(
        fontSize: 32,
        fontWeight: FontWeight.w600,
        color: AppColors.text,
      ),
      headlineMedium: fontFamily.copyWith(
        fontSize: 28,
        fontWeight: FontWeight.w600,
        color: AppColors.text,
      ),
      headlineSmall: fontFamily.copyWith(
        fontSize: 24,
        fontWeight: FontWeight.w600,
        color: AppColors.text,
      ),
      titleLarge: fontFamily.copyWith(
        fontSize: 22,
        fontWeight: FontWeight.w600,
        color: AppColors.text,
      ),
      titleMedium: fontFamily.copyWith(
        fontSize: 16,
        fontWeight: FontWeight.w600,
        letterSpacing: 0.15,
        color: AppColors.text,
      ),
      titleSmall: fontFamily.copyWith(
        fontSize: 14,
        fontWeight: FontWeight.w600,
        letterSpacing: 0.1,
        color: AppColors.text,
      ),
      bodyLarge: fontFamily.copyWith(
        fontSize: 16,
        fontWeight: FontWeight.w400,
        letterSpacing: 0.5,
        color: AppColors.text,
      ),
      bodyMedium: fontFamily.copyWith(
        fontSize: 14,
        fontWeight: FontWeight.w400,
        letterSpacing: 0.25,
        color: AppColors.text,
      ),
      bodySmall: fontFamily.copyWith(
        fontSize: 12,
        fontWeight: FontWeight.w400,
        letterSpacing: 0.4,
        color: AppColors.outline,
      ),
      labelLarge: fontFamily.copyWith(
        fontSize: 14,
        fontWeight: FontWeight.w600,
        letterSpacing: 0.1,
        color: AppColors.text,
      ),
      labelMedium: fontFamily.copyWith(
        fontSize: 12,
        fontWeight: FontWeight.w600,
        letterSpacing: 0.5,
        color: AppColors.text,
      ),
      labelSmall: fontFamily.copyWith(
        fontSize: 11,
        fontWeight: FontWeight.w600,
        letterSpacing: 0.5,
        color: AppColors.text,
      ),
    );
  }
}
`, FontName, FontName, FontName, FontName, FontName, FontName, FontName, FontName, FontName)

	if err := os.WriteFile(filepath.Join(base, "theme", "app_theme.dart"), []byte(themeContent), 0644); err != nil {
		return err
	}

	// Create router files
	if err := createRouterFiles(); err != nil {
		return err
	}

	// Create common widgets
	if err := createCommonWidgets(); err != nil {
		return err
	}

	return nil
}

func createRouterFiles() error {
	base := filepath.Join(ProjectName, "lib", "src", "core", "router")

	routesContent := `abstract class AppRoutes {
  AppRoutes._();

  static const String splash = '/';
  static const String onboarding = '/onboarding';
  static const String login = '/login';
  static const String signup = '/signup';
  static const String forgotPassword = '/forgot-password';
  static const String home = '/home';
}
`

	routerContent := fmt.Sprintf(`import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/signup_screen.dart';
import '../../features/auth/presentation/screens/forgot_password_screen.dart';
import '../../features/auth/presentation/providers/auth_provider.dart';
import '../../features/onboarding/presentation/screens/onboarding_screen.dart';
import '../../features/splash/presentation/screens/splash_screen.dart';
import '../../features/home/presentation/screens/home_screen.dart';
import 'app_routes.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);

  return GoRouter(
    initialLocation: AppRoutes.splash,
    debugLogDiagnostics: true,
    redirect: (context, state) {
      final isLoggedIn = authState.valueOrNull != null;
      final isOnAuthRoute = state.matchedLocation == AppRoutes.login ||
          state.matchedLocation == AppRoutes.signup ||
          state.matchedLocation == AppRoutes.forgotPassword;
      final isOnboarding = state.matchedLocation == AppRoutes.onboarding;

      // If not logged in and not on auth route, redirect to login
      if (!isLoggedIn && !isOnAuthRoute && !isOnboarding) {
        return AppRoutes.login;
      }

      // If logged in and on auth route, redirect to home
      if (isLoggedIn && isOnAuthRoute) {
        return AppRoutes.home;
      }

      return null;
    },
    routes: [
      GoRoute(
        path: AppRoutes.splash,
        name: 'splash',
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(
        path: AppRoutes.onboarding,
        name: 'onboarding',
        builder: (context, state) => const OnboardingScreen(),
      ),
      GoRoute(
        path: AppRoutes.login,
        name: 'login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: AppRoutes.signup,
        name: 'signup',
        builder: (context, state) => const SignupScreen(),
      ),
      GoRoute(
        path: AppRoutes.forgotPassword,
        name: 'forgotPassword',
        builder: (context, state) => const ForgotPasswordScreen(),
      ),
      GoRoute(
        path: AppRoutes.home,
        name: 'home',
        builder: (context, state) => const HomeScreen(),
      ),
    ],
    errorBuilder: (context, state) => Scaffold(
      body: Center(
        child: Text('Page not found: ${state.error}'),
      ),
    ),
  );
});
`)

	if err := os.WriteFile(filepath.Join(base, "app_routes.dart"), []byte(routesContent), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(base, "app_router.dart"), []byte(routerContent), 0644); err != nil {
		return err
	}

	return nil
}

func createCommonWidgets() error {
	base := filepath.Join(ProjectName, "lib", "src", "core", "widgets")

	loadingWidget := `import 'package:flutter/material.dart';
import '../constants/app_colors.dart';

class LoadingWidget extends StatelessWidget {
  final String? message;
  
  const LoadingWidget({super.key, this.message});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const CircularProgressIndicator(
            color: AppColors.primary,
          ),
          if (message != null) ...[
            const SizedBox(height: 16),
            Text(
              message!,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
          ],
        ],
      ),
    );
  }
}
`

	errorWidget := `import 'package:flutter/material.dart';
import '../constants/app_colors.dart';

class ErrorWidget extends StatelessWidget {
  final String message;
  final VoidCallback? onRetry;

  const ErrorWidget({
    super.key,
    required this.message,
    this.onRetry,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.error_outline,
              size: 64,
              color: AppColors.error,
            ),
            const SizedBox(height: 16),
            Text(
              message,
              style: Theme.of(context).textTheme.bodyLarge,
              textAlign: TextAlign.center,
            ),
            if (onRetry != null) ...[
              const SizedBox(height: 24),
              ElevatedButton.icon(
                onPressed: onRetry,
                icon: const Icon(Icons.refresh),
                label: const Text('Retry'),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
`

	emptyStateWidget := `import 'package:flutter/material.dart';
import 'package:gap/gap.dart';
import '../constants/app_colors.dart';

class EmptyStateWidget extends StatelessWidget {
  final IconData icon;
  final String title;
  final String? subtitle;
  final String? actionLabel;
  final VoidCallback? onAction;

  const EmptyStateWidget({
    super.key,
    required this.icon,
    required this.title,
    this.subtitle,
    this.actionLabel,
    this.onAction,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              icon,
              size: 80,
              color: AppColors.outline,
            ),
            const Gap(16),
            Text(
              title,
              style: Theme.of(context).textTheme.headlineSmall,
              textAlign: TextAlign.center,
            ),
            if (subtitle != null) ...[
              const Gap(8),
              Text(
                subtitle!,
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: AppColors.outline,
                ),
                textAlign: TextAlign.center,
              ),
            ],
            if (actionLabel != null && onAction != null) ...[
              const Gap(24),
              ElevatedButton(
                onPressed: onAction,
                child: Text(actionLabel!),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
`

	widgets := map[string]string{
		"loading_widget.dart":     loadingWidget,
		"error_widget.dart":       errorWidget,
		"empty_state_widget.dart": emptyStateWidget,
	}

	for filename, content := range widgets {
		if err := os.WriteFile(filepath.Join(base, filename), []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

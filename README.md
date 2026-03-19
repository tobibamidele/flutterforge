# FlutterForge

FlutterForge is a CLI tool that scaffolds Flutter projects with a feature-first architecture, saving developers hours of repetitive setup work.

## Features

- Feature-first directory structure following Clean Architecture principles
- go_router integration with typed routes and navigation
- Riverpod state management with providers and notifiers pre-configured
- Firebase and Supabase authentication scaffolding
- Firebase Cloud Messaging (FCM) setup for push notifications
- Material Design 3 theming with customizable colors and typography
- Google Fonts integration
- Complete authentication flow (login, signup, password reset)
- Onboarding flow with page indicators
- Splash screen with auth state detection
- Common widgets (loading, error, empty states)

## Installation

### From Source

```bash
git clone https://github.com/tobibamidele/flutterforge.git
cd flutterforge
go build -o flutterforge .
sudo mv flutterforge /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/tobibamidele/flutterforge@latest
```

## Usage

### Basic Usage

```bash
flutterforge create my_app
```

This creates a Flutter project with:
- Feature-first directory structure
- go_router and Riverpod pre-configured
- Dark theme with Material 3
- Google Fonts (Roboto)
- Authentication screens
- Onboarding flow
- Splash screen

### With Firebase

```bash
flutterforge create my_app --backend firebase
```

### With Supabase

```bash
flutterforge create my_app --backend supabase
```

### With FCM (Push Notifications)

```bash
flutterforge create my_app --backend firebase --fcm
```

### With Custom Colors

```bash
flutterforge create my_app \
  --background-color #0D0D0D \
  --primary-color #6C63FF \
  --text-color #F5F5F5
```

### With Custom Font

```bash
flutterforge create my_app --font Inter
```

## Command Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--backend` | - | Backend integration: `firebase`, `supabase`, or `none` | `none` |
| `--font` | `-f` | Google Font family to use | `Roboto` |
| `--background-color` | - | App background color (hex) | `#121212` |
| `--primary-color` | - | Primary theme color (hex) | `#6750A4` |
| `--secondary-color` | - | Secondary theme color (hex) | `#CCC2DC` |
| `--tertiary-color` | - | Tertiary theme color (hex) | `#EFB8C8` |
| `--text-color` | - | Primary text color (hex) | `#FFFFFF` |
| `--modal-color` | - | Modal/dialog background color (hex) | `#1E1E1E` |
| `--elevated-button-color` | - | Elevated button color (hex) | `#6750A4` |
| `--outlined-button-color` | - | Outlined button color (hex) | `#938F99` |
| `--org` | - | Organization identifier for Flutter | `com.example` |
| `--onboarding` | - | Include onboarding flow | `true` |
| `--fcm` | - | Include Firebase Cloud Messaging | `false` |

## Project Structure

After scaffolding, your project will have the following structure:

```
lib/src/
├── core/
│   ├── constants/
│   │   ├── app_colors.dart
│   │   └── app_strings.dart
│   ├── errors/
│   ├── router/
│   │   ├── app_routes.dart
│   │   └── app_router.dart
│   ├── services/
│   ├── theme/
│   │   └── app_theme.dart
│   ├── utils/
│   └── widgets/
│       ├── loading_widget.dart
│       ├── error_widget.dart
│       └── empty_state_widget.dart
└── features/
    ├── auth/
    │   ├── data/
    │   │   ├── datasources/
    │   │   ├── models/
    │   │   └── repositories/
    │   ├── domain/
    │   │   ├── entities/
    │   │   ├── repositories/
    │   │   └── usecases/
    │   ├── presentation/
    │   │   ├── providers/
    │   │   ├── screens/
    │   │   └── widgets/
    │   └── auth_router.dart
    ├── home/
    │   ├── data/
    │   ├── domain/
    │   └── presentation/
    ├── onboarding/
    │   ├── data/
    │   ├── domain/
    │   └── presentation/
    └── splash/
        └── presentation/
```

## Generated Files

### Core Files

- **app_colors.dart**: Color constants for the theme
- **app_strings.dart**: String constants for localization-ready text
- **app_theme.dart**: Material 3 theme with all TextStyles (displayLarge, headlineMedium, bodyLarge, etc.)
- **app_routes.dart**: Route path constants
- **app_router.dart**: go_router configuration with auth-aware redirects

### Authentication Feature

- **User entity**: Domain model with id, email, displayName, photoUrl, emailVerified
- **AuthRepository**: Interface for auth operations
- **AuthRepositoryImpl**: Implementation for Firebase, Supabase, or mock
- **AuthNotifier**: StateNotifier for auth state management
- **Login screen**: Email/password form with validation
- **Signup screen**: Registration form with confirm password
- **Forgot password screen**: Password reset flow

### State Management

Riverpod providers are generated for:

- Auth state stream
- Sign in/sign up futures
- Sign out
- Password reset
- Home navigation state
- Splash initialization state

## Post-Scaffolding Steps

1. Navigate to the project directory:
   ```bash
   cd my_app
   ```

2. Fetch dependencies:
   ```bash
   flutter pub get
   ```

3. For Firebase projects:
   - Add your `google-services.json` to `android/app/`
   - Add your `GoogleService-Info.plist` to `ios/Runner/`

4. For Supabase projects:
   - Update the URL and anon key in `lib/src/core/services/supabase_service.dart`

5. Run the app:
   ```bash
   flutter run
   ```

## Requirements

- Flutter SDK 3.7.0 or later
- Go 1.21 or later (for building from source)

## License

Apache-2.0 License

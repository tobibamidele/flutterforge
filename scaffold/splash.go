package scaffold

import (
	"os"
	"path/filepath"
)

func CreateSplashFeature() error {
	base := filepath.Join(ProjectName, "lib", "src", "features", "splash")

	providerContent := `import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

enum SplashStatus {
  initial,
  checkingAuth,
  authenticated,
  unauthenticated,
  onboardingRequired,
}

class SplashState {
  final SplashStatus status;
  final String? message;

  const SplashState({
    this.status = SplashStatus.initial,
    this.message,
  });

  SplashState copyWith({
    SplashStatus? status,
    String? message,
  }) {
    return SplashState(
      status: status ?? this.status,
      message: message ?? this.message,
    );
  }
}

class SplashNotifier extends StateNotifier<SplashState> {
  SplashNotifier() : super(const SplashState());

  Future<void> checkAuthStatus() async {
    state = state.copyWith(status: SplashStatus.checkingAuth);

    try {
      // Simulate splash delay
      await Future.delayed(const Duration(seconds: 2));
      
      // Check onboarding status
      final prefs = await SharedPreferences.getInstance();
      final hasSeenOnboarding = prefs.getBool('has_seen_onboarding') ?? false;
      
      // TODO: Check actual auth status from AuthRepository
      // For now, we'll simulate based on onboarding status
      if (!hasSeenOnboarding) {
        state = state.copyWith(status: SplashStatus.onboardingRequired);
      } else {
        // Check if user is logged in
        // final user = await ref.read(authRepositoryProvider).getCurrentUser();
        // if (user != null) {
        //   state = state.copyWith(status: SplashStatus.authenticated);
        // } else {
        state = state.copyWith(status: SplashStatus.unauthenticated);
        // }
      }
    } catch (e) {
      state = state.copyWith(
        status: SplashStatus.unauthenticated,
        message: e.toString(),
      );
    }
  }

  Future<void> completeOnboarding() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool('has_seen_onboarding', true);
  }
}

final splashNotifierProvider = StateNotifierProvider<SplashNotifier, SplashState>((ref) {
  return SplashNotifier();
});
`
	if err := os.WriteFile(filepath.Join(base, "presentation", "providers", "splash_provider.dart"), []byte(providerContent), 0644); err != nil {
		return err
	}

	screenContent := `import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/router/app_routes.dart';
import '../providers/splash_provider.dart';

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(splashNotifierProvider.notifier).checkAuthStatus();
    });
  }

  @override
  Widget build(BuildContext context) {
    final splashState = ref.watch(splashNotifierProvider);

    ref.listen<SplashState>(splashNotifierProvider, (previous, next) {
      switch (next.status) {
        case SplashStatus.onboardingRequired:
          context.go(AppRoutes.onboarding);
          break;
        case SplashStatus.authenticated:
          context.go(AppRoutes.home);
          break;
        case SplashStatus.unauthenticated:
          context.go(AppRoutes.login);
          break;
        default:
          break;
      }
    });

    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [
              AppColors.primary.withValues(alpha: 0.3),
              AppColors.background,
            ],
          ),
        ),
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: 120,
                height: 120,
                decoration: BoxDecoration(
                  color: AppColors.primary,
                  borderRadius: BorderRadius.circular(30),
                ),
                child: const Icon(
                  Icons.flutter_dash,
                  size: 80,
                  color: Colors.white,
                ),
              ),
              const SizedBox(height: 24),
              Text(
                'App Name',
                style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'Loading...',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: AppColors.outline,
                ),
              ),
              const SizedBox(height: 48),
              const CircularProgressIndicator(),
            ],
          ),
        ),
      ),
    );
  }
}
`
	if err := os.WriteFile(filepath.Join(base, "presentation", "screens", "splash_screen.dart"), []byte(screenContent), 0644); err != nil {
		return err
	}

	return nil
}

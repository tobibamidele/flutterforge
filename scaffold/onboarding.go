package scaffold

import (
	"os"
	"path/filepath"
)

func CreateOnboardingFeature() error {
	base := filepath.Join(ProjectName, "lib", "src", "features", "onboarding")

	entityContent := `import 'package:equatable/equatable.dart';

class OnboardingPage extends Equatable {
  final String title;
  final String description;
  final String imagePath;

  const OnboardingPage({
    required this.title,
    required this.description,
    required this.imagePath,
  });

  @override
  List<Object?> get props => [title, description, imagePath];
}
`
	if err := os.WriteFile(filepath.Join(base, "domain", "entities", "onboarding_page.dart"), []byte(entityContent), 0644); err != nil {
		return err
	}

	providerContent := `import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../domain/entities/onboarding_page.dart';

final onboardingPagesProvider = Provider<List<OnboardingPage>>((ref) {
  return const [
    OnboardingPage(
      title: 'Welcome to the App',
      description: 'Discover amazing features that will help you achieve your goals faster.',
      imagePath: 'assets/images/onboarding_1.svg',
    ),
    OnboardingPage(
      title: 'Stay Organized',
      description: 'Keep track of everything in one place. Never miss an important task again.',
      imagePath: 'assets/images/onboarding_2.svg',
    ),
    OnboardingPage(
      title: 'Get Started',
      description: 'Ready to begin? Create your account and start exploring today!',
      imagePath: 'assets/images/onboarding_3.svg',
    ),
  ];
});

class OnboardingState {
  final int currentPage;
  final bool isCompleted;

  const OnboardingState({
    this.currentPage = 0,
    this.isCompleted = false,
  });

  OnboardingState copyWith({
    int? currentPage,
    bool? isCompleted,
  }) {
    return OnboardingState(
      currentPage: currentPage ?? this.currentPage,
      isCompleted: isCompleted ?? this.isCompleted,
    );
  }
}

class OnboardingNotifier extends StateNotifier<OnboardingState> {
  OnboardingNotifier() : super(const OnboardingState());

  void nextPage(int totalPages) {
    if (state.currentPage < totalPages - 1) {
      state = state.copyWith(currentPage: state.currentPage + 1);
    } else {
      state = state.copyWith(isCompleted: true);
    }
  }

  void previousPage() {
    if (state.currentPage > 0) {
      state = state.copyWith(currentPage: state.currentPage - 1);
    }
  }

  void goToPage(int page) {
    state = state.copyWith(currentPage: page);
  }

  void complete() {
    state = state.copyWith(isCompleted: true);
  }

  void reset() {
    state = const OnboardingState();
  }
}

final onboardingNotifierProvider = StateNotifierProvider<OnboardingNotifier, OnboardingState>((ref) {
  return OnboardingNotifier();
});
`
	if err := os.WriteFile(filepath.Join(base, "presentation", "providers", "onboarding_provider.dart"), []byte(providerContent), 0644); err != nil {
		return err
	}

	screenContent := `import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:smooth_page_indicator/smooth_page_indicator.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/router/app_routes.dart';
import '../providers/onboarding_provider.dart';
import '../widgets/onboarding_page_widget.dart';

class OnboardingScreen extends ConsumerStatefulWidget {
  const OnboardingScreen({super.key});

  @override
  ConsumerState<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends ConsumerState<OnboardingScreen> {
  final PageController _pageController = PageController();

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final pages = ref.watch(onboardingPagesProvider);
    final onboardingState = ref.watch(onboardingNotifierProvider);

    ref.listen<OnboardingState>(onboardingNotifierProvider, (previous, next) {
      if (next.isCompleted) {
        context.go(AppRoutes.signup);
      }
    });

    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  TextButton(
                    onPressed: () => context.go(AppRoutes.signup),
                    child: const Text(AppStrings.skip),
                  ),
                ],
              ),
            ),
            Expanded(
              child: PageView.builder(
                controller: _pageController,
                itemCount: pages.length,
                onPageChanged: (index) {
                  ref.read(onboardingNotifierProvider.notifier).goToPage(index);
                },
                itemBuilder: (context, index) {
                  return OnboardingPageWidget(page: pages[index]);
                },
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(24),
              child: Column(
                children: [
                  SmoothPageIndicator(
                    controller: _pageController,
                    count: pages.length,
                    effect: WormEffect(
                      dotHeight: 8,
                      dotWidth: 8,
                      activeDotColor: AppColors.primary,
                      dotColor: AppColors.surfaceVariant,
                    ),
                  ),
                  const SizedBox(height: 32),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      if (onboardingState.currentPage > 0)
                        OutlinedButton(
                          onPressed: () {
                            _pageController.previousPage(
                              duration: const Duration(milliseconds: 300),
                              curve: Curves.easeInOut,
                            );
                          },
                          child: const Text(AppStrings.previous),
                        )
                      else
                        const SizedBox(width: 100),
                      ElevatedButton(
                        onPressed: () {
                          if (onboardingState.currentPage == pages.length - 1) {
                            ref.read(onboardingNotifierProvider.notifier).complete();
                          } else {
                            _pageController.nextPage(
                              duration: const Duration(milliseconds: 300),
                              curve: Curves.easeInOut,
                            );
                          }
                        },
                        child: Text(
                          onboardingState.currentPage == pages.length - 1
                              ? AppStrings.getStarted
                              : AppStrings.next,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
`
	if err := os.WriteFile(filepath.Join(base, "presentation", "screens", "onboarding_screen.dart"), []byte(screenContent), 0644); err != nil {
		return err
	}

	pageWidget := `import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:gap/gap.dart';
import '../../../../core/constants/app_colors.dart';
import '../../domain/entities/onboarding_page.dart';

class OnboardingPageWidget extends StatelessWidget {
  final OnboardingPage page;

  const OnboardingPageWidget({
    super.key,
    required this.page,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(32),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            height: 250,
            width: 250,
            decoration: BoxDecoration(
              color: AppColors.surfaceVariant,
              borderRadius: BorderRadius.circular(20),
            ),
            child: Center(
              child: Icon(
                Icons.image_outlined,
                size: 100,
                color: AppColors.primary,
              ),
            ),
          ),
          const Gap(48),
          Text(
            page.title,
            style: Theme.of(context).textTheme.headlineMedium,
            textAlign: TextAlign.center,
          ),
          const Gap(16),
          Text(
            page.description,
            style: Theme.of(context).textTheme.bodyLarge?.copyWith(
              color: AppColors.outline,
            ),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}
`
	if err := os.WriteFile(filepath.Join(base, "presentation", "widgets", "onboarding_page_widget.dart"), []byte(pageWidget), 0644); err != nil {
		return err
	}

	return nil
}

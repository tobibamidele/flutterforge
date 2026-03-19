package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

func SetupFirebase() error {
	base := filepath.Join(ProjectName, "lib", "src")

	// Create Firebase service
	firebaseService := `import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';

class FirebaseService {
  static FirebaseService? _instance;
  static FirebaseService get instance => _instance ??= FirebaseService._();

  FirebaseService._();

  FirebaseAuth get auth => FirebaseAuth.instance;
  FirebaseMessaging get messaging => FirebaseMessaging.instance;

  Future<void> initialize() async {
    await Firebase.initializeApp();
    await _setupMessaging();
    _setupAuthInterceptors();
  }

  Future<void> _setupMessaging() async {
    final settings = await messaging.requestPermission(
      alert: true,
      announcement: false,
      badge: true,
      carPlay: false,
      criticalAlert: false,
      provisional: false,
      sound: true,
    );

    if (settings.authorizationStatus == AuthorizationStatus.authorized) {
      String? token = await messaging.getToken();
      if (kDebugMode) {
        print('FCM Token: $token');
      }
      
      messaging.onTokenRefresh.listen((newToken) {
        if (kDebugMode) {
          print('New FCM Token: $newToken');
        }
        _sendTokenToServer(newToken);
      });
    }
  }

  void _setupAuthInterceptors() {
    auth.authStateChanges().listen((User? user) {
      if (user == null) {
        if (kDebugMode) {
          print('User is currently signed out!');
        }
      } else {
        if (kDebugMode) {
          print('User is signed in!');
          print('User ID: ${user.uid}');
          print('User Email: ${user.email}');
        }
      }
    });
  }

  Future<void> _sendTokenToServer(String token) async {
    // TODO: Implement sending token to your backend server
    if (kDebugMode) {
      print('Sending token to server: $token');
    }
  }
}

@pragma('vm:entry-point')
Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  await Firebase.initializeApp();
  if (kDebugMode) {
    print('Handling a background message: ${message.messageId}');
    print('Message data: ${message.data}');
    print('Message notification: ${message.notification?.title}');
    print('Message notification: ${message.notification?.body}');
  }
}
`
	if err := os.WriteFile(filepath.Join(base, "core", "services", "firebase_service.dart"), []byte(firebaseService), 0644); err != nil {
		return err
	}

	// Create main.dart with Firebase setup
	mainContent := fmt.Sprintf(`import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'src/core/theme/app_theme.dart';
import 'src/core/router/app_router.dart';
import 'src/core/services/firebase_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Set preferred orientations
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);
  
  // Initialize Firebase
  await Firebase.initializeApp();
  
  // Set up FCM background handler
  FirebaseMessaging.onBackgroundMessage(FirebaseMessaging.instance);
  
  // Initialize other services
  await FirebaseService.instance.initialize();
  
  runApp(
    const ProviderScope(
      child: MyApp(),
    ),
  );
}

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: '%s',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.darkTheme,
      routerConfig: router,
    );
  }
}
`, ProjectName)

	if err := os.WriteFile(filepath.Join(ProjectName, "lib", "main.dart"), []byte(mainContent), 0644); err != nil {
		return err
	}

	// Create google-services.json placeholder
	jsonPlaceholder := `{
  "placeholder": "Replace this with your google-services.json from Firebase Console"
}
`
	if err := os.WriteFile(filepath.Join(ProjectName, "android", "app", "google-services.json"), []byte(jsonPlaceholder), 0644); err != nil {
		fmt.Printf("Warning: Could not create google-services.json placeholder: %v\n", err)
	}

	// Update Android build.gradle
	androidBuildGradle := filepath.Join(ProjectName, "android", "app", "build.gradle.kts")
	if _, err := os.Stat(androidBuildGradle); err == nil {
		content, _ := os.ReadFile(androidBuildGradle)
		newContent := string(content)
		if !contains(newContent, "com.google.gms.google-services") {
			newContent = newContent + `
plugins.apply(com.google.gms.google-services)
`
			os.WriteFile(androidBuildGradle, []byte(newContent), 0644)
		}
	}

	return nil
}

func SetupSupabase() error {
	base := filepath.Join(ProjectName, "lib", "src")

	// Create Supabase service
	supabaseService := `import 'package:supabase_flutter/supabase_flutter.dart';
import 'package:flutter/foundation.dart';

class SupabaseService {
  static SupabaseService? _instance;
  static SupabaseService get instance => _instance ??= SupabaseService._();

  SupabaseService._();

  SupabaseClient get client => Supabase.instance.client;

  Future<void> initialize({
    required String url,
    required String anonKey,
  }) async {
    await Supabase.initialize(
      url: url,
      anonKey: anonKey,
    );
    
    if (kDebugMode) {
      print('Supabase initialized successfully');
    }
  }

  Future<AuthResponse> signInWithEmail(String email, String password) async {
    return await client.auth.signInWithPassword(
      email: email,
      password: password,
    );
  }

  Future<AuthResponse> signUpWithEmail(String email, String password) async {
    return await client.auth.signUp(
      email: email,
      password: password,
    );
  }

  Future<void> signOut() async {
    await client.auth.signOut();
  }

  User? get currentUser => client.auth.currentUser;
  
  Session? get currentSession => client.auth.currentSession;

  Stream<AuthState> get authStateChanges => client.auth.onAuthStateChange;
}
`
	if err := os.WriteFile(filepath.Join(base, "core", "services", "supabase_service.dart"), []byte(supabaseService), 0644); err != nil {
		return err
	}

	// Create main.dart with Supabase setup
	mainContent := fmt.Sprintf(`import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:supabase_flutter/supabase_flutter.dart';
import 'src/core/theme/app_theme.dart';
import 'src/core/router/app_router.dart';
import 'src/core/services/supabase_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Set preferred orientations
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);
  
  // Initialize Supabase
  // TODO: Replace with your actual Supabase URL and anon key
  await SupabaseService.instance.initialize(
    url: 'YOUR_SUPABASE_URL',
    anonKey: 'YOUR_SUPABASE_ANON_KEY',
  );
  
  runApp(
    const ProviderScope(
      child: MyApp(),
    ),
  );
}

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: '%s',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.darkTheme,
      routerConfig: router,
    );
  }
}
`, ProjectName)

	if err := os.WriteFile(filepath.Join(ProjectName, "lib", "main.dart"), []byte(mainContent), 0644); err != nil {
		return err
	}

	return nil
}

func SetupFcm() error {
	base := filepath.Join(ProjectName, "lib", "src")

	// Create notification service
	notificationService := `import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';

class NotificationService {
  static NotificationService? _instance;
  static NotificationService get instance => _instance ??= NotificationService._();

  NotificationService._();

  final FirebaseMessaging _messaging = FirebaseMessaging.instance;

  Function(RemoteMessage)? onMessageHandler;
  Function(RemoteMessage)? onMessageOpenedHandler;

  Future<void> initialize() async {
    await _requestPermission();
    await _getToken();
    _setupMessageHandlers();
  }

  Future<void> _requestPermission() async {
    final settings = await _messaging.requestPermission(
      alert: true,
      announcement: false,
      badge: true,
      carPlay: false,
      criticalAlert: false,
      provisional: false,
      sound: true,
    );

    if (kDebugMode) {
      print('Notification permission status: ${settings.authorizationStatus}');
    }
  }

  Future<void> _getToken() async {
    String? token = await _messaging.getToken();
    if (kDebugMode) {
      print('FCM Token: $token');
    }
    
    _messaging.onTokenRefresh.listen((newToken) {
      if (kDebugMode) {
        print('New FCM Token: $newToken');
      }
      _sendTokenToServer(newToken);
    });
  }

  void _setupMessageHandlers() {
    FirebaseMessaging.onMessage.listen((RemoteMessage message) {
      if (kDebugMode) {
        print('Received foreground message:');
        print('Title: ${message.notification?.title}');
        print('Body: ${message.notification?.body}');
        print('Data: ${message.data}');
      }
      onMessageHandler?.call(message);
    });

    FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) {
      if (kDebugMode) {
        print('Message opened from terminated state:');
        print('Title: ${message.notification?.title}');
        print('Body: ${message.notification?.body}');
        print('Data: ${message.data}');
      }
      onMessageOpenedHandler?.call(message);
    });
  }

  Future<void> _sendTokenToServer(String token) async {
    // TODO: Implement sending token to your backend server
    if (kDebugMode) {
      print('Sending token to server: $token');
    }
  }

  Future<RemoteMessage?> getInitialMessage() async {
    return await _messaging.getInitialMessage();
  }

  Future<void> subscribeToTopic(String topic) async {
    await _messaging.subscribeToTopic(topic);
    if (kDebugMode) {
      print('Subscribed to topic: $topic');
    }
  }

  Future<void> unsubscribeFromTopic(String topic) async {
    await _messaging.unsubscribeFromTopic(topic);
    if (kDebugMode) {
      print('Unsubscribed from topic: $topic');
    }
  }
}
`
	if err := os.WriteFile(filepath.Join(base, "core", "services", "notification_service.dart"), []byte(notificationService), 0644); err != nil {
		return err
	}

	return nil
}

func UpdateMainDart() error {
	mainContent := fmt.Sprintf(`import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'src/core/theme/app_theme.dart';
import 'src/core/router/app_router.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Set preferred orientations
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);
  
  // Set system UI overlay style
  SystemChrome.setSystemUIOverlayStyle(
    const SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness: Brightness.light,
      systemNavigationBarColor: Color(0xFF121212),
      systemNavigationBarIconBrightness: Brightness.light,
    ),
  );
  
  runApp(
    const ProviderScope(
      child: MyApp(),
    ),
  );
}

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: '%s',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.darkTheme,
      routerConfig: router,
    );
  }
}
`, ProjectName)

	if err := os.WriteFile(filepath.Join(ProjectName, "lib", "main.dart"), []byte(mainContent), 0644); err != nil {
		return err
	}

	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

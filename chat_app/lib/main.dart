import 'package:chat_app/core/config/router.dart';
import 'package:chat_app/core/storage/storage_service.dart';
import 'package:chat_app/shared/utils/logger.dart';
import 'package:chat_app/theme/app_theme.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

void main() async {
  // Ensure Flutter is initialized
  WidgetsFlutterBinding.ensureInitialized();

  // Set preferred orientations
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
    DeviceOrientation.landscapeLeft,
    DeviceOrientation.landscapeRight,
  ]);

  // Initialize storage services
  try {
    await StorageService.init();
  } catch (e, stackTrace) {
    AppLogger.error('Failed to initialize app', e, stackTrace);
  }

  // Run app
  runApp(
    const ProviderScope(
      child: MyLittlePriceApp(),
    ),
  );
}

class MyLittlePriceApp extends ConsumerWidget {
  const MyLittlePriceApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // TODO: Add theme mode provider to switch between light/dark
    const themeMode = ThemeMode.system;

    return MaterialApp.router(
      title: 'MyLittlePrice',
      debugShowCheckedModeBanner: false,

      // Theme
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,

      // Routing
      routerConfig: AppRouter.router,

      // Locale
      // TODO: Add localization support
      // localizationsDelegates: AppLocalizations.localizationsDelegates,
      // supportedLocales: AppLocalizations.supportedLocales,
    );
  }
}

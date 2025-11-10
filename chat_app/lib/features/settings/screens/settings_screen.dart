import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../theme/app_colors.dart';
import '../providers/settings_provider.dart';
import '../providers/settings_state.dart';
import '../widgets/widgets.dart';

/// Settings screen
class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final settings = ref.watch(settingsProvider);
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings'),
      ),
      body: ListView(
        children: [
          // Appearance Section
          _buildSectionHeader(context, 'Appearance'),
          _buildThemeModeCard(context, ref, settings),
          const SizedBox(height: 16),

          // Preferences Section
          _buildSectionHeader(context, 'Preferences'),
          _buildPreferencesTile(
            context,
            icon: Icons.public,
            title: 'Country',
            subtitle: _getCountryName(settings.country),
            onTap: () => _showCountryPicker(context, ref),
          ),
          _buildPreferencesTile(
            context,
            icon: Icons.language,
            title: 'Language',
            subtitle: _getLanguageName(settings.language),
            onTap: () => _showLanguagePicker(context, ref),
          ),
          _buildPreferencesTile(
            context,
            icon: Icons.attach_money,
            title: 'Currency',
            subtitle: _getCurrencyName(settings.currency),
            onTap: () => _showCurrencyPicker(context, ref),
          ),
          const SizedBox(height: 16),

          // Notifications Section
          _buildSectionHeader(context, 'Notifications'),
          _buildSwitchTile(
            context,
            icon: Icons.notifications_outlined,
            title: 'Push Notifications',
            subtitle: 'Receive notifications about search results',
            value: settings.notificationsEnabled,
            onChanged: (value) {
              ref.read(settingsProvider.notifier).toggleNotifications();
            },
          ),
          _buildSwitchTile(
            context,
            icon: Icons.volume_up_outlined,
            title: 'Sound',
            subtitle: 'Play sound for notifications',
            value: settings.soundEnabled,
            onChanged: (value) {
              ref.read(settingsProvider.notifier).toggleSound();
            },
          ),
          const SizedBox(height: 16),

          // About Section
          _buildSectionHeader(context, 'About'),
          _buildPreferencesTile(
            context,
            icon: Icons.info_outline,
            title: 'About',
            subtitle: 'Version 1.0.0',
            onTap: () => _showAboutDialog(context),
          ),
          _buildPreferencesTile(
            context,
            icon: Icons.privacy_tip_outlined,
            title: 'Privacy Policy',
            subtitle: 'Read our privacy policy',
            onTap: () {
              // TODO: Open privacy policy
            },
          ),
          _buildPreferencesTile(
            context,
            icon: Icons.description_outlined,
            title: 'Terms of Service',
            subtitle: 'Read our terms of service',
            onTap: () {
              // TODO: Open terms
            },
          ),
          const SizedBox(height: 16),

          // Reset Settings
          _buildSectionHeader(context, 'Advanced'),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: OutlinedButton.icon(
              onPressed: () => _confirmResetSettings(context, ref),
              icon: const Icon(Icons.restore),
              label: const Text('Reset to Defaults'),
              style: OutlinedButton.styleFrom(
                foregroundColor: AppColors.warning,
              ),
            ),
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(BuildContext context, String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Text(
        title,
        style: Theme.of(context).textTheme.titleSmall?.copyWith(
              color: AppColors.primary,
              fontWeight: FontWeight.w600,
            ),
      ),
    );
  }

  Widget _buildThemeModeCard(
    BuildContext context,
    WidgetRef ref,
    SettingsState settings,
  ) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(
                  Icons.palette_outlined,
                  color: AppColors.primary,
                ),
                const SizedBox(width: 12),
                Text(
                  'Theme',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            SegmentedButton<ThemeMode>(
              segments: const [
                ButtonSegment(
                  value: ThemeMode.light,
                  icon: Icon(Icons.light_mode),
                  label: Text('Light'),
                ),
                ButtonSegment(
                  value: ThemeMode.system,
                  icon: Icon(Icons.brightness_auto),
                  label: Text('Auto'),
                ),
                ButtonSegment(
                  value: ThemeMode.dark,
                  icon: Icon(Icons.dark_mode),
                  label: Text('Dark'),
                ),
              ],
              selected: {settings.themeMode},
              onSelectionChanged: (Set<ThemeMode> newSelection) {
                ref
                    .read(settingsProvider.notifier)
                    .setThemeMode(newSelection.first);
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPreferencesTile(
    BuildContext context, {
    required IconData icon,
    required String title,
    required String subtitle,
    required VoidCallback onTap,
  }) {
    return ListTile(
      leading: Icon(icon),
      title: Text(title),
      subtitle: Text(subtitle),
      trailing: const Icon(Icons.chevron_right),
      onTap: onTap,
    );
  }

  Widget _buildSwitchTile(
    BuildContext context, {
    required IconData icon,
    required String title,
    required String subtitle,
    required bool value,
    required ValueChanged<bool> onChanged,
  }) {
    return SwitchListTile(
      secondary: Icon(icon),
      title: Text(title),
      subtitle: Text(subtitle),
      value: value,
      onChanged: onChanged,
    );
  }

  String _getCountryName(String code) {
    final country = availableCountries.firstWhere(
      (c) => c.code == code,
      orElse: () => availableCountries.first,
    );
    return '${country.flag} ${country.name}';
  }

  String _getLanguageName(String code) {
    final language = availableLanguages.firstWhere(
      (l) => l.code == code,
      orElse: () => availableLanguages.first,
    );
    return language.nativeName;
  }

  String _getCurrencyName(String code) {
    final currency = availableCurrencies.firstWhere(
      (c) => c.code == code,
      orElse: () => availableCurrencies.first,
    );
    return '${currency.symbol} ${currency.code}';
  }

  void _showCountryPicker(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      builder: (context) => CountryPicker(
        onCountrySelected: (country) {
          ref.read(settingsProvider.notifier).setCountry(country.code);
          Navigator.pop(context);
        },
      ),
    );
  }

  void _showLanguagePicker(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      builder: (context) => LanguagePicker(
        onLanguageSelected: (language) {
          ref.read(settingsProvider.notifier).setLanguage(language.code);
          Navigator.pop(context);
        },
      ),
    );
  }

  void _showCurrencyPicker(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      builder: (context) => CurrencyPicker(
        onCurrencySelected: (currency) {
          ref.read(settingsProvider.notifier).setCurrency(currency.code);
          Navigator.pop(context);
        },
      ),
    );
  }

  void _showAboutDialog(BuildContext context) {
    showAboutDialog(
      context: context,
      applicationName: 'MyLittlePrice',
      applicationVersion: '1.0.0',
      applicationIcon: const Icon(
        Icons.shopping_bag,
        size: 48,
        color: AppColors.primary,
      ),
      children: [
        const Text(
          'AI-powered product search assistant that helps you find products through conversational chat.',
        ),
      ],
    );
  }

  Future<void> _confirmResetSettings(BuildContext context, WidgetRef ref) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Reset Settings'),
        content: const Text(
          'Are you sure you want to reset all settings to their default values?',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Reset'),
          ),
        ],
      ),
    );

    if (confirmed == true && context.mounted) {
      await ref.read(settingsProvider.notifier).resetToDefaults();
    }
  }
}

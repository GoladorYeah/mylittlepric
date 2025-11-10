import 'package:flutter/material.dart';
import '../providers/settings_state.dart';

/// Language picker bottom sheet
class LanguagePicker extends StatelessWidget {
  final Function(Language) onLanguageSelected;

  const LanguagePicker({
    super.key,
    required this.onLanguageSelected,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Handle bar
          Container(
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.onSurface.withOpacity(0.2),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          const SizedBox(height: 16),
          // Title
          Text(
            'Select Language',
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
          ),
          const SizedBox(height: 16),
          // Language list
          Flexible(
            child: ListView.builder(
              shrinkWrap: true,
              itemCount: availableLanguages.length,
              itemBuilder: (context, index) {
                final language = availableLanguages[index];
                return ListTile(
                  leading: const Icon(Icons.language),
                  title: Text(language.nativeName),
                  subtitle: Text(language.name),
                  trailing: Text(
                    language.code.toUpperCase(),
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: Theme.of(context)
                              .colorScheme
                              .onSurface
                              .withOpacity(0.5),
                        ),
                  ),
                  onTap: () => onLanguageSelected(language),
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

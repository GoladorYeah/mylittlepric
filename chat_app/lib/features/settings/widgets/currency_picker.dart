import 'package:flutter/material.dart';
import '../providers/settings_state.dart';

/// Currency picker bottom sheet
class CurrencyPicker extends StatelessWidget {
  final Function(Currency) onCurrencySelected;

  const CurrencyPicker({
    super.key,
    required this.onCurrencySelected,
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
            'Select Currency',
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
          ),
          const SizedBox(height: 16),
          // Currency list
          Flexible(
            child: ListView.builder(
              shrinkWrap: true,
              itemCount: availableCurrencies.length,
              itemBuilder: (context, index) {
                final currency = availableCurrencies[index];
                return ListTile(
                  leading: Text(
                    currency.symbol,
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  title: Text(currency.name),
                  subtitle: Text(currency.code),
                  onTap: () => onCurrencySelected(currency),
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

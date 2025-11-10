import 'package:flutter/material.dart';

/// Data class for parsed quick reply
class ParsedQuickReply {
  final String text;
  final String? price;

  const ParsedQuickReply({
    required this.text,
    this.price,
  });

  /// Parse quick reply to separate text and price
  /// Examples:
  /// - "Option (≈CHF 100-200)" -> text: "Option", price: "≈CHF 100-200"
  /// - "Option (CHF 100–200)" -> text: "Option", price: "CHF 100–200"
  /// - "Option (≈$100)" -> text: "Option", price: "≈$100"
  /// - "Option (CHF 500–1500+)" -> text: "Option", price: "CHF 500–1500+"
  static ParsedQuickReply parse(String reply) {
    // Match price patterns with various dash types: - – — (hyphen, en-dash, em-dash)
    final pricePattern = RegExp(
      r'\(([≈~]?[A-Z$€£¥]{1,4}[\s]?[\d,.\-–—]+[\+]?(?:[\s]?[kK]|[\s]?[\-–—][\s]?[\d,.\-–—]+[\+]?(?:[kK])?)?)\)$',
    );

    final match = pricePattern.firstMatch(reply);
    if (match != null) {
      final text = reply.substring(0, match.start).trim();
      final price = match.group(1);
      return ParsedQuickReply(text: text, price: price);
    }

    return ParsedQuickReply(text: reply);
  }
}

/// Widget displaying quick reply buttons with optional price badges
class QuickRepliesWidget extends StatelessWidget {
  final List<String> quickReplies;
  final ValueChanged<String> onReplyTap;

  const QuickRepliesWidget({
    super.key,
    required this.quickReplies,
    required this.onReplyTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Wrap(
      spacing: 8,
      runSpacing: 8,
      children: quickReplies.map((reply) {
        final parsed = ParsedQuickReply.parse(reply);

        return InkWell(
          onTap: () => onReplyTap(reply),
          borderRadius: BorderRadius.circular(8),
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(
              color: theme.colorScheme.secondaryContainer,
              borderRadius: BorderRadius.circular(8),
              border: Border.all(
                color: theme.colorScheme.outline.withOpacity(0.5),
              ),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Text
                Text(
                  parsed.text,
                  style: theme.textTheme.bodyMedium?.copyWith(
                    fontWeight: FontWeight.w500,
                    color: theme.colorScheme.onSecondaryContainer,
                  ),
                ),

                // Price badge (if exists)
                if (parsed.price != null) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 2,
                    ),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.primary.withOpacity(0.15),
                      borderRadius: BorderRadius.circular(4),
                      border: Border.all(
                        color: theme.colorScheme.primary.withOpacity(0.2),
                      ),
                    ),
                    child: Text(
                      parsed.price!,
                      style: theme.textTheme.labelSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: theme.colorScheme.primary,
                      ),
                    ),
                  ),
                ],
              ],
            ),
          ),
        );
      }).toList(),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:chat_app/features/settings/providers/settings_provider.dart';
import 'package:chat_app/core/config/constants.dart';

/// Chat input field with country selector and send button
class ChatInputWidget extends ConsumerStatefulWidget {
  final ValueChanged<String> onSend;
  final bool isLoading;
  final bool isConnected;
  final String connectionStatus;

  const ChatInputWidget({
    super.key,
    required this.onSend,
    required this.isLoading,
    required this.isConnected,
    required this.connectionStatus,
  });

  @override
  ConsumerState<ChatInputWidget> createState() => _ChatInputWidgetState();
}

class _ChatInputWidgetState extends ConsumerState<ChatInputWidget> {
  final _controller = TextEditingController();
  final _focusNode = FocusNode();

  @override
  void dispose() {
    _controller.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  void _handleSend() {
    final text = _controller.text.trim();
    if (text.isNotEmpty) {
      widget.onSend(text);
      _controller.clear();
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final country = ref.watch(countryProvider);

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: theme.colorScheme.background,
        border: Border(
          top: BorderSide(
            color: theme.colorScheme.outline.withOpacity(0.2),
          ),
        ),
      ),
      child: SafeArea(
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.end,
          children: [
            // Country selector dropdown
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                color: theme.colorScheme.secondaryContainer,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: theme.colorScheme.outline.withOpacity(0.3),
                ),
              ),
              child: DropdownButton<String>(
                value: country,
                underline: const SizedBox.shrink(),
                isDense: true,
                items: AppConstants.availableCountries.entries.map((entry) {
                  return DropdownMenuItem(
                    value: entry.key,
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          entry.value['flag'] ?? '',
                          style: const TextStyle(fontSize: 18),
                        ),
                        const SizedBox(width: 6),
                        Text(
                          entry.key,
                          style: theme.textTheme.bodySmall?.copyWith(
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  );
                }).toList(),
                onChanged: widget.isConnected
                    ? (value) {
                        if (value != null) {
                          ref
                              .read(settingsProvider.notifier)
                              .updateCountry(value);
                        }
                      }
                    : null,
              ),
            ),
            const SizedBox(width: 8),

            // Input field
            Expanded(
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                decoration: BoxDecoration(
                  color: theme.colorScheme.secondaryContainer,
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(
                    color: _focusNode.hasFocus
                        ? theme.colorScheme.primary
                        : theme.colorScheme.outline.withOpacity(0.3),
                  ),
                ),
                child: TextField(
                  controller: _controller,
                  focusNode: _focusNode,
                  enabled: widget.isConnected && !widget.isLoading,
                  maxLines: null,
                  textCapitalization: TextCapitalization.sentences,
                  decoration: InputDecoration(
                    hintText: widget.isConnected
                        ? 'Type your message...'
                        : '${widget.connectionStatus}...',
                    border: InputBorder.none,
                    hintStyle: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onSecondaryContainer
                          .withOpacity(0.5),
                    ),
                  ),
                  onSubmitted: (_) => _handleSend(),
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: theme.colorScheme.onSecondaryContainer,
                  ),
                ),
              ),
            ),
            const SizedBox(width: 8),

            // Send button
            Material(
              color: widget.isConnected && !widget.isLoading && _controller.text.trim().isNotEmpty
                  ? theme.colorScheme.primary
                  : theme.colorScheme.primary.withOpacity(0.5),
              borderRadius: BorderRadius.circular(12),
              child: InkWell(
                onTap: widget.isConnected && !widget.isLoading
                    ? _handleSend
                    : null,
                borderRadius: BorderRadius.circular(12),
                child: Container(
                  width: 48,
                  height: 48,
                  alignment: Alignment.center,
                  child: Icon(
                    Icons.send,
                    color: theme.colorScheme.onPrimary,
                    size: 20,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

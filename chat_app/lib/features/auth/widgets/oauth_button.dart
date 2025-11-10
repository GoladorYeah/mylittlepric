import 'package:flutter/material.dart';

/// OAuth provider button
class OAuthButton extends StatelessWidget {
  final String provider;
  final IconData icon;
  final Color backgroundColor;
  final Color textColor;
  final VoidCallback onPressed;

  const OAuthButton({
    super.key,
    required this.provider,
    required this.icon,
    required this.backgroundColor,
    required this.textColor,
    required this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton.icon(
      onPressed: onPressed,
      icon: Icon(icon, color: textColor),
      label: Text(
        'Continue with $provider',
        style: TextStyle(
          color: textColor,
          fontWeight: FontWeight.w600,
        ),
      ),
      style: ElevatedButton.styleFrom(
        backgroundColor: backgroundColor,
        foregroundColor: textColor,
        minimumSize: const Size(double.infinity, 56),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
          side: backgroundColor == Colors.white
              ? BorderSide(
                  color: Colors.grey.shade300,
                  width: 1,
                )
              : BorderSide.none,
        ),
        elevation: backgroundColor == Colors.white ? 0 : 2,
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:chat_app/core/config/constants.dart';

class AppTextStyles {
  // Display Styles
  static const TextStyle displayLarge = TextStyle(
    fontSize: AppConstants.fontSizeXxl,
    fontWeight: FontWeight.bold,
    height: 1.2,
    letterSpacing: -0.5,
  );

  static const TextStyle displayMedium = TextStyle(
    fontSize: AppConstants.fontSizeXl,
    fontWeight: FontWeight.bold,
    height: 1.2,
    letterSpacing: -0.25,
  );

  static const TextStyle displaySmall = TextStyle(
    fontSize: AppConstants.fontSizeL,
    fontWeight: FontWeight.bold,
    height: 1.3,
  );

  // Headline Styles
  static const TextStyle headlineLarge = TextStyle(
    fontSize: AppConstants.fontSizeL,
    fontWeight: FontWeight.w600,
    height: 1.3,
  );

  static const TextStyle headlineMedium = TextStyle(
    fontSize: 18.0,
    fontWeight: FontWeight.w600,
    height: 1.4,
  );

  static const TextStyle headlineSmall = TextStyle(
    fontSize: AppConstants.fontSizeM,
    fontWeight: FontWeight.w600,
    height: 1.4,
  );

  // Body Styles
  static const TextStyle bodyLarge = TextStyle(
    fontSize: AppConstants.fontSizeM,
    fontWeight: FontWeight.normal,
    height: 1.5,
  );

  static const TextStyle bodyMedium = TextStyle(
    fontSize: AppConstants.fontSizeS,
    fontWeight: FontWeight.normal,
    height: 1.5,
  );

  static const TextStyle bodySmall = TextStyle(
    fontSize: AppConstants.fontSizeXs,
    fontWeight: FontWeight.normal,
    height: 1.5,
  );

  // Label Styles
  static const TextStyle labelLarge = TextStyle(
    fontSize: AppConstants.fontSizeM,
    fontWeight: FontWeight.w500,
    height: 1.4,
  );

  static const TextStyle labelMedium = TextStyle(
    fontSize: AppConstants.fontSizeS,
    fontWeight: FontWeight.w500,
    height: 1.4,
  );

  static const TextStyle labelSmall = TextStyle(
    fontSize: AppConstants.fontSizeXs,
    fontWeight: FontWeight.w500,
    height: 1.4,
  );

  // Chat Specific
  static const TextStyle chatMessage = TextStyle(
    fontSize: AppConstants.fontSizeS,
    fontWeight: FontWeight.normal,
    height: 1.5,
  );

  static const TextStyle chatTimestamp = TextStyle(
    fontSize: 11.0,
    fontWeight: FontWeight.normal,
    height: 1.3,
  );

  static const TextStyle quickReply = TextStyle(
    fontSize: AppConstants.fontSizeS,
    fontWeight: FontWeight.w500,
    height: 1.4,
  );

  // Product Specific
  static const TextStyle productTitle = TextStyle(
    fontSize: AppConstants.fontSizeS,
    fontWeight: FontWeight.w600,
    height: 1.4,
  );

  static const TextStyle productPrice = TextStyle(
    fontSize: AppConstants.fontSizeM,
    fontWeight: FontWeight.bold,
    height: 1.3,
  );

  static const TextStyle productSource = TextStyle(
    fontSize: AppConstants.fontSizeXs,
    fontWeight: FontWeight.normal,
    height: 1.3,
  );

  // Button Styles
  static const TextStyle button = TextStyle(
    fontSize: AppConstants.fontSizeS,
    fontWeight: FontWeight.w600,
    height: 1.2,
    letterSpacing: 0.5,
  );

  static const TextStyle buttonSmall = TextStyle(
    fontSize: AppConstants.fontSizeXs,
    fontWeight: FontWeight.w600,
    height: 1.2,
    letterSpacing: 0.5,
  );
}

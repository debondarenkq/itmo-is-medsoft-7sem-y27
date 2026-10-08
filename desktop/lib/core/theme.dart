import 'package:flutter/material.dart';

abstract final class MedicalColors {
  static const teal = Color(0xFF087F83);
  static const ink = Color(0xFF183D43);
  static const muted = Color(0xFF6B8388);
  static const background = Color(0xFFF3F8F8);
  static const line = Color(0xFFE0ECEC);
  static const mint = Color(0xFFE7F5EE);
  static const green = Color(0xFF25865F);
  static const red = Color(0xFFBD555D);
}

ThemeData medicalTheme() => ThemeData(
  useMaterial3: true,
  fontFamily: 'HospitalSans',
  colorScheme: ColorScheme.fromSeed(
    seedColor: MedicalColors.teal,
    primary: MedicalColors.teal,
    surface: Colors.white,
    onSurface: MedicalColors.ink,
    onSurfaceVariant: MedicalColors.muted,
  ),
  scaffoldBackgroundColor: MedicalColors.background,
  dividerColor: MedicalColors.line,
  textTheme: const TextTheme(
    bodyMedium: TextStyle(fontSize: 14, color: MedicalColors.ink),
    bodySmall: TextStyle(fontSize: 12, color: MedicalColors.muted),
    titleLarge: TextStyle(
      fontSize: 24,
      fontWeight: FontWeight.w700,
      color: MedicalColors.ink,
    ),
  ),
  inputDecorationTheme: InputDecorationTheme(
    filled: true,
    fillColor: Colors.white,
    contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 17),
    border: OutlineInputBorder(
      borderRadius: BorderRadius.circular(10),
      borderSide: const BorderSide(color: MedicalColors.line),
    ),
    enabledBorder: OutlineInputBorder(
      borderRadius: BorderRadius.circular(10),
      borderSide: const BorderSide(color: MedicalColors.line),
    ),
    focusedBorder: OutlineInputBorder(
      borderRadius: BorderRadius.circular(10),
      borderSide: const BorderSide(color: MedicalColors.teal, width: 1.5),
    ),
  ),
  filledButtonTheme: FilledButtonThemeData(
    style: FilledButton.styleFrom(
      padding: const EdgeInsets.symmetric(horizontal: 22, vertical: 17),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
    ),
  ),
  outlinedButtonTheme: OutlinedButtonThemeData(
    style: OutlinedButton.styleFrom(
      padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 17),
      side: const BorderSide(color: MedicalColors.line),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
    ),
  ),
  cardTheme: CardThemeData(
    elevation: 0,
    color: Colors.white,
    shape: RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(16),
      side: const BorderSide(color: MedicalColors.line),
    ),
  ),
  snackBarTheme: SnackBarThemeData(
    elevation: 0,
    behavior: SnackBarBehavior.floating,
    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
  ),
  dataTableTheme: DataTableThemeData(
    headingRowColor: WidgetStateProperty.all(const Color(0xFFF7FAFA)),
    headingRowHeight: 52,
    dataRowMinHeight: 66,
    dataRowMaxHeight: 82,
    horizontalMargin: 22,
    columnSpacing: 28,
    headingTextStyle: const TextStyle(
      fontSize: 12,
      fontWeight: FontWeight.w600,
      color: MedicalColors.muted,
    ),
  ),
);

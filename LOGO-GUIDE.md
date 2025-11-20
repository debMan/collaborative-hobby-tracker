# Stash - Logo Design Guide

## Logo Concept

The Stash logo represents the core values of the Collaborative Hobby Tracker:

### Design Elements

1. **The Box/Container** - Represents "stashing" or collecting hobbies
   - Rounded corners for a friendly, modern feel
   - Blue color (#1a73e8) from the app's design system

2. **Checklist Items** - Symbolizes hobby tracking and organization
   - One checked item (completed hobby)
   - Two unchecked items (pending hobbies)
   - Clean, minimal lines reflecting the app's distraction-free philosophy

3. **Collaborative Dots** - Three small circles representing:
   - Different user circles (Partner, Friends, Family)
   - The collaborative nature of the app
   - Community and sharing features

### Color Palette

- **Primary Blue**: `#1a73e8` - Main brand color
- **Dark Blue**: `#1967d2` - Accents and depth
- **Light Blue**: `#4285f4` - Highlights
- **White**: `#ffffff` - Contrast and clarity

## Logo Variations

### 1. `logo.svg` (200×200)
- **Use for**: Website headers, documentation, social media
- **Features**: Full detail, circular background, all elements visible
- **Format**: SVG (scalable)

### 2. `icon.svg` (64×64)
- **Use for**: Favicon, app icons, small UI elements
- **Features**: Compact, high contrast, simplified details
- **Format**: SVG (can be converted to ICO/PNG)

### 3. `logo-with-text.svg` (400×120)
- **Use for**: Main branding, marketing materials, splash screens
- **Features**: Logo + "Stash" wordmark + tagline
- **Format**: SVG (scalable)

## Usage Guidelines

### Do's
✅ Use the logo on white or light backgrounds
✅ Maintain the aspect ratio when scaling
✅ Ensure minimum size of 32×32 pixels for icon
✅ Use official blue colors for brand consistency

### Don'ts
❌ Don't distort or stretch the logo
❌ Don't change the colors (except for monochrome versions)
❌ Don't add effects like shadows or gradients
❌ Don't place on busy backgrounds

## Converting for Different Platforms

### Favicon (ICO)
```bash
# Convert icon.svg to favicon
convert icon.svg -define icon:auto-resize=64,48,32,16 favicon.ico
```

### PNG Export
```bash
# Export logo.svg to PNG at different sizes
convert logo.svg -resize 512x512 logo-512.png
convert logo.svg -resize 256x256 logo-256.png
convert logo.svg -resize 128x128 logo-128.png
```

### iOS App Icons
```bash
# iOS requires specific sizes
convert icon.svg -resize 180x180 ios-icon-180.png  # iPhone @3x
convert icon.svg -resize 120x120 ios-icon-120.png  # iPhone @2x
convert icon.svg -resize 167x167 ios-icon-167.png  # iPad @2x
```

### Android App Icons
```bash
# Android adaptive icon
convert icon.svg -resize 432x432 android-icon-432.png  # xxxhdpi
convert icon.svg -resize 324x324 android-icon-324.png  # xxhdpi
convert icon.svg -resize 216x216 android-icon-216.png  # xhdpi
```

## Design Philosophy Alignment

The logo embodies the app's core principles:

- **Simplicity**: Clean, minimal design without clutter
- **Content-First**: The checklist items are the focal point
- **List-Based**: Linear, organized representation
- **Professional**: Mature color scheme and balanced composition
- **Collaborative**: Visual representation of multiple users/circles

## File Formats

All logos are provided as SVG (Scalable Vector Graphics) for:
- Infinite scalability without quality loss
- Small file size
- Easy editing and customization
- Browser-native support

## Brand Integration

These logos are designed to work seamlessly with:
- Material Design principles
- The app's blue color scheme
- List-based UI components
- Minimal, distraction-free interface

---

**Created for**: Stash - Collaborative Hobby Tracker
**Design Date**: 2025-11-20
**Version**: 1.0

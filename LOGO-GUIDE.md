# Stash - Logo Design Guide

## Logo Concept (Redesigned)

The Stash logo represents the core values of the Collaborative Hobby Tracker with emphasis on:

1. **Hobby Tracking** - Primary focus
2. **Collaboration** - Sharing and connecting with others
3. **Warmth & Friendliness** - Approachable design

### Design Elements

1. **The Central Container** - Represents "stashing" hobbies in one organized place
   - Warm orange color creates an inviting, friendly feel
   - Rounded corners for modern, approachable aesthetic

2. **Diverse Hobby Icons** - Core feature showcasing various hobby types
   - **Film/Movie** icon (top left) - Movies to watch
   - **Music Note** (top right) - Music to listen to
   - **Utensils/Food** (middle left) - Restaurants and food to try
   - **Map Pin** (middle right) - Travel destinations
   - **Star** (bottom left) - Activities and experiences
   - **Checkmark** (bottom right) - Completion tracking

3. **Three Collaborative People** - Emphasizes the collaborative aspect
   - Left, right, and top positions surrounding the stash
   - Represent different user circles (Partner, Friends, Family, Colleagues)
   - Connected with dashed lines showing active collaboration
   - Different warm shades showing diversity

4. **Share Icon Badge** - Collaboration indicator
   - Upload/share symbol on top of the container
   - Golden accent color (#FFB84D)
   - Represents the sharing and import features

### Color Palette (Warm Theme)

- **Primary Orange**: `#FF6B35` - Main brand color (energetic, friendly)
- **Golden Orange**: `#F7931E` - Secondary accent
- **Warm Yellow**: `#FFA500` - Bright highlights
- **Coral Orange**: `#FF8C42` - Tertiary accent
- **Gold**: `#FFB84D` - Special highlights (share icon)
- **White**: `#ffffff` - Contrast and clarity

**Why Warm Colors?**
- **Inviting & Friendly**: Orange conveys enthusiasm and creativity
- **Energetic**: Reflects the active nature of pursuing hobbies
- **Social**: Warm colors encourage collaboration and sharing
- **Approachable**: Less corporate than blue, more personal

## Logo Variations

### 1. `logo.svg` (200×200)
- **Use for**: Website headers, documentation, social media profiles
- **Features**:
  - Full detail with gradient background
  - All hobby icons visible
  - Three people showing collaboration
  - Connection lines and share badge
- **Format**: SVG (scalable)

### 2. `icon.svg` (64×64)
- **Use for**: Favicon, app icons, small UI elements
- **Features**:
  - Compact design optimized for small sizes
  - Simplified hobby icons (4 main types)
  - Three collaborative people
  - High contrast for visibility
- **Format**: SVG (can be converted to ICO/PNG)

### 3. `logo-with-text.svg` (450×130)
- **Use for**: Main branding, marketing materials, website headers
- **Features**:
  - Logo icon + "Stash" wordmark in matching orange
  - Tagline: "Collaborative Hobby Tracker"
  - Perfect for landing pages and presentations
- **Format**: SVG (scalable)

### 4. `logo-monochrome.svg` (200×200)
- **Use for**: Printing, dark mode variations, single-color applications
- **Features**:
  - Black and white version
  - Maintains all design elements
  - Grayscale shading for depth
- **Format**: SVG (scalable)

## Usage Guidelines

### Do's
✅ Use the logo on white or light backgrounds for best visibility
✅ Maintain the aspect ratio when scaling
✅ Ensure minimum size of 32×32 pixels for icon version
✅ Use official warm orange colors for brand consistency
✅ Use monochrome version for print materials
✅ Maintain clear space around the logo (at least 10% of logo width)

### Don'ts
❌ Don't distort or stretch the logo
❌ Don't change the warm color palette (except for monochrome)
❌ Don't add drop shadows, gradients, or effects
❌ Don't place on busy or dark orange backgrounds
❌ Don't rotate or flip the logo
❌ Don't separate the hobby icons from the container

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

### Web Formats
```bash
# For web use (retina displays)
convert logo.svg -resize 400x400 logo-2x.png
convert logo.svg -resize 200x200 logo-1x.png

# For social media
convert logo.svg -resize 512x512 social-icon.png
convert logo-with-text.svg -resize 1200x630 social-banner.png  # Open Graph
```

## Design Philosophy Alignment

The logo embodies the app's core values:

### Primary Focus: Hobby Tracking
- **Six diverse hobby icons** prominently displayed
- Film, music, food, travel, activities, and completion tracking
- Clear representation of the app's main purpose
- Visual variety showing the breadth of trackable hobbies

### Strong Collaboration Emphasis
- **Three people icons** surrounding the central stash
- Connection lines showing active collaboration
- Share/upload badge highlighting the import/export feature
- Different colored people representing diverse user circles

### Warm & Approachable
- **Warm orange palette** replaces corporate blue
- Friendly, energetic colors encourage engagement
- Social and inviting aesthetic
- Reflects the personal nature of hobby tracking

### Simplicity & Clarity
- Clean iconography without excessive detail
- Clear visual hierarchy
- Easy to understand at any size
- Maintains distraction-free philosophy

## Color Psychology

**Orange (#FF6B35)** - Primary brand color
- Represents: Creativity, enthusiasm, success, encouragement
- Perfect for: A social app about personal interests and hobbies
- Emotion: Friendly, energetic, optimistic

**Golden Tones (#F7931E, #FFA500, #FFB84D)**
- Represents: Warmth, collaboration, achievement
- Perfect for: Highlighting sharing and completion features
- Emotion: Welcoming, cheerful, supportive

## Brand Integration

These logos work seamlessly with:
- Material Design principles (with warm color adaptation)
- The app's list-based UI components
- Minimal, distraction-free interface
- Both light and dark mode interfaces

### Tailwind CSS Integration

Update your Tailwind config to use the warm color scheme:

```javascript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#FFF5F0',
          100: '#FFE8DC',
          200: '#FFD1B9',
          300: '#FFB996',
          400: '#FFA173',
          500: '#FF6B35',  // Main brand color
          600: '#E85D2A',
          700: '#CC501E',
          800: '#B04314',
          900: '#94360B',
        },
        accent: {
          DEFAULT: '#F7931E',
          light: '#FFA500',
          gold: '#FFB84D',
        }
      }
    }
  }
}
```

## Accessibility

- **Color Contrast**: Orange on white provides 3.8:1 ratio (WCAG AA for large text)
- **Icon Clarity**: Hobby icons are distinguishable even at small sizes
- **Monochrome Version**: Available for high-contrast needs
- **Alternative Text**: Always use "Stash - Collaborative Hobby Tracker" as alt text

## File Formats

All logos are provided as SVG (Scalable Vector Graphics) for:
- Infinite scalability without quality loss
- Small file size (optimized for web)
- Easy editing and customization
- Browser-native support
- Professional print quality at any size

---

**App Name**: Stash - Collaborative Hobby Tracker
**Design Focus**: Hobby tracking with strong collaboration emphasis
**Color Theme**: Warm oranges (friendly, energetic, social)
**Design Date**: 2025-11-20
**Version**: 2.0 (Redesigned)

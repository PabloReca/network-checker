#!/bin/bash

# Configuration
LOGO="build/appicon.png"
BG_COLOR="#EDDE29"
OUTPUT_DIR="build/darwin/icon.iconset"
FINAL_ICNS="build/darwin/icon.icns"
FINAL_ICO="build/windows/icon.ico"

echo "🎨 Generating icons with background $BG_COLOR..."

# Create temp directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Clean the logo: ensure it has transparency and remove potential white backgrounds
# We'll create a temp clean logo
CLEAN_LOGO="build/temp_logo_clean.png"
magick "$LOGO" -trim -fuzz 5% -transparent white "$CLEAN_LOGO"

# Function to create a square with the logo overlayed
generate_step() {
    SIZE=$1
    NAME=$2
    
    # Calculate dimensions
    # macOS squircle radius is approx 22%
    RADIUS=$(awk -v s=$SIZE 'BEGIN { print s * 0.22 }')
    # Logo should be around 70% of the size
    INNER_SIZE=$(awk -v s=$SIZE 'BEGIN { print s * 0.7 }')
    
    # 1. Create a solid background colored square
    # 2. Apply a rounded rectangle mask (Squircle-ish)
    # 3. Composite the clean logo on top
    magick -size "${SIZE}x${SIZE}" xc:none \
        -fill "$BG_COLOR" -draw "roundrectangle 0,0,$SIZE,$SIZE,$RADIUS,$RADIUS" \
        \( "$CLEAN_LOGO" -resize "${INNER_SIZE}x${INNER_SIZE}" -gravity center \) \
        -gravity center -composite \
        "$OUTPUT_DIR/$NAME.png"
}

# Generate all sizes for macOS
for size in 16 32 128 256 512; do
    generate_step $size "icon_${size}x${size}"
    generate_step $((size * 2)) "icon_${size}x${size}@2x"
done

# Create ICNS
iconutil -c icns "$OUTPUT_DIR"

# Keep the largest PNG as a preview
cp "$OUTPUT_DIR/icon_512x512@2x.png" "build/icon_preview.png"

# Create ICO for Windows
# We'll use a version without the squircle mask for Windows if preferred, 
# but usually matching branding is better. Let's keep the squircle for consistency.
magick "$OUTPUT_DIR/icon_256x256.png" -define icon:auto-resize=256,128,64,48,32,16 "$FINAL_ICO"

# Cleanup
rm -rf "$OUTPUT_DIR"
rm "$CLEAN_LOGO"

echo "✅ Icons generated successfully!"
echo "📍 macOS: $FINAL_ICNS (Check if this looks yellow/EDDE29)"
echo "📍 Windows: $FINAL_ICO"

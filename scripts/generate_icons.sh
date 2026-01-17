#!/bin/bash

# Configuration
LOGO="build/appicon.png"
BG_COLOR="#EDDE29"
OUTPUT_DIR="build/darwin/icon.iconset"
FINAL_ICNS="build/darwin/icon.icns"
FINAL_ICO="build/windows/icon.ico"

echo "🎨 Generating icons with background $BG_COLOR..."

# Create temp directory
mkdir -p "$OUTPUT_DIR"

# Function to create a square with the logo overlayed
generate_step() {
    SIZE=$1
    NAME=$2
    # Create the background squircle
    # We use a rounded rectangle that covers approx 80% of the area to look like a macOS icon
    # or we can just fill the whole square if the user wants it full-bleed.
    # macOS icons usually have some padding.
    
    # 1. Create a canvas of SIZE x SIZE with BG_COLOR
    # 2. Apply a rounded corner mask (radius is roughly 22% of size)
    # 3. Overlay the logo (scaled to roughly 60-70% of the size)
    
    RADIUS=$(echo "$SIZE * 0.22" | bc)
    INNER_SIZE=$(echo "$SIZE * 0.7" | bc)
    
    magick -size "${SIZE}x${SIZE}" xc:none \
        -fill "$BG_COLOR" -draw "roundrectangle 0,0,$SIZE,$SIZE,$RADIUS,$RADIUS" \
        \( "$LOGO" -resize "${INNER_SIZE}x${INNER_SIZE}" \) \
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

# Create ICO for Windows (using the 256px version as base)
magick "$OUTPUT_DIR/icon_256x256.png" -define icon:auto-resize=256,128,64,48,32,16 "$FINAL_ICO"

# Cleanup
rm -rf "$OUTPUT_DIR"

echo "✅ Icons generated successfully!"
echo "📍 macOS: $FINAL_ICNS"
echo "📍 Windows: $FINAL_ICO"

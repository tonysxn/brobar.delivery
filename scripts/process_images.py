import os
import subprocess
from pathlib import Path

# Configuration
SOURCE_DIR = 'uploads_backup'
DEST_DIR = 'uploads'
MAX_SIZE = 800
QUALITY = 80

def get_base_name(filename):
    # Removes extensions and hashes if present
    # Strategy: split by dot.
    # If the file is "batat-fri.jpg", base is "batat-fri".
    # If "batat-fri.697b8d8a119085.79514815.jpg", base is "batat-fri".
    parts = filename.split('.')
    if len(parts) > 2: 
        # Heuristic: if more than 2 parts, likely "slug.hash.hash.ext"
        # We assume the first part is the slug.
        # But wait, what if slug has dots? Unlikely for slugs.
        return parts[0]
    return os.path.splitext(filename)[0]

def optimize_image(source_path, dest_path):
    # Use ffmpeg to resize and convert
    # scale=800:-1:flags=lanczos: sets width to 800, height auto, high quality scaling
    # If using 'force_original_aspect_ratio', we can make sure it fits within box.
    # command = f'ffmpeg -y -i "{source_path}" -vf "scale=w={MAX_SIZE}:h={MAX_SIZE}:force_original_aspect_ratio=decrease" -q:v {QUALITY} "{dest_path}"'
    
    # We use qscale:v for quality control in ffmpeg for webp? 
    # For libwebp, -qscale is mapped to quality 0-100?
    # ffmpeg docs say -q:v is for quality.
    
    cmd = [
        'ffmpeg', '-y', '-v', 'error',
        '-i', source_path,
        '-vf', f'scale=\'min({MAX_SIZE},iw)\':-1', # Only downscale
        '-c:v', 'libwebp',
        '-q:v', str(QUALITY),
        dest_path
    ]
    
    subprocess.run(cmd, check=True)

def main():
    if not os.path.exists(DEST_DIR):
        os.makedirs(DEST_DIR)

    processed_slugs = set()
    updates = []
    
    # Check what files we have
    files = sorted(os.listdir(SOURCE_DIR))
    
    # Filter for image files
    image_files = [f for f in files if f.lower().endswith(('.jpg', '.jpeg', '.png'))]
    
    # Prefer non-hashed files if duplicates exist?
    # Actually, if we have 'batat-fri.jpg' and 'batat-fri.hash.jpg', they are likely the same or 'batat-fri.jpg' is better named.
    # Let's simple check if we already processed a slug.
    
    # We want to match what's in the DB.
    # The DB has "batat-fri.697b8d8a119085.79514815.jpg".
    # But we want to UPDATE the DB to use "batat-fri.webp".
    # So we just need to produce "batat-fri.webp" from WHATEVER source image we find for that slug.
    
    count = 0
    for filename in image_files:
        slug = get_base_name(filename)
        
        if slug in processed_slugs:
            continue
            
        source_path = os.path.join(SOURCE_DIR, filename)
        dest_filename = f"{slug}.webp"
        dest_path = os.path.join(DEST_DIR, dest_filename)
        
        print(f"Processing {filename} -> {dest_filename}")
        
        try:
            optimize_image(source_path, dest_path)
            processed_slugs.add(slug)
            
            # Generate SQL update
            # careful with quoting
            updates.append(f"UPDATE products SET image = '{dest_filename}' WHERE slug = '{slug}';")
            count += 1
        except Exception as e:
            print(f"Error processing {filename}: {e}")

    # Write SQL script
    with open('scripts/update_images.sql', 'w') as f:
        f.write('\n'.join(updates))
        
    print(f"Finished. Processed {count} images. SQL script generated at scripts/update_images.sql")

if __name__ == "__main__":
    main()

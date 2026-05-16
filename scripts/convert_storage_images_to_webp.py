"""
One-time script to convert all images in the Supabase `event-media` bucket
to WebP. Downloads each image, re-encodes it as WebP, and re-uploads to the
SAME path with `Content-Type: image/webp`.

Filenames, locations, and extensions stay exactly the same (a `.jpg` file
keeps the name `.jpg`, only its bytes become WebP). Browsers render it fine
because Supabase serves the stored content-type, not the extension, so
nothing in the DB needs to be re-added or renamed.

JPEG/JPG -> lossy WebP (quality 80, same as the compress script).
PNG      -> lossless WebP (keeps transparency and sharp edges, still smaller).

Images are downscaled (never upscaled) so the longest edge fits the cap:
files under /_gallery/ -> 1280px, everything else (cover/location/etc) ->
1920px. EXIF orientation is baked into the pixels first so phone photos
don't come out rotated.

Note: do NOT re-run compress_storage_images.py after this. That script picks
its codec from the extension, so it would re-encode these WebP-in-.jpg files
back to JPEG and undo the conversion.

Usage:
    pip install supabase Pillow
    python convert_storage_images_to_webp.py

Reads SUPABASE_URL and SUPABASE_SERVICE_KEY from .env via python-dotenv,
or you can pass them as env vars.
"""

import os
import sys
from io import BytesIO
from pathlib import Path

from dotenv import load_dotenv
from PIL import Image, ImageOps

from supabase import create_client

load_dotenv()

BUCKET = "event-media"
WEBP_QUALITY = 80
WEBP_LOSSLESS_QUALITY = 100
WEBP_METHOD = 6
SUPPORTED_EXTENSIONS = {".jpg", ".jpeg", ".png"}
MIN_SAVINGS_RATIO = 0.05
GALLERY_DIR_MARKER = "/_gallery/"
GALLERY_MAX_EDGE = 1280
DEFAULT_MAX_EDGE = 1920


def format_bytes(size: int) -> str:
    """Returns a human-readable file size string."""
    if size < 1024:
        return f"{size} B"
    if size < 1024 * 1024:
        return f"{size / 1024:.1f} KB"
    return f"{size / (1024 * 1024):.2f} MB"


def is_image(filepath: str) -> bool:
    """Checks if a file path has a supported (convertible) image extension."""
    return Path(filepath).suffix.lower() in SUPPORTED_EXTENSIONS


def list_all_files(storage, prefix: str = "") -> list[str]:
    """Recursively lists every file path in the bucket under `prefix`."""
    paths: list[str] = []
    result = storage.from_(BUCKET).list(
        prefix, {"limit": 1000, "sortBy": {"column": "name", "order": "asc"}}
    )

    for item in result:
        full_path = f"{prefix}/{item['name']}" if prefix else item["name"]
        if item.get("id"):
            paths.append(full_path)
        else:
            paths.extend(list_all_files(storage, full_path))

    return paths


def convert_to_webp(data: bytes, ext: str, max_edge: int) -> bytes | None:
    """Re-encodes image bytes as WebP, downscaling so the longest edge is
    at most `max_edge` (never upscales).

    EXIF orientation is baked into the pixels first so phone photos don't
    come out rotated. JPEGs go lossy (quality 80), PNGs go lossless so
    transparency and sharp edges survive. Returns None if the format isn't
    convertible.
    """
    img = Image.open(BytesIO(data))
    img = ImageOps.exif_transpose(img)

    if max(img.size) > max_edge:
        img.thumbnail((max_edge, max_edge), Image.LANCZOS)

    buf = BytesIO()

    if ext in (".jpg", ".jpeg"):
        # JPEGs have no alpha, drop palette/alpha modes WebP can't take as RGB.
        if img.mode not in ("RGB", "L"):
            img = img.convert("RGB")
        img.save(buf, format="WEBP", quality=WEBP_QUALITY, method=WEBP_METHOD)
    elif ext == ".png":
        # Keep alpha for lossless WebP; WebP only takes RGB/RGBA.
        has_alpha = img.mode in ("RGBA", "LA", "PA") or (
            img.mode == "P" and "transparency" in img.info
        )
        img = img.convert("RGBA" if has_alpha else "RGB")
        img.save(
            buf,
            format="WEBP",
            lossless=True,
            quality=WEBP_LOSSLESS_QUALITY,
            method=WEBP_METHOD,
        )
    else:
        return None

    return buf.getvalue()


def main() -> None:
    url = os.getenv("SUPABASE_URL")
    key = os.getenv("SUPABASE_SERVICE_KEY")

    if not url or not key:
        print("Missing SUPABASE_URL or SUPABASE_SERVICE_KEY.")
        print("Pass them as env vars or add them to your .env file.")
        sys.exit(1)

    supabase = create_client(url, key)
    storage = supabase.storage

    print(f'\nScanning bucket "{BUCKET}"...\n')

    all_files = list_all_files(storage)
    image_files = [f for f in all_files if is_image(f)]
    ignored = len(all_files) - len(image_files)

    print(f"Found {len(all_files)} files total, {len(image_files)} images (jpg/png).")
    if ignored > 0:
        print(f"Ignoring {ignored} non-convertible files (webp, mp4, mp3, etc.).")
    print()

    if not image_files:
        print("Nothing to convert.")
        return

    total_original = 0
    total_converted = 0
    converted_count = 0
    skipped_count = 0

    for i, filepath in enumerate(image_files, 1):
        ext = Path(filepath).suffix.lower()
        progress = f"[{i}/{len(image_files)}]"
        max_edge = (
            GALLERY_MAX_EDGE if GALLERY_DIR_MARKER in filepath else DEFAULT_MAX_EDGE
        )

        try:
            data = storage.from_(BUCKET).download(filepath)
        except Exception as e:
            print(f"{progress} SKIP {filepath} — download failed: {e}")
            skipped_count += 1
            continue

        original_size = len(data)
        total_original += original_size

        try:
            webp_data = convert_to_webp(data, ext, max_edge)

            savings = 1 - (len(webp_data) / original_size) if webp_data else 0
            if webp_data is None or savings < MIN_SAVINGS_RATIO:
                print(
                    f"{progress} SKIP {filepath} ({format_bytes(original_size)}) — already optimal"
                )
                total_converted += original_size
                skipped_count += 1
                continue

            # Same path, same name. Only the bytes and content-type change.
            storage.from_(BUCKET).update(
                filepath,
                webp_data,
                {"content-type": "image/webp", "upsert": "true"},
            )

            saved = original_size - len(webp_data)
            pct = (saved / original_size) * 100
            total_converted += len(webp_data)
            converted_count += 1

            print(
                f"{progress} OK   {filepath}  "
                f"{format_bytes(original_size)} → {format_bytes(len(webp_data))}  "
                f"(−{pct:.1f}%)"
            )
        except Exception as e:
            print(f"{progress} FAIL {filepath} — {e}")
            total_converted += original_size
            skipped_count += 1

    total_saved = total_original - total_converted
    total_pct = (total_saved / total_original * 100) if total_original > 0 else 0

    print("\n--- Summary ---")
    print(f"Converted: {converted_count} files")
    print(f"Skipped:   {skipped_count} files (already optimal or errors)")
    print(f"Before:    {format_bytes(total_original)}")
    print(f"After:     {format_bytes(total_converted)}")
    print(f"Saved:     {format_bytes(total_saved)} ({total_pct:.1f}%)")
    print()


if __name__ == "__main__":
    main()

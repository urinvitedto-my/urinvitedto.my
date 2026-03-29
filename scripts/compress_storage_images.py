"""
One-time script to compress all images in the Supabase `event-media` bucket.
Downloads each image, compresses with Pillow, and re-uploads to the same path.
Names, locations, and formats stay the same — only file size changes.

Usage:
    pip install supabase Pillow
    python compress_storage_images.py

Reads SUPABASE_URL and SUPABASE_SERVICE_KEY from .env via python-dotenv,
or you can pass them as env vars.
"""

import os
import sys
from io import BytesIO
from pathlib import Path

from dotenv import load_dotenv
from PIL import Image

from supabase import create_client

load_dotenv()

BUCKET = "event-media"
JPEG_QUALITY = 80
PNG_COMPRESS_LEVEL = 9
SUPPORTED_EXTENSIONS = {".jpg", ".jpeg", ".png"}
MIN_SAVINGS_RATIO = 0.05


def format_bytes(size: int) -> str:
    """Returns a human-readable file size string."""
    if size < 1024:
        return f"{size} B"
    if size < 1024 * 1024:
        return f"{size / 1024:.1f} KB"
    return f"{size / (1024 * 1024):.2f} MB"


def is_image(filepath: str) -> bool:
    """Checks if a file path has a supported image extension."""
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


def compress_image(data: bytes, ext: str) -> bytes | None:
    """Compresses image bytes with Pillow, keeping the same format."""
    img = Image.open(BytesIO(data))

    buf = BytesIO()
    if ext in (".jpg", ".jpeg"):
        if img.mode in ("RGBA", "P"):
            img = img.convert("RGB")
        img.save(buf, format="JPEG", quality=JPEG_QUALITY, optimize=True)
    elif ext == ".png":
        img.save(buf, format="PNG", compress_level=PNG_COMPRESS_LEVEL, optimize=True)
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
        print(f"Ignoring {ignored} non-image files (mp4, mp3, etc.).")
    print()

    if not image_files:
        print("Nothing to compress.")
        return

    total_original = 0
    total_compressed = 0
    compressed_count = 0
    skipped_count = 0

    for i, filepath in enumerate(image_files, 1):
        ext = Path(filepath).suffix.lower()
        progress = f"[{i}/{len(image_files)}]"
        content_type = "image/jpeg" if ext in (".jpg", ".jpeg") else "image/png"

        try:
            data = storage.from_(BUCKET).download(filepath)
        except Exception as e:
            print(f"{progress} SKIP {filepath} — download failed: {e}")
            skipped_count += 1
            continue

        original_size = len(data)
        total_original += original_size

        try:
            compressed_data = compress_image(data, ext)

            savings = 1 - (len(compressed_data) / original_size) if compressed_data else 0
            if compressed_data is None or savings < MIN_SAVINGS_RATIO:
                print(
                    f"{progress} SKIP {filepath} ({format_bytes(original_size)}) — already optimized"
                )
                total_compressed += original_size
                skipped_count += 1
                continue

            storage.from_(BUCKET).update(
                filepath,
                compressed_data,
                {"content-type": content_type, "upsert": "true"},
            )

            saved = original_size - len(compressed_data)
            pct = (saved / original_size) * 100
            total_compressed += len(compressed_data)
            compressed_count += 1

            print(
                f"{progress} OK   {filepath}  "
                f"{format_bytes(original_size)} → {format_bytes(len(compressed_data))}  "
                f"(−{pct:.1f}%)"
            )
        except Exception as e:
            print(f"{progress} FAIL {filepath} — {e}")
            total_compressed += original_size
            skipped_count += 1

    total_saved = total_original - total_compressed
    total_pct = (total_saved / total_original * 100) if total_original > 0 else 0

    print("\n--- Summary ---")
    print(f"Compressed: {compressed_count} files")
    print(f"Skipped:    {skipped_count} files (already optimized or errors)")
    print(f"Before:     {format_bytes(total_original)}")
    print(f"After:      {format_bytes(total_compressed)}")
    print(f"Saved:      {format_bytes(total_saved)} ({total_pct:.1f}%)")
    print()


if __name__ == "__main__":
    main()

#!/usr/bin/env python3
"""
image_to_rom_csv.py

Reads an image, resamples it down to a grid, and writes the pixel data
to a CSV file as a single packed 24-bit value per address:

    packed = R * (2**0) + G * (2**8) + B * (2**16)

i.e. R is the low byte, G is the middle byte, B is the high byte.

All settings are configured in the CONFIG section below — just edit the
values and run the script (no command-line arguments needed).
"""

from PIL import Image
import csv


# =============================================================
# CONFIG — edit these values directly
# =============================================================

INPUT_IMAGE   = "D:\\Golang\\GeareoDecodingV2\\in.png"    # path to the source image
OUTPUT_CSV    = "D:\\Golang\\GeareoDecodingV2\\out.csv"   # path to write the CSV to

GRID_WIDTH    = 25             # number of columns to sample down to
GRID_HEIGHT   = 20             # number of rows to sample down to

# Resampling filter used when resizing to the grid.
# "nearest"  -> blocky, best for pixel-art style sampling (default)
# "bilinear" / "bicubic" / "lanczos" -> smoother, blends pixels
RESAMPLE = "nearest"

# Address format for the CSV:
# "linear" -> one "address" column (bottom row left->right, then upward)
# "rowcol" -> separate "row" and "col" columns instead of one address
ADDRESS_FORMAT = "linear"

INVERT       = False  # if True, invert each channel (255 - value) before packing
WRITE_HEADER = False   # if True, write a CSV header row

# =============================================================


RESAMPLE_FILTERS = {
    "nearest": Image.NEAREST,
    "bilinear": Image.BILINEAR,
    "bicubic": Image.BICUBIC,
    "lanczos": Image.LANCZOS,
}


def load_grid(path, width, height, resample):
    img = Image.open(path).convert("RGB")
    img = img.resize((width, height), RESAMPLE_FILTERS[resample])
    return img


def pack_rgb(r, g, b, invert):
    if invert:
        r, g, b = 255 - r, 255 - g, 255 - b
    return r * (2 ** 0) + g * (2 ** 8) + b * (2 ** 16)


def main():
    img = load_grid(INPUT_IMAGE, GRID_WIDTH, GRID_HEIGHT, RESAMPLE)

    with open(OUTPUT_CSV, "w", newline="") as f:
        writer = csv.writer(f)

        if WRITE_HEADER:
            if ADDRESS_FORMAT == "linear":
                writer.writerow(["address", "packed_rgb"])
            else:
                writer.writerow(["row", "col", "packed_rgb"])

        address = 0
        for y in range(GRID_HEIGHT - 1, -1, -1):
            for x in range(GRID_WIDTH):
                r, g, b = img.getpixel((x, y))
                packed = pack_rgb(r, g, b, INVERT)

                if ADDRESS_FORMAT == "linear":
                    writer.writerow([packed])
                else:
                    writer.writerow([packed])

                address += 1

    total = GRID_WIDTH * GRID_HEIGHT
    print(f"Wrote {total} entries ({GRID_WIDTH}x{GRID_HEIGHT} grid) to {OUTPUT_CSV}")


if __name__ == "__main__":
    main()
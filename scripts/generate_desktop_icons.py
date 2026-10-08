#!/usr/bin/env python3
"""Render the app's geometric medical cross into native icons (stdlib only)."""
import math
from pathlib import Path
import struct
import zlib

ROOT = Path(__file__).resolve().parent.parent / 'desktop'


def rounded_distance(x, y, cx, cy, half_x, half_y, radius):
    dx, dy = abs(x - cx) - half_x + radius, abs(y - cy) - half_y + radius
    return math.hypot(max(dx, 0), max(dy, 0)) + min(max(dx, dy), 0) - radius


def png(size):
    rows = bytearray()
    for iy in range(size):
        rows.append(0)
        for ix in range(size):
            x, y = (ix + .5) / size, (iy + .5) / size
            background = rounded_distance(x, y, .5, .5, .46, .46, .15)
            cross = min(rounded_distance(x, y, .5, .5, .075, .25, .025),
                        rounded_distance(x, y, .5, .5, .25, .075, .025))
            alpha = max(0, min(1, .5 - background * size))
            white = max(0, min(1, .5 - cross * size))
            color = [round(c + (255 - c) * white) for c in (8, 127, 131)]
            rows.extend((*color, round(255 * alpha)))
    def chunk(kind, data):
        return struct.pack('>I', len(data)) + kind + data + struct.pack('>I', zlib.crc32(kind + data))
    return b'\x89PNG\r\n\x1a\n' + chunk(b'IHDR', struct.pack('>IIBBBBB', size, size, 8, 6, 0, 0, 0)) + chunk(b'IDAT', zlib.compress(rows, 9)) + chunk(b'IEND', b'')


def main():
    folder = ROOT / 'macos/Runner/Assets.xcassets/AppIcon.appiconset'
    for size in (16, 32, 64, 128, 256, 512, 1024):
        (folder / f'app_icon_{size}.png').write_bytes(png(size))
    sizes = (16, 32, 48, 64, 128, 256)
    images = [png(size) for size in sizes]
    offset = 6 + 16 * len(sizes)
    header = bytearray(struct.pack('<HHH', 0, 1, len(sizes)))
    for size, image in zip(sizes, images):
        header.extend(struct.pack('<BBBBHHII', size % 256, size % 256, 0, 0, 1, 32, len(image), offset))
        offset += len(image)
    (ROOT / 'windows/runner/resources/app_icon.ico').write_bytes(header + b''.join(images))


if __name__ == '__main__':
    main()

#!/usr/bin/env python3
"""Generate a QR code image for the Sets in Go speakerdeck URL."""

import qrcode

URL = "https://speakerdeck.com/ramalho/sets-in-go"
OUTPUT_FILE = "sets-in-go-qr.png"

img = qrcode.make(URL)
img.save(OUTPUT_FILE)
print(f"QR code saved to {OUTPUT_FILE}")

#!/bin/bash
# Launch SimplePDFCompress in the background
# Output is discarded and process is detached from terminal

nohup ./simplepdfcompress >/dev/null 2>&1 &

echo "SimplePDFCompress launched in background."

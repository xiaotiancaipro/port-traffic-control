#!/bin/bash

echo "[$(date +'%Y-%m-%d %H:%M:%S')] Start building multi-platform binaries"

APP_NAME="port-traffic-control"
OUTPUT_PATH="./dist"
platforms=(
  "linux/amd64"
  "linux/arm64"
)

mkdir -p $OUTPUT_PATH
echo "[$(date +'%Y-%m-%d %H:%M:%S')] Create Output Directory: $OUTPUT_PATH"

total_platforms=${#platforms[@]}
success_count=0
fail_count=0

echo "[$(date +'%Y-%m-%d %H:%M:%S')] Start building $total_platforms Platforms..."

for platform in "${platforms[@]}"; do

    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    OUTPUT="$OUTPUT_PATH/$APP_NAME-$GOOS-$GOARCH"

    echo "[$(date +'%Y-%m-%d %H:%M:%S')] ▶ Build $GOOS/$GOARCH ..."

    build_log=$(GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT 2>&1)
    exit_status=$?

    if [ $exit_status -eq 0 ]; then
        echo "[$(date +'%Y-%m-%d %H:%M:%S')] ✓ $GOOS/$GOARCH build successfully "
        ((success_count++))
    else
        echo "[$(date +'%Y-%m-%d %H:%M:%S')] ✗ $GOOS/$GOARCH build failed (exit code: $exit_status)"
        echo "[$(date +'%Y-%m-%d %H:%M:%S')] Build error messages:"
        echo "[$(date +'%Y-%m-%d %H:%M:%S')]     $build_log"
        ((fail_count++))
    fi

done

echo "[$(date +'%Y-%m-%d %H:%M:%S')] Build Complete!"
echo "[$(date +'%Y-%m-%d %H:%M:%S')] Total: $total_platforms"
echo "[$(date +'%Y-%m-%d %H:%M:%S')] Successful: $success_count"
echo "[$(date +'%Y-%m-%d %H:%M:%S')] Failed: $fail_count"
echo "[$(date +'%Y-%m-%d %H:%M:%S')] Output path: $OUTPUT_PATH"

if [ $fail_count -gt 0 ]; then
    exit 1
else
    exit 0
fi

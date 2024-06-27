#! /bin/bash
OUTPUT_DIR=./bin

# Array of directories where your Go modules are located
declare -a goModuleDirs=("internal/cmd/generator" "playerData")

mkdir -p "$OUTPUT_DIR"

# Loop through each directory
for dir in "${goModuleDirs[@]}"; do
    echo "Building module $dir"
    (
        # cd "$dir" || exit # Change to the directory, exit if it fails
        go build -o "$OUTPUT_DIR/$(basename "$dir")" $dir/main.go > "$OUTPUT_DIR/$(basename "$dir").log" 2>&1
    )
done

wait

echo "Finished running build script."

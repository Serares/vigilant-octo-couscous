#! /bin/bash

declare -a goModuleDirs=("hfDescr" "materialize" "playerData" "internal/cmd/generator")

# Loop through each directory
for dir in "${goModuleDirs[@]}"; do
    echo "Starting service in $dir"
    (
        cd "$dir" || exit # Change to the directory, exit if it fails
        service_name="/" read -ra ADDR <<<"$dir"
        go run main.go >log.log 2>&1 & # Run the Go server and redirect output to logs.log
        # give them time to bootstrap
        sleep 1
    )
done

echo "All services started."

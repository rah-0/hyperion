#!/bin/bash
cd ..

# Function to run benchmarks and generate profiling outputs
run_benchmarks() {
  local profile_name="$1" # Third parameter: name for pprof output files (e.g., "custom_profile")
  local bench_time="5s"

  # Construct profiling output file paths based on profile_name
  local mem_profile="./pprof_svg/${profile_name}_mem.out"
  local cpu_profile="./pprof_svg/${profile_name}_cpu.out"
  local mem_svg="./pprof_svg/${profile_name}_mem.svg"
  local cpu_svg="./pprof_svg/${profile_name}_cpu.svg"

  # Ensure the output directory exists
  mkdir -p ./pprof_svg

  # Run the benchmark and generate profiling outputs
  go test -run=^$ -bench=. -benchmem ./ -benchtime="${bench_time}" -timeout=0 -memprofile="${mem_profile}" -cpuprofile="${cpu_profile}" \
    && go tool pprof -svg -output="${mem_svg}" "${mem_profile}" \
    && go tool pprof -svg -output="${cpu_svg}" "${cpu_profile}"

  # Check for errors
  if [[ $? -eq 0 ]]; then
    echo "Benchmarks and profiles generated successfully."
    echo "Memory profile saved to: ${mem_svg}"
    echo "CPU profile saved to: ${cpu_svg}"
  else
    echo "An error occurred during benchmarks or profile generation."
  fi
}

rm -rf ./pprof_svg/*

run_benchmarks "node"

rm hyperion.test

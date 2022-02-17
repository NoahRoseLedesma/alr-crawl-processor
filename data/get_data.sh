#!/bin/bash

# Script to facilitate downloading WARC files from common crawl.

# Which web crawl to download. A list of the available crawls can be found here:
# https://commoncrawl.org/the-data/get-started/ 
CRAWL_NAME="CC-MAIN-2022-05"

# Number of WARC files to download.
NUM_WARC_FILES=60

# HTTP Server hosting the crawling data
BASE_URL="https://commoncrawl.s3.amazonaws.com/"

download_paths() {
	# Download the WARC paths manifest for the crawl
	local dl_file=$(mktemp)
	wget "${BASE_URL}crawl-data/${CRAWL_NAME}/warc.paths.gz" -O "${dl_file}"
	# Extract the paths.gz file
	warc_paths=$(zcat ${dl_file} | head -n ${NUM_WARC_FILES})
}

download_paths
total_paths=0
for path in $warc_paths; do
	let total_paths++
done

current_idx=1
for path in $warc_paths; do
	echo Downloading "${current_idx}"/"${total_paths}"
	let current_idx++
	# Download the WARC file. Retrying if the server returns HTTP 503
	# (rate limer error)
	wget -cq --show-progress "${BASE_URL}${path}" --retry-on-http-error=503
done


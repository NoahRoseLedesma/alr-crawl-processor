# ARIA Landmark Web Crawl Processing

This application searches for HTML documents in web crawls performed by
[CommonCrawl](https://commoncrawl.org/) that contain ARIA Landmark Roles. It
reads web crawling data in the WARC format and saves all of the documents
containing
[valid landmark roles](https://www.washington.edu/accessibility/web/landmarks/)
to separate WARC files.

## Downloading the crawl data

The script `data/get_data.sh` facilitates the downloading a percentage of a
CommonCrawl crawl.

Modify the `NUM_WARC_FILES` variable to change how many WARC files are
downloaded. Each file is approximately one gigabyte.

Modify the `CRAWL_NAME` variable to change which crawl is
downloaded. By default the January 2022 (`CC-MAIN-2022-05`) is used.

## Running the landmark scanner

The script `main.go` scans all of the gzipped WARC files in `data/` and writes
the records containing HTML documents with valid landmark roles to WARC files
in the `output/` directory.

To run this program use

```
go install
```

followed by

```
go run main.go
```

Or, use the provided Dockerfile. Ensure to mount your data directory to the
`/app/data/` folder and output directory to the `/app/output` folder in the
container like so:

```
docker run -it -v /path/to/data:/app/data -v /path/to/output:/app/output <image name>
```
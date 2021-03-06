# Mani

Create AWS Redshift manifest files from a file list.

## Application

I wanted to a simple command line application which I could embed into an ETL
workflow without the need for any external dependencies.

## Installation

This will build with the standard [Go tools](http://golang.org/). No further dependancies required unless you want to build cross platform binaries, in which case the `make dist` action will use [gox](https://github.com/mitchellh/gox) for a few common target environments.

You can download binaries from the [github release pages](https://github.com/psmithuk/mani/releases/tag/v0.0.1)

## Usage

```bash
usage: mani [flags]
  -i="": input filename (stdout if none provided)
  -o="": output filename (stdout if none provided)
  -p="": s3 bucketname or prefix
  -version=false: print version string
```

For example, given the file `files.txt` with the following content:

```text
2013-10-04-custdata
2013-10-05-custdata
2013-10-04-custdata
2013-10-05-custdata
```

and the command:

```bash
mani -i=files.txt -o=manifest.json -p="s3://mybucket-alpha/"
```

the file `manifest.json` with the following json will be created:

```json
{
  "entries": [
    {
      "url": "s3://mybucket-alpha/2013-10-04-custdata",
      "mandatory": true
    },
    {
      "url": "s3://mybucket-alpha/2013-10-05-custdata",
      "mandatory": true
    },
    {
      "url": "s3://mybucket-alpha/2013-10-04-custdata",
      "mandatory": true
    },
    {
      "url": "s3://mybucket-alpha/2013-10-05-custdata",
      "mandatory": true
    }
  ]
}
```

A more complex usage example would be:

```bash
aws s3 ls s3://my-bucket/customers-sales/ --recursive \
  | awk '{print $4}' \
  | grep -v '^$' \
  | grep '2013' \
  | mani -p="s3://my-bucket/" \
  > manifest.json
```

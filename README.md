# Imagelinter

Imagelinter is used to verify and validate docker images.Docker images can be from dockerfiles or yaml | yml files.

## The problem

In order to adhere to open source compliance, some repos should not contain Alpine based images in any Docker or Image files. Hence, we make sure that all docker images (from Dockerfiles or from yaml|yml files) should not be Alpine-based images.

## Solution

There is no straight command or a way to identify Linux OS from Docker perspective. After a deep analysis, we have found a few solutions that can apparently identify Linux OS of a given image.The below are steps that are executed when the linter is run.

1. Read Image metadata from the registry/Image libraries
2. Pull image and analyze Image meta-data (History)
3. Create a container and read /etc/os-release file
4. Copy etc/os-release to the local path and then analyze
5. Copy /usr/lib/os-release to the local path and then analyze
6. Check if there is any License file inside the container
7. Copy a simple binary to the container
8. Nothing determines the OS then reject the image

## Configuration

imagelinter supports few configurations.The below is the default configuration file.

```yaml
---
includeExts:
- ".yaml"
- ".sh"
- ".yml"
includeFiles:
- README.md
includeLines:
- 'image:'
- FROM
excludeLines:
- "#"
- "//"
excludeFiles:
- "cli/"
- .git/
- docs/
- ".gitignore"
- ".github/"
- "*.md"
- "*.sh"
succesValidators:
- apt-get
- apt
- yum
- "/lib/x86_64-linux-gnu"
- "/usr/lib/x86_64-linux-gnu"
- "imgpkg"
failureValidators:
- Alpine
```

## How to install imagelinter

To download source code and install imagelinter 
```git clone https://github.com/JitenPalaparthi/imagelinter.git```

cd to the imagelinter directory and run 
```go install github.com/JitenPalaparthi/imagelinter```

The above command generates the single binary file.Default configuration is embed into it.

## How to run imagelinter

- ```go run main.go --path <any valid path or it takes only current working directory> -- config <provide config path or default path will be taken> --summary=true --details=fail```

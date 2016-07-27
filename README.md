# nvd-mirror

A simple shell script to mirror the CPE/CVE XML data from NIST.

The purpose of nvd-mirror is to be able to replicate the NIST vulnerabiity 
data inside a company, so that local access to NIST data can be achieved.

If you need a version that works on windows or for any reason you cannot use
wget or gzip (used in the bash script), you can build the golang client.

# Usage

    nvd-mirror.sh <output_directory>
    
# Build utility in golang
	go get github.com/ralberto/nvd-mirror
    go build
    
    
# Similar projects
- [nist-data-mirror](https://github.com/stevespringett/nist-data-mirror) - Utility in Java


# License
Please eee the [LICENSE](./LICENSE) file.


#!/bin/bash

BASE_URL="https://nvd.nist.gov"
START_YEAR=2002
END_YEAR=`date "+%Y"`

if [ "$#" -ne 1 ]
then
  echo "Usage: $0 <outputDir>"
  exit 1
fi

OUTPUT_DIR=$1

function download {
  urlToFetch=$1
  outDir=$2

  gzFile=$(echo $urlToFetch | rev | cut -d "/" -f 1 | rev)
  uncompressedFile=$(echo $gzFile | rev | cut -d "." -f 2- | rev)

  echo "Conditional download: $urlToFetch"
  fetchOutput=$(wget -P $outDir -N -c -nv $urlToFetch 2>&1)
  if [ $? -eq 0 ]; then
    if [ ! -z "$fetchOutput" ]; then
      echo "Uncompressing file $gzFile"
      gzip -d < $outDir/$gzFile > $outDir/$uncompressedFile
    else
      echo "Using cached version of $urlToFetch"
    fi
  else
    echo "Unable to fetch $urlToFetch. $fetchOutput"
    exit $?
  fi
}


download $BASE_URL/download/nvdcve-Modified.xml.gz $OUTPUT_DIR
download $BASE_URL/feeds/xml/cve/nvdcve-2.0-Modified.xml.gz $OUTPUT_DIR

for year in $(seq $START_YEAR $END_YEAR); do
  download $BASE_URL/download/nvdcve-"$year".xml.gz $OUTPUT_DIR
  download $BASE_URL/feeds/xml/cve/nvdcve-2.0-"$year".xml.gz $OUTPUT_DIR
done


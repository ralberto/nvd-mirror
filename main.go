/**
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
* THE SOFTWARE.
*
* This is a golang application that can be called from the command line to
* download contents from NIST CPE/CVE XML data to the specified output path.
*
* It's an alternative to nvd-mirror.sh if for any reason you can't use wget or
* gzip. You'll need to build the binary from this source to the appropriate
* target architecture.
*/
package main

import (
	"net/http"
	"strings"
	"fmt"
	"os"
	"strconv"
	"io"
	"compress/gzip"
	"time"
)

const (
	CVE_12_MODIFIED_URL = "https://nvd.nist.gov/download/nvdcve-Modified.xml.gz";
	CVE_20_MODIFIED_URL = "https://nvd.nist.gov/feeds/xml/cve/nvdcve-2.0-Modified.xml.gz";
	CVE_12_BASE_URL = "https://nvd.nist.gov/download/nvdcve-%d.xml.gz";
	CVE_20_BASE_URL = "https://nvd.nist.gov/feeds/xml/cve/nvdcve-2.0-%d.xml.gz";
	START_YEAR = 2002;
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Usage: %s <outputDir>\n", args[0]);
		os.Exit(-2);
	}
	outDir := args[1]

	err := Get(CVE_12_MODIFIED_URL, outDir)
	if (err != nil) {
		fmt.Printf("%s\n", err.Error())
		os.Exit(-1)
	}

	err = Get(CVE_20_MODIFIED_URL, outDir)
	if (err != nil) {
		fmt.Printf("%s\n", err.Error())
		os.Exit(-1)
	}

	currentYear,_,_ := time.Now().Date()

	for year := START_YEAR; year <= currentYear; year++ {
		err = Get(strings.Replace(CVE_12_BASE_URL, "%d", strconv.Itoa(year), 1), outDir)
		if (err != nil) {
			fmt.Printf("%s\n", err.Error())
			os.Exit(-1)
		}

		err = Get(strings.Replace(CVE_20_BASE_URL, "%d", strconv.Itoa(year), 1), outDir)
		if (err != nil) {
			fmt.Printf("%s\n", err.Error())
			os.Exit(-1)
		}
	}
}

func Get(url, outDir string) error {
	idx := strings.LastIndex(url, "/")
	fileName := url[idx + 1:]
	absOutputFile := outDir + "/" + fileName

	remoteContentLength, err := contentLength(url)
	if (err != nil) {
		return err
	}

	localFileContentLength, err := fileSize(absOutputFile)
	if (err != nil) {
		return err
	}

	if (remoteContentLength == localFileContentLength) {
		fmt.Printf("%s => Cached\n", url)
	} else {
		fmt.Printf("Downloading %s ...", url)
		bytes, err := download(url, absOutputFile)
		if (err != nil) {
			return err
		}
		fmt.Printf(" => %d bytes downloaded\n", bytes)
		err = unzip(absOutputFile)
		if (err != nil) {
			return err
		}
	}
	return nil
}

func fileSize(filepath string) (int64, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return -1, nil
	}

	file, err := os.Open(filepath)
	defer file.Close()

	if err != nil {
		return -1, err
	}

	fi, err := file.Stat()
	if err != nil {
		return -1, err
	}

	return fi.Size(), nil
}

func contentLength(url string) (size int64, err error) {
	resp, err := http.Head(url)
	if (err != nil) {
		return -1, err
	}

	return resp.ContentLength, nil
}


//Output is written into the same directory
func unzip(absInputFile string) error {
	reader, err := os.Open(absInputFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()
	target := strings.Replace(absInputFile, ".gz", "", 1)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = io.Copy(writer, archive)
	return err

}

func download(url, absOutputFile string) (bytes int64, err error) {
	out, err := os.Create(absOutputFile)
	if (err != nil ) {
		return -1, err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if (err != nil) {
		return -1, err
	}

	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)

	if (err != nil) {
		return n, err
	}
	return n, nil
}

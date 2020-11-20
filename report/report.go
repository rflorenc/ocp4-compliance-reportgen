package report

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	rgerr "github.com/rflorenc/ocp4-compliance-reportgen/error"
	rgexec "github.com/rflorenc/ocp4-compliance-reportgen/exec"
	"github.com/rflorenc/ocp4-compliance-reportgen/templates"
)

// RootHTMLDir is the NGINX server root
const RootHTMLDir = "/opt/nginx/html"

const fileFinalMode = 444
const fileAppendMode = 0644
const maxScans = 19

var data = templates.TemplateData{
	Title:      "Scan Results",
	ScanNumber: 0,
	TemplateMap: map[string]string{
		"d1": templates.D1,
		"d2": templates.D2,
		"d3": templates.D3,
		"d4": templates.D4,
		"d5": templates.D5,
		"d6": templates.D6,
	},
}

var indexNewFileHandle *os.File

// ParseExecTemplate Parses and Executes templates
func ParseExecTemplate(absPath string, templateMapKey string, tm map[string]string) (*os.File, error) {
	_, exists := tm[templateMapKey]
	if exists {
		f, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileAppendMode)
		if err != nil {
			log.Fatalf("ParseExecTemplate os.OpenFile: %v ", err)
			return nil, err
		}
		defer f.Close()

		t, terr := template.New(templateMapKey).Parse(tm[templateMapKey])
		rgerr.CheckError(fmt.Sprintf("template.New(%s)", templateMapKey), terr)

		terr = t.Execute(f, data)
		rgerr.CheckError(fmt.Sprintf("t.Execute(%s)", templateMapKey), terr)

		return f, terr
	}
	log.Fatalf("Non existing key: %v", templateMapKey)
	return nil, nil
}

// GenerateReportHTML generates the HTML reports by using the defined templates and the oscap report results
func GenerateReportHTML() {
	// Create first HTML block in RootHTMLDir/index_new.html
	ParseExecTemplate(filepath.Join(RootHTMLDir, "index_new.html"), "d1", data.TemplateMap)

	for data.ScanNumber <= maxScans {
		log.Printf("Processing scan number: %d\n", data.ScanNumber)
		// Scans are created / mounted at /opt/nginx/html/?
		fi, err := os.Lstat(filepath.Join(RootHTMLDir, strconv.Itoa(data.ScanNumber)))
		if err != nil {
			log.Printf("os.Lstat")
			break
		}

		// Create first HTML block in RootHTMLDir/scannumber/index_new.html
		indexNew := filepath.Join(RootHTMLDir, strconv.Itoa(data.ScanNumber), "index_new.html")
		indexNewFileHandle, _ = ParseExecTemplate(indexNew, "d2", data.TemplateMap)

		bzipList, err := filepath.Glob(filepath.Join(RootHTMLDir, strconv.Itoa(data.ScanNumber), "*.bzip2"))
		if err != nil {
			log.Fatal(err)
		}
		for _, bzipFile := range bzipList {
			name := strings.TrimRight(bzipFile, "-pod.xml.bzip2")
			aux := strings.TrimPrefix(name, filepath.Join(RootHTMLDir, strconv.Itoa(data.ScanNumber)))
			trimmedName := strings.TrimLeft(aux, string(os.PathSeparator))
			fmt.Printf("name: %s", name)

			exists := fileExists(name + "result.html")
			if !exists {
				outXML, err := os.OpenFile(name+"-result.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileAppendMode)
				if err != nil {
					panic(err)
				}
				defer outXML.Close()

				log.Printf("Extracting %s ...", bzipFile)
				bunzipCmd := exec.Command("bunzip2", bzipFile, "-c")
				bunzipCmd.Stdout = outXML
				err = bunzipCmd.Start()
				if err != nil {
					panic(err)
				}
				bunzipCmd.Wait()

				log.Printf("Converting %s to HTML ...", outXML.Name())
				outHTML, err := os.OpenFile(name+"-result.html", os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileAppendMode)
				if err != nil {
					panic(err)
				}
				defer outHTML.Close()

				oscapCmd := exec.Command("oscap", "xccdf", "generate", "report", outXML.Name(), "-c")
				oscapCmd.Stdout = outHTML
				err = oscapCmd.Start()
				if err != nil {
					panic(err)
				}
				oscapCmd.Wait()
			}

			data.UpdateScanNameResult(trimmedName + "-result")
			data.UpdateScanNameResultHTML(trimmedName + "-result.html")
			data.UpdateScanNameModifiedTime(fi.ModTime().String())

			// Append per scan oscap HTML to RootHTMLDir/scannumber/index_new.html
			indexNewFileHandle, _ = ParseExecTemplate(indexNew, "d3", data.TemplateMap)
		}

		// Append closing HTML to RootHTMLDir/scannumber/index_new.html
		indexNewFileHandle, _ = ParseExecTemplate(indexNew, "d4", data.TemplateMap)

		prefixScanPath := filepath.Join(RootHTMLDir, strconv.Itoa(data.ScanNumber))
		log.Printf("prefixScanPath: %s\n", prefixScanPath)

		src := filepath.Join(prefixScanPath, "index_new.html")
		dest := filepath.Join(prefixScanPath, "index.html")
		cleanup(src, dest, strconv.Itoa(fileFinalMode))

		// Append closing HTML to RootHTMLDir/index_new.html
		_, _ = ParseExecTemplate(filepath.Join(RootHTMLDir, "index_new.html"), "d5", data.TemplateMap)

		data.IncrementScanNumber()
	}

	// Append closing HTML to RootHTMLDir/index_new.html
	_, _ = ParseExecTemplate(filepath.Join(RootHTMLDir, "index_new.html"), "d6", data.TemplateMap)

	src := filepath.Join(RootHTMLDir, "index_new.html")
	dest := filepath.Join(RootHTMLDir, "index.html")
	cleanup(src, dest, strconv.Itoa(fileFinalMode))

}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func cleanup(src, dest, fileMode string) {
	log.Printf("src: %s, dest: %s, fileMode: %s", src, dest, fileMode)
	rgexec.ExecuteCmd("rm", "-f", dest)
	rgexec.ExecuteCmd("mv", src, dest)
	rgexec.ExecuteCmd("chmod", fileMode, dest)
}

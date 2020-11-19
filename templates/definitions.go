package templates

// TemplateData defines fields to be used during HTML templating
type TemplateData struct {
	Title                string
	ScanNumber           int
	ScanNameResult       string
	ScanNameResultHTML   string
	ScanNameModifiedTime string

	TemplateMap map[string]string
}

// IncrementScanNumber increments ScanNumber counter
func (td *TemplateData) IncrementScanNumber() {
	td.ScanNumber++
}

// D1 is the first block of templated HTML to be applied to the index.html located in the root folder
// by default at /opt/nginx/html/
// index.html simply contains a sequence of scans (0..N) containing an href
// to the actual per SCAN page containing the oscap html results, which is built dinamically using templates D2 to D5.
//
// Scan Results
// Scan 0 HTML Reports
// Scan 1 HTML Reports
// Scan 2 HTML Reports
const D1 = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<link rel="stylesheet" href="css/bootstrap.min.css">
		<title>{{.Title}}</title>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
			<meta http-equiv="refresh" content="30">
		<style>
			a.list-group-item:hover{
			background: #d5d5ff;
			}
		</style>
	</head>
	<body>
	<h1> {{.Title}} </h1>
	<div class="list-group">
`

// D2 is initial portion of the the per scan html report,
// with added metadata, such as modification time and scan name.
//
// Scan 0 Results
// Scan Name	                     Scan Time	                            Report
// rhcos4-e8-master-example.f.q.d.n	 2020-10-27 15:25:41.076401137 +0000	HTML Report
// rhcos4-e8-node-example.f.q.d.n	 2020-10-27 15:25:37.794438632 +0000	HTML Report
const D2 = `
	<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link rel="stylesheet" href="../css/bootstrap.min.css">
		<title>Scan {{.ScanNumber}} Results</title>
		<style>
			.table-striped tbody tr:nth-of-type(odd) {
				background-color: #e4e4e4;
		}
			.table tbody tr:hover {background-color: #d5d5ff;}
		</style>
	</head>
	<body>
	  <h1> Scan {{.ScanNumber}} Results</h1>
	  <table class="table table-striped table-bordered table-hover">
		<thead>
		  <tr>
			<th scope="col">Scan Name</th>
			<th scope="col">Scan Time</th>
			<th scope="col">Report</th>
		  </tr>
		</thead>
		<tbody>
`

// D3 is a template block that is appended to D2
// Adds Scan Metadata
const D3 = `          <tr><th scope="row">{{.ScanNameResult}}</th><td>{{.ScanNameModifiedTime}}</td><td><a href="{{.ScanNameResultHTML}}">HTML Report</a></td></tr>`

// D4 is a template block to be appended to D2
// Closes D2 HTML
const D4 = `      </tbody></table>
</body>
</html>
`

// D5 is a template block to be appended to D1
// Adds href to the final html reports
const D5 = `        <a href="{{.ScanNumber}}/index.html" class="list-group-item list-group-item-action">Scan {{.ScanNumber}} HTML Reports</a>"`

// D6 is a template block to be appended to D1
// Closes D1 HTML
const D6 = `      </div>
</body>
</html>`

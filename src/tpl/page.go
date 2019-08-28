package tpl

import (
	"../serverError"
	"../util"
	"net/url"
	"path"
	"text/template"
)

const pageTplStr = `
<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	<meta http-equiv="X-UA-Compatible" content="IE=edge"/>
	<meta name="viewport" content="initial-scale=1"/>
	<meta name="format-detection" content="telephone=no"/>
	<meta name="renderer" content="webkit"/>
	<meta name="wap-font-scale" content="no"/>
	<base href="{{.Scheme}}//{{.Host}}/{{if .Path}}{{.Path}}/{{end}}"/>

	<title>{{.Path}}</title>

	<style type="text/css">
		html, body {
			margin: 0;
			padding: 0;
			background: #fff;
		}

		html {
			font-family: "roboto_condensedbold", "Helvetica Neue", Helvetica, Arial, sans-serif;
		}

		body {
			color: #333;
			font-size: 0.625em;
			font-family: Consolas, Monaco, "Andale Mono", "DejaVu Sans Mono", monospace;
		}

		a {
			display: block;
			padding: 0.25em 0.5em;
			color: inherit;
			text-decoration: none;
		}

		a:hover {
			color: #000;
			background: #f5f5f5;
		}

		input, button {
			margin: 0;
			padding: 0.25em 0;
		}

		.path-list {
			font-size: 1.4em;
			overflow: hidden;
			border-bottom: 1px #999 solid;
		}

		.path-list a {
			position: relative;
			float: left;
			padding-right: 1.2em;
			text-align: center;
			white-space: nowrap;
			min-width: 1em;
		}

		.path-list a:after {
			content: '';
			position: absolute;
			top: 0.6em;
			right: 0.4em;
			width: 0.4em;
			height: 0.4em;
			border: 1px solid;
			border-color: #ccc #ccc transparent transparent;
			transform: rotate(45deg);
		}

		.path-list a:last-child {
			padding-right: 0.5em;
		}

		.path-list a:last-child:after {
			display: none;
		}

		.upload {
			position: relative;
			margin: 1em;
			padding: 1em;
			background: #f7f7f7;
		}

		.upload.dragging::before {
			content: '';
			position: absolute;
			left: 0;
			top: 0;
			right: 0;
			bottom: 0;
			opacity: 0.7;
			background: #c9c;
		}

		.upload form {
			margin: 0;
			padding: 0;
		}

		.upload input {
			display: block;
			width: 100%;
			box-sizing: border-box;
		}

		.upload input + input {
			margin-top: 0.5em;
		}

		.item-list {
			margin: 1em;
		}

		.item-list a {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			border-bottom: 1px #f5f5f5 solid;
		}

		.item-list span {
			margin-left: 1em;
			flex-shrink: 0;
		}

		.item-list .name {
			flex: 1 1 0;
			margin-left: 0;
			font-size: 1.4em;
			word-break: break-all;
		}

		.item-list .size {
			white-space: nowrap;
			text-align: right;
			color: #666;
		}

		.item-list .time {
			width: 10em;
			color: #999;
			text-align: right;
			white-space: nowrap;
			overflow: hidden;
		}

		.error {
			margin: 1em;
			padding: 1em;
			background: #ffc;
		}
	</style>
</head>
<body>

<div class="path-list">
	<a href="/">/</a>
    {{range $path := .Paths}}
		<a href="{{$path.Path}}">{{html $path.Name}}</a>
    {{end}}
</div>

{{if .CanUpload}}
	<div class="upload">
		<form method="POST" enctype="multipart/form-data">
			<input type="file" name="files" class="files" multiple="multiple" accept="*/*"/>
			<input type="submit" value="Upload"/>
		</form>
	</div>
{{end}}

<div class="item-list">
	<a href="../">
		<span class="name">../</span>
		<span class="size"></span>
		<span class="time"></span>
	</a>
    {{range .SubItems}}
        {{$isDir := .IsDir}}
		<a href="./{{path .Name}}" class="item {{if $isDir}}item-dir{{else}}item-file{{end}}">
			<span class="name">{{html .Name}}{{if $isDir}}/{{end}}</span>
			<span class="size">{{if not $isDir}}{{fmtSize .Size}}{{end}}</span>
			<span class="time">{{fmtTime .ModTime}}</span>
		</a>
    {{end}}
</div>

{{range $error := .Errors}}
	<div class="error">{{$error}}</div>
{{end}}

<script type="text/javascript">
    (function () {
        if (!document.querySelector) {
            return;
        }

        var upload = document.querySelector('.upload');
        if (!upload || !upload.addEventListener) {
            return;
        }
        var fileInput = upload.querySelector('.files');

        var addClass = function (ele, className) {
            ele && ele.classList && ele.classList.add(className)
        };

        var removeClass = function (ele, className) {
            ele && ele.classList && ele.classList.remove(className)
        };

        var onDragEnterOver = function (e) {
            e.stopPropagation();
            e.preventDefault();
            addClass(e.currentTarget, 'dragging');
        };

        var onDragLeave = function (e) {
            removeClass(e.currentTarget, 'dragging');
        };

        var onDrop = function (e) {
            e.stopPropagation();
            e.preventDefault();
            removeClass(e.currentTarget, 'dragging');

            if (!e.dataTransfer.files) {
                return;
            }
            fileInput.files = e.dataTransfer.files;
        };

        upload.addEventListener('dragenter', onDragEnterOver);
        upload.addEventListener('dragover', onDragEnterOver);
        upload.addEventListener('dragleave', onDragLeave);
        upload.addEventListener('drop', onDrop);
    })()
</script>
</body>
</html>
`

var defaultPage *template.Template

func init() {
	tplObj := template.New("page")
	tplObj = addFuncMap(tplObj)

	var err error
	defaultPage, err = tplObj.Parse(pageTplStr)
	if serverError.CheckError(err) {
		defaultPage = template.Must(tplObj.Parse("Builtin Template Error"))
	}
}

func LoadPage(tplPath string) *template.Template {
	var tplObj *template.Template
	var err error

	if len(tplPath) > 0 {
		tplObj = template.New(path.Base(tplPath))
		tplObj = addFuncMap(tplObj)
		tplObj, err = tplObj.ParseFiles(tplPath)
		serverError.CheckError(err)
	}
	if err != nil || len(tplPath) == 0 {
		tplObj = defaultPage
	}

	return tplObj
}

func addFuncMap(tpl *template.Template) *template.Template {
	return tpl.Funcs(template.FuncMap{
		"path":    url.PathEscape,
		"fmtSize": util.FormatSize,
		"fmtTime": util.FormatTimeMinute,
	})
}
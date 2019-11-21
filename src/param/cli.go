package param

import (
	"../goNixArgParser"
	"../serverErrHandler"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var cliParams []*Param
var cliCmd *goNixArgParser.Command

func init() {
	cliCmd = goNixArgParser.NewSimpleCommand(os.Args[0], "Simple command line based HTTP file server to share local file system")
	options := cliCmd.Options()

	// define option
	var err error
	err = options.AddFlagsValue("root", []string{"-r", "--root"}, "GHFS_ROOT", ".", "root directory of server")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("fallbackproxies", []string{"-x", "--fallback-proxy"}, "", nil, "reverse proxy to target if local resource not found, <sep><url><sep><target>, e.g. :/doc:http://remote/doc")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("alwaysproxies", []string{"-X", "--always-proxy"}, "", nil, "reverse proxy to target shadows local resource, <sep><url><sep><target>, e.g. :/doc:http://remote/doc")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("ignorebadproxycert", "--ignore-proxy-target-bad-cert", "", "ignore proxy target bad certificate")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("aliases", []string{"-a", "--alias"}, "", nil, "set alias path, <sep><url><sep><path>, e.g. :/doc:/usr/share/doc")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlags("globalupload", []string{"-U", "--global-upload"}, "", "allow upload files for all url paths")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("uploadurls", []string{"-u", "--upload"}, "", nil, "url path that allow upload files")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("uploaddirs", []string{"-p", "--upload-dir"}, "", nil, "file system path that allow upload files")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlags("globalarchive", []string{"-A", "--global-archive"}, "GHFS_GLOBAL_ARCHIVE", "enable download archive for all directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("archiveurls", "--archive", "", nil, "url path that enable download as archive for specific directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("archivedirs", "--archive-dir", "", nil, "file system path that enable download as archive for specific directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("globalcors", "--global-cors", "GHFS_GLOBAL_CORS", "enable CORS headers for all directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("corsurls", "--cors", "", nil, "url path that enable CORS headers")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("corsdirs", "--cors-dir", "", nil, "file system path that enable CORS headers")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlag("globalauth", "--global-auth", "GHFS_GLOBAL_AUTH", "require Basic Auth for all directories")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("authurls", "--auth", "", nil, "url path that require Basic Auth")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("authdirs", "--auth-dir", "", nil, "file system path that require Basic Auth")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("users", "--user", "", nil, "user info: <username>:<password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("usersbase64", "--user-base64", "", nil, "user info: <username>:<base64-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("usersmd5", "--user-md5", "", nil, "user info: <username>:<md5-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("userssha1", "--user-sha1", "", nil, "user info: <username>:<sha1-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("userssha256", "--user-sha256", "", nil, "user info: <username>:<sha256-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("userssha512", "--user-sha512", "", nil, "user info: <username>:<sha512-password>")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValue("key", []string{"-k", "--key"}, "GHFS_KEY", "", "TLS certificate key path")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValue("cert", []string{"-c", "--cert"}, "GHFS_CERT", "", "TLS certificate path")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("listen", []string{"-l", "--listen"}, "GHFS_LISTEN", nil, "address and port to listen")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("listenplain", "--listen-plain", "GHFS_LISTEN_PLAIN", nil, "address and port to listen, force plain http protocol")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("listentls", "--listen-tls", "GHFS_LISTEN_TLS", nil, "address and port to listen, force https protocol")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValues("hostnames", "--hostname", "", nil, "hostname for the virtual host")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValue("template", []string{"-t", "--template"}, "GHFS_TEMPLATE", "", "custom template file for page")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("shows", []string{"-S", "--show"}, "GHFS_SHOW", nil, "show directories or files match wildcard")
	serverErrHandler.CheckFatal(err)
	err = options.AddFlagsValues("showdirs", []string{"-SD", "--show-dir"}, "GHFS_SHOW_DIR", nil, "show directories match wildcard")
	serverErrHandler.CheckFatal(err)
	err = options.AddFlagsValues("showfiles", []string{"-SF", "--show-file"}, "GHFS_SHOW_FILE", nil, "show files match wildcard")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValues("hides", []string{"-H", "--hide"}, "GHFS_HIDE", nil, "hide directories or files match wildcard")
	serverErrHandler.CheckFatal(err)
	err = options.AddFlagsValues("hidedirs", []string{"-HD", "--hide-dir"}, "GHFS_HIDE_DIR", nil, "hide directories match wildcard")
	serverErrHandler.CheckFatal(err)
	err = options.AddFlagsValues("hidefiles", []string{"-HF", "--hide-file"}, "GHFS_HIDE_FILE", nil, "hide files match wildcard")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValue("accesslog", []string{"-L", "--access-log"}, "GHFS_ACCESS_LOG", "", "access log file, use \"-\" for stdout")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagsValue("errorlog", []string{"-E", "--error-log"}, "GHFS_ERROR_LOG", "-", "error log file, use \"-\" for stderr")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlagValue("config", "--config", "", "", "print this help")
	serverErrHandler.CheckFatal(err)

	err = options.AddFlags("help", []string{"-h", "--help"}, "", "print this help")
	serverErrHandler.CheckFatal(err)
}

func doParseCli() []*Param {
	params := []*Param{}

	args := os.Args

	// parse option
	results := cliCmd.ParseGroups(args, nil)
	configs := []string{}
	groupSeps := cliCmd.Options().GroupSeps()[0]
	foundConfig := false
	for i, length := 0, len(results); i < length; i++ {
		result := results[i]

		// help
		if result.HasFlagKey("help") {
			cliCmd.PrintHelp()
			os.Exit(0)
		}

		configs = append(configs, groupSeps)

		// config file
		config, _ := result.GetString("config")
		if len(config) == 0 {
			continue
		}

		configStr, err := ioutil.ReadFile(config)
		if serverErrHandler.CheckError(err) || len(configStr) == 0 {
			continue
		}

		configArgs := strings.Fields(string(configStr))
		if len(configArgs) == 0 {
			continue
		}

		foundConfig = true
		configs = append(configs, configArgs...)
	}

	if foundConfig {
		configs = configs[1:]
		results = cliCmd.ParseGroups(args, configs)
	}

	for _, result := range results {
		param := &Param{}

		// normalize option
		param.Root, _ = result.GetString("root")
		param.IgnoreProxyTargetBadCert = result.HasKey("ignorebadproxycert")
		param.GlobalUpload = result.HasKey("globalupload")
		param.GlobalArchive = result.HasKey("globalarchive")
		param.GlobalCors = result.HasKey("globalcors")
		param.GlobalAuth = result.HasKey("globalauth")
		param.Key, _ = result.GetString("key")
		param.Cert, _ = result.GetString("cert")
		param.Hostnames, _ = result.GetStrings("hostnames")
		param.Template, _ = result.GetString("template")
		param.AccessLog, _ = result.GetString("accesslog")
		param.ErrorLog, _ = result.GetString("errorlog")

		// normalize listen
		listen, _ := result.GetStrings("listen")
		param.Listen = append(param.Listen, listen...)

		listenRests := result.GetRests()
		param.Listen = append(param.Listen, listenRests...)

		param.ListenPlain, _ = result.GetStrings("listenplain")

		param.ListenTLS, _ = result.GetStrings("listentls")

		// normalize proxies
		arrFallbackProxies, _ := result.GetStrings("fallbackproxies")
		param.FallbackProxies = normalizeProxyMaps(arrFallbackProxies)
		arrAlwaysProxies, _ := result.GetStrings("alwaysproxies")
		param.AlwaysProxies = normalizeProxyMaps(arrAlwaysProxies)

		// normalize aliases
		arrAlias, _ := result.GetStrings("aliases")
		param.Aliases = normalizePathMaps(arrAlias)

		// normalize upload urls
		arrUploadUrls, _ := result.GetStrings("uploadurls")
		param.UploadUrls = normalizeUrlPaths(arrUploadUrls)

		// normalize upload dirs
		arrUploadDirs, _ := result.GetStrings("uploaddirs")
		param.UploadDirs = normalizeFsPaths(arrUploadDirs)

		// normalize archive urls
		arrArchiveUrls, _ := result.GetStrings("archiveurls")
		param.ArchiveUrls = normalizeUrlPaths(arrArchiveUrls)

		// normalize archive dirs
		arrArchiveDirs, _ := result.GetStrings("archivedirs")
		param.ArchiveDirs = normalizeFsPaths(arrArchiveDirs)

		// normalize cors urls
		arrCorsUrls, _ := result.GetStrings("corsurls")
		param.CorsUrls = normalizeUrlPaths(arrCorsUrls)

		// normalize cors dirs
		arrCorsDirs, _ := result.GetStrings("corsdirs")
		param.CorsDirs = normalizeFsPaths(arrCorsDirs)

		// normalize auth urls
		arrAuthUrls, _ := result.GetStrings("authurls")
		param.AuthUrls = normalizeUrlPaths(arrAuthUrls)

		// normalize auth dirs
		arrAuthDirs, _ := result.GetStrings("authdirs")
		param.AuthDirs = normalizeFsPaths(arrAuthDirs)

		// normalize users
		arrUsersPlain, _ := result.GetStrings("users")
		param.UsersPlain = getUsers(arrUsersPlain)
		arrUsersBase64, _ := result.GetStrings("usersbase64")
		param.UsersBase64 = getUsers(arrUsersBase64)
		arrUsersMd5, _ := result.GetStrings("usersmd5")
		param.UsersMd5 = getUsers(arrUsersMd5)
		arrUsersSha1, _ := result.GetStrings("userssha1")
		param.UsersSha1 = getUsers(arrUsersSha1)
		arrUsersSha256, _ := result.GetStrings("userssha256")
		param.UsersSha256 = getUsers(arrUsersSha256)
		arrUsersSha512, _ := result.GetStrings("userssha512")
		param.UsersSha512 = getUsers(arrUsersSha512)

		dupUsers := getDupUserNames(
			param.UsersPlain,
			param.UsersBase64,
			param.UsersMd5,
			param.UsersSha1,
			param.UsersSha256,
			param.UsersSha512,
		)
		if len(dupUsers) > 0 {
			serverErrHandler.CheckFatal(fmt.Errorf("duplicated usernames: %q", dupUsers))
		}

		// shows
		shows, err := getWildcardRegexp(result.GetStrings("shows"))
		serverErrHandler.CheckFatal(err)
		param.Shows = shows

		showDirs, err := getWildcardRegexp(result.GetStrings("showdirs"))
		serverErrHandler.CheckFatal(err)
		param.ShowDirs = showDirs

		showFiles, err := getWildcardRegexp(result.GetStrings("showfiles"))
		serverErrHandler.CheckFatal(err)
		param.ShowFiles = showFiles

		// hides
		hides, err := getWildcardRegexp(result.GetStrings("hides"))
		serverErrHandler.CheckFatal(err)
		param.Hides = hides

		hideDirs, err := getWildcardRegexp(result.GetStrings("hidedirs"))
		serverErrHandler.CheckFatal(err)
		param.HideDirs = hideDirs

		hideFiles, err := getWildcardRegexp(result.GetStrings("hidefiles"))
		serverErrHandler.CheckFatal(err)
		param.HideFiles = hideFiles

		params = append(params, param)
	}

	return params
}

func ParseCli() []*Param {
	if cliParams == nil {
		cliParams = doParseCli()
	}

	return cliParams
}

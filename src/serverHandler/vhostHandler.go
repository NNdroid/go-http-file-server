package serverHandler

import (
	"mjpclab.dev/ghfs/src/param"
	"mjpclab.dev/ghfs/src/serverError"
	"mjpclab.dev/ghfs/src/serverLog"
	"mjpclab.dev/ghfs/src/tpl/theme"
	"mjpclab.dev/ghfs/src/user"
	"net/http"
	"regexp"
)

type vhostContext struct {
	users  *user.List
	theme  theme.Theme
	logger *serverLog.Logger

	shows     *regexp.Regexp
	showDirs  *regexp.Regexp
	showFiles *regexp.Regexp
	hides     *regexp.Regexp
	hideDirs  *regexp.Regexp
	hideFiles *regexp.Regexp

	restrictAccess     bool
	restrictAccessUrls []pathStrings
	restrictAccessDirs []pathStrings

	headersUrls []pathHeaders
	headersDirs []pathHeaders

	vary string
}

func NewVhostHandler(
	p *param.Param,
	logger *serverLog.Logger,
	theme theme.Theme,
) (handler http.Handler, errs []error) {
	// users
	users := user.NewList(p.UserMatchCase)
	for _, u := range p.UsersPlain {
		errs = serverError.AppendError(errs, users.AddPlain(u[0], u[1]))
	}
	for _, u := range p.UsersBase64 {
		errs = serverError.AppendError(errs, users.AddBase64(u[0], u[1]))
	}
	for _, u := range p.UsersMd5 {
		errs = serverError.AppendError(errs, users.AddMd5(u[0], u[1]))
	}
	for _, u := range p.UsersSha1 {
		errs = serverError.AppendError(errs, users.AddSha1(u[0], u[1]))
	}
	for _, u := range p.UsersSha256 {
		errs = serverError.AppendError(errs, users.AddSha256(u[0], u[1]))
	}
	for _, u := range p.UsersSha512 {
		errs = serverError.AppendError(errs, users.AddSha512(u[0], u[1]))
	}

	// show/hide
	shows, err := wildcardToRegexp(p.Shows)
	errs = serverError.AppendError(errs, err)
	showDirs, err := wildcardToRegexp(p.ShowDirs)
	errs = serverError.AppendError(errs, err)
	showFiles, err := wildcardToRegexp(p.ShowFiles)
	errs = serverError.AppendError(errs, err)
	hides, err := wildcardToRegexp(p.Hides)
	errs = serverError.AppendError(errs, err)
	hideDirs, err := wildcardToRegexp(p.HideDirs)
	errs = serverError.AppendError(errs, err)
	hideFiles, err := wildcardToRegexp(p.HideFiles)
	errs = serverError.AppendError(errs, err)

	if len(errs) > 0 {
		return nil, errs
	}

	// restrict access
	restrictAccessUrls := newRestrictAccesses(p.RestrictAccessUrls)
	restrictAccessDirs := newRestrictAccesses(p.RestrictAccessDirs)
	restrictAccess := hasRestrictAccess(p.GlobalRestrictAccess, restrictAccessUrls, restrictAccessDirs)

	// `Vary` header
	vary := "accept-encoding"
	if restrictAccess {
		vary += ", referer, origin"
	}

	// alias param
	vhostCtx := &vhostContext{
		users:  users,
		theme:  theme,
		logger: logger,

		shows:     shows,
		showDirs:  showDirs,
		showFiles: showFiles,
		hides:     hides,
		hideDirs:  hideDirs,
		hideFiles: hideFiles,

		restrictAccess:     restrictAccess,
		restrictAccessUrls: restrictAccessUrls,
		restrictAccessDirs: restrictAccessDirs,

		headersUrls: newPathHeaders(p.HeadersUrls),
		headersDirs: newPathHeaders(p.HeadersDirs),

		vary: vary,
	}

	handler = newMultiplexHandler(p, vhostCtx)
	handler = newPreprocessHandler(logger, p.PreMiddlewares, handler)
	handler = newPathTransformHandler(p.PrefixUrls, handler)
	return
}

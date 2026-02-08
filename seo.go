package seo

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	gtmScriptTemplate = `<!-- Google Tag Manager -->
<script>(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':
new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0],
j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src=
'https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);
})(window,document,'script','dataLayer','%s');</script>
<!-- End Google Tag Manager -->`

	gtmNoscriptTemplate = `<!-- Google Tag Manager (noscript) -->
<noscript><iframe src="https://www.googletagmanager.com/ns.html?id=%s"
height="0" width="0" style="display:none;visibility:hidden"></iframe></noscript>
<!-- End Google Tag Manager (noscript) -->`
)

var (
	localePrefixRe = regexp.MustCompile(`^/(mobile/)?(ru|en)(/|$)`)
)

type Config struct {
	SitemapPath   string   `json:"sitemapPath,omitempty"`
	RobotsPath    string   `json:"robotsPath,omitempty"`
	Ignore        []string `json:"ignore,omitempty"`
	GTMID         string   `json:"gtmID,omitempty"`
	DefaultLang   string   `json:"defaultLang,omitempty"`   // for x-default, e.g. "en"
	SupportedLangs []string `json:"supportedLangs,omitempty"` // e.g. ["ru", "en"]
}

func CreateConfig() *Config {
	return &Config{
		SitemapPath: "/sitemap.xml",
		RobotsPath:  "/robots.txt",
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

func (w *statusWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
	}
	w.ResponseWriter.WriteHeader(code)
}

type modifyingWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (w *modifyingWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.body.Write(b)
}

func (w *modifyingWriter) WriteHeader(code int) {
	w.status = code
}

type sitemapGenerator struct {
	next          http.Handler
	name          string
	sitemapPath   string
	robotsPath    string
	ignores       []*regexp.Regexp
	paths         map[string]struct{}
	gtmID         string
	defaultLang   string
	supportedLangs []string
	mu            sync.Mutex
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.SitemapPath == "" {
		config.SitemapPath = "/sitemap.xml"
	}
	if config.RobotsPath == "" {
		config.RobotsPath = "/robots.txt"
	}

	var ignores []*regexp.Regexp
	for _, pattern := range config.Ignore {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid ignore regex %s: %v", pattern, err)
		}
		ignores = append(ignores, re)
	}

	defaultPatterns := []string{
		// Config, env, backup
		`(?i)\.env`,
		`(?i)\.bak`,
		`(?i)\.old`,
		`(?i)\.example`,
		`(?i)\.exmaple`,
		`(?i)\.sample`,
		`(?i)\.tmpl`,
		`(?i)\.tpl`,
		`(?i)\.dist`,
		`(?i)\.~`,
		`(?i)config`,
		`(?i)wp-`,
		`(?i)sitemap`,
		`(?i)undefined`,
		`^/_next/`,
		// Защита от сканирования: VCS, служебные каталоги
		`(?i)/\.git/`,
		`(?i)/\.git$`,
		`(?i)/\.svn/`,
		`(?i)/\.hg/`,
		`(?i)/\.bzr/`,
		`(?i)/\.env\b`,
		`(?i)/\.htaccess`,
		`(?i)/\.htpasswd`,
		`(?i)/\.well-known/`,
		`(?i)/wp-admin`,
		`(?i)/wp-login`,
		`(?i)/wp-content/`,
		`(?i)/wp-includes/`,
		`(?i)/administrator`,
		`(?i)/phpmyadmin`,
		`(?i)/phpinfo`,
		`(?i)/server-status`,
		`(?i)/server-info`,
		`(?i)/actuator`,
		`(?i)/jmx`,
		`(?i)/jolokia`,
		`(?i)/heapdump`,
		`(?i)/threaddump`,
		`(?i)/trace\.axd`,
		`(?i)/elmah\.axd`,
		`(?i)/node_modules/`,
		`(?i)/vendor/`,
		`(?i)/\.idea/`,
		`(?i)/\.vscode/`,
		`(?i)/cgi-bin/`,
		`(?i)/\.aws/`,
		`(?i)/\.ssh/`,
		`(?i)/\.docker/`,
		`(?i)/debug`,
		`(?i)/\.DS_Store`,
		`(?i)/crossdomain\.xml`,
		`(?i)/web\.config`,
		`(?i)/Dockerfile`,
		`(?i)/docker-compose`,
		`(?i)/Vagrantfile`,
		`(?i)/composer\.(json|lock)`,
		`(?i)/package(-lock)?\.json`,
		`(?i)/Gemfile`,
		`(?i)/pom\.xml`,
		`(?i)/build\.(gradle|xml)`,
		`(?i)/\.terraform/`,
		`(?i)/\.vagrant/`,
		`(?i)/\.travis\.yml`,
		`(?i)/\.github/`,
		`(?i)/\.gitlab`,
		`(?i)/Jenkinsfile`,
		`(?i)/Procfile`,
		`(?i)/\.bash_history`,
		`(?i)/\.mysql_history`,
		`(?i)/\.pgpass`,
		`(?i)/\.netrc`,
		`(?i)/\.my\.cnf`,
		`(?i)/api/v[0-9]+/`,
		`(?i)/graphql`,
		`(?i)/swagger`,
		`(?i)/openapi`,
		`(?i)/api-docs`,
		`(?i)/health`,
		`(?i)/metrics`,
		`(?i)/status`,
		`(?i)/__pycache__/`,
		`(?i)/\.pytest_cache/`,
		`(?i)/\.mypy_cache/`,
		`(?i)/\.cache/`,
		`(?i)/\.next/`,
		`(?i)/\.nuxt/`,
		`(?i)/\.output/`,
		`(?i)/target/`,
		`(?i)/dist/`,
		`(?i)/build/`,
		// Расширения файлов
		`(?i)\.(jpg|jpeg|jpe|png|gif|webp|svg|svgz|bmp|tif|tiff|ico|icns|heic|heif|avif|raw|cr2|nef|arw|dng|psd|ai|eps|orf|rw2|sr2|x3f)$`,
		`(?i)\.(html|htm|shtml|xhtml|css|scss|sass|less|js|mjs|cjs|ts|tsx|jsx|map|json|xml|xaml|yaml|yml|csv|tsv)$`,
		`(?i)\.(pdf|doc|docx|xls|xlsx|ppt|pptx|odt|ods|odp|docm|xlsm|pptm|rtf|tex|bib|pages|numbers|keynote)$`,
		`(?i)\.(zip|rar|7z|tar|gz|gzip|bz2|bzip2|xz|lz|lzma|lzo|z|sz|br|zst|lz4)$`,
		`(?i)\.(mp3|mp4|m4a|wav|wma|ogg|oga|opus|flac|aac|webm|mkv|avi|mov|wmv|flv|mpg|mpeg|m4v|3gp|3g2|ts|mts|m2ts|aiff|au|mid|midi|ra|rm|rmvb)$`,
		`(?i)\.(ttf|otf|woff|woff2|eot)$`,
		`(?i)\.(csv|db|sqlite|sqlite3|dat|dbf|accdb|mdb)$`,
		`(?i)\.(env|ini|conf|config|cfg|properties|log|tmp|temp)$`,
		`(?i)\.(php|php3|php4|php5|phtml|asp|aspx|aspx|jsp|jspx|cgi|py|pyc|pyo|rb|jar|war|class)$`,
		`(?i)\.(exe|dll|so|dylib|a|bin|sh|bash|bat|cmd|ps1|msi|deb|rpm|apk|app|dmg|pkg|run)$`,
		`(?i)\.(swf|wasm|data|sql|dump|bak|backup|old|orig|save|swp|tmp|temp|lock)$`,
	}
	for _, pattern := range defaultPatterns {
		re := regexp.MustCompile(pattern)
		ignores = append(ignores, re)
	}

	defaultLang := "en"
	if config.DefaultLang != "" {
		defaultLang = config.DefaultLang
	}
	supportedLangs := []string{"ru", "en"}
	if len(config.SupportedLangs) > 0 {
		supportedLangs = config.SupportedLangs
	}

	sg := &sitemapGenerator{
		next:           next,
		name:           name,
		sitemapPath:    config.SitemapPath,
		robotsPath:     config.RobotsPath,
		ignores:        ignores,
		paths:          make(map[string]struct{}),
		gtmID:          config.GTMID,
		defaultLang:    defaultLang,
		supportedLangs: supportedLangs,
	}

	return sg, nil
}

func (sg *sitemapGenerator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == sg.sitemapPath {
		xmlContent := sg.buildSitemapXML(req)
		if xmlContent == nil {
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/xml")
		rw.WriteHeader(http.StatusOK)
		rw.Write(xmlContent)
		return
	}

	if path == sg.robotsPath {
		robotsContent := sg.buildRobotsTxt(req)
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(robotsContent))
		return
	}

	ignored := false
	for _, re := range sg.ignores {
		if re.MatchString(path) {
			ignored = true
			break
		}
	}

	scheme := req.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = req.URL.Scheme
	}
	host := req.Host
	fullURL := scheme + "://" + host + strings.TrimSuffix(path, "/")

	mw := &modifyingWriter{
		ResponseWriter: rw,
		body:           bytes.NewBuffer([]byte{}),
	}
	sg.next.ServeHTTP(mw, req)

	if mw.status == 0 {
		mw.status = http.StatusOK
	}

	contentType := rw.Header().Get("Content-Type")
	contentEncoding := rw.Header().Get("Content-Encoding")

	var bodyBytes []byte = mw.body.Bytes()
	var isGzipped bool = strings.EqualFold(contentEncoding, "gzip")

	if isGzipped {
		reader, err := gzip.NewReader(bytes.NewReader(bodyBytes))
		if err != nil {
			rw.WriteHeader(mw.status)
			rw.Write(bodyBytes)
			return
		}
		defer reader.Close()
		decompressed, err := io.ReadAll(reader)
		if err != nil {
			rw.WriteHeader(mw.status)
			rw.Write(bodyBytes)
			return
		}
		bodyBytes = decompressed
	}

	bodyStr := string(bodyBytes)

	if strings.HasPrefix(strings.ToLower(contentType), "text/html") && mw.status == http.StatusOK {
		modified := bodyStr

		// Inject canonical and alternate links for locale pages (/ru/, /en/, /mobile/ru/, /mobile/en/)
		if seoLinks := sg.buildSEOLinks(req); seoLinks != "" {
			modified = strings.Replace(modified, "</head>", seoLinks+"</head>", 1)
		}

		// Inject GTM
		if sg.gtmID != "" {
			gtmScript := fmt.Sprintf(gtmScriptTemplate, sg.gtmID)
			gtmNoscript := fmt.Sprintf(gtmNoscriptTemplate, sg.gtmID)
			modified = strings.Replace(modified, "</head>", gtmScript+"</head>", 1)
			re := regexp.MustCompile(`(?i)<body\b[^>]*>`)
			match := re.FindStringIndex(modified)
			if match != nil {
				insertPos := match[1]
				modified = modified[:insertPos] + gtmNoscript + modified[insertPos:]
			}
		}

		bodyBytes = []byte(modified)

		if isGzipped {
			var gzippedBuf bytes.Buffer
			writer := gzip.NewWriter(&gzippedBuf)
			_, err := writer.Write(bodyBytes)
			if err != nil {
				rw.WriteHeader(mw.status)
				rw.Write(mw.body.Bytes())
				return
			}
			writer.Close()
			bodyBytes = gzippedBuf.Bytes()
		}

		rw.Header().Del("Content-Length")
		rw.WriteHeader(mw.status)
		rw.Write(bodyBytes)
	} else {
		rw.WriteHeader(mw.status)
		rw.Write(mw.body.Bytes())
	}

	if !ignored && mw.status == http.StatusOK {
		sg.mu.Lock()
		sg.paths[fullURL] = struct{}{}
		sg.mu.Unlock()
	}
}

type pathInfo struct {
	loc string
}

func (sg *sitemapGenerator) buildSitemapXML(req *http.Request) []byte {
	sg.mu.Lock()
	infos := make([]pathInfo, 0, len(sg.paths))
	for p := range sg.paths {
		infos = append(infos, pathInfo{loc: p})
	}
	sg.mu.Unlock()

	var filteredInfos []pathInfo
	var base string
	if req != nil {
		scheme := req.Header.Get("X-Forwarded-Proto")
		if scheme == "" {
			scheme = req.URL.Scheme
		}
		base = scheme + "://" + req.Host
		hasRoot := false
		for _, info := range infos {
			if strings.HasPrefix(info.loc, base+"/") || info.loc == base {
				if info.loc == base {
					hasRoot = true
				}
				filteredInfos = append(filteredInfos, info)
			}
		}
		if !hasRoot {
			filteredInfos = append(filteredInfos, pathInfo{loc: base})
		}
	} else {
		filteredInfos = infos
	}

	sort.Slice(filteredInfos, func(i, j int) bool {
		return filteredInfos[i].loc < filteredInfos[j].loc
	})

	type URL struct {
		Loc      string  `xml:"loc"`
		Lastmod  string  `xml:"lastmod"`
		Priority float64 `xml:"priority"`
	}
	type URLSet struct {
		XMLName xml.Name `xml:"urlset"`
		Xmlns   string   `xml:"xmlns,attr"`
		URLs    []URL    `xml:"url"`
	}

	urlset := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	now := time.Now().UTC()
	lastmodStr := now.Format("2006-01-02T15:04:05Z")

	for _, info := range filteredInfos {
		priority := 0.8
		if base != "" && info.loc == base {
			priority = 1.0
		}
		urlset.URLs = append(urlset.URLs, URL{Loc: info.loc, Lastmod: lastmodStr, Priority: priority})
	}

	output, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return nil
	}

	return []byte(xml.Header + string(output))
}

func (sg *sitemapGenerator) buildRobotsTxt(req *http.Request) string {
	scheme := req.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = req.URL.Scheme
	}
	host := req.Host
	sitemapURL := scheme + "://" + host + sg.sitemapPath

	return fmt.Sprintf("User-agent: *\nSitemap: %s\n", sitemapURL)
}

// buildSEOLinks returns canonical and alternate link tags for pages with /ru/ or /en/ locale
func (sg *sitemapGenerator) buildSEOLinks(req *http.Request) string {
	path := req.URL.Path
	match := localePrefixRe.FindStringSubmatch(path)
	if match == nil {
		return ""
	}

	mobilePrefix := match[1] // "mobile/" or ""
	currentLang := match[2] // "ru" or "en"
	isMobile := mobilePrefix != ""

	scheme := req.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = req.URL.Scheme
	}
	if scheme == "" {
		scheme = "https"
	}
	host := req.Host
	baseURL := scheme + "://" + host

	// Desktop path (strip /mobile/ for canonical on mobile pages)
	desktopPath := path
	if isMobile {
		desktopPath = "/" + strings.TrimPrefix(path, "/mobile/")
	}

	// Canonical: always desktop URL
	canonicalPath := strings.TrimSuffix(desktopPath, "/")
	if canonicalPath == "" {
		canonicalPath = "/"
	}
	canonicalURL := baseURL + canonicalPath

	var links []string
	links = append(links, fmt.Sprintf(`<link rel="canonical" href="%s" />`, canonicalURL))

	// Hreflang alternates for each supported language
	for _, lang := range sg.supportedLangs {
		var altPath string
		if isMobile {
			altPath = "/mobile/" + lang + strings.TrimPrefix(desktopPath, "/"+currentLang)
		} else {
			altPath = "/" + lang + strings.TrimPrefix(desktopPath, "/"+currentLang)
		}
		altPath = strings.TrimSuffix(altPath, "/")
		if altPath == "" {
			altPath = "/"
		}
		altURL := baseURL + altPath
		links = append(links, fmt.Sprintf(`<link rel="alternate" hreflang="%s" href="%s" />`, lang, altURL))
	}

	// x-default: homepage in default language (desktop only)
	if !isMobile {
		defaultHome := baseURL + "/" + sg.defaultLang + "/"
		links = append(links, fmt.Sprintf(`<link rel="alternate" hreflang="x-default" href="%s" />`, defaultHome))
	}

	// Mobile alternate with media query (desktop pages only)
	if !isMobile {
		mobileURL := baseURL + "/mobile" + desktopPath
		mobileURL = strings.TrimSuffix(mobileURL, "/")
		if mobileURL == baseURL+"/mobile" {
			mobileURL = baseURL + "/mobile/" + currentLang + "/"
		}
		links = append(links, fmt.Sprintf(`<link rel="alternate" media="only screen and (max-width: 640px)" href="%s" />`, mobileURL))
	}

	return strings.Join(links, "")
}

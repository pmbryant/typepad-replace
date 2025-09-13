package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// matchAndReplaceRegex performs match and replace against a slice of regex patterns across all html files
// baseDir: directory to scan
// patterns: slice of pairs (regex pattern, replacement string)
func matchAndReplaceRegex(baseDir string, patterns []struct {
	pattern *regexp.Regexp
	replace string
}) {
	var files []string
	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})

	for _, path := range files {
		input, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", path, err)
			continue
		}
		content := string(input)
		newContent := content
		for _, pr := range patterns {
			newContent = pr.pattern.ReplaceAllString(newContent, pr.replace)
		}
		if newContent != content {
			err = os.WriteFile(path, []byte(newContent), 0644)
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", path, err)
				continue
			}
			fmt.Printf("Updated: %s\n", path)
		}
	}
}

func main() {
	//rootDir := "/Users/peterbryant/Documents/Codebase/typepad-replace/lyg-blog/pmbryant.typepad.com/letyourselfgo/linda-darnell"
	//rootDir := "/Users/peterbryant/Documents/Codebase/typepad-replace/bb-blog/pmbryant.typepad.com/b_and_b/2012/12"
	rootDir := "/Users/peterbryant/Documents/Codebase/typepad-replace"

	// Define match and replace pairs
	matchReplace := []struct {
		matchStr   string
		replaceStr string
	}{
		{`href="\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../index.html`}, // These are files in 2021/01 subdir
		{`href="../01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../index.html`},
		{`href="2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="index.html`},
		{`href="../2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../index.html`},
		{`href="../../2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../index.html`},
		{`href="../../../2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../../index.html`},
		{`href="../../../../2021/01/\&quot;6Ldg1s4SAAAAAEvvZX2ILFkWp7KB-jjdL4v0JV2e\&quot;.html`, `href="../../../../index.html`},
		{`href="atom.xml"`, ``},
		{`href="../atom.xml"`, ``},
		{`href="../../atom.xml"`, ``},
		{`href="../../../atom.xml"`, ``},
		{`href="../../../../atom.xml"`, ``},
		{`href="../../../../../atom.xml"`, ``},
		{`href="https://pmbryant.typepad.com/letyourselfgo/atom.xml"`, ``}, // These only show up in atom.xml files so no replacement happens
		{`href="https://pmbryant.typepad.com/b_and_b/atom.xml"`, ``},       // These only show up in atom.xml files so no replacement happens
		{`<a href="https://www.typepad.com/" title="Blog">Blog</a>`, `<a  title="Blog">Blog</a>`},
		{`<a href="https://www.typepad.com/" title="TypePad">Typepad</a>`, `<a  title="TypePad">Typepad</a>`},
		{`<a href="http://www.typepad.com/">Powered by Typepad</a>`, `<a >Powered by Typepad</a>`},
		{`Blog powered by Typepad`, `<a  title="Blog">Blog</a> powered by <a  title="Typepad">Typepad</a>`},                                      // Put back links (as disabled) that I removed entirely earlier (B and B only)
		{`<input type="submit" value="Search" />`, `<input type="submit" value="Search" disabled/>`},                                             // Search submit button for LYG blog
		{`<input type="submit" name="btnG" value=" Google Search " /> `, `<input type="submit" name="btnG" value=" Google Search " disabled/> `}, // Search submit button for BandB blog
		{`<script type="text/javascript">
<!--
var extra_happy = Math.floor(1000000000 * Math.random());
document.write('<img src="https://www.typepad.com/t/stats?blog_id=124080652610737932&amp;user_id=100838&amp;page=' + escape(location.href) + '&amp;referrer=' + escape(document.referrer) + '&amp;i=' + extra_happy + '" width="1" height="1" alt="" style="position: absolute; top: 0; left: 0;" />');
// -->
</script>`, ``}, // remove stats calls for LYG blog
		{`<script type="text/javascript">
<!--
var extra_happy = Math.floor(1000000000 * Math.random());
document.write('<img src="https://www.typepad.com/t/stats?blog_id=50890&amp;user_id=100838&amp;page=' + escape(location.href) + '&amp;referrer=' + escape(document.referrer) + '&amp;i=' + extra_happy + '" width="1" height="1" alt="" style="position: absolute; top: 0; left: 0;" />');
// -->
</script>`, ``}, // remove stats calls for BandB blog
		{`<!--WEBBOT bot="Script" startspan PREVIEW="Site Meter" --><script type="text/javascript" language="JavaScript">var site="sm3BandBBlog"</script><script type="text/javascript" language="JavaScript1.2" src="https://sm3.sitemeter.com/js/counter.js?site=sm3BandBBlog"></script><noscript><a href="http://sm3.sitemeter.com/stats.asp?site=sm3BandBBlog" target="_top"><img src="http://sm3.sitemeter.com/meter.asp?site=sm3BandBBlog" alt="Site Meter" border="0"/></a></noscript><!-- Copyright (c)2005 Site Meter --><!--WEBBOT bot="Script" Endspan -->`, ``}, // Disable for BandB blog
		{`onclick="b=document.body; TYPEPAD___bookmarklet_domain='https://www.typepad.com/'; TYPEPAD___reblog_entryxid='6a00d83451ffec69e202c8d3bac21f200b'; TYPEPAD___is_reblog = 1; if (b &amp;&amp; !document.xmlVersion) { void ( z=document.createElement('script')); void(z.type='text/javascript'); void(z.src='https://static.typepad.com/.shared/js/qp/loader-combined-min.js'); void(b.appendChild(z));}else{}"`, ``},                                                                                                                                             // Disable 'Reblog' script
		{`onclick="b=document.body; TYPEPAD___bookmarklet_domain='https://www.typepad.com/'; TYPEPAD___reblog_entryxid='6a00d83451ffec69e201b8d2b4251e970c'; TYPEPAD___is_reblog = 1; if (b &amp;&amp; !document.xmlVersion) { void ( z=document.createElement('script')); void(z.type='text/javascript'); void(z.src='https://static.typepad.com/.shared/js/qp/loader-combined-min.js'); void(b.appendChild(z));}else{}"`, ``},                                                                                                                                             // And there are apparenlty over 150 other other these entry xIDs!
		{`TYPEPAD___bookmarklet_domain='https://www.typepad.com/'`, `TYPEPAD___bookmarklet_domain='https://defunct.domain.goes.here'`},                                         // Give up on prior approach and just get rid of the typepad.com call entirely.
		{`void(z.src='https://static.typepad.com/.shared/js/qp/loader-combined-min.js'`, `void(z.src='https://defunct.domain.goes.here/.shared/js/qp/loader-combined-min.js'`}, // Do same for the statis.typepad.com call in that line
		{`TPApp.app_uri = "https://www.typepad.com/";`, `TPApp.app_uri = "https://defunct.domain.goes.here/";`},
		{`<meta property="og:url" content="https://pmbryant.typepad.com/`, `<meta property="og:url" content="https://pmbryant.com/x/`},
		{`<meta property="og:image" content="https://up4.typepad.com/`, `<meta property="og:image" content="https://pmbryant.typepad.com/up4.typepad.com/`}, // Ooops ! Fixed with below line.
		{`<meta property="og:image" content="https://pmbryant.typepad.com/up4.typepad.com/`, `<meta property="og:image" content="https://pmbryant.com/up4.typepad.com/`},
		{`<link rel="openid.server" href="https://www.typepad.com/services/openid/server" />`, `<link rel="openid.server" href="https://defunct.domain.goes.here/services/openid/server" />`},
		{`<link rel="EditURI" type="application/rsd+xml" title="RSD" href="https://www.typepad.com/`, `<link rel="EditURI" type="application/rsd+xml" title="RSD" href="https://defunct.domain.goes.here/`},
		{`<link rel="meta" type="application/rdf+xml" title="FOAF" href="https://pmbryant.typepad.com/foaf.rdf" />`, `<link rel="meta" type="application/rdf+xml" title="FOAF" href="https://x.defunct.domain.goes.here/foaf.rdf" />`},
		{`href="../page/30/pzmyers@pharyngula.org.html"`, `href="../index.html"`},
		{`href="../../page/30/pzmyers@pharyngula.org.html"`, `href="../../index.html"`},
		{`href="../../../page/30/pzmyers@pharyngula.org.html"`, `href="../../../index.html"`},
		{`href="../../../../page/30/pzmyers@pharyngula.org.html"`, `href="../../../../index.html"`},
		{`href="../30/pzmyers@pharyngula.org.html"`, `href="../30/index.html"`},
		{`href="page/30/pzmyers@pharyngula.org.html"`, `href="index.html"`},
		{`href="pzmyers@pharyngula.org.html"`, `href="../..index.html"`},
		{`href="https://cache.blogads.com/261914443/feed.css"`, ``}, // blogads is defunct and irrelevant
		{`src="https://cache.blogads.com/261914443/feed.js"`, ``},
		{`src="http://pharyngula.org/images/ttbbadge.gif"`, ``},                                                        // pharyngula.org is defunct
		{`src="https://www.leftyblogs.com/cgi-bin/blogwire.cgi?feed=texas&site=pmbryant.typepad.com&tz=cst&n=60"`, ``}, // leftyblogs.com is defunct
		{`href="https://pmbryant.typepad.com/about.html"`, `href="https://pmbryant.com/x/b_and_b/about.html"`},         // Will not work locally but impatient and don't want to deal with different folder depths.
		{`href="https://profile.typepad.com/pmbryant"`, `href="https://pmbryant.com/profile.typepad.com/pmbryant.html"`},
		{`value="pmbryant.typepad.com"`, `value="pmbryant.defunct.domain.goeshere"`}, // typepad refs in search button for BandB (search button already disabled)
	}

	var files []string
	// Walk all html files under rootDir
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})

	for _, path := range files {
		input, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", path, err)
			continue
		}
		content := string(input)
		newContent := content
		for _, mr := range matchReplace {
			newContent = strings.ReplaceAll(newContent, mr.matchStr, mr.replaceStr)
		}
		if newContent != content {
			err = os.WriteFile(path, []byte(newContent), 0644)
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", path, err)
				continue
			}
			fmt.Printf("Updated: %s\n", path)
		}
	}

	patterns := []struct {
		pattern *regexp.Regexp
		replace string
	}{
		{regexp.MustCompile(`<a href="https://www\.typekey\.com/.*">Typepad</a>`), `<a title="disabled-as-defunct">Typepad</a>`},
		{regexp.MustCompile(`<a href="https://www\.typekey\.com/.*">Facebook</a>`), `<a title="disabled-as-defunct">Facebook</a>`},
		{regexp.MustCompile(`<a href="https://www\.typekey\.com/.*">Twitter</a>`), `<a title="disabled-as-defunct">Twitter</a>`},
		{regexp.MustCompile(`<a href="https://www\.typekey\.com/.*">more...</a>`), `<a title="disabled-as-defunct">more...</a>`},
		{regexp.MustCompile(`<a href="https://pmbryant\.typepad\.com/\.services/sitelogout.*">Sign Out</a>`), `<a title="disabled-as-defunct">Sign Out</a>`}, // Another typekey ref in the URI
		{regexp.MustCompile(`<script type="text/javascript">
\(function\(i,s,o,g,r,a,m\)\{i\['GoogleAnalyticsObject'\][\s\S]*ga\('Typepad.send', 'pageview'\);
</script>`), `<!-- google - analytics REMOVED -->`},
		{regexp.MustCompile(`<iframe src="https://www\.typepad\.com/services/connect/profile_module.*" width="100%" height="20" frameborder="0" scrolling="no" allowtransparency="true"></iframe>`), `<iframe  width="100%" height="20" frameborder="0" scrolling="no" allowtransparency="true"></iframe>`},
	}

	matchAndReplaceRegex(rootDir, patterns)
}

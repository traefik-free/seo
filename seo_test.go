package seo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateConfig(t *testing.T) {
	cfg := CreateConfig()
	if cfg.SitemapPath != "/sitemap.xml" {
		t.Errorf("expected SitemapPath /sitemap.xml, got %s", cfg.SitemapPath)
	}
	if cfg.RobotsPath != "/robots.txt" {
		t.Errorf("expected RobotsPath /robots.txt, got %s", cfg.RobotsPath)
	}
}

func TestNew_ValidConfig(t *testing.T) {
	cfg := &Config{
		SitemapPath: "/sitemap.xml",
		RobotsPath:  "/robots.txt",
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h, err := New(context.Background(), next, cfg, "test")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if h == nil {
		t.Fatal("handler is nil")
	}
}

func TestNew_InvalidIgnoreRegex(t *testing.T) {
	cfg := &Config{
		Ignore: []string{`[invalid`},
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	_, err := New(context.Background(), next, cfg, "test")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestServeHTTP_RobotsTxt(t *testing.T) {
	cfg := &Config{SitemapPath: "/sitemap.xml", RobotsPath: "/robots.txt"}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	h, _ := New(context.Background(), next, cfg, "test")

	req := httptest.NewRequest("GET", "/robots.txt", nil)
	req.Host = "example.com"
	req.Header.Set("X-Forwarded-Proto", "https")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "User-agent: *") {
		t.Error("robots.txt should contain User-agent: *")
	}
	if !strings.Contains(body, "Sitemap: https://example.com/sitemap.xml") {
		t.Errorf("robots.txt should contain sitemap URL, got: %s", body)
	}
}

func TestServeHTTP_SitemapXML(t *testing.T) {
	cfg := &Config{SitemapPath: "/sitemap.xml", RobotsPath: "/robots.txt"}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	h, _ := New(context.Background(), next, cfg, "test")

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<?xml") {
		t.Error("sitemap should be valid XML")
	}
	if !strings.Contains(body, "urlset") {
		t.Error("sitemap should contain urlset")
	}
}

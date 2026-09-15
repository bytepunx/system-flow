// Package config reads and writes the flai user configuration file.
//
// The file lives at ~/.flai/config.json by default (ADR-0008). It is plain
// JSON with a fixed schema; unknown keys are rejected so typos surface early.
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// EnvVar overrides the config file path when set.
const EnvVar = "FLAI_CONFIG"

// Config is the schema of ~/.flai/config.json.
type Config struct {
	Template  Template  `json:"template"`
	Dashboard Dashboard `json:"dashboard"`
	CacheDir  string    `json:"cache_dir"`
	Author    string    `json:"author"`
}

// Template is where flai fetches the project template from.
type Template struct {
	Repo string `json:"repo"`
	Ref  string `json:"ref"`
}

// Dashboard is how flai runs the flaiover image.
type Dashboard struct {
	Image string `json:"image"`
	Tag   string `json:"tag"`
	Port  int    `json:"port"`
}

// Default returns the configuration written on first run.
func Default() Config {
	return Config{
		Template:  Template{Repo: "https://github.com/bytepunx/system-flow-template", Ref: "main"},
		Dashboard: Dashboard{Image: "ghcr.io/bytepunx/flaiover", Tag: "latest", Port: 4242},
		CacheDir:  "~/.flai/cache",
		Author:    currentUser(),
	}
}

// DefaultPath is ~/.flai/config.json, or a relative fallback if the home
// directory cannot be determined.
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".flai", "config.json")
	}
	return filepath.Join(home, ".flai", "config.json")
}

// ResolvePath picks the config path: the explicit override, then FLAI_CONFIG,
// then the default.
func ResolvePath(override string) string {
	if override != "" {
		return ExpandHome(override)
	}
	if env := os.Getenv(EnvVar); env != "" {
		return ExpandHome(env)
	}
	return DefaultPath()
}

// Load reads the config at path. If the file does not exist it is created
// with defaults and created reports true.
func Load(path string) (cfg Config, created bool, err error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg = Default()
		if err := Save(path, cfg); err != nil {
			return Config{}, false, err
		}
		return cfg, true, nil
	}
	if err != nil {
		return Config{}, false, fmt.Errorf("read config: %w", err)
	}
	cfg, err = Parse(data)
	if err != nil {
		return Config{}, false, fmt.Errorf("parse %s: %w", path, err)
	}
	return cfg, false, nil
}

// Parse decodes JSON into a Config, rejecting unknown keys.
func Parse(data []byte) (Config, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Save writes cfg to path atomically, creating parent directories.
func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(path), ".config-*.json")
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write config: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("write config: %w", err)
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("write config: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// Keys lists every settable dotted key in schema order.
func Keys() []string {
	return []string{
		"template.repo", "template.ref",
		"dashboard.image", "dashboard.tag", "dashboard.port",
		"cache_dir", "author",
	}
}

// Get returns the value at a dotted key. Nested objects are returned as
// map[string]any.
func (c Config) Get(key string) (any, error) {
	m := c.toMap()
	var cur any = m
	for _, part := range strings.Split(key, ".") {
		obj, ok := cur.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unknown key %q", key)
		}
		cur, ok = obj[part]
		if !ok {
			return nil, fmt.Errorf("unknown key %q", key)
		}
	}
	return cur, nil
}

// Set assigns raw (a string from the command line) to a dotted key,
// converting it to the field's type. It returns the updated Config.
func (c Config) Set(key, raw string) (Config, error) {
	m := c.toMap()
	parts := strings.Split(key, ".")
	cur := m
	for _, part := range parts[:len(parts)-1] {
		next, ok := cur[part].(map[string]any)
		if !ok {
			return Config{}, fmt.Errorf("unknown key %q", key)
		}
		cur = next
	}
	last := parts[len(parts)-1]
	existing, ok := cur[last]
	if !ok {
		return Config{}, fmt.Errorf("unknown key %q", key)
	}
	switch existing.(type) {
	case float64: // json numbers
		n, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("%s expects an integer, got %q", key, raw)
		}
		cur[last] = n
	case bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, fmt.Errorf("%s expects true or false, got %q", key, raw)
		}
		cur[last] = b
	case string:
		cur[last] = raw
	default:
		return Config{}, fmt.Errorf("%s is not a settable value", key)
	}
	data, err := json.Marshal(m)
	if err != nil {
		return Config{}, err
	}
	return Parse(data)
}

// ExpandHome replaces a leading ~ with the user's home directory.
func ExpandHome(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(p, "~"))
		}
	}
	return p
}

func (c Config) toMap() map[string]any {
	data, _ := json.Marshal(c)
	var m map[string]any
	_ = json.Unmarshal(data, &m)
	return m
}

func currentUser() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	return os.Getenv("USER")
}

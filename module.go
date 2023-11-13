package hitcounter

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/templates"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(HitCounter{})
}

// HitCounter implements a simple early-Web hit counter.
type HitCounter struct {
	// The style of digit/counter to use.
	// Supported values are bright_green, green, odometer, or yellow.
	// Default: green.
	// (Styles and default are subject to change.)
	Style string `json:"style,omitempty"`

	// How many digits wide to make the counter. If zero/unset,
	// padding is disabled.
	PadDigits int `json:"pad_digits,omitempty"`

	// Pre-generated <img> tags with base64 data URIs for portability.
	// The index is the digit.
	imgTags [10]string

	// Counter states.
	counters   map[string]uint64
	countersMu *sync.Mutex

	// Time counters were last persisted.
	lastStore   time.Time
	lastStoreMu *sync.Mutex

	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (HitCounter) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.templates.functions.hitCounter",
		New: func() caddy.Module { return new(HitCounter) },
	}
}

func (hc *HitCounter) Provision(ctx caddy.Context) error {
	hc.counters = make(map[string]uint64)
	hc.countersMu = new(sync.Mutex)
	hc.lastStoreMu = new(sync.Mutex)
	hc.logger = ctx.Logger()

	if hc.Style == "" {
		hc.Style = "green"
	}

	// generate digit HTML; converting the embedded images to
	// base64 data URIs is more portable and self-contained for
	// small, sub-KB images like this
	for i := 0; i <= 9; i++ {
		fpath := fmt.Sprintf("digits/%s/%d.png", hc.Style, i)
		mimeType := mime.TypeByExtension(path.Ext(fpath))

		file, err := digits.ReadFile(fpath)
		if err != nil {
			return fmt.Errorf("unable to load digits: %v", err)
		}
		b64 := base64.StdEncoding.EncodeToString(file)

		hc.imgTags[i] = fmt.Sprintf(`<img src="data:%s;base64,%s">`, mimeType, b64)
	}

	if err := hc.restore(); err != nil {
		hc.logger.Error("restoring hit counters", zap.Error(err))
	}

	return nil
}

func (hc *HitCounter) CustomTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"hitCounter": func(key string) (string, error) {
			// get and increment the count
			hc.countersMu.Lock()
			count := hc.counters[key]
			count++
			hc.counters[key] = count
			hc.countersMu.Unlock()

			// store the updated count if it's been some time
			if err := hc.persist(); err != nil {
				hc.logger.Error("persisting hit counter data", zap.Error(err))
			}

			// convert the count to a string
			var countStr string
			if hc.PadDigits > 0 {
				formatString := "%0" + strconv.Itoa(hc.PadDigits) + "d"
				countStr = fmt.Sprintf(formatString, count)
			} else {
				countStr = strconv.FormatUint(count, 10)
			}

			// generate the HTML to display the count
			var sb strings.Builder
			for _, digit := range countStr {
				sb.WriteString(hc.imgTags[int(digit)-'0'])
			}

			return sb.String(), nil
		},
	}
}

// persist writes the counts to storage if a certain duration has passed
// since the last persist.
func (hc *HitCounter) persist() error {
	hc.lastStoreMu.Lock()

	if time.Since(hc.lastStore) < 30*time.Second {
		hc.lastStoreMu.Unlock()
		return nil
	}

	// we're optimistic it'll succeed; but we also want to
	// avoid hammering disk for every request if it's erroring,
	// so just assuming it succeeds also does that for us
	hc.lastStore = time.Now()
	hc.lastStoreMu.Unlock()

	file, err := os.Create(persistencePath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)

	hc.countersMu.Lock()
	defer hc.countersMu.Unlock()
	payload := persistedCounters{
		Timestamp: time.Now(),
		Counts:    hc.counters,
	}

	return enc.Encode(payload)
}

// restore loads the last-persisted counter state, if any.
func (hc *HitCounter) restore() error {
	file, err := os.Open(persistencePath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	var pc persistedCounters
	if err := dec.Decode(&pc); err != nil {
		return err
	}

	hc.countersMu.Lock()
	hc.lastStoreMu.Lock()
	hc.counters = pc.Counts
	hc.lastStore = pc.Timestamp
	hc.countersMu.Unlock()
	hc.lastStoreMu.Unlock()

	return nil
}

type persistedCounters struct {
	Timestamp time.Time
	Counts    map[string]uint64
}

var persistencePath = filepath.Join(caddy.AppDataDir(), "hitcounters.json")

//go:embed digits
var digits embed.FS

// Interface guards
var (
	_ caddy.Provisioner         = (*HitCounter)(nil)
	_ templates.CustomFunctions = (*HitCounter)(nil)
)

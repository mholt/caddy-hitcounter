Hit counter / visitor counter Caddy plugin
===========================================

Remember those hit counters that were so fun to have on sites from the 90s and early 2000s?

This plugin adds a retro hit counter component to Caddy templates.

In a modern Web infested with bots, hit counters are terribly inaccurate these days... but that's part of the charm.

**NOTE:** This is still in development, and currently requires an upstream patch in Caddy. This module is subject to change.

## Usage

In the Caddyfile:

```
templates {
	extensions {
		hitCounter {
			style <style>
			pad_digits <num>
		}
	}
}
```

(The `hitCounter` block is optional, but at least the function name must appear in the config.)

Possible styles are green, bright_green, odometer, or yellow.

If you want your hit counter to have a fixed size / number of digits, you can set pad_digits > 0.

Then in your template:

```
{{ hitCounter "key" }}
```

where `key` is a string identifier for this particular counter. If you want to share the same hit counter for the whole site, use `.Req.URL.Host`. Or if you want it to be for just the page, use `.Req.URL.Path` as the key.

## How it works

Every time the hit counter is shown, it increments the count.

This module embeds several styles of digit images. At provision-time, the images of the configured style are encoded as base64 so that they can be written out as self-contained `<img>` tags with data URIs.

Then the `hitCounter` template action returns consecutive HTML `<img>` tags to produce a hit counter.

Although traditional hit counters generated images dynamically for each page load, this module displays a static image for each digit to improve performance and portability and to simplify things. The resulting look is still classic nonetheless.

## Tips

- Make sure your hit counter has enough room. As each digit is its own image, you want to avoid wrapping. (CSS can help with this.)
- Since the output is multiple `<img>` tags and not a single image source, make sure your page's CSS is compatible with the look you want for your hit counter (i.e. see if your existing CSS affects these `<img>` tags).

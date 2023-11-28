Hit counter / visitor counter Caddy plugin
===========================================

Remember those hit counters that were so fun to have on sites from the 90s and early 2000s?

This plugin adds a retro hit counter component to Caddy templates.

In a modern Web infested with bots, hit counters are terribly inaccurate these days... but that's part of the charm.

> [!WARNING]
> This is still in development, and currently requires [an upstream patch in Caddy](https://github.com/caddyserver/caddy/pull/5939). This module is subject to change.

> [!NOTE]
> This is not an official repository of the [Caddy Web Server](https://github.com/caddyserver) organization.


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

Possible styles are:

- <img src="https://github.com/mholt/caddy-hitcounter/assets/1128849/0ece69c9-4e5a-43e9-a826-34f8d15bbda5" height="24"> `green` (default)
- <img src="https://github.com/mholt/caddy-hitcounter/assets/1128849/df9b6f07-9c8d-43ef-9235-fd57d0f13af0" height="22"> `bright_green`
- <img src="https://github.com/mholt/caddy-hitcounter/assets/1128849/31736f9b-dee3-4670-8e38-b66b5514053c" height="18"> `odometer`
- <img src="https://github.com/mholt/caddy-hitcounter/assets/1128849/aa0ee1f3-5dc6-4be4-a911-a5281618ace6" height="22"> `yellow`

If you want your hit counter to have a fixed size / number of digits, you can set pad_digits > 0.

Then in your template:

```
{{ hitCounter "key" }}
```

where `key` is a string identifier for this particular counter. If you want to share the same hit counter for the whole site, use `.Req.URL.Host`. Or if you want it to be for just the page, use `.Req.URL.Path` as the key.

If your hit counter is optional, and you want your template to work for other Caddy users that may not have this module installed, you can invoke it "softly":

```
{{ maybe "hitCounter" "key" }}
```

The `maybe` template action is like the built-in `call` function, except it gracefully no-ops if the requested function is not plugged in.

## How it works

Every time the hit counter is shown, it increments the count.

This module embeds several styles of digit images. At provision-time, the images of the configured style are encoded as base64 so that they can be written out as self-contained `<img>` tags with data URIs.

Then the `hitCounter` template action returns consecutive HTML `<img>` tags to produce a hit counter.

Although traditional hit counters generated images dynamically for each page load, this module displays a static image for each digit to improve performance and portability and to simplify things. The resulting look is still classic nonetheless.

## Tips

- Make sure your hit counter has enough room. As each digit is its own image, you want to avoid wrapping. (CSS can help with this.)
- Since the output is multiple `<img>` tags and not a single image source, make sure your page's CSS is compatible with the look you want for your hit counter (i.e. see if your existing CSS affects these `<img>` tags).


## Why

I grew up on a more whimsical, fun Web. So basically, [purely for nostalgia's sake](https://twitter.com/mholt6/status/1723538541505106343) (and the idea got a whole 10 likes on Twitter). There are almost DOZENS of us who want this!

# caddy-markdown-ex

This module adds a template function `markdown_ex` to caddy (v2) which can be
used with the `templates` directive. `markdown_ex` uses a custom markdown
renderer that, in addition to the default markdown renderer (template function
`markdown`)

* adds the html class "task-list-item" for items of a
  [TaskList](https://github.blog/2014-04-28-task-lists-in-all-markdown-documents/)
  to add an easy style option via CSS to disable the normal list-style-type and
  use the checkbox instead like this:
  
  ```css
  .task-list-item {
    list-style-type: none;
  }

  .task-list-item input {
    margin: 0 .2em .25em -1.6em;
    vertical-align: middle;
  }
  ```

  which results in lists like this: (the github renderer should support this 😅)

  * [ ] not checked
  * [x] checked
  * [ ] also not checked

  instead of looking like this with caddy's default renderer

  * * [ ] not checked
  * * [x] checked
  * * [ ] also not checked

* support for [mermaidJS syntax](https://github.com/mermaid-js/mermaid) by using
  [goldmark-mermaid](https://github.com/abhinav/goldmark-mermaid)

## Installation

```sh
xcaddy build --with github.com/ueffel/caddy-markdown-ex
```

The latest version of this module needs caddy v2.7.6 or above.

## Configuration

As of <https://github.com/caddyserver/caddy/pull/5939> to use a new template
function it needs to be configured as extension.

```caddy-d
templates {
    extensions {
        markdown_ex
    }
}
```

Mermaid support needs a javascript file in order to render the diagrams
client-side. This script is loaded by default from
<https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js>

There is the option to change the source path of this script to serve the file
from your own domain with the following in a Caddyfile:

```caddy-d
templates {
    extensions {
        markdown_ex {
            mermaid_js /mermaid.min.js
        }
    }
}
```

The file then should be available @ <https://[[your.domain]]/mermaid.min.js>.
For example by a `file_server` route:

```caddy-d
[your.domain] {
    file_server /mermaid.min.js {
        root /var/www/html
    }
}
```

## Usage

Use the template function `markdown_ex` instead of `markdown` to render your
markdown within the template. See [templates
directive](https://caddyserver.com/docs/caddyfile/directives/templates#templates)

## Note

I created this module just for myself but feel free to use it.

You can also fork the repository and configure your own markdown renderer to
your liking in `unmarshalCaddyfile`. The [goldmark
renderer](https://github.com/yuin/goldmark), which is used by this module and
caddy's default `markdown` function, has excellent extensibility.

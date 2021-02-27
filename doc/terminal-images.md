
## Map view (v3)

In the experimental map view, that is available under the view code `v3`,
weather information about a geographical region is available:

```
    $ curl v3.wttr.in/Bayern.sxl
```

![v3.wttr.in/Bayern](https://v3.wttr.in/Bayern.png)

or directly in browser:

*   https://v3.wttr.in/Bayern

The map view currently supports three formats:

* PNG (for browser and messangers);
* Sixel (terminal inline images support);
* IIP (terminal with iterm2 inline images protocol support).

## Terminal with images support


| Terminal              | Environment    | Images support | Protocol |
| --------------------- | --------- | ------------- | --------- |
| uxterm                |   X11     |   yes         |   Sixel   |
| mlterm                |   X11     |   yes         |   Sixel   |
| kitty                 |   X11     |   yes         |   Kitty   |
| wezterm               |   X11     |   yes         |   IIP     |
| aminal                |   X11     |   yes         |   Sixel   |
| Jexer                 |   X11     |   yes         |   Sixel   |
| GNOME Terminal        |   X11     |   [in-progress](https://gitlab.gnome.org/GNOME/vte/-/issues/253) |   Sixel   |
| alacritty             |   X11     |   [in-progress](https://github.com/alacritty/alacritty/issues/910) |  Sixel   |
| st                    |   X11     | [stixel](https://github.com/vizs/stixel) or [st-sixel](https://github.com/galatolofederico/st-sixel)     |   Sixel   |
| Konsole               |   X11     |   [requested](https://bugs.kde.org/show_bug.cgi?id=391781) | Sixel   |
| DomTerm               |   Web     |   yes         |   Sixel   |
| Yaft                  |   FB      |   yes         |   Sixel   |
| iTerm2                |   Mac OS X|   yes         |   IIP     |
| mintty                | Windows   |   yes         |   Sixel   |
| Windows Terminal  |   Windows     |   [in-progress](https://github.com/microsoft/terminal/issues/448) |   Sixel   |
| [RLogin](http://nanno.dip.jp/softlib/man/rlogin/) | Windows | yes         |   Sixel   |   |

Support in all VTE-based terminals: termite, terminator, etc is more or less the same as in the GNOME Terminal

## Notes

### xterm/uxterm

To start xterm/uxterm with Sixel support:

```
uxterm -ti vt340
```

### Kitty

To view images in kitty:

```
curl -s v3.wttr.in/Tabasco.png | kitty icat --align=left
```

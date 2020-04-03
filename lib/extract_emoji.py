#!/usr/bin/env python
#vim: fileencoding=utf-8

"""

At the moment, Pillow library does not support colorful emojis,
that is why emojis must be extracted to external files first,
and then they must be handled as usual graphical objects
and not as text.

The files are extracted using Imagemagick.

Usage:

    ve/bi/python lib/extract_emoji.py
"""

import subprocess

EMOJIS = [
    "✨",
    "☁️",
    "🌫",
    "🌧",
    "🌧",
    "❄️",
    "❄️",
    "🌦",
    "🌦",
    "🌧",
    "🌧",
    "🌨",
    "🌨",
    "⛅️",
    "☀️",
    "🌩",
    "⛈",
    "⛈",
    "☁️",
    "🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"
]

def extract_emojis_to_directory(dirname):
    """
    Extract emoji from an emoji font, to separate files.
    """

    emoji_font = "Noto Color Emoji"
    emoji_size = 30

    for emoji in EMOJIS:
        filename = "%s/%s.png" % (dirname, emoji)
        convert_string = [
            "convert", "-background", "black", "-size", "%sx%s" % (emoji_size, emoji_size),
            "-set", "colorspace", "sRGB",
            "pango:<span font=\"%s\" size=\"20000\">%s</span>" % (emoji_font, emoji),
            filename
        ]
        subprocess.Popen(convert_string)

if __name__ == '__main__':
    extract_emojis_to_directory("share/emoji")

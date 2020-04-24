#!/usr/bin/python
#vim: encoding=utf-8
# pylint: disable=wrong-import-position,wrong-import-order,redefined-builtin

"""
This module is used to generate png-files for wttr.in queries.
The only exported function is:

* render_ansi(png_file, text, options=None)

`render_ansi` is the main function of the module,
which does rendering of stream into a PNG-file.

The module uses PIL for graphical tasks, and pyte for rendering
of ANSI stream into terminal representation.
"""

from __future__ import print_function

import sys
import io
import os
import glob

from PIL import Image, ImageFont, ImageDraw
import pyte.screens
import emoji
import grapheme

from . import unicodedata2

sys.path.insert(0, "..")
import constants
import globals

COLS = 180
ROWS = 100
CHAR_WIDTH = 9
CHAR_HEIGHT = 18
FONT_SIZE = 15
FONT_CAT = {
    'default':      "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
    'Cyrillic':     "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
    'Greek':        "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
    'Arabic':       "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
    'Hebrew':       "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
    'Han':          "/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc",
    'Hiragana':     "/usr/share/fonts/truetype/motoya-l-cedar/MTLc3m.ttf",
    'Katakana':     "/usr/share/fonts/truetype/motoya-l-cedar/MTLc3m.ttf",
    'Hangul':       "/usr/share/fonts/truetype/lexi/LexiGulim.ttf",
    'Braille':      "/usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf",
    'Emoji':        "/usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf",
}

#
# How to find font for non-standard scripts:
#
#   $ fc-list :lang=ja
#
# GNU/Debian packages, that the fonts come from:
#
#   * fonts-dejavu-core
#   * fonts-wqy-zenhei (Han)
#   * fonts-motoya-l-cedar (Hiragana/Katakana)
#   * fonts-lexi-gulim (Hangul)
#   * fonts-symbola (Braille/Emoji)
#

def render_ansi(text, options=None):
    """Render `text` (terminal sequence) in a PNG file
    paying attention to passed command line `options`.

    Return: file content
    """

    screen = pyte.screens.Screen(COLS, ROWS)
    screen.set_mode(pyte.modes.LNM)
    stream = pyte.Stream(screen)

    text, graphemes = _fix_graphemes(text)
    stream.feed(text)

    buf = sorted(screen.buffer.items(), key=lambda x: x[0])
    buf = [[x[1] for x in sorted(line[1].items(), key=lambda x: x[0])] for line in buf]

    return _gen_term(buf, graphemes, options=options)

def _color_mapping(color):
    """Convert pyte color to PIL color

    Return: tuple of color values (R,G,B)
    """

    if color == 'default':
        return 'lightgray'
    if color in ['green', 'black', 'cyan', 'blue', 'brown']:
        return color
    try:
        return (
            int(color[0:2], 16),
            int(color[2:4], 16),
            int(color[4:6], 16))
    except (ValueError, IndexError):
        # if we do not know this color and it can not be decoded as RGB,
        # print it and return it as it is (will be displayed as black)
        # print color
        return color
    return color

def _strip_buf(buf):
    """Strips empty spaces from behind and from the right side.
    (from the right side is not yet implemented)
    """

    def empty_line(line):
        "Returns True if the line consists from spaces"
        return all(x.data == ' ' for x in line)

    def line_len(line):
        "Returns len of the line excluding spaces from the right"

        last_pos = len(line)
        while last_pos > 0 and line[last_pos-1].data == ' ':
            last_pos -= 1
        return last_pos

    number_of_lines = 0
    for line in buf[::-1]:
        if not empty_line(line):
            break
        number_of_lines += 1

    if number_of_lines:
        buf = buf[:-number_of_lines]

    max_len = max(line_len(x) for x in buf)
    buf = [line[:max_len] for line in buf]

    return buf

def _script_category(char):
    """Returns category of a Unicode character

    Possible values:
        default, Cyrillic, Greek, Han, Hiragana
    """

    if char in emoji.UNICODE_EMOJI:
        return "Emoji"

    cat = unicodedata2.script_cat(char)[0]
    if char == u'：':
        return 'Han'
    if cat in ['Latin', 'Common']:
        return 'default'
    return cat

def _load_emojilib():
    """Load known emojis from a directory, and return dictionary
    of PIL Image objects correspodent to the loaded emojis.
    Each emoji is resized to the CHAR_HEIGHT size.
    """

    emojilib = {}
    for filename in glob.glob("share/emoji/*.png"):
        character = os.path.basename(filename)[:-4]
        emojilib[character] = \
            Image.open(filename).resize((CHAR_HEIGHT, CHAR_HEIGHT))
    return emojilib

# pylint: disable=too-many-locals,too-many-branches,too-many-statements
def _gen_term(buf, graphemes, options=None):
    """Renders rendered pyte buffer `buf` and list of workaround `graphemes`
    to a PNG file, and return its content
    """

    if not options:
        options = {}

    current_grapheme = 0

    buf = _strip_buf(buf)
    cols = max(len(x) for x in buf)
    rows = len(buf)

    image = Image.new('RGB', (cols * CHAR_WIDTH, rows * CHAR_HEIGHT))

    buf = buf[-ROWS:]

    draw = ImageDraw.Draw(image)
    font = {}
    for cat in FONT_CAT:
        font[cat] = ImageFont.truetype(FONT_CAT[cat], FONT_SIZE)

    emojilib = _load_emojilib()

    x_pos = 0
    y_pos = 0
    for line in buf:
        x_pos = 0
        for char in line:
            current_color = _color_mapping(char.fg)
            if char.bg != 'default':
                draw.rectangle(
                    ((x_pos, y_pos),
                     (x_pos+CHAR_WIDTH, y_pos+CHAR_HEIGHT)),
                    fill=_color_mapping(char.bg))

            if char.data == "!":
                try:
                    data = graphemes[current_grapheme]
                except IndexError:
                    pass
                current_grapheme += 1
            else:
                data = char.data

            if data:
                cat = _script_category(data[0])
                if cat not in font:
                    globals.log("Unknown font category: %s" % cat)
                if cat == 'Emoji' and emojilib.get(data):
                    image.paste(emojilib.get(data), (x_pos, y_pos))
                else:
                    draw.text(
                        (x_pos, y_pos),
                        data,
                        font=font.get(cat, font.get('default')),
                        fill=current_color)

            x_pos += CHAR_WIDTH * constants.WEATHER_SYMBOL_WIDTH_VTE.get(data, 1)
        y_pos += CHAR_HEIGHT

    if 'transparency' in options:
        transparency = options.get('transparency', '255')
        try:
            transparency = int(transparency)
        except ValueError:
            transparency = 255

        if transparency < 0:
            transparency = 0

        if transparency > 255:
            transparency = 255

        image = image.convert("RGBA")
        datas = image.getdata()

        new_data = []
        for item in datas:
            new_item = tuple(list(item[:3]) + [transparency])
            new_data.append(new_item)

        image.putdata(new_data)

    img_bytes = io.BytesIO()
    image.save(img_bytes, format="png")
    return img_bytes.getvalue()

def _fix_graphemes(text):
    """
    Extract long graphemes sequences that can't be handled
    by pyte correctly because of the bug pyte#131.
    Graphemes are omited and replaced with placeholders,
    and returned as a list.

    Return:
        text_without_graphemes, graphemes
    """

    output = ""
    graphemes = []

    for gra in grapheme.graphemes(text):
        if len(gra) > 1:
            character = "!"
            graphemes.append(gra)
        else:
            character = gra
        output += character

    return output, graphemes

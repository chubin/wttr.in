import re

def metric_or_imperial(query, lang, us_ip=False):
    """
    """

    # what units should be used
    # metric or imperial
    # based on query and location source (imperial for US by default)
    if query.get('use_metric', False) and not query.get('use_imperial', False):
        query['use_imperial'] = False
        query['use_metric'] = True
    elif query.get('use_imperial', False) and not query.get('use_metric', False):
        query['use_imperial'] = True
        query['use_metric'] = False
    elif lang == 'us':
        # slack uses m by default, to override it speciy us.wttr.in
        query['use_imperial'] = True
        query['use_metric'] = False
    else:
        if us_ip:
            query['use_imperial'] = True
            query['use_metric'] = False
        else:
            query['use_imperial'] = False
            query['use_metric'] = True

    return query

def parse_query(args):
    result = {}

    reserved_args = ["lang"]
    #q = "&".join(x for x in args.keys() if x not in reserved_args)

    q = ""

    for key, val in args.items():
        if len(val) == 0:
            q += key
            continue
        if val == 'True':
            val = True
        if val == 'False':
            val = False
        result[key] = val

    if q is None:
        return result
    if 'A' in q:
        result['force-ansi'] = True
    if 'n' in q:
        result['narrow'] = True
    if 'm' in q:
        result['use_metric'] = True
    if 'M' in q:
        result['use_ms_for_wind'] = True
    if 'u' in q:
        result['use_imperial'] = True
    if 'I' in q:
        result['inverted_colors'] = True
    if 't' in q:
        result['transparency'] = '150'
    if 'T' in q:
        result['no-terminal'] = True
    if 'p' in q:
        result['padding'] = True

    for days in "0123":
        if days in q:
            result['days'] = days

    result['no-caption'] = False
    result['no-city'] = False
    if 'q' in q:
        result['no-caption'] = True
    if 'Q' in q:
        result['no-city'] = True
    if 'F' in q:
        result['no-follow-line'] = True

    for key, val in args.items():
        if val == 'True':
            val = True
        if val == 'False':
            val = False
        result[key] = val

    return result

def parse_wttrin_png_name(name):
    """
    Parse the PNG filename and return the result as a dictionary.
    For example:
        input = City_200x_lang=ru.png
        output = {
            "lang": "ru",
            "width": "200",
            "filetype": "png",
            "location": "City"
        }
    """

    parsed = {}
    to_be_parsed = {}

    if name.lower()[-4:] == '.png':
        parsed['filetype'] = 'png'
        name = name[:-4]

    parts = name.split('_')
    parsed['location'] = parts[0]

    for part in parts[1:]:
        if re.match('(?:[0-9]+)x', part):
            parsed['width'] = part[:-1]
        elif re.match('x(?:[0-9]+)', part):
            parsed['height'] = part[1:]
        elif re.match(part, '(?:[0-9]+)x(?:[0-9]+)'):
            parsed['width'], parsed['height'] = part.split('x', 1)
        elif '=' in part:
            arg, val = part.split('=', 1)
            to_be_parsed[arg] = val
        else:
            to_be_parsed[part] = ''

    parsed.update(parse_query(to_be_parsed))

    return parsed


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
    if 'n' in q:
        result['narrow'] = True
    if 'm' in q:
        result['use_metric'] = True
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

    for key, val in args.items():
        if val == 'True':
            val = True
        if val == 'False':
            val = False
        result[key] = val

    return result


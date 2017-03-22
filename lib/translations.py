# vim: set encoding=utf-8

MESSAGE = {
    'NOT_FOUND_MESSAGE': {
        'en': u"""
We were unable to find your location,
so we have brought you to Oymyakon,
one of the coldest permanently inhabited locales on the planet.
""",

        'cs': u"""
Nepodařilo se nám najít vaši polohu,
takže jsme vás přivedl do Ojmjakonu.
Je to jedno z nejchladnějších trvale obydlených míst na planetě.
Doufáme, že budete mít lepší počasí!
""",

        'de': u"""
Wir konnten Ihren Standort nicht finden,
so haben wir Sie nach Oimjakon gebracht,
einer der kältesten dauerhaft bewohnten Orte auf dem Planeten.
Wir hoffen, dass Sie besseres Wetter haben!
""",

        'el': u"""
Δεν μπορέσαμε να βρούμε την τοποθεσία σου,
για αυτό διαλέξαμε το Οϊμιάκον για εσένα,
μία από τις πιο κρύες μόνιμα κατοικημένες περιοχές στον πλανήτη.
Ελπίζουμε να έχεις καλύτερο καιρό!
""",

        'fi': u"""
Emme löytänyt sijaintiasi, joten toimme sinut Oimjakoniin,
yhteen maailman kylmimmistä pysyvästi asutetuista paikoista.
Toivottavasti sinulla on parempi sää!
""",

        'fr': u"""
Nous n'avons pas pu trouver votre position,
Nous vous avons donc amenés à Oïmiakon,
L'un des endroits les plus froids habités en permanence sur la planète.
Nous espérons que vous avez une meilleure météo
""",

        'ja': u"""
指定された場所が見つかりませんでした。
代わりにオイミャコンの天気予報を表示しています。
オイミャコンは地球上で最も寒い居住地の一つです。
""",

        'ro': u"""
Nu v-am putut identifica locația, prin urmare va aratam vremea din Oimiakon,
una dintre cele mai reci localități permanent locuite de pe planetă.
Sperăm că aveți vreme mai bună!        
""",

        'ru': u"""
Ваше местоположение определить не удалось,
поэтому мы покажем вам погоду в Оймяконе,
самом холодном населённом пункте на планете.
Будем надеяться, что у вас сегодня погода лучше!
""",

        'sk': u"""
Nepodarilo sa nám nájsť vašu polohu,
takže sme vás priviedli do Ojmiakonu.
Je to jedno z najchladnejších trvale obývaných miest na planéte.
Dúfame, že budete mať lepšie počasie!
""",

        'sr': u"""
Нисмо успели да пронађемо Вашу локацију,
па смо Вас довели у Ојмјакон,
једно од најхладнијих стално насељених места на планети.
Надамо се да је време код Вас боље него што је то случај овде!
""",

        'sv': u"""
Vi lyckades inte hitta er plats så vi har istället tagit er till Ojmjakon,
en av planetens kallaste platser med permanent bosättning.
Vi hoppas att vädret är bättre hos dig!
""",

        'nn': u"""
Me klarte ikkje å finna din stad, så me har i staden teke deg til Ojmjakon,
ein av dei kaldaste stadane på kloten med permanent busetnad.
Me håper vêret er betre hos deg!
""",

        'tr': u"""
Aradığınız bölge bulunamadı. O yüzden sizi dünyadaki en soğuk sürekli
yerleşim yerlerinden biri olan Oymyakon'e getirdik.
Umarız sizin olduğunuz yerde havalar daha iyidir!
""",

        'uk': u"""
Ваше місце розташування визначити не вдалося,
тому ми покажемо вам погоду в Оймяконе,
найхолоднішому населеному пункті на планеті.
Будемо сподіватися, що у вас сьогодні погода краще!
"""
    },

    'UNKNOWN_LOCATION': {
        'en': u'Unknown location',
        'cs': u'Neznámá poloha',
        'de': u'Unbekannter Ort',
        'el': u'Άνγωστη τοποθεσία',
        'fi': u'Tuntematon sijainti',
        'fr': u'Emplacement inconnu',
        'ja': u'未知の場所です',
        'ro': u'Locaţie necunoscută',
        'ru': u'Неизвестное местоположение',
        'sk': u'Neznáma poloha',
        'sr': u'Непозната локација',
        'sv': u'Okänd plats',
        'nn': u'Ukjend stad',
        'tr': u'Bölge bulunamadı',
        'ua': u'Невідоме місце',
    },

    'LOCATION': {
        'en': u'Location',
        'cs': u'Poloha',
        'de': u'Ort',
        'el': u'Τοποθεσία',
        'fi': u'Tuntematon sijainti',
        'fr': u'Emplacement',
        'ja': u'未知の場所です',
        'nn': u'Stad',
        'ro': u'Locaţie',
        'ru': u'Местоположение',
        'sk': u'Poloha',
        'sr': u'Локација',
        'sv': u'Plats',
        'tr': u'Bölge bulunamadı',
        'ua': u'Місце',
    }
}

def get_message(message_name, lang):
    if message_name not in MESSAGE:
        return ''
    message_dict = MESSAGE[message_name]
    return message_dict.get(lang, message_dict.get('en', ''))


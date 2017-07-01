# vim: set encoding=utf-8

FULL_TRANSLATION = [
    	"de",
]

PARTIAL_TRANSLATION = [
        "az", "be", "bg", "bs", "ca", "cy", "cs",
        "da", "el", "eo", "es", "et", "fi", "fr",
        "hi", "hr", "hu", "hy", "is", "it", "ja",
        "jv", "ka", "kk", "ko", "ky", "lt", "lv",
        "mk", "ml", "nl", "nn", "pt", "pl", "ro",
        "ru", "sk", "sl", "sr", "sr-lat", "sv",
        "sw", "th", "tr", "uk", "uz", "vi", "zh",
        "zu",
]


SUPPORTED_LANGS = FULL_TRANSLATION + PARTIAL_TRANSLATION

MESSAGE = {
    'NOT_FOUND_MESSAGE': {
        'en': u"""
We were unable to find your location
so we have brought you to Oymyakon,
one of the coldest permanently inhabited locales on the planet.
""",
        'bs': u"""
Nismo mogli pronaći vašu lokaciju,
tako da smo te doveli do Oymyakon,
jedan od najhladnijih stalno naseljena mjesta na planeti.
Nadamo se da ćete imati bolje vreme!       
""",
        'ca': u"""
Hem estat incapaços de trobar la seva ubicació,
és per aquest motiu que l'hem portat fins Oymyakon,
un dels llocs més freds inhabitats de manera permanent al planeta.
""",
            
        'cs': u"""
Nepodařilo se nám najít vaši polohu,
takže jsme vás přivedl do Ojmjakonu.
Je to jedno z nejchladnějších trvale obydlených míst na planetě.
Doufáme, že budete mít lepší počasí!
""",
            
        'cy': u"""
Ni darganfyddwyd eich lleoliad,
felly rydym wedi dod â chi i Oymyakon,
un o'r llefydd oerach ar y blaned lle mae pobl yn fyw!
""",
            
        'de': u"""
Wir konnten Ihren Standort nicht finden,
also haben wir Sie nach Oimjakon gebracht,
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

        'hy': u"""
Ձեր գտնվելու վայրը չհաջողվեց որոշել,
այդ պատճառով մենք ձեզ կցուցադրենք եղանակը Օյմյակոնում.
երկրագնդի ամենասառը բնակավայրում։
Հույս ունենք որ ձեր եղանակը այսօր ավելի լավն է։
""",

        'is': u"""
Við finnum ekki staðsetninguna þína og vísum þér þar með á Ojmjakon,
ein af köldustu byggðum jarðar.
Vonandi er betra veður hjá þér.
""",

        'ja': u"""
指定された場所が見つかりませんでした。
代わりにオイミャコンの天気予報を表示しています。
オイミャコンは地球上で最も寒い居住地の一つです。
""",

        'ko': u"""
지정된 장소를 찾을 수 없습니다,
대신 오이먀콘의 일기 예보를 표시합니다,
오이먀콘은 지구상에서 가장 추운 곳에 위치한 마을입니다!
""",

        'mk': u"""
Неможевме да ја пронајдеме вашата локација,
затоа ве однесовме во Ојмајкон,
еден од најладните трајно населени места на планетата.
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
""",

        'uz': u"""
Sizning joylashuvingizni aniqlay olmadik,
shuning uchun sizga sayyoramizning eng sovuq aholi punkti - Oymyakondagi ob-havo haqida ma'lumot beramiz.
Umid qilamizki, sizda bugungi ob-havo bundan yaxshiroq!
"""
   
		},


    'UNKNOWN_LOCATION': {
        'en': u'Unknown location',
        'bs': u'Nepoznatoja lokacija',
        'ca': u'Localització desconeguda',
        'cs': u'Neznámá poloha',
        'cy': u'Lleoliad anhysbys',
        'de': u'Unbekannter Ort',
        'el': u'Άνγωστη τοποθεσία',
        'fi': u'Tuntematon sijainti',
        'fr': u'Emplacement inconnu',
        'hy': u'Անհայտ գտնվելու վայր',
        'is': u'Óþekkt staðsetning',
        'ja': u'未知の場所です',
        'ko': u'알 수 없는 장소',
        'mk': u'Непозната локација',
        'ro': u'Locaţie necunoscută',
        'ru': u'Неизвестное местоположение',
        'sk': u'Neznáma poloha',
        'sl': u'Neznano lokacijo',
        'sr': u'Непозната локација',
        'sv': u'Okänd plats',
        'tr': u'Bölge bulunamadı',
        'ua': u'Невідоме місце',
        'uz': u'Аникланмаган худуд',
    },

    'LOCATION': {
        'en': u'Location',
        'bs': u'Lokacija',
        'ca': u'Localització',
        'cs': u'Poloha',
        'cy': u'Lleoliad',
        'de': u'Ort',
        'el': u'Τοποθεσία',
        'fi': u'Tuntematon sijainti',
        'fr': u'Emplacement',
        'hy': u'Դիրք',
        'is': u'Staðsetning',
        'ja': u'位置情報',
        'ko': u'위치',
        'mk': u'Локација',
        'ro': u'Locaţie',
        'ru': u'Местоположение',
        'sk': u'Poloha',
        'sl': u'Lokacijo',
        'sr': u'Локација',
        'sv': u'Plats',
        'tr': u'Bölge bulunamadı',
        'ua': u'Місце',
    },

    'CAPACITY_LIMIT_REACHED': {
        'en': u"""
Sorry, we are running out of queries to the weather service at the moment.
Here is the weather report for the default city (just to show you, how it looks like).
We will get new queries as soon as possible.
You can follow https://twitter.com/igor_chubin for the updates.
======================================================================================
""",
        'bs': u"""
Žao mi je, mi ponestaje upita i vremenska prognoza u ovom trenutku.
Ovdje je izvještaj o vremenu za default grada (samo da vam pokažem kako to izgleda).
Mi ćemo dobiti nove upite u najkraćem mogućem roku.
Možete pratiti https://twitter.com/igor_chubin za ažuriranja.
======================================================================================
""",
        'ca': u"""
Disculpi'ns, ens hem quedat sense consultes al servei meteorològic momentàniament.
Aquí li oferim l'informe del temps a la ciutat per defecte (només per mostrar, quin aspecte té).
Obtindrem noves consultes tan aviat com ens sigui possible.
Pot seguir https://twitter.com/igor_chubin per noves actualitzacions.
======================================================================================
""",
        'de': u"""
Entschuldigung, wir können momentan den Wetterdienst nicht erreichen.
Dafür zeigen wir Ihnen das Wetter an einem Beispielort, damit Sie sehen wie die Seite das Wetter anzeigt.
Wir werden versuchen das Problem so schnell wie möglich zu beheben.
Folgen Sie https://twitter.com/igor_chubin für Updates.
======================================================================================
""",
        'hy': u"""
Կներեք, այս պահին մենք գերազանցել ենք եղանակային տեսության կայանին հարցումների քանակը.
Կարող եք տեսնել տիպային եղանակը զեկուցում հիմնական քաղաքի համար (Ուղղակի որպես նմուշ):
Մենք մշտապես աշխատում ենք հարցումների քանակը բարելավելու ուղղությամբ:
Կարող եք հետևել մեզ https://twitter.com/igor_chubin թարմացումների համար.
======================================================================================
""",
        'ko': u"""
죄송합니다. 현재 날씨 정보를 가져오는 쿼리 요청이 한도에 도달했습니다.
대신 기본으로 설정된 도시에 대한 일기 예보를 보여드리겠습니다. (이는 단지 어떻게 보이는지 알려주기 위함입니다).
쿼리 요청이 가능한 한 빨리 이루어질 수 있도록 하겠습니다.
업데이트 소식을 원하신다면 https://twitter.com/igor_chubin 을 팔로우 해주세요.
======================================================================================
""",
        'mk': u"""
Извинете, ни снемуваат барања за до сервисот кој ни нуди временска прогноза во моментот.
Еве една временска прогноза за град (за да видите како изгледа).
Ќе добиеме нови барања најбрзо што можеме.
Следете го https://twitter.com/igor_chubin за известувања.
======================================================================================
""",

    },

    #'Check new Feature: \033[92mwttr.in/Moon\033[0m or \033[92mwttr.in/Moon@2016-Mar-23\033[0m to see the phase of the Moon'
    #'New feature: \033[92mwttr.in/Rome?lang=it\033[0m or \033[92mcurl -H "Accept-Language: it" wttr.in/Rome\033[0m for the localized version. Your lang instead of "it"'

    'NEW_FEATURE': {
        'en': u'New feature: multilingual location names \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) and location search \033[92mwttr.in/~Kilimanjaro\033[0m (just add ~ before)',
        'bs': u'XXXXXXXXXXXXXXXXXXXX: XXXXXXXXXXXXXXXXXXXXXXXXXXXXX\ 033[92mwttr.in/станция+Восток\033 [0m (XX UTF-8) XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
        'ca': u'Noves funcionalitats: noms d\'ubicació multilingües \033[92mwttr.in/станция+Восток\033 [0m (en UTF-8) i la ubicació de recerca \ 033 [92mwttr.in/~Kilimanjaro\033 [0m (només cal afegir ~ abans)',
        'cy': u'Nodwedd newydd: enwau lleoliad amlieithog \033[92mwttr.in/станция+Восток\033[0m (yn UTF-8) a chwilio lleoliad \033[92mwttr.in/~Kilimanjaro\033[0m (ychwanegwch ~ yn gyntaf)',
        'de': u'Neue Funktion: mehrsprachige Ortsnamen \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) und Ortssuche \033[92mwttr.in/~Kilimanjaro\033[0m (fügen Sie ein ~ vor dem Ort ein)',
        'hy': u'Փորձարկեք: տեղամասերի անունները կամայական լեզվով \033[92mwttr.in/Դիլիջան\033[0m (в UTF-8) և տեղանքի որոնում \033[92mwttr.in/~Kilimanjaro\033[0m (հարկավոր է ~ ավելացնել դիմացից)',
        'ko': u'새로운 기능: 다국어로 대응된 위치 \033[92mwttr.in/서울\033[0m (UTF-8에서) 장소 검색 \033[92mwttr.in/~Kilimanjaro\033[0m (앞에 ~를 붙이세요)',
        'mk': u'Нова функција: повеќе јазично локациски имиња \033[92mwttr.in/станция+Восток\033[0m (во UTF-8) и локациско пребарување \033[92mwttr.in/~Kilimanjaro\033[0m (just add ~ before)',
        'ru': u'Попробуйте: названия мест на любом языке \033[92mwttr.in/станция+Восток\033[0m (в UTF-8) и поиск мест \033[92mwttr.in/~Kilimanjaro\033[0m (нужно добавить ~ спереди)',
    },

    'FOLLOW_ME': {
        'en': u'Follow \033[46m\033[30m@igor_chubin\033[0m for wttr.in updates',
        'bs': u'XXXXXX \033[46m\033[30m@igor_chubin\033[0m XXXXXXXXXXXXXXXXXXX',   
        'ca': u'Seguiu \033[46m\033[30m@igor_chubin\033[0m per actualitzacions de wttr.in',
        'cy': u'Dilyn \033[46m\033[30m@igor_Chubin\033[0m am diweddariadau wttr.in',
        'de': u'Folgen Sie \033[46m\033[30mhttps://twitter.com/igor_chubin\033[0m für wttr.in Updates';
        'hy': u'Նոր ֆիչռների համար հետևեք՝ \033[46m\033[30m@igor_chubin\033[0m',
        'ko': u'wttr.in의 업데이트 소식을 원하신다면 \033[46m\033[30m@igor_chubin\033[0m 을 팔로우 해주세요',
        'mk': u'Следете \033[46m\033[30m@igor_chubin\033[0m за wttr.in новости',
        'ru': u'Все новые фичи публикуются здесь: \033[46m\033[30m@igor_chubin\033[0m',
    },
}

def get_message(message_name, lang):
    if message_name not in MESSAGE:
        return ''
    message_dict = MESSAGE[message_name]
    return message_dict.get(lang, message_dict.get('en', ''))

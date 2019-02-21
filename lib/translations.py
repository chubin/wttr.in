# vim: fileencoding=utf-8

"""
Translation of almost everything.
"""

FULL_TRANSLATION = [
    "af", "da", "de", "fr", "fa", "id", "it", "nb", "nl", "pl", "ro", "ru",
]

PARTIAL_TRANSLATION = [
    "az", "be", "bg", "bs", "ca", "cy", "cs",
    "el", "eo", "es", "et", "fi",
    "hi", "hr", "hu", "hy", "is",
    "ja", "jv", "ka", "kk", "ko", "ky", "lt",
    "lv", "mk", "ml", "nl", "nn", "pt",
    "sk", "sl", "sr", "sr-lat",
    "sv", "sw", "th", "tr", "te", "uk", "uz", "vi",
    "zh", "zu",
    "he",
]

PROXY_LANGS = [
    'af', 'az', 'be', 'bs', 'ca', 'cy', 'eo', 'fa',
    'he', 'hr', 'hy', 'id', 'is', 'it', 'ja',
    'kk', 'lv', 'mk', 'nb', 'nn', 'ro', 'sl', 'uz'
]

SUPPORTED_LANGS = FULL_TRANSLATION + PARTIAL_TRANSLATION

MESSAGE = {
    'NOT_FOUND_MESSAGE': {
        'en': u"""
We were unable to find your location
so we have brought you to Oymyakon,
one of the coldest permanently inhabited locales on the planet.
""",
        'af': u"""
Ons kon nie u ligging opspoor nie
gevolglik het ons vir u na Oymyakon geneem,
een van die koudste permanent bewoonde plekke op aarde.
""",
        'be': u"""
Ваша месцазнаходжанне вызначыць не атрымалася,
таму мы пакажам вам надвор'е ў Аймяконе,
самым халодным населеным пункце на планеце.
Будзем спадзявацца, што ў вас сёння надвор'е лепей!
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
un o'r llefydd oeraf ar y blaned ble mae pobl yn dal i fyw!
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
        'es': u"""
No hemos logrado encontrar tu ubicación,
asi que hemos decidido enseñarte el tiempo en Oymyakon,
uno de los sitios más fríos y permanentemente deshabitados del planeta.
""",
        'fa': u"""
ما نتونستیم مکان شما رو پیدا کنیم. به همین خاطر شما رو به om بردیم
، یکی از سردترین مکان های روی زمین که اصلا قابل سکونت نیست!
""",
        'fi': u"""
Emme löytänyt sijaintiasi, joten toimme sinut Oimjakoniin,
yhteen maailman kylmimmistä pysyvästi asutetuista paikoista.
Toivottavasti sinulla on parempi sää!
""",
        'fr': u"""
Nous n'avons pas pu déterminer votre position,
Nous vous avons donc amenés à Oïmiakon,
l'un des endroits les plus froids habités en permanence sur la planète.
Nous espérons qu'il fait meilleur chez vous !
""",
        'hy': u"""
Ձեր գտնվելու վայրը չհաջողվեց որոշել,
այդ պատճառով մենք ձեզ կցուցադրենք եղանակը Օյմյակոնում.
երկրագնդի ամենասառը բնակավայրում։
Հույս ունենք որ ձեր եղանակը այսօր ավելի լավն է։
""",
        'id': u"""
Kami tidak dapat menemukan lokasi anda,
jadi kami membawa anda ke Oymyakon,
salah satu tempat terdingin yang selalu dihuni di planet ini!
""",
        'is': u"""
Við finnum ekki staðsetninguna þína og vísum þér þar með á Ojmjakon,
ein af köldustu byggðum jarðar.
Vonandi er betra veður hjá þér.
""",
        'it': u"""
Non siamo riusciti a trovare la sua posizione
quindi la abbiamo portato a Oymyakon,
uno dei luoghi abitualmente abitati più freddi del pianeta.
Ci auguriamo che le condizioni dove lei si trova siano migliori!
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
        'lv': u"""
Mēs nevarējām atrast jūsu atrašanās vietu tādēļ nogādājām jūs Oimjakonā,
vienā no aukstākajām apdzīvotajām vietām uz planētas.
""",
        'mk': u"""
Неможевме да ја пронајдеме вашата локација,
затоа ве однесовме во Ојмајкон,
еден од најладните трајно населени места на планетата.
""",
        'nb': u"""
Vi kunne ikke finne din lokasjon,
så her får du Ojmjakon, et av de kaldeste bebodde stedene på planeten.
Vi håper været er bedre hos deg!
""",
        'nl': u"""
Wij konden uw locatie niet vaststellen
dus hebben we u naar Ojmjakon gebracht,
één van de koudste permanent bewoonde gebieden op deze planeet.
""",
        'pt': u"""
Não conseguimos encontrar a sua localização,
então decidimos te mostrar o tempo em Oymyakon,
um dos lugares mais frios e permanentemente desabitados do planeta.
""",
        'pl': u"""
Nie udało nam się znaleźć podanej przez Ciebie lokalizacji,
więc zabraliśmy Cię do Ojmiakonu,
jednego z najzimniejszych, stale zamieszkanych miejsc na Ziemi.
Mamy nadzieję, że u Ciebie jest cieplej!
""",
        'ro': u"""
Nu v-am putut identifica localitatea, prin urmare vă arătăm vremea din Oimiakon,
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
Aradığınız konum bulunamadı. O yüzden sizi dünyadaki en soğuk sürekli
yerleşim yerlerinden biri olan Oymyakon'e getirdik.
Umarız sizin olduğunuz yerde havalar daha iyidir!
""",
        'te': u"""
మేము మీ స్థానాన్ని కనుగొనలేకపోయాము
కనుక మనం "ఓమాయకాన్కు" తీసుకొని వచ్చాము,
భూమిపై అత్యల్ప శాశ్వతంగా నివసించే స్థానిక ప్రదేశాలలో ఒకటి.
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
""",
        'da': u"""
Vi kunne desværre ikke finde din lokation
så vi har bragt dig til Oymyakon,
En af koldeste og helt ubolige lokationer på planeten.
""",
        'et': u"""
Me ei suutnud tuvastada teie asukohta
ning seetõttu paigutasime teid Oymyakoni,
mis on üks kõige külmemaid püsivalt asustatud paiku planeedil.
""",
    },

    'UNKNOWN_LOCATION': {
        'en': u'Unknown location',
        'af': u'Onbekende ligging',
        'be': u'Невядомае месцазнаходжанне',
        'bs': u'Nepoznatoja lokacija',
        'ca': u'Localització desconeguda',
        'cs': u'Neznámá poloha',
        'cy': u'Lleoliad anhysbys',
        'de': u'Unbekannter Ort',
        'da': u'Ukendt lokation',
        'el': u'Άνγωστη τοποθεσία',
        'es': u'Ubicación desconocida',
        'et': u'Tundmatu asukoht',
        'fa': u'مکان نامعلوم',
        'fi': u'Tuntematon sijainti',
        'fr': u'Emplacement inconnu',
        'hy': u'Անհայտ գտնվելու վայր',
        'id': u'Lokasi tidak diketahui',
        'is': u'Óþekkt staðsetning',
        'it': u'Località sconosciuta',
        'ja': u'未知の場所です',
        'ko': u'알 수 없는 장소',
        'kk': u'',
        'lv': u'Nezināma atrašanās vieta',
        'mk': u'Непозната локација',
        'nb': u'Ukjent sted',
        'nl': u'Onbekende locatie',
        'pl': u'Nieznana lokalizacja',
        'pt': u'Localização desconhecida',
        'ro': u'Localitate necunoscută',
        'ru': u'Неизвестное местоположение',
        'sk': u'Neznáma poloha',
        'sl': u'Neznano lokacijo',
        'sr': u'Непозната локација',
        'sv': u'Okänd plats',
        'te': u'తెలియని ప్రదేశం',
        'tr': u'Bilinmeyen konum',
        'ua': u'Невідоме місце',
        'uz': u'Аникланмаган худуд',
    },

    'LOCATION': {
        'en': u'Location',
        'af': u'Ligging',
        'be': u'Месцазнаходжанне',
        'bs': u'Lokacija',
        'ca': u'Localització',
        'cs': u'Poloha',
        'cy': u'Lleoliad',
        'de': u'Ort',
        'da': u'Lokation',
        'el': u'Τοποθεσία',
        'es': u'Ubicación',
        'et': u'Asukoht',
        'fa': u'مکان',
        'fi': u'Tuntematon sijainti',
        'fr': u'Emplacement',
        'hy': u'Դիրք',
        'id': u'Lokasi',
        'is': u'Staðsetning',
        'it': u'Località',
        'ja': u'位置情報',
        'ko': u'위치',
        'kk': u'',
        'lv': u'Atrašanās vieta',
        'mk': u'Локација',
        'nb': u'Sted',
        'nl': u'Locatie',
        'pl': u'Lokalizacja',
        'pt': u'Localização',
        'ro': u'Localitate',
        'ru': u'Местоположение',
        'sk': u'Poloha',
        'sl': u'Lokacijo',
        'sr': u'Локација',
        'sv': u'Plats',
        'te': u'స్థానము',
        'tr': u'Konum',
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
        'af': u"""
Verskoning, ons oorskry tans die vermoë om navrae aan die weerdiens te rig.
Hier is die weerberig van 'n voorbeeld ligging (bloot om aan u te wys hoe dit lyk).
Ons sal weereens nuwe navrae kan hanteer so gou as moontlik.
U kan vir https://twitter.com/igor_chubin volg vir opdaterings.
======================================================================================
""",
        'be': u"""
Прабачце, мы выйшлі за ліміты колькасці запытаў да службы надвор'я ў дадзены момант.
Вось прагноз надвор'я для горада па змаўчанні (толькі, каб паказаць вам, як гэта выглядае).
Мы вернемся як мага хутчэй.
Вы можаце сачыць на https://twitter.com/igor_chubin за абнаўленнямі.
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
        'cy': u"""
Rydym yn brin o ymholiadau i'r gwasanaeth tywydd ar hyn o bryd.
Felly dyma'r adroddiad tywydd ar gyfer y ddinas ragosod (er mwyn arddangos sut mae'n edrych).
Byddwn gyda ymholiadau newydd yn fuan.
Gellir dilyn https://twitter.com/igor_chubin i gael newyddion pellach.
======================================================================================
""",
        'es': u"""
Lo siento, hemos alcanzado el límite de peticiones al servicio de previsión del tiempo en este momento.
A continuación, la previsión del tiempo para una ciudad estándar (solo para que puedas ver que aspecto tiene el informe).
Muy pronto volveremos a tener acceso a las peticiones.
Puedes seguir https://twitter.com/igor_chubin para estar al tanto de la situación.
======================================================================================
""",
        'fa': u"""
متأسفانه در حال حاضر ظرفیت ما برای درخواست به سرویس هواشناسی به اتمام رسیده.
اینجا می تونید گزارش هواشناسی برای شهر پیش فرض رو ببینید (فقط برای اینه که بهتون نشون بدیم چه شکلی هست)
ما تلاش میکنیم در اسرع وقت ظرفیت جدید به دست بیاریم.
برای دنبال کردن اخبار جدید میتونید https://twitter.com/igor_chubin رو فالو کنید.
======================================================================================
""",
        'fr': u"""
Désolé, nous avons épuisé les requêtes vers le service météo.
Voici un bulletin météo de l'emplacement par défaut (pour vous donner un aperçu).
Nous serons très bientôt en mesure de faire de nouvelles requêtes.
Vous pouvez suivre https://twitter.com/igor_chubin pour rester informé.
======================================================================================
""",
        'hy': u"""
Կներեք, այս պահին մենք գերազանցել ենք եղանակային տեսության կայանին հարցումների քանակը.
Կարող եք տեսնել տիպային եղանակը զեկուցում հիմնական քաղաքի համար (Ուղղակի որպես նմուշ):
Մենք մշտապես աշխատում ենք հարցումների քանակը բարելավելու ուղղությամբ:
Կարող եք հետևել մեզ https://twitter.com/igor_chubin թարմացումների համար.
======================================================================================
""",
        'id': u"""
Maaf, kami kehabian permintaan ke layanan cuaca saat ini.
Ini adalah laporan cuaca dari kota standar (hanya untuk menunjukkan kepada anda bagaimana tampilannya).
Kami akan mencoba permintaan baru lagi sesegera mungkin.
Anda dapat mengikuti https://twitter.com/igor_chubin untuk informasi terbaru.
======================================================================================
""",
        'it': u"""
Scusate, attualmente stiamo esaurendo le risorse a disposizione del servizio meteo.
Qui trovate il bollettino del tempo per la città di default (solo per mostrarvi come si presenta).
Potremo elaborare nuove richieste appena possibile.
Potete seguire https://twitter.com/igor_chubin per gli aggiornamenti.
======================================================================================
""",
        'ko': u"""
죄송합니다. 현재 날씨 정보를 가져오는 쿼리 요청이 한도에 도달했습니다.
대신 기본으로 설정된 도시에 대한 일기 예보를 보여드리겠습니다. (이는 단지 어떻게 보이는지 알려주기 위함입니다).
쿼리 요청이 가능한 한 빨리 이루어질 수 있도록 하겠습니다.
업데이트 소식을 원하신다면 https://twitter.com/igor_chubin 을 팔로우 해주세요.
======================================================================================
""",
        'lv': u"""
Atvainojiet, uz doto brīdi mēs esam mazliet noslogoti.
Šeit ir laika ziņas noklusējuma pilsētai (lai parādītu jums, kā izskatās izveidotais ziņojums).
Mēs atsāksim darbu cik ātri vien varēsim.
Jūs varat sekot https://twitter.com/igor_chubin lai redzētu visus jaunumus.
======================================================================================
""",
        'mk': u"""
Извинете, ни снемуваат барања за до сервисот кој ни нуди временска прогноза во моментот.
Еве една временска прогноза за град (за да видите како изгледа).
Ќе добиеме нови барања најбрзо што можеме.
Следете го https://twitter.com/igor_chubin за известувања
======================================================================================
""",
        'nb': u"""
Beklager, vi kan ikke nå værtjenesten for øyeblikket.
Her er værmeldingen for standardbyen så du får se hvordan tjenesten ser ut.
Vi vil forsøke å fikse problemet så snart som mulig.
Du kan følge https://twitter.com/igor_chubin for oppdateringer.
======================================================================================
""",
        'nl': u"""
Excuse, wij kunnen u op dit moment dit weerbericht niet laten zien.
Hier is het weerbericht voor de standaard stad(zodat u weet hoe het er uitziet)
Wij lossen dit probleem zo snel mogelijk op.
voor updates kunt u ons op https://twitter.com/igor_chubin volgen.
======================================================================================
""",
        'pl': u"""
Bardzo nam przykro, ale chwilowo wykorzystaliśmy limit zapytań do serwisu pogodowego.
To, co widzisz jest przykładowym raportem pogodowym dla domyślnego miasta.
Postaramy się przywrócić funkcjonalność tak szybko, jak to tylko możliwe.
Możesz śledzić https://twitter.com/igor_chubin na Twitterze, aby być na bieżąco.
======================================================================================
""",
        'pt': u"""
Desculpe-nos, estamos atingindo o limite de consultas ao serviço de previsão do tempo neste momento.
Veja a seguir a previsão do tempo para uma cidade padrão (apenas para você ver que aspecto o relatório tem).
Em breve voltaremos a ter acesso às consultas.
Você pode seguir https://twitter.com/igor_chubin para acompanhar a situação.
======================================================================================
""",
        'ro': u"""
Ne pare rău, momentan am epuizat cererile alocate de către serviciul de prognoză meteo.
Vă arătăm prognoza meteo pentru localitatea implicită (ca exemplu, să vedeți cum arată).
Vom obține alocarea de cereri noi cât de curând posibil.
Puteți urmări https://twitter.com/igor_chubin pentru actualizări.
======================================================================================
""",
        'te': u"""
క్షమించండి, ప్రస్తుతానికి మేము వాతావరణ సేవకు ప్రశ్నలను గడుపుతున్నాం.
ఇక్కడ డిఫాల్ట్ నగరం కోసం వాతావరణ నివేదిక (కేవలం మీకు చూపించడానికి, ఇది ఎలా కనిపిస్తుంది).
సాధ్యమైనంత త్వరలో కొత్త ప్రశ్నలను పొందుతారు.
నవీకరణల కోసం https://twitter.com/igor_chubin ను మీరు అనుసరించవచ్చు.
======================================================================================
""",
        'tr': u"""
Üzgünüz, an itibariyle hava durumu servisine yapabileceğimiz sorgu limitine ulaştık.
Varsayılan şehir için hava durumu bilgisini görüyorsunuz (neye benzediğini gösterebilmek için).
Mümkün olan en kısa sürede servise yeniden sorgu yapmaya başlayacağız.
Gelişmeler için https://twitter.com/igor_chubin adresini takip edebilirsiniz.
======================================================================================
""",
        'da': u"""
Beklager, men vi er ved at løbe tør for forespørgsler til vejr-servicen lige nu.
Her er vejr rapporten for standard byen (bare så du ved hvordan det kan se ud).
Vi får nye forespørsler hurtigst muligt.
Du kan følge https://twitter.com/igor_chubin for at få opdateringer.
======================================================================================
""",
        'et': u"""
Vabandage, kuid hetkel on päringud ilmateenusele piiratud.
Selle asemel kuvame hetkel näidislinna ilmaprognoosi (näitamaks, kuidas see välja näeb).
Üritame probleemi lahendada niipea kui võimalik.
Jälgige https://twitter.com/igor_chubin värskenduste jaoks.
======================================================================================
""",

    },

    # Historical messages:
    #     'Check new Feature: \033[92mwttr.in/Moon\033[0m or \033[92mwttr.in/Moon@2016-Mar-23\033[0m to see the phase of the Moon'
    #     'New feature: \033[92mwttr.in/Rome?lang=it\033[0m or \033[92mcurl -H "Accept-Language: it" wttr.in/Rome\033[0m for the localized version. Your lang instead of "it"'

    'NEW_FEATURE': {
        'en': u'New feature: multilingual location names \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) and location search \033[92mwttr.in/~Kilimanjaro\033[0m (just add ~ before)',
        'af': u'Nuwe eienskap: veeltalige name vir liggings \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) en ligging soek \033[92mwttr.in/~Kilimanjaro\033[0m (plaas net ~ vooraan)',
        'be': u'Новыя магчымасці: назвы месц на любой мове \033[92mwttr.in/станция+Восток\033[0m (в UTF-8) i пошук месц \033[92mwttr.in/~Kilimanjaro\033[0m (трэба дадаць ~ ў пачатак)',
        'bs': u'XXXXXXXXXXXXXXXXXXXX: XXXXXXXXXXXXXXXXXXXXXXXXXXXXX\033[92mwttr.in/станция+Восток\033[0m (XX UTF-8) XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
        'ca': u'Noves funcionalitats: noms d\'ubicació multilingües \033[92mwttr.in/станция+Восток\033[0m (en UTF-8) i la ubicació de recerca \033[92mwttr.in/~Kilimanjaro\033[0m (només cal afegir ~ abans)',
        'es': u'Nuevas funcionalidades: los nombres de las ubicaciones en vários idiomas \033[92mwttr.in/станция+Восток\033[0m (em UTF-8) y la búsqueda por ubicaciones \033[92mwttr.in/~Kilimanjaro\033[0m (tan solo inserte ~ en frente)',
        'fa': u'قابلیت جدید: پشتیبانی از نام چند زبانه مکانها \033[92mwttr.in/станция+Восток\033[0m (در فرمت UTF-8) و جسجتوی مکان ها \033[92mwttr.in/~Kilimanjaro\033[0m (فقط قبل از اون ~ اضافه کنید)',
        'fr': u'Nouvelles fonctionnalités: noms d\'emplacements multilingues \033[92mwttr.in/станция+Восток\033[0m (en UTF-8) et recherche d\'emplacement \033[92mwttr.in/~Kilimanjaro\033[0m (ajouter ~ devant)',
        'mk': u'Нова функција: повеќе јазично локациски имиња \033[92mwttr.in/станция+Восток\033[0m (во UTF-8) и локациско пребарување \033[92mwttr.in/~Kilimanjaro\033[0m (just add ~ before)',
        'nb': u'Ny funksjon: flerspråklige stedsnavn \033[92mwttr.in/станция+Восток\033[0m (i UTF-8) og lokasjonssøk \033[92mwttr.in/~Kilimanjaro\033[0m (bare legg til ~ foran)',
        'nl': u'Nieuwe functie: tweetalige locatie namen \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) en locatie zoeken \033[92mwttr.in/~Kilimanjaro\033[0m (zet er gewoon een ~ voor)',
        'cy': u'Nodwedd newydd: enwau lleoliadau amlieithog \033[92mwttr.in/станция+Восток\033[0m (yn UTF-8) a chwilio am leoliad \033[92mwttr.in/~Kilimanjaro\033[0m (ychwanegwch ~ yn gyntaf)',
        'de': u'Neue Funktion: mehrsprachige Ortsnamen \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) und Ortssuche \033[92mwttr.in/~Kilimanjaro\033[0m (fügen Sie ein ~ vor dem Ort ein)',
        'hy': u'Փորձարկեք: տեղամասերի անունները կամայական լեզվով \033[92mwttr.in/Դիլիջան\033[0m (в UTF-8) և տեղանքի որոնում \033[92mwttr.in/~Kilimanjaro\033[0m (հարկավոր է ~ ավելացնել դիմացից)',
        'id': u'Fitur baru: nama lokasi dalam multibahasa \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) dan pencarian lokasi \033[92mwttr.in/~Kilimanjaro\033[0m (hanya tambah tanda ~ sebelumnya)',
        'it': u'Nuove funzionalità: nomi delle località multilingue \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) e ricerca della località \033[92mwttr.in/~Kilimanjaro\033[0m (basta premettere ~)',
        'ko': u'새로운 기능: 다국어로 대응된 위치 \033[92mwttr.in/서울\033[0m (UTF-8에서) 장소 검색 \033[92mwttr.in/~Kilimanjaro\033[0m (앞에 ~를 붙이세요)',
        'kk': u'',
        'lv': u'Jaunums: Daudzvalodu atrašanās vietu nosaukumi \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) un dabas objektu meklēšana \033[92mwttr.in/~Kilimanjaro\033[0m (tikai priekšā pievieno ~)',
        'mk': u'Нова функција: повеќе јазично локациски имиња \033[92mwttr.in/станция+Восток\033[0m (во UTF-8) и локациско пребарување \033[92mwttr.in/~Kilimanjaro\033[0m (just add ~ before)',
        'pl': u'Nowa funkcjonalność: wielojęzyczne nazwy lokalizacji \033[92mwttr.in/станция+Восток\033[0m (w UTF-8) i szukanie lokalizacji \033[92mwttr.in/~Kilimanjaro\033[0m (poprzedź zapytanie ~ - znakiem tyldy)',
        'pt': u'Nova funcionalidade: nomes de localidades em várias línguas \033[92mwttr.in/станция+Восток\033[0m (em UTF-8) e procura por localidades \033[92mwttr.in/~Kilimanjaro\033[0m (é só colocar ~ antes)',
        'ro': u'Funcționalitate nouă: nume de localități multilingve \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) și căutare de localități \033[92mwttr.in/~Kilimanjaro\033[0m (adăuați ~ în față)',
        'ru': u'Попробуйте: названия мест на любом языке \033[92mwttr.in/станция+Восток\033[0m (в UTF-8) и поиск мест \033[92mwttr.in/~Kilimanjaro\033[0m (нужно добавить ~ спереди)',
        'tr': u'Yeni özellik: çok dilli konum isimleri \033[92mwttr.in/станция+Восток\033[0m (UTF-8 ile) ve konum arama \033[92mwttr.in/~Kilimanjaro\033[0m (sadece önüne ~ ekleyin)',
        'te': u'క్రొత్త లక్షణం: బహుభాషా స్థాన పేర్లు \ 033 [92mwttr.in/stancelя+Vostок\033 [0 U (UTF-8 లో) మరియు స్థానం శోధన \ 033 [92mwttr.in/~kilimanjaro\033 [0m (కేవలం ~ ముందుకి జోడించండి)',
        'da': u'Ny funktion: flersprogede lokationsnavne \033[92mwttr.in/станция+Восток\033[0m (som UTF-8) og lokations søgning \033[92mwttr.in/~Kilimanjaro\033[0m (bare tilføj ~ inden)',
        'et': u'Uus funktsioon: mitmekeelsed asukohanimed \033[92mwttr.in/станция+Восток\033[0m (UTF-8 vormingus) ja asukoha otsing \033[92mwttr.in/~Kilimanjaro\033[0m (lisa ~ enne)',
    },

    'FOLLOW_ME': {
        'en': u'Follow \033[46m\033[30m@igor_chubin\033[0m for wttr.in updates',
        'af': u'Volg \033[46m\033[30m@igor_chubin\033[0m vir wttr.in opdaterings',
        'be': u'Сачыце за \033[46m\033[30m@igor_chubin\033[0m за навінамі wttr.in',
        'bs': u'XXXXXX \033[46m\033[30m@igor_chubin\033[0m XXXXXXXXXXXXXXXXXXX',
        'ca': u'Seguiu \033[46m\033[30m@igor_chubin\033[0m per actualitzacions de wttr.in',
        'es': u'Seguir \033[46m\033[30m@igor_chubin\033[0m para recibir las novedades de wttr.in',
        'cy': u'Dilyner \033[46m\033[30m@igor_Chubin\033[0m am diweddariadau wttr.in',
        'fa': u'برای دنبال کردن خبرهای wttr.in شناسه \033[46m\033[30m@igor_chubin\033[0m رو فالو کنید.',
        'fr': u'Suivez \033[46m\033[30m@igor_Chubin\033[0m pour rester informé sur wttr.in',
        'de': u'Folgen Sie \033[46m\033[30mhttps://twitter.com/igor_chubin\033[0m für wttr.in Updates',
        'hy': u'Նոր ֆիչռների համար հետևեք՝ \033[46m\033[30m@igor_chubin\033[0m',
        'id': u'Ikuti \033[46m\033[30m@igor_chubin\033[0m untuk informasi wttr.in terbaru',
        'it': u'Seguite \033[46m\033[30m@igor_chubin\033[0m per aggiornamenti a wttr.in',
        'ko': u'wttr.in의 업데이트 소식을 원하신다면 \033[46m\033[30m@igor_chubin\033[0m 을 팔로우 해주세요',
        'kk': u'',
        'lv': u'Seko \033[46m\033[30m@igor_chubin\033[0m , lai uzzinātu wttr.in jaunumus',
        'mk': u'Следете \033[46m\033[30m@igor_chubin\033[0m за wttr.in новости',
        'nb': u'Følg \033[46m\033[30m@igor_chubin\033[0m for wttr.in oppdateringer',
        'nl': u'Volg \033[46m\033[30m@igor_chubin\033[0m voor wttr.in updates',
        'pl': u'Śledź \033[46m\033[30m@igor_chubin\033[0m aby być na bieżąco z nowościami dotyczącymi wttr.in',
        'pt': u'Seguir \033[46m\033[30m@igor_chubin\033[0m para as novidades de wttr.in',
        'ro': u'Urmăriți \033[46m\033[30m@igor_chubin\033[0m pentru actualizări despre wttr.in',
        'ru': u'Все новые фичи публикуются здесь: \033[46m\033[30m@igor_chubin\033[0m',
        'te': u'అనుసరించండి \ 033 [46m \ 033 [30m @ igor_chubin \ 033 [wttr.in నవీకరణలను కోసం',
        'tr': u'wttr.in ile ilgili gelişmeler için \033[46m\033[30m@igor_chubin\033[0m adresini takip edin',
        'da': u'Følg \033[46m\033[30m@igor_chubin\033[0m for at få wttr.in opdateringer',
        'et': u'Jälgi \033[46m\033[30m@igor_chubin\033[0m wttr.in uudiste tarbeks',
    },
}

def get_message(message_name, lang):
    if message_name not in MESSAGE:
        return ''
    message_dict = MESSAGE[message_name]
    return message_dict.get(lang, message_dict.get('en', ''))

# vim: fileencoding=utf-8

"""
Translation of almost everything.
"""

FULL_TRANSLATION = [
    "am",
    "ar",
    "af",
    "be",
    "bn",
    "ca",
    "da",
    "de",
    "el",
    "et",
    "fr",
    "fa",
    "gl",
    "hi",
    "hu",
    "ia",
    "id",
    "it",
    "lt",
    "mg",
    "nb",
    "nl",
    "oc",
    "pl",
    "pt-br",
    "ro",
    "ru",
    "ta",
    "tr",
    "th",
    "uk",
    "vi",
    "zh-cn",
    "zh-tw",
]

PARTIAL_TRANSLATION = [
    "az",
    "bg",
    "bs",
    "cy",
    "cs",
    "eo",
    "es",
    "eu",
    "fi",
    "ga",
    "hi",
    "hr",
    "hy",
    "is",
    "ja",
    "jv",
    "ka",
    "kk",
    "ko",
    "ky",
    "lv",
    "mk",
    "ml",
    "mr",
    "nl",
    "fy",
    "nn",
    "pt",
    "pt-br",
    "sk",
    "sl",
    "sr",
    "sr-lat",
    "sv",
    "sw",
    "te",
    "uz",
    "zh",
    "zu",
    "he",
]

PROXY_LANGS = [
    "af",
    "am",
    "ar",
    "az",
    "be",
    "bn",
    "bs",
    "ca",
    "cy",
    "de",
    "el",
    "eo",
    "et",
    "eu",
    "fa",
    "fr",
    "fy",
    "ga",
    "gl",
    "he",
    "hr",
    "hu",
    "hy",
    "ia",
    "id",
    "is",
    "it",
    "ja",
    "kk",
    "lt",
    "lv",
    "mg",
    "mk",
    "mr",
    "nb",
    "nn",
    "oc",
    "ro",
    "ru",
    "sl",
    "th",
    "pt-br",
    "uk",
    "uz",
    "vi",
    "zh-cn",
    "zh-tw",
]

SUPPORTED_LANGS = FULL_TRANSLATION + PARTIAL_TRANSLATION


MESSAGE = {
    "NOT_FOUND_MESSAGE": {
        "en": """
We were unable to find your location
so we have brought you to Oymyakon,
one of the coldest permanently inhabited locales on the planet.
""",
        "af": """
Ons kon nie u ligging opspoor nie
gevolglik het ons vir u na Oymyakon geneem,
een van die koudste permanent bewoonde plekke op aarde.
""",
        "am": """
አካባቢዎን ማግኘት ስላልቻልን ወደ ኦይማያኮን አምጥተንዎታል፤ 
በአለም ላይ በጣም ቀዝቃዛ እና ሰው አልባ ከሆኑት ቦታዎች አንዱ።
""",
        "ar": """
تعذر علينا العثور على موقعك
لذلك قمنا بجلبك إلي أويمياكون,
 إحدى الأماكن المَأْهُولة الأكثر برودة علي الإطلاق في هذا الكوكب.
""",
        "be": """
Ваша месцазнаходжанне вызначыць не атрымалася,
таму мы пакажам вам надвор'е ў Аймяконе,
самым халодным населеным пункце на планеце.
Будзем спадзявацца, што ў вас сёння надвор'е лепей!
""",
        "bg": """
Не успяхме да открием вашето местоположение
така че ви доведохме в Оймякон,
едно от най-студените постоянно обитавани места на планетата.
""",
        "bn": """
দুঃখিত, আপনার অবস্থান আমরা খুঁজে পাইনি।
তাই, আমরা আপনাকে নিয়ে এসেছি ওয়মিয়াকনে, 
যা পৃথিবীর শীতলতম স্থায়ী জন-বসতিগুলোর একটি। 
""",
        "bs": """
Nismo mogli pronaći vašu lokaciju,
tako da smo te doveli do Oymyakon,
jedan od najhladnijih stalno naseljena mjesta na planeti.
Nadamo se da ćete imati bolje vreme!
""",
        "ca": """
Hem estat incapaços de trobar la seva ubicació,
per això l'hem portat fins Oymyakon,
un dels llocs més freds i permanentment deshabitats del planeta.
""",
        "cs": """
Nepodařilo se nám najít vaši polohu,
takže jsme vás přivedli do Ojmjakonu.
Je to jedno z nejchladnějších trvale obydlených míst na planetě.
Doufáme, že budete mít lepší počasí!
""",
        "cy": """
Ni darganfyddwyd eich lleoliad,
felly rydym wedi dod â chi i Oymyakon,
un o'r llefydd oeraf ar y blaned ble mae pobl yn dal i fyw!
""",
        "de": """
Wir konnten Ihren Standort nicht finden,
also haben wir Sie nach Oimjakon gebracht,
einer der kältesten dauerhaft bewohnten Orte auf dem Planeten.
Wir hoffen, dass Sie besseres Wetter haben!
""",
        "el": """
Δεν μπορέσαμε να βρούμε την τοποθεσία σου,
για αυτό διαλέξαμε το Οϊμιάκον για εσένα,
μία από τις πιο κρύες μόνιμα κατοικημένες περιοχές στον πλανήτη.
Ελπίζουμε να έχεις καλύτερο καιρό!
""",
        "es": """
No hemos logrado encontrar tu ubicación,
así que hemos decidido enseñarte el tiempo en Oimiakón,
uno de los sitios más fríos y permanentemente deshabitados del planeta.
""",
        "eu": """
Ezin izan dugu zure kokapena aurkitu,
beraz eguraldia Oymyakonen erakustea erabaki dugu,
munduko lekurik hotz eta hutseneraiko bat
""",
        "fa": """
ما نتونستیم مکان شما رو پیدا کنیم. به همین خاطر شما رو به اویمیاکن بردیم
، یکی از سردترین مکان های روی زمین که اصلا قابل سکونت نیست!
""",
        "fi": """
Emme löytänyt sijaintiasi, joten toimme sinut Oimjakoniin,
yhteen maailman kylmimmistä pysyvästi asutetuista paikoista.
Toivottavasti sinulla on parempi sää!
""",
        "fr": """
Nous n'avons pas pu déterminer votre position,
Nous vous avons donc amenés à Oïmiakon,
l'un des endroits les plus froids habités en permanence sur la planète.
Nous espérons qu'il fait meilleur chez vous !
""",
        "ga": """
Ní rabhamar ábalta do cheantar a aimsiú
mar sin thugamar go dtí Oymyakon,
tú ceann do na ceantair bhuanáitrithe is fuaire ar domhan.
""",
        "gl": """
Non logramos atopar a túa localización
polo que te trouxemos até Oimiakón,
un dos lugares máis fríos e permamentemente deshabitados do planeta.
""",
        "hi": """
हम आपका स्थान खोजने में असमर्थ है,
इसलिए हम आपको ओयमयाकोन पर ले आए है,
जो ग्रह के सबसे ठंडे स्थानों में से एक है|
""",
        "hu": """
Nem sikerült megtalálni a pozíciódat,
így elhoztunk Ojmjakonba;
az egyik leghidegebb állandóan lakott településre a bolygón.
""",
        "hy": """
Ձեր գտնվելու վայրը չհաջողվեց որոշել,
այդ պատճառով մենք ձեզ կցուցադրենք եղանակը Օյմյակոնում.
երկրագնդի ամենասառը բնակավայրում։
Հույս ունենք որ ձեր եղանակը այսօր ավելի լավն է։
""",
        "ia": """
Nos non trovate su location,
assi nos su apporte a Oymyakon,
un del plus frigide locos habita super le planeta!
""",
        "id": """
Kami tidak dapat menemukan lokasi anda,
jadi kami membawa anda ke Oymyakon,
salah satu tempat terdingin yang selalu dihuni di planet ini!
""",
        "is": """
Við finnum ekki staðsetninguna þína og vísum þér þar með á Ojmjakon,
ein af köldustu byggðum jarðar.
Vonandi er betra veður hjá þér.
""",
        "it": """
Non siamo riusciti a trovare la sua posizione
quindi la abbiamo portato a Oymyakon,
uno dei luoghi abitualmente abitati più freddi del pianeta.
Ci auguriamo che le condizioni dove lei si trova siano migliori!
""",
        "ja": """
指定された場所が見つかりませんでした。
代わりにオイミャコンの天気予報を表示しています。
オイミャコンは地球上で最も寒い居住地の一つです。
""",
        "ko": """
지정된 장소를 찾을 수 없습니다,
대신 오이먀콘의 일기 예보를 표시합니다,
오이먀콘은 지구상에서 가장 추운 곳에 위치한 마을입니다!
""",
        "lt": """
Mums nepavyko rasti jūsų vietovės,
todėl mes nukreipėme jus į Omjakoną,
vieną iš šalčiausių nuolatinių gyvenviečių planetoje.
""",
        "lv": """
Mēs nevarējām noteikt Tavu atrašanās vietu, tādēļ nogādājām Tevi Oimjakonā,
vienā no planētas aukstākajām pastāvīgi apdzīvotajām vietām.
""",
        "mk": """
Неможевме да ја пронајдеме вашата локација,
затоа ве однесовме во Ојмајкон,
еден од најладните трајно населени места на планетата.
""",
        "mr": """
आमहाला तुमचे स्थळ सापडले नाही.
म्हणून आम्ही तुम्हाला ओयम्याकोन येथे आणले आहे,
जे कि आपल्या ग्रहावरील सर्वात थंड वस्तिस्थानांपैकी एक आहे.
""",
        "nb": """
Vi kunne ikke finne din lokasjon,
så her får du Ojmjakon, et av de kaldeste bebodde stedene på planeten.
Vi håper været er bedre hos deg!
""",
        "nl": """
Wij konden uw locatie niet vaststellen
dus hebben we u naar Ojmjakon gebracht,
één van de koudste permanent bewoonde gebieden op deze planeet.
""",
        "fy": """
Wy koenen jo lokaasje net fêststelle
dus wy ha jo nei Ojmjakon brocht,
ien fan de kâldste permanent bewenbere plakken op ierde.
""",
        "oc": """
Avèm pas pogut determinar vòstra posicion,
Vos avèm doncas menat a Oïmiakon,
un dels endreches mai freds abitat permanéncia del monde.
Esperam que fa melhor en çò vòstre !
""",
        "pt": """
Não conseguimos encontrar a sua localização,
então decidimos te mostrar o tempo em Oymyakon,
um dos lugares mais frios e permanentemente desabitados do planeta.
""",
        "pt-br": """
Não conseguimos encontrar a sua localização,
então decidimos te mostrar o tempo em Oymyakon,
um dos lugares mais frios e permanentemente desabitados do planeta.
""",
        "pl": """
Nie udało nam się znaleźć podanej przez Ciebie lokalizacji,
więc zabraliśmy Cię do Ojmiakonu,
jednego z najzimniejszych, stale zamieszkanych miejsc na Ziemi.
Mamy nadzieję, że u Ciebie jest cieplej!
""",
        "ro": """
Nu v-am putut identifica localitatea, prin urmare vă arătăm vremea din Oimiakon,
una dintre cele mai reci localități permanent locuite de pe planetă.
Sperăm că aveți vreme mai bună!
""",
        "ru": """
Ваше местоположение определить не удалось,
поэтому мы покажем вам погоду в Оймяконе,
самом холодном населённом пункте на планете.
Будем надеяться, что у вас сегодня погода лучше!
""",
        "sk": """
Nepodarilo sa nám nájsť vašu polohu,
takže sme vás priviedli do Ojmiakonu.
Je to jedno z najchladnejších trvale obývaných miest na planéte.
Dúfame, že budete mať lepšie počasie!
""",
        "sr": """
Нисмо успели да пронађемо Вашу локацију,
па смо Вас довели у Ојмјакон,
једно од најхладнијих стално насељених места на планети.
Надамо се да је време код Вас боље него што је то случај овде!
""",
        "sv": """
Vi lyckades inte hitta er plats så vi har istället tagit er till Ojmjakon,
en av planetens kallaste platser med permanent bosättning.
Vi hoppas att vädret är bättre hos dig!
""",
        "tr": """
Aradığınız konum bulunamadı. O yüzden sizi dünyadaki en soğuk sürekli
yerleşim yerlerinden biri olan Oymyakon'e getirdik.
Umarız sizin olduğunuz yerde havalar daha iyidir!
""",
        "te": """
మేము మీ స్థానాన్ని కనుగొనలేకపోయాము
కనుక మనం "ఓమాయకాన్కు" తీసుకొని వచ్చాము,
భూమిపై అత్యల్ప శాశ్వతంగా నివసించే స్థానిక ప్రదేశాలలో ఒకటి.
""",
        "th": """
เราไม่สามารถหาตำแหน่งของคุณได้เราจึงนำคุณไปสู่ Oymyakon หมู่บ้านที่หนาวที่สุดในโลก!
""",
        "uk": """
Ми не змогли визначити Ваше місцезнаходження,
тому покажемо Вам погоду в Оймяконі —
найхолоднішому населеному пункті на планеті.
Будемо сподіватися, що у Вас сьогодні погода краще!
""",
        "uz": """
Sizning joylashuvingizni aniqlay olmadik,
shuning uchun sizga sayyoramizning eng sovuq aholi punkti - Oymyakondagi ob-havo haqida ma'lumot beramiz.
Umid qilamizki, sizda bugungi ob-havo bundan yaxshiroq!
""",
        "zh": """
我们无法找到您的位置,
当前显示奥伊米亚康(Oymyakon)，这个星球上最冷的人类定居点。
""",
        "da": """
Vi kunne desværre ikke finde din lokation
så vi har bragt dig til Oymyakon,
En af koldeste og helt ubolige lokationer på planeten.
""",
        "et": """
Me ei suutnud tuvastada teie asukohta
ning seetõttu paigutasime teid Oymyakoni,
mis on üks kõige külmemaid püsivalt asustatud paiku planeedil.
""",
        "vi": """
Chúng tôi không tìm thấy địa điểm của bạn
vì vậy chúng tôi đưa bạn đến Oymyakon,
một trong những nơi lạnh nhất có người sinh sống trên trái đất.
""",
        "zh-tw": """
我們找不到您的位置
所以我們帶您到奧伊米亞康，
這個星球上有人類定居最冷之處。
""",
        "mg": """
Tsy hita ny toerana misy anao koa nentinay tany Oymyakon ianao,
iray amin'ireo toerana mangatsiaka indrindra tsisy mponina eto an-tany.
""",
        "ta": """
உங்கள் இருப்பிடத்தை எங்களால் கண்டுபிடிக்க முடியவில்லை
எனவே நாங்கள் உங்களை ஓமியாகோனுக்கு அழைத்து வந்தோம்.
கிரகத்தின் குளிர்ந்த நிரந்தரமாக வசிக்கும் இடங்களில் ஒன்று.
""",
    },
    "UNKNOWN_LOCATION": {
        "en": "Unknown location",
        "af": "Onbekende ligging",
        "am": "ያልታወቀ ቦታ",
        "ar": "موقع غير معروف",
        "be": "Невядомае месцазнаходжанне",
        "bg": "Неизвестно местоположение",
        "bn": "অজানা অবস্থান",
        "bs": "Nepoznatoja lokacija",
        "ca": "Ubicació desconeguda",
        "cs": "Neznámá poloha",
        "cy": "Lleoliad anhysbys",
        "de": "Unbekannter Ort",
        "da": "Ukendt lokation",
        "el": "Άνγωστη τοποθεσία",
        "es": "Ubicación desconocida",
        "et": "Tundmatu asukoht",
        "eu": "Kokapen ezezaguna",
        "fa": "مکان نامعلوم",
        "fi": "Tuntematon sijainti",
        "fr": "Emplacement inconnu",
        "ga": "Ceantar anaithnid",
        "gl": "Localización descoñecida",
        "hi": "अज्ञात स्थान",
        "hu": "Ismeretlen lokáció",
        "hy": "Անհայտ գտնվելու վայր",
        "id": "Lokasi tidak diketahui",
        "ia": "Location incognite",
        "is": "Óþekkt staðsetning",
        "it": "Località sconosciuta",
        "ja": "未知の場所です",
        "ko": "알 수 없는 장소",
        "kk": "",
        "lt": "Nežinoma vietovė",
        "lv": "Nezināma atrašanās vieta",
        "mk": "Непозната локација",
        "mr": "अज्ञात स्थळ",
        "nb": "Ukjent sted",
        "nl": "Onbekende locatie",
        "oc": "Emplaçament desconegut",
        "fy": "Ûnbekende lokaasje",
        "pl": "Nieznana lokalizacja",
        "pt": "Localização desconhecida",
        "pt-br": "Localização desconhecida",
        "ro": "Localitate necunoscută",
        "ru": "Неизвестное местоположение",
        "sk": "Neznáma poloha",
        "sl": "Neznano lokacijo",
        "sr": "Непозната локација",
        "sv": "Okänd plats",
        "te": "తెలియని ప్రదేశం",
        "tr": "Bilinmeyen konum",
        "th": "ไม่สามารถระบุตำแหน่งได้",
        "uk": "Невідоме місце",
        "uz": "Аникланмаган худуд",
        "zh": "未知地点",
        "vi": "Địa điểm không xác định",
        "zh-tw": "未知位置",
        "mg": "Toerana tsy fantatra",
        "ta": "தெரியாத இடம்",
    },
    "LOCATION": {
        "en": "Location",
        "af": "Ligging",
        "ar": "الموقع",
        "be": "Месцазнаходжанне",
        "bg": "Местоположение",
        "bn": "অবস্থান",
        "bs": "Lokacija",
        "ca": "Ubicació",
        "cs": "Poloha",
        "cy": "Lleoliad",
        "de": "Ort",
        "da": "Lokation",
        "el": "Τοποθεσία",
        "es": "Ubicación",
        "et": "Asukoht",
        "eu": "Kokaena",
        "fa": "مکان",
        "fi": "Tuntematon sijainti",
        "fr": "Emplacement",
        "ga": "Ceantar",
        "gl": "Localización",
        "hi": "स्थान",
        "hu": "Lokáció",
        "hy": "Դիրք",
        "ia": "Location",
        "id": "Lokasi",
        "is": "Staðsetning",
        "it": "Località",
        "ja": "位置情報",
        "ko": "위치",
        "kk": "",
        "lt": "Vietovė",
        "lv": "Atrašanās vieta",
        "mk": "Локација",
        "mr": "स्थळ",
        "nb": "Sted",
        "nl": "Locatie",
        "oc": "Emplaçament",
        "fy": "Lokaasje",
        "pl": "Lokalizacja",
        "pt": "Localização",
        "pt-br": "Localização",
        "ro": "Localitate",
        "ru": "Местоположение",
        "sk": "Poloha",
        "sl": "Lokacijo",
        "sr": "Локација",
        "sv": "Plats",
        "zh": "地点",
        "te": "స్థానము",
        "tr": "Konum",
        "th": "ตำแหน่ง",
        "uk": "Місцезнаходження",
        "vi": "Địa điểm",
        "zh-tw": "位置",
        "mg": "Toerana",
        "ta": "இடம்",
    },
    "CAPACITY_LIMIT_REACHED": {
        "en": """
Sorry, we are running out of queries to the weather service at the moment.
Here is the weather report for the default city (just to show you what it looks like).
We will get new queries as soon as possible.
You can follow https://x.com/igor_chubin for the updates.
======================================================================================
""",
        "af": """
Verskoning, ons oorskry tans die vermoë om navrae aan die weerdiens te rig.
Hier is die weerberig van 'n voorbeeld ligging (bloot om aan u te wys hoe dit lyk).
Ons sal weereens nuwe navrae kan hanteer so gou as moontlik.
U kan vir https://x.com/igor_chubin volg vir opdaterings.
======================================================================================
""",
        "am": """
ይቅርታ ፤ ለጊዜው መጠይቆችን ማስተናገድ አልቻልንም፡፡
የነባሪ ከተማው የአየር ሁኔታ ሪፖርት ይኸውልዎ (ምን እንደሚመስል ለማሳየት)።
አዳዲስ መጠይቆችን በተቻለ ፍጥነት ለማስተናገድ እንሞክራለን።
ለበለጠ መረጃ እባክዎ https://x.com/igor_chubin ን ይከተሉ።
======================================================================================
""",
        "ar": """
عذرًا ، استعلامات خدمة الطقس نفذت في الوقت الحالي.
هذا هو تقرير الطقس للمدينة الافتراضية (فقط لتظهر لك كيف تبدو).
سوف نتلقى استعلامات جديدة في أقرب وقت ممكن.
يمكنك متابعة https://x.com/igor_chubin لآخر المستجدات.
======================================================================================
""",
        "be": """
Прабачце, мы выйшлі за ліміты колькасці запытаў да службы надвор'я ў дадзены момант.
Вось прагноз надвор'я для горада па змаўчанні (толькі, каб паказаць вам, як гэта выглядае).
Мы вернемся як мага хутчэй.
Вы можаце сачыць на https://x.com/igor_chubin за абнаўленнямі.
======================================================================================
""",
        "bg": """
Съжаляваме, количеството заявки към услугата за предсказване на време е почти изчерпано.
Ето доклад за града по подразбиране (просто да видите как изглежда).
Ще осогурим допълнителни заявки максимално бързо.
Може да последвате https://x.com/igor_chubin за обновления.
""",
        "bn": """
দুঃখিত, এই মুহুর্তে আবহাওয়া পরিসেবাতে  আমাদের কুইরী শেষ হয়ে আসছে। 
এখানে ডিফল্ট শহরের আবহাওয়ার প্রতিবেদন রয়েছে (এটি দেখতে কেমন তা আপনাকে দেখানোর জন্য)। 
আমরা খুব দ্রুত নতুন কুইরী পাওয়ার ব্যবস্থা করছি। 
আপডেটের জন্য আপনি https://x.com/igor_chubin অনুসরণ করতে পারেন। 
""",
        "bs": """
Žao mi je, mi ponestaje upita i vremenska prognoza u ovom trenutku.
Ovdje je izvještaj o vremenu za default grada (samo da vam pokažem kako to izgleda).
Mi ćemo dobiti nove upite u najkraćem mogućem roku.
Možete pratiti https://x.com/igor_chubin za ažuriranja.
======================================================================================
""",
        "ca": """
Disculpa'ns, ens hem quedat momentàniament sense consultes al servei meteorològic.
Aquí t'oferim l'informe del temps a la ciutat per defecte (només per mostrar quin aspecte té).
Obtindrem noves consultes tan aviat com ens sigui possible.
Pots seguir https://x.com/igor_chubin per noves actualitzacions.
======================================================================================
""",
        "cs": """
Omlouváme se, ale aktuálně nám dochází počet dotazů na službu počasí.
Zde je předpověď počasí pro výchozí město (abyste viděli, jak to vypadá).
Problém bude vyřešen co nejdříve.
Aktualizace můžete následovat na https://x.com/igor_chubin.
======================================================================================
""",
        "de": """
Entschuldigung, wir können momentan den Wetterdienst nicht erreichen.
Dafür zeigen wir Ihnen das Wetter an einem Beispielort, damit Sie sehen, wie die Seite das Wetter anzeigt.
Wir werden versuchen das Problem so schnell wie möglich zu beheben.
Folgen Sie https://x.com/igor_chubin für Updates.
======================================================================================
""",
        "cy": """
Rydym yn brin o ymholiadau i'r gwasanaeth tywydd ar hyn o bryd.
Felly dyma'r adroddiad tywydd ar gyfer y ddinas ragosod (er mwyn arddangos sut mae'n edrych).
Byddwn gyda ymholiadau newydd yn fuan.
Gellir dilyn https://x.com/igor_chubin i gael newyddion pellach.
======================================================================================
""",
        "es": """
Lo siento, hemos alcanzado el límite de peticiones al servicio de previsión del tiempo en este momento.
A continuación, la previsión del tiempo para una ciudad estándar (solo para que puedas ver el formato del reporte).
Muy pronto volveremos a tener acceso a las peticiones.
Puedes seguir https://x.com/igor_chubin para estar al tanto de la situación.
======================================================================================
""",
        "fa": """
متأسفانه در حال حاضر ظرفیت ما برای درخواست به سرویس هواشناسی به اتمام رسیده.
اینجا می تونید گزارش هواشناسی برای شهر پیش فرض رو ببینید (فقط برای اینه که بهتون نشون بدیم چه شکلی هست)
ما تلاش میکنیم در اسرع وقت ظرفیت جدید به دست بیاریم.
برای دنبال کردن اخبار جدید میتونید https://x.com/igor_chubin رو فالو کنید.
======================================================================================
""",
        "fr": """
Désolé, nous avons épuisé les requêtes vers le service météo.
Voici un bulletin météo de l'emplacement par défaut (pour vous donner un aperçu).
Nous serons très bientôt en mesure de faire de nouvelles requêtes.
Vous pouvez suivre https://x.com/igor_chubin pour rester informé.
======================================================================================
""",
        "ga": """
Tá brón orainn, níl mórán iarratas le fail chuig seirbhís na haimsire faoi láthair.
Seo duit réamhaisnéis na haimsire don chathair réamhshocraithe (chun é a thaispeaint duit).
Gheobhaimid iarratais nua chomh luath agus is feidir.
Lean orainn ar https://x.com/igor_chubin don eolas is déanaí.
======================================================================================
""",
        "gl": """
Desculpa, estamos a chegar ao límite de peticións ao servizo meteorolóxico neste momento.
Aquí está a previsión do tempo para a cidade por defecto (tan só para amosarche un exemplo).
Imos obter máis peticións tan pronto como poidamos.
Podes seguir https://x.com/igor_chubin para estares actualizada.
======================================================================================
""",
        "hi": """
क्षमा करें, इस समय हम मौसम सेवा से संपर्क नहीं कर पा रहे है।
यहा पूर्व निर्धारित शहर के लिए मौसम की जानकारी (आपको यह दिखाने के लिए कि यह कैसा दिखता है)।
हम जल्द से जल्द मौसम सेवा से संपर्क करने नई कोशिश करेंगे।
आप अपडेट के लिए https://x.com/igor_chubin का अनुसरण कर सकते हैं।
======================================================================================
""",
        "hu": """
Sajnáljuk, kifogytunk az időjárási szolgáltatásra fordított erőforrásokból.
Itt van az alapértelmezett város időjárási jelentése (hogy lásd, hogyan néz ki).
A lehető leghamarabb új erőforrásokat fogunk kapni.
A frissítésekért tekintsd meg a https://x.com/igor_chubin oldalt.
======================================================================================
""",
        "hy": """
Կներեք, այս պահին մենք գերազանցել ենք եղանակային տեսության կայանին հարցումների քանակը.
Կարող եք տեսնել տիպային եղանակը զեկուցում հիմնական քաղաքի համար (Ուղղակի որպես նմուշ):
Մենք մշտապես աշխատում ենք հարցումների քանակը բարելավելու ուղղությամբ:
Կարող եք հետևել մեզ https://x.com/igor_chubin թարմացումների համար.
======================================================================================
""",
        "ia": """
Pardono, nos ha exhaurite inquisitione del servicio tempore alora.
Ecce es le reporto tempore por qualque citate (demonstrar le reporte).
Nos recipera inquisitione nove tosto.
Tu pote abonar a https://x.com/igor_chubin por nove information.
======================================================================================
""",
        "id": """
Maaf, kami kehabian permintaan ke layanan cuaca saat ini.
Ini adalah laporan cuaca dari kota standar (hanya untuk menunjukkan kepada anda bagaimana tampilannya).
Kami akan mencoba permintaan baru lagi sesegera mungkin.
Anda dapat mengikuti https://x.com/igor_chubin untuk informasi terbaru.
======================================================================================
""",
        "it": """
Scusate, attualmente stiamo esaurendo le risorse a disposizione del servizio meteo.
Qui trovate il bollettino del tempo per la città di default (solo per mostrarvi come si presenta).
Potremo elaborare nuove richieste appena possibile.
Potete seguire https://x.com/igor_chubin per gli aggiornamenti.
======================================================================================
""",
        "ko": """
죄송합니다. 현재 날씨 정보를 가져오는 쿼리 요청이 한도에 도달했습니다.
대신 기본으로 설정된 도시에 대한 일기 예보를 보여드리겠습니다. (이는 단지 어떻게 보이는지 알려주기 위함입니다).
쿼리 요청이 가능한 한 빨리 이루어질 수 있도록 하겠습니다.
업데이트 소식을 원하신다면 https://x.com/igor_chubin 을 팔로우 해주세요.
======================================================================================
""",
        "lt": """
Atsiprašome, šiuo metu pasiekėme orų prognozės paslaugos užklausų ribą.
Štai orų prognozė numatomam miestui (tam, kad parodytume, kaip ji atrodo).
Naujas užklausas priimsime, kai tik galėsime.
Atnaujinimus galite sekti https://x.com/igor_chubin
======================================================================================
""",
        "lv": """
Atvainojamies, šobrīd laikapstākļu pakalpojums ir noslogots!
Šeit ir laika ziņas noklusējuma pilsētai (lai parādītu, kā tās izskatās).
Mēs atsāksim darbu, cik ātri vien varēsim.
Var sekot https://x.com/igor_chubin, lai uzzinātu jaunumus.
======================================================================================
""",
        "mk": """
Извинете, ни снемуваат барања за до сервисот кој ни нуди временска прогноза во моментот.
Еве една временска прогноза за град (за да видите како изгледа).
Ќе добиеме нови барања најбрзо што можеме.
Следете го https://x.com/igor_chubin за известувања
======================================================================================
""",
        "mr": """
क्षमस्व, याक्षणी आम्ही हवामान सेवेशी संपर्क करू शकत नाही.
हा एका पूर्वनिर्धारित शहराचा हवामान अहवाल आहे (केवळ तो कसा दिसतो हे दाखवण्याकरिता).
आम्ही लवकरात लवकर सेवा पुनः चालू करण्याचा प्रयत्न करू.
अद्यावत माहितीसाठी तुम्ही https://x.com/igor_chubin चे अनुसरण करू शकता.
======================================================================================
""",
        "nb": """
Beklager, vi kan ikke nå værtjenesten for øyeblikket.
Her er værmeldingen for standardbyen så du får se hvordan tjenesten ser ut.
Vi vil forsøke å fikse problemet så snart som mulig.
Du kan følge https://x.com/igor_chubin for oppdateringer.
======================================================================================
""",
        "zh": """
抱歉，当前天气服务不可用。
以下显示默认城市（只对您可见）。
我们将会尽快获取新数据。
您可以通过 https://x.com/igor_chubin 获取最新动态。
======================================================================================
""",
        "nl": """
Excuse, wij kunnen u op dit moment dit weerbericht niet laten zien.
Hier is het weerbericht voor de standaard stad(zodat u weet hoe het er uitziet)
Wij lossen dit probleem zo snel mogelijk op.
voor updates kunt u ons op https://x.com/igor_chubin volgen.
======================================================================================
""",
        "oc": """
O planhèm, avèm pas mai de requèstas cap al servici metèo.
Vaquí las prevision metèo de l'emplaçament per defaut (per vos donar un apercebut).
Poirem lèu ne far de novèlas.
Podètz seguir https://x.com/igor_chubin per demorar informat.
======================================================================================
""",
        "fy": """
Excuses, wy kinne op dit moment 't waarberjocht net sjin litte.
Hjir is 't waarberjocht foar de standaard stêd.
Wy losse dit probleem sa gau mooglik op.
Foar updates kinne jo ús op https://x.com/igor_chubin folgje.
======================================================================================
""",
        "pl": """
Bardzo nam przykro, ale chwilowo wykorzystaliśmy limit zapytań do serwisu pogodowego.
To, co widzisz jest przykładowym raportem pogodowym dla domyślnego miasta.
Postaramy się przywrócić funkcjonalność tak szybko, jak to tylko możliwe.
Możesz śledzić https://x.com/igor_chubin na Twitterze, aby być na bieżąco.
======================================================================================
""",
        "pt": """
Desculpe-nos, estamos atingindo o limite de consultas ao serviço de previsão do tempo neste momento.
Veja a seguir a previsão do tempo para uma cidade padrão (apenas para você ver que aspecto o relatório tem).
Em breve voltaremos a ter acesso às consultas.
Você pode seguir https://x.com/igor_chubin para acompanhar a situação.
======================================================================================
""",
        "pt-br": """
Desculpe-nos, atingimos o limite de consultas ao serviço de previsão do tempo neste momento.
Veja a seguir a previsão do tempo para uma cidade padrão (apenas para você ver que aspecto o relatório tem).
Em breve voltaremos a ter acesso às consultas.
Você pode seguir https://x.com/igor_chubin para acompanhar a situação.
======================================================================================
""",
        "ro": """
Ne pare rău, momentan am epuizat cererile alocate de către serviciul de prognoză meteo.
Vă arătăm prognoza meteo pentru localitatea implicită (ca exemplu, să vedeți cum arată).
Vom obține alocarea de cereri noi cât de curând posibil.
Puteți urmări https://x.com/igor_chubin pentru actualizări.
======================================================================================
""",
        "te": """
క్షమించండి, ప్రస్తుతానికి మేము వాతావరణ సేవకు ప్రశ్నలను గడుపుతున్నాం.
ఇక్కడ డిఫాల్ట్ నగరం కోసం వాతావరణ నివేదిక (కేవలం మీకు చూపించడానికి, ఇది ఎలా కనిపిస్తుంది).
సాధ్యమైనంత త్వరలో కొత్త ప్రశ్నలను పొందుతారు.
నవీకరణల కోసం https://x.com/igor_chubin ను మీరు అనుసరించవచ్చు.
======================================================================================
""",
        "tr": """
Üzgünüz, an itibariyle hava durumu servisine yapabileceğimiz sorgu limitine ulaştık.
Varsayılan şehir için hava durumu bilgisini görüyorsunuz (neye benzediğini gösterebilmek için).
Mümkün olan en kısa sürede servise yeniden sorgu yapmaya başlayacağız.
Gelişmeler için https://x.com/igor_chubin adresini takip edebilirsiniz.
======================================================================================
""",
        "th": """
ขออภัย การเรียกค้นหาระบบสภาพอากาศของเราหมดชั่วคราว
เราจึงแสดงข้อมูลของเมืองตัวอย่าง (เพื่อที่จะแสดงให้คุณเห็นว่าหน้าตาเป็นยังไง)
เราจะเรียกค้นหาใหม่โดยเร็วที่สุด
คุณสามารถติดตามการอัพเดทได้ที่ https://x.com/igor_chubin
======================================================================================
""",
        "da": """
Beklager, men vi er ved at løbe tør for forespørgsler til vejr-servicen lige nu.
Her er vejr rapporten for standard byen (bare så du ved hvordan det kan se ud).
Vi får nye forespørsler hurtigst muligt.
Du kan følge https://x.com/igor_chubin for at få opdateringer.
======================================================================================
""",
        "et": """
Vabandage, kuid hetkel on päringud ilmateenusele piiratud.
Selle asemel kuvame hetkel näidislinna ilmaprognoosi (näitamaks, kuidas see välja näeb).
Üritame probleemi lahendada niipea kui võimalik.
Jälgige https://x.com/igor_chubin värskenduste jaoks.
======================================================================================
""",
        "uk": """
Вибачте, ми перевищили максимальну кількість запитів до сервісу погоди.
Ось прогноз погоди у нашому місті (просто показати Вам як це виглядає).
Ми відновимо роботу як тільки зможемо.
Ви можете підписатися на https://x.com/igor_chubin для отримання новин.
======================================================================================
""",
        "vi": """
Xin lỗi, hiện tại chúng tôi đã hết lượt yêu cầu thông tin thời tiết.
Đây là dự báo thời tiết cho thành phố mặc định (chỉ để cho bạn thấy nó trông như thế nào).
Chung tôi sẽ có thêm lượt truy vấn sớm nhất có thể
Bạn có thể theo dõi https://x.com/igor_chubin để cập nhật thông tin mới nhất.
======================================================================================
""",
        "zh-tw": """
抱歉，目前天氣服務的查詢請求太多了。
這裡是預設城市的天氣報告（只是為您展示它的外觀）。
我們將盡快取得新的資料。
您可以追蹤 https://x.com/igor_chubin 以取得更新。
======================================================================================
""",
        "mg": """
Miala tsiny fa misedra olana ny sampan-draharaha momba ny toetrandro amin'izao fotoana izao.
Ity ny tatitra momba ny toetr'andro ho an'ny tanàna mahazatra (mba hampisehoana anao ny endriny).
Haivaly aminao haingana ny fangatahanao.
Azonao atao ny manaraka ny pejy https://x.com/igor_chubin.
""",
        "ta": """
மன்னிக்கவும், தற்போது வானிலை சேவைக்கான வினவல்கள் எங்களிடம் இல்லை.
இயல்புநிலை நகரத்திற்கான வானிலை அறிக்கை இதோ (அது எப்படி இருக்கும் என்பதை உங்களுக்குக் காண்பிப்பதற்காக).
கூடிய விரைவில் புதிய வினவல்களைப் பெறுவோம்.
புதுப்பிப்புகளுக்கு நீங்கள் https://x.com/igor_chubin ஐப் பின்தொடரலாம்.
""",
    },
    # Historical messages:
    #     'Check new Feature: \033[92mwttr.in/Moon\033[0m or \033[92mwttr.in/Moon@2016-Mar-23\033[0m to see the phase of the Moon'
    #     'New feature: \033[92mwttr.in/Rome?lang=it\033[0m or \033[92mcurl -H "Accept-Language: it" wttr.in/Rome\033[0m for the localized version. Your lang instead of "it"'
    "NEW_FEATURE": {
        "en": "New feature: multilingual location names \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) and location search \033[92mwttr.in/~Kilimanjaro\033[0m (just add ~ before)",
        "am": "አዲስ ባህሪ: የአካባቢ ስሞች በተለያዩ ቋንቋዎች \033[92mwttr.in/ኢትዮጵያ\033[0m (በ UTF-8) እና የአካባቢ ፍለጋ \033[92mwttr.in/~Kilimanjaro\033[0m (ከአካባቢ ስም በፊት ~ ያክሉ)",
        "ar": "ميزة جديدة : أسماء الأماكن بلغات متعددة \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) والبحث عن الأماكن \033[92mwttr.in/~Kilimanjaro\033[0m (فقط أضف ~ قبل)",
        "af": "Nuwe eienskap: veeltalige name vir liggings \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) en ligging soek \033[92mwttr.in/~Kilimanjaro\033[0m (plaas net ~ vooraan)",
        "be": "Новыя магчымасці: назвы месц на любой мове \033[92mwttr.in/станция+Восток\033[0m (в UTF-8) i пошук месц \033[92mwttr.in/~Kilimanjaro\033[0m (трэба дадаць ~ ў пачатак)",
        "bg": "Нова функционалност: многоезични имена на места\033[92mwttr.in/станция+Восток\033[0m (в UTF-8) и в търсенето \033[92mwttr.in/~Kilimanjaro\033[0m (добавете ~ преди)",
        "bn": "নতুন ফিচার : বহুভাষিক অবস্থানের নাম \ 033 [92mwttr.in/станция+Восток\033 [0m (UTF-8)] এবং অবস্থান অনুসন্ধান \ 033 [92mwttr.in/~Kilimanjaro\033 [0m (শুধু আগে ~ যোগ করুন)",
        "bs": "XXXXXXXXXXXXXXXXXXXX: XXXXXXXXXXXXXXXXXXXXXXXXXXXXX\033[92mwttr.in/станция+Восток\033[0m (XX UTF-8) XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
        "ca": "Noves funcionalitats: noms d'ubicació multilingües \033[92mwttr.in/станция+Восток\033[0m (en UTF-8) i la ubicació de recerca \033[92mwttr.in/~Kilimanjaro\033[0m (només cal afegir ~ abans)",
        "cs": "Nová funkce: vícejazyčná jména poloh  \033[92mwttr.in/станция+Восток\033[0m (v UTF-8) a vyhledávání polohy \033[92mwttr.in/~Kilimanjaro\033[0m (pouze přidejte ~ před pozici)",
        "es": "Nuevas funcionalidades: los nombres de las ubicaciones en varios idiomas \033[92mwttr.in/станция+Восток\033[0m (em UTF-8) y la búsqueda por ubicaciones \033[92mwttr.in/~Kilimanjaro\033[0m (tan solo inserte ~ al principio)",
        "fa": "قابلیت جدید: پشتیبانی از نام چند زبانه مکانها \033[92mwttr.in/станция+Восток\033[0m (در فرمت UTF-8) و جسجتوی مکان ها \033[92mwttr.in/~Kilimanjaro\033[0m (فقط قبل از اون ~ اضافه کنید)",
        "fr": "Nouvelles fonctionnalités: noms d'emplacements multilingues \033[92mwttr.in/станция+Восток\033[0m (en UTF-8) et recherche d'emplacement \033[92mwttr.in/~Kilimanjaro\033[0m (ajouter ~ devant)",
        "gl": "Nova funcionalidade: nomes de localizacións en varios idiomas\033[92mwttr.in/станция+Восток\033[0m (en UTF-8) e procuras de localizacións \033[92mwttr.in/~Kilimanjaro\033[0m (engade ~ antes)",
        "mk": "Нова функција: повеќе јазично локациски имиња \033[92mwttr.in/станция+Восток\033[0m (во UTF-8) и локациско пребарување \033[92mwttr.in/~Kilimanjaro\033[0m (just add ~ before)",
        "nb": "Ny funksjon: flerspråklige stedsnavn \033[92mwttr.in/станция+Восток\033[0m (i UTF-8) og lokasjonssøk \033[92mwttr.in/~Kilimanjaro\033[0m (bare legg til ~ foran)",
        "nl": "Nieuwe functie: tweetalige locatie namen \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) en locatie zoeken \033[92mwttr.in/~Kilimanjaro\033[0m (zet er gewoon een ~ voor)",
        "fy": "Nije funksje: twatalige lokaasje nammen \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) en lokaasje sykjen \033[92mwttr.in/~Kilimanjaro\033[0m (set er gewoan in ~ foar)",
        "cy": "Nodwedd newydd: enwau lleoliadau amlieithog \033[92mwttr.in/станция+Восток\033[0m (yn UTF-8) a chwilio am leoliad \033[92mwttr.in/~Kilimanjaro\033[0m (ychwanegwch ~ yn gyntaf)",
        "de": "Neue Funktion: mehrsprachige Ortsnamen \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) und Ortssuche \033[92mwttr.in/~Kilimanjaro\033[0m (fügen Sie ein ~ vor dem Ort ein)",
        "hi": "नई सुविधा: बहुभाषी स्थान के नाम \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) और स्थान खोज \033[92mwttr.in/~Kilimanjaro\033[0m (बस ~ आगे लगाये)",
        "hu": "Új funkcinalitás: többnyelvű helynevek \033[92mwttr.in/станция+Восток\033[0m (UTF-8-ban) és pozíció keresés \033[92mwttr.in/~Kilimanjaro\033[0m (csak adj egy ~ jelet elé)",
        "hy": "Փորձարկեք: տեղամասերի անունները կամայական լեզվով \033[92mwttr.in/Դիլիջան\033[0m (в UTF-8) և տեղանքի որոնում \033[92mwttr.in/~Kilimanjaro\033[0m (հարկավոր է ~ ավելացնել դիմացից)",
        "ia": "Nove functione: location nomine multilingue \033[92mwttr.in/станция+Восток\033[0m (a UTF-8) e recerca de location\033[92mwttr.in/~Kilimanjaro\033[0m (solo adde ~ ante)",
        "id": "Fitur baru: nama lokasi dalam multibahasa \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) dan pencarian lokasi \033[92mwttr.in/~Kilimanjaro\033[0m (hanya tambah tanda ~ sebelumnya)",
        "it": "Nuove funzionalità: nomi delle località multilingue \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) e ricerca della località \033[92mwttr.in/~Kilimanjaro\033[0m (basta premettere ~)",
        "ko": "새로운 기능: 다국어로 대응된 위치 \033[92mwttr.in/서울\033[0m (UTF-8에서) 장소 검색 \033[92mwttr.in/~Kilimanjaro\033[0m (앞에 ~를 붙이세요)",
        "kk": "",
        "lt": "Naujiena: daugiakalbiai vietovių pavadinimai \033[92mwttr.in/станция+Восток\033[0m (UTF-8) ir vietovių paieška \033[92mwttr.in/~Kilimanjaro\033[0m (tiesiog priekyje pridėkite ~)",
        "lv": "Jauna iespēja: daudzvalodu atrašanās vietu nosaukumi \033[92mwttr.in/станция+Восток\033[0m (UTF-8) un atrašanās vietu meklēšana \033[92mwttr.in/~Kilimanjaro\033[0m (sākumā tikai jāpievieno ~)",
        "mk": "Нова функција: повеќе јазично локациски имиња \033[92mwttr.in/станция+Восток\033[0m (во UTF-8) и локациско пребарување \033[92mwttr.in/~Kilimanjaro\033[0m (just add ~ before)",
        "mr": "नवीन वैशिष्ट्य: स्थळांची बहुभाषिक नावे \033[92mwttr.in/станция+Восток\033[0m (UTF-8 मध्ये) आणि स्थळ शोध \033[92mwttr.in/~Kilimanjaro\033[0m (फक्त आधी ~ जोडा)",
        "oc": "Novèla foncionalitat : nom de lòc multilenga \033[92mwttr.in/станция+Восток\033[0m (en UTF-8) e recèrca de lòc \033[92mwttr.in/~Kilimanjaro\033[0m (solament ajustatz ~ abans)",
        "pl": "Nowa funkcjonalność: wielojęzyczne nazwy lokalizacji \033[92mwttr.in/станция+Восток\033[0m (w UTF-8) i szukanie lokalizacji \033[92mwttr.in/~Kilimanjaro\033[0m (poprzedź zapytanie ~ - znakiem tyldy)",
        "pt": "Nova funcionalidade: nomes de localidades em várias línguas \033[92mwttr.in/станция+Восток\033[0m (em UTF-8) e procura por localidades \033[92mwttr.in/~Kilimanjaro\033[0m (é só colocar ~ antes)",
        "pt-br": "Nova funcionalidade: nomes de localidades em várias línguas \033[92mwttr.in/станция+Восток\033[0m (em UTF-8) e procura por localidades \033[92mwttr.in/~Kilimanjaro\033[0m (é só colocar ~ antes)",
        "ro": "Funcționalitate nouă: nume de localități multilingve \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) și căutare de localități \033[92mwttr.in/~Kilimanjaro\033[0m (adăuați ~ în față)",
        "ru": "Попробуйте: названия мест на любом языке \033[92mwttr.in/станция+Восток\033[0m (в UTF-8) и поиск мест \033[92mwttr.in/~Kilimanjaro\033[0m (нужно добавить ~ спереди)",
        "zh": "新功能：多语言地点名称 \033[92mwttr.in/станция+Восток\033[0m (in UTF-8) 及地点搜索\033[92mwttr.in/~Kilimanjaro\033[0m （只需在名称前加~）",
        "tr": "Yeni özellik: çok dilli konum isimleri \033[92mwttr.in/станция+Восток\033[0m (UTF-8 ile) ve konum arama \033[92mwttr.in/~Kilimanjaro\033[0m (sadece önüne ~ ekleyin)",
        "te": "క్రొత్త లక్షణం: బహుభాషా స్థాన పేర్లు \ 033 [92mwttr.in/stancelя+Vostок\033 [0 U (UTF-8 లో) మరియు స్థానం శోధన \ 033 [92mwttr.in/~kilimanjaro\033 [0m (కేవలం ~ ముందుకి జోడించండి)",
        "th": "ฟีเจอร์ใหม่: แสดงที่ตั้งได้หลายภาษา \033[92mwttr.in/станция+Восток\033[0m (ใน UTF-8) และการค้นหาที่ตั้ง \033[92mwttr.in/~Kilimanjaro\033[0m (เพียงเพิ่ม ~ ข้างหน้า)",
        "da": "Ny funktion: flersprogede lokationsnavne \033[92mwttr.in/станция+Восток\033[0m (som UTF-8) og lokations søgning \033[92mwttr.in/~Kilimanjaro\033[0m (bare tilføj ~ inden)",
        "et": "Uus funktsioon: mitmekeelsed asukohanimed \033[92mwttr.in/станция+Восток\033[0m (UTF-8 vormingus) ja asukoha otsing \033[92mwttr.in/~Kilimanjaro\033[0m (lisa ~ enne)",
        "uk": "Спробуйте: назви місць будь-якою мовою \033[92mwttr.in/станція+Восток\033[0m (в UTF-8) та пошук місць \033[92mwttr.in/~Kilimanjaro\033[0m (потрібно додати ~ спочатку)",
        "vi": "Chức năng mới: tên địa điểm đa ngôn ngữ \033[92mwttr.in/станция+Восток\033[0m (dùng UTF-8) và tìm kiếm địa điểm \033[92mwttr.in/~Kilimanjaro\033[0m (chỉ cần thêm ~ phía trước)",
        "zh-tw": "新功能：多語言地點名稱 \033[92mwttr.in/станция+Восток\033[0m （使用 UTF-8 編碼）與位置搜尋 \033[92mwttr.in/~Kilimanjaro\033[0m （只要在地點前加 ~ 就可以了）",
        "mg": "Fanatsrana vaovao: anarana toerana amin'ny fiteny maro\033[92mwttr.in/станция+Восток\033[0m (en UTF-8) sy fitadiavana toerana \033[92mwttr.in/~Kilimanjaro\033[0m (ampio ~ fotsiny eo aloha)",
        "ta": "புதிய அம்சம்: பன்மொழி இருப்பிடப் பெயர்கள் \033[92mwttr.in/станция+Восток\033[0m (UTF-8 இல்) மற்றும் இருப்பிடத் தேடல் \033[92mwttr.in/~Kilimanjaro\033[0m (முன் ~ஐச் சேர்க்கவும்)",
    },
    "FOLLOW_ME": {
        "en": "Follow \033[46m\033[30m@igor_chubin\033[0m for wttr.in updates",
        "ar": "لآخر المستجدات تابع \033[46m\033[30m@igor_chubin\033[0m",
        "af": "Volg \033[46m\033[30m@igor_chubin\033[0m vir wttr.in opdaterings",
        "am": "ለተጨማሪ wttr.in ዜና እና መረጃ \033[46m\033[30m@igor_chubin\033[0m ን ይከተሉ",
        "be": "Сачыце за \033[46m\033[30m@igor_chubin\033[0m за навінамі wttr.in",
        "bg": "Последвай \033[46m\033[30m@igor_chubin\033[0m за обновления свързани с wttr.in",
        "bn": "wttr.in আপডেটের জন্য \033[46m\033[30m@igor_chubin\033[0m কে অনুসরণ করুন",
        "bs": "XXXXXX \033[46m\033[30m@igor_chubin\033[0m XXXXXXXXXXXXXXXXXXX",
        "ca": "Segueix \033[46m\033[30m@igor_chubin\033[0m per actualitzacions de wttr.in",
        "cs": "Sledujte \033[46m\033[30m@igor_chubin\033[0m pro aktualizace pro wttr.in",
        "es": "Sigue a \033[46m\033[30m@igor_chubin\033[0m para enterarte de las novedades de wttr.in",
        "eu": "\033[46m\033[30m@igor_chubin\033[0m jarraitu wttr.in berriak jasotzeko",
        "cy": "Dilyner \033[46m\033[30m@igor_Chubin\033[0m am diweddariadau wttr.in",
        "fa": "برای دنبال کردن خبرهای wttr.in شناسه \033[46m\033[30m@igor_chubin\033[0m رو فالو کنید.",
        "fr": "Suivez \033[46m\033[30m@igor_Chubin\033[0m pour rester informé sur wttr.in",
        "de": "Folgen Sie \033[46m\033[30mhttps://x.com/igor_chubin\033[0m für wttr.in Updates",
        "ga": "Lean \033[46m\033[30m@igor_chubin\033[0m don wttr.in eolas is deanaí",
        "gl": "Segue a \033[46m\033[30m@igor_chubin\033[0m para actualizacións sobre wttr.in",
        "hi": "अपडेट के लिए फॉलो करें \033[46m\033[30m@igor_chubin\033[0m",
        "hu": "Kövesd \033[46m\033[30m@igor_chubin\033[0m-t további wttr.in információkért",
        "hy": "Նոր ֆիչռների համար հետևեք՝ \033[46m\033[30m@igor_chubin\033[0m",
        "ia": "Seque \033[46m\033[30m@igor_chubin\033[0m por nove information de wttr.in",
        "id": "Ikuti \033[46m\033[30m@igor_chubin\033[0m untuk informasi wttr.in terbaru",
        "it": "Seguite \033[46m\033[30m@igor_chubin\033[0m per aggiornamenti a wttr.in",
        "ko": "wttr.in의 업데이트 소식을 원하신다면 \033[46m\033[30m@igor_chubin\033[0m 을 팔로우 해주세요",
        "kk": "",
        "lt": "wttr.in atnaujinimus sekite \033[46m\033[30m@igor_chubin\033[0m",
        "lv": "Seko \033[46m\033[30m@igor_chubin\033[0m, lai uzzinātu jaunumus par wttr.in",
        "mk": "Следете \033[46m\033[30m@igor_chubin\033[0m за wttr.in новости",
        "mr": "wttr.in च्या अद्यावत माहितीसाठी  \033[46m\033[30m@igor_chubin\033[0m चे अनुसरण करा",
        "nb": "Følg \033[46m\033[30m@igor_chubin\033[0m for wttr.in oppdateringer",
        "nl": "Volg \033[46m\033[30m@igor_chubin\033[0m voor wttr.in updates",
        "oc": "Seguissètz \033[46m\033[30m@igor_Chubin\033[0m per demorar informat sus wttr.in",
        "fy": "Folgje \033[46m\033[30m@igor_chubin\033[0m foar wttr.in updates",
        "pl": "Śledź \033[46m\033[30m@igor_chubin\033[0m aby być na bieżąco z nowościami dotyczącymi wttr.in",
        "pt": "Seguir \033[46m\033[30m@igor_chubin\033[0m para as novidades de wttr.in",
        "pt-br": "Seguir \033[46m\033[30m@igor_chubin\033[0m para as novidades de wttr.in",
        "ro": "Urmăriți \033[46m\033[30m@igor_chubin\033[0m pentru actualizări despre wttr.in",
        "ru": "Все новые фичи публикуются здесь: \033[46m\033[30m@igor_chubin\033[0m",
        "zh": "关注 \033[46m\033[30m@igor_chubin\033[0m 获取 wttr.in 动态",
        "te": "అనుసరించండి \ 033 [46m \ 033 [30m @ igor_chubin \ 033 [wttr.in నవీకరణలను కోసం",
        "tr": "wttr.in ile ilgili gelişmeler için \033[46m\033[30m@igor_chubin\033[0m adresini takip edin",
        "th": "ติดตาม \033[46m\033[30m@igor_chubin\033[0m สำหรับการอัพเดท wttr.in",
        "da": "Følg \033[46m\033[30m@igor_chubin\033[0m for at få wttr.in opdateringer",
        "et": "Jälgi \033[46m\033[30m@igor_chubin\033[0m wttr.in uudiste tarbeks",
        "uk": "Нові можливості wttr.in публікуються тут: \033[46m\033[30m@igor_chubin\033[0m",
        "vi": "Theo dõi \033[46m\033[30m@igor_chubin\033[0m để cập nhật thông tin về wttr.in",
        "zh-tw": "追蹤 \033[46m\033[30m@igor_chubin\033[0m 以取得更多 wttr.in 的動態",
        "mg": "Araho ao ny pejy \033[46m\033[30m@igor_Chubin\033[0m raha toa ka te hahazo vaovao momban'ny wttr.in",
        "ta": "wttr.in புதுப்பிப்புகளுக்கு \033[46m\033[30m@igor_chubin\033[0m ஐப் பின்தொடரவும்",
    },
}
CAPTION = {
    "af": "Weer verslag vir:",
    "am": "የአየር ሁኔታ ሪፖርት ለ",
    "ar": "تقرير حالة الطقس",
    "az": "Hava proqnozu:",
    "be": "Прагноз надвор'я для:",
    "bg": "Прогноза за времето в:",
    "bn": "আবহাওয়ার প্রতিবেদন:",
    "bs": "Vremenske prognoze za:",
    "ca": "Informe del temps per a:",
    "cs": "Předpověď počasí pro:",
    "cy": "Adroddiad tywydd ar gyfer:",
    "da": "Vejret i:",
    "de": "Wetterbericht für:",
    "el": "Πρόγνωση καιρού για:",
    "en": "Weather report for:",
    "eo": "Veterprognozo por:",
    "es": "Pronóstico del tiempo en:",
    "et": "Ilmaprognoos:",
    "eu": "Eguraldiaren iragarpena:",
    "fa": "گزارش آب و هئا برای شما:",
    "fi": "Säätiedotus:",
    "fr": "Prévisions météo pour:",
    "fy": "Waarberjocht foar:",
    "ga": "Réamhaisnéis na haimsire do:",
    "he": ":ריוואה גזמ תיזחת",
    "hi": "मौसम की जानकारी",
    "hr": "Vremenska prognoza za:",
    "hu": "Időjárás előrejelzés:",
    "hy": "Եղանակի տեսություն:",
    "ia": "Reporto tempore por:",
    "id": "Prakiraan cuaca:",
    "it": "Previsioni meteo:",
    "is": "Veðurskýrsla fyrir:",
    "ja": "天気予報：",
    "jv": "Weather forecast for:",
    "ka": "ამინდის პროგნოზი:",
    "kk": "Ауа райы:",
    "ko": "일기 예보：",
    "ky": "Аба ырайы:",
    "lt": "Orų prognozė:",
    "lv": "Laika ziņas:",
    "mk": "Прогноза за времето во:",
    "ml": "കാലാവസ്ഥ റിപ്പോർട്ട്:",
    "mr": "हवामान अहवाल:",
    "nb": "Værmelding for:",
    "nl": "Weerbericht voor:",
    "nn": "Vêrmelding for:",
    "oc": "Previsions metèo per :",
    "pl": "Pogoda w:",
    "pt": "Previsão do tempo para:",
    "pt-br": "Previsão do tempo para:",
    "ro": "Prognoza meteo pentru:",
    "ru": "Прогноз погоды:",
    "sk": "Predpoveď počasia pre:",
    "sl": "Vremenska napoved za",
    "sr": "Временска прогноза за:",
    "sr-lat": "Vremenska prognoza za:",
    "sv": "Väderleksprognos för:",
    "sw": "Ripoti ya hali ya hewa, jiji la:",
    "te": "వాతావరణ సమాచారము:",
    "th": "รายงานสภาพอากาศ:",
    "tr": "Hava beklentisi:",
    "uk": "Прогноз погоди для:",
    "uz": "Ob-havo bashorati:",
    "vi": "Báo cáo thời tiết:",
    "zh": "天气预报：",
    "zu": "Isimo sezulu:",
    "zh-tw": "天氣報告：",
    "mg": "Toetr'andro any :",
    "ta": "வானிலை அறிக்கை:",
}


def get_message(message_name, lang):
    if lang == "zh-cn":
        lang = "zh"
    if message_name not in MESSAGE:
        return ""
    message_dict = MESSAGE[message_name]
    return message_dict.get(lang, message_dict.get("en", ""))

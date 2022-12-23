package v1

func getIcon(name string) []string {
	icon := map[string][]string{
		"iconUnknown": {
			"    .-.      ",
			"     __)     ",
			"    (        ",
			"     `-’     ",
			"      •      ",
		},

		"iconSunny": {
			"\033[38;5;226m    \\   /    \033[0m",
			"\033[38;5;226m     .-.     \033[0m",
			"\033[38;5;226m  ― (   ) ―  \033[0m",
			"\033[38;5;226m     `-’     \033[0m",
			"\033[38;5;226m    /   \\    \033[0m",
		},

		"iconPartlyCloudy": {
			"\033[38;5;226m   \\  /\033[0m      ",
			"\033[38;5;226m _ /\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m   \\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"             ",
		},

		"iconCloudy": {
			"             ",
			"\033[38;5;250m     .--.    \033[0m",
			"\033[38;5;250m  .-(    ).  \033[0m",
			"\033[38;5;250m (___.__)__) \033[0m",
			"             ",
		},

		"iconVeryCloudy": {
			"             ",
			"\033[38;5;240;1m     .--.    \033[0m",
			"\033[38;5;240;1m  .-(    ).  \033[0m",
			"\033[38;5;240;1m (___.__)__) \033[0m",
			"             ",
		},

		"iconLightShowers": {
			"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"\033[38;5;111m     ‘ ‘ ‘ ‘ \033[0m",
			"\033[38;5;111m    ‘ ‘ ‘ ‘  \033[0m",
		},

		"iconHeavyShowers": {
			"\033[38;5;226m _`/\"\"\033[38;5;240;1m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;240;1m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;240;1m(___(__) \033[0m",
			"\033[38;5;21;1m   ‚‘‚‘‚‘‚‘  \033[0m",
			"\033[38;5;21;1m   ‚’‚’‚’‚’  \033[0m",
		},

		"iconLightSnowShowers": {
			"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"\033[38;5;255m     *  *  * \033[0m",
			"\033[38;5;255m    *  *  *  \033[0m",
		},

		"iconHeavySnowShowers": {
			"\033[38;5;226m _`/\"\"\033[38;5;240;1m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;240;1m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;240;1m(___(__) \033[0m",
			"\033[38;5;255;1m    * * * *  \033[0m",
			"\033[38;5;255;1m   * * * *   \033[0m",
		},

		"iconLightSleetShowers": {
			"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"\033[38;5;111m     ‘ \033[38;5;255m*\033[38;5;111m ‘ \033[38;5;255m* \033[0m",
			"\033[38;5;255m    *\033[38;5;111m ‘ \033[38;5;255m*\033[38;5;111m ‘  \033[0m",
		},

		"iconThunderyShowers": {
			"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"\033[38;5;228;5m    ⚡\033[38;5;111;25m‘‘\033[38;5;228;5m⚡\033[38;5;111;25m‘‘ \033[0m",
			"\033[38;5;111m    ‘ ‘ ‘ ‘  \033[0m",
		},

		"iconThunderyHeavyRain": {
			"\033[38;5;240;1m     .-.     \033[0m",
			"\033[38;5;240;1m    (   ).   \033[0m",
			"\033[38;5;240;1m   (___(__)  \033[0m",
			"\033[38;5;21;1m  ‚‘\033[38;5;228;5m⚡\033[38;5;21;25m‘‚\033[38;5;228;5m⚡\033[38;5;21;25m‚‘ \033[0m",
			"\033[38;5;21;1m  ‚’‚’\033[38;5;228;5m⚡\033[38;5;21;25m’‚’  \033[0m",
		},

		"iconThunderySnowShowers": {
			"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"\033[38;5;255m     *\033[38;5;228;5m⚡\033[38;5;255;25m*\033[38;5;228;5m⚡\033[38;5;255;25m* \033[0m",
			"\033[38;5;255m    *  *  *  \033[0m",
		},

		"iconLightRain": {
			"\033[38;5;250m     .-.     \033[0m",
			"\033[38;5;250m    (   ).   \033[0m",
			"\033[38;5;250m   (___(__)  \033[0m",
			"\033[38;5;111m    ‘ ‘ ‘ ‘  \033[0m",
			"\033[38;5;111m   ‘ ‘ ‘ ‘   \033[0m",
		},

		"iconHeavyRain": {
			"\033[38;5;240;1m     .-.     \033[0m",
			"\033[38;5;240;1m    (   ).   \033[0m",
			"\033[38;5;240;1m   (___(__)  \033[0m",
			"\033[38;5;21;1m  ‚‘‚‘‚‘‚‘   \033[0m",
			"\033[38;5;21;1m  ‚’‚’‚’‚’   \033[0m",
		},

		"iconLightSnow": {
			"\033[38;5;250m     .-.     \033[0m",
			"\033[38;5;250m    (   ).   \033[0m",
			"\033[38;5;250m   (___(__)  \033[0m",
			"\033[38;5;255m    *  *  *  \033[0m",
			"\033[38;5;255m   *  *  *   \033[0m",
		},

		"iconHeavySnow": {
			"\033[38;5;240;1m     .-.     \033[0m",
			"\033[38;5;240;1m    (   ).   \033[0m",
			"\033[38;5;240;1m   (___(__)  \033[0m",
			"\033[38;5;255;1m   * * * *   \033[0m",
			"\033[38;5;255;1m  * * * *    \033[0m",
		},

		"iconLightSleet": {
			"\033[38;5;250m     .-.     \033[0m",
			"\033[38;5;250m    (   ).   \033[0m",
			"\033[38;5;250m   (___(__)  \033[0m",
			"\033[38;5;111m    ‘ \033[38;5;255m*\033[38;5;111m ‘ \033[38;5;255m*  \033[0m",
			"\033[38;5;255m   *\033[38;5;111m ‘ \033[38;5;255m*\033[38;5;111m ‘   \033[0m",
		},

		"iconFog": {
			"             ",
			"\033[38;5;251m _ - _ - _ - \033[0m",
			"\033[38;5;251m  _ - _ - _  \033[0m",
			"\033[38;5;251m _ - _ - _ - \033[0m",
			"             ",
		},
	}

	return icon[name]
}

func codes() map[int][]string {
	return map[int][]string{
		113: getIcon("iconSunny"),
		116: getIcon("iconPartlyCloudy"),
		119: getIcon("iconCloudy"),
		122: getIcon("iconVeryCloudy"),
		143: getIcon("iconFog"),
		176: getIcon("iconLightShowers"),
		179: getIcon("iconLightSleetShowers"),
		182: getIcon("iconLightSleet"),
		185: getIcon("iconLightSleet"),
		200: getIcon("iconThunderyShowers"),
		227: getIcon("iconLightSnow"),
		230: getIcon("iconHeavySnow"),
		248: getIcon("iconFog"),
		260: getIcon("iconFog"),
		263: getIcon("iconLightShowers"),
		266: getIcon("iconLightRain"),
		281: getIcon("iconLightSleet"),
		284: getIcon("iconLightSleet"),
		293: getIcon("iconLightRain"),
		296: getIcon("iconLightRain"),
		299: getIcon("iconHeavyShowers"),
		302: getIcon("iconHeavyRain"),
		305: getIcon("iconHeavyShowers"),
		308: getIcon("iconHeavyRain"),
		311: getIcon("iconLightSleet"),
		314: getIcon("iconLightSleet"),
		317: getIcon("iconLightSleet"),
		320: getIcon("iconLightSnow"),
		323: getIcon("iconLightSnowShowers"),
		326: getIcon("iconLightSnowShowers"),
		329: getIcon("iconHeavySnow"),
		332: getIcon("iconHeavySnow"),
		335: getIcon("iconHeavySnowShowers"),
		338: getIcon("iconHeavySnow"),
		350: getIcon("iconLightSleet"),
		353: getIcon("iconLightShowers"),
		356: getIcon("iconHeavyShowers"),
		359: getIcon("iconHeavyRain"),
		362: getIcon("iconLightSleetShowers"),
		365: getIcon("iconLightSleetShowers"),
		368: getIcon("iconLightSnowShowers"),
		371: getIcon("iconHeavySnowShowers"),
		374: getIcon("iconLightSleetShowers"),
		377: getIcon("iconLightSleet"),
		386: getIcon("iconThunderyShowers"),
		389: getIcon("iconThunderyHeavyRain"),
		392: getIcon("iconThunderySnowShowers"),
		395: getIcon("iconHeavySnowShowers"),
	}
}

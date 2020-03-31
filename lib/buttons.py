TWITTER_BUTTON = """
<a href="https://twitter.com/igor_chubin?ref_src=twsrc%5Etfw" class="twitter-follow-button" data-show-count="false">Follow @igor_chubin</a><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
"""

GITHUB_BUTTON = """
<a aria-label="Star chubin/wttr.in on GitHub" data-count-aria-label="# stargazers on GitHub" data-count-api="/repos/chubin/wttr.in#stargazers_count" data-count-href="/chubin/wttr.in/stargazers" data-show-count="true" data-icon="octicon-star" href="https://github.com/chubin/wttr.in" class="github-button">wttr.in</a>
"""

GITHUB_BUTTON_2 = """
<!-- Place this tag where you want the button to render. -->
<a aria-label="Star schachmat/wego on GitHub" data-count-aria-label="# stargazers on GitHub" data-count-api="/repos/schachmat/wego#stargazers_count" data-count-href="/schachmat/wego/stargazers" data-show-count="true" data-icon="octicon-star" href="https://github.com/schachmat/wego" class="github-button">wego</a>
"""

GITHUB_BUTTON_3 = """
<!-- Place this tag where you want the button to render. -->
<a aria-label="Star chubin/pyphoon on GitHub" data-count-aria-label="# stargazers on GitHub" data-count-api="/repos/chubin/pyphoon#stargazers_count" data-count-href="/chubin/pyphoon/stargazers" data-show-count="true" data-icon="octicon-star" href="https://github.com/chubin/pyphoon" class="github-button">pyphoon</a>
"""

GITHUB_BUTTON_FOOTER = """
<!-- Place this tag right after the last button or just before your close body tag. -->
<script async defer id="github-bjs" src="https://buttons.github.io/buttons.js"></script>
"""

def add_buttons(output):
    """
    Add buttons to html output
    """

    return output.replace('</body>',
                          (TWITTER_BUTTON
                           + GITHUB_BUTTON
                           + GITHUB_BUTTON_3
                           + GITHUB_BUTTON_2
                           + GITHUB_BUTTON_FOOTER) + '</body>')

queries=(
    /
    /Kiev
    /Kiev.png
    /?T
    /Киев
    /Kiev?2
    "/Kiev?format=1"
    "/Kiev?format=2"
    "/Kiev?format=3"
    "/Kiev?format=4"
    "/Kiev?format=v2"
    "/:help"
    "/Kiev?T"
    "/Kiev?p"
    "/Kiev?q"
    "/Kiev?Q"
    "/Kiev_text=no_view=v2.png"
    "/Kiev.png?1nqF"
    "/Kiev_1nqF.png"
)

options=$(cat <<EOF

-A firefox
-H Accept-Language:ru
-H X-Forwarded-For:1.1.1.1
EOF
)

server="http://127.0.0.1:8002"

if [ "$1" = update ]; then
  UPDATE=yes
fi

if [[ $UPDATE = yes ]]; then
  true > test-data/signatures
fi

result_tmp=$(mktemp wttrin-test-XXXXX)

while read -r -a args
do
  for q in "${queries[@]}"; do
    signature=$(echo "${args[@]}" "$q" | sha1sum | awk '{print $1}')
    curl -ks "${args[@]}" "$server$q" > "$result_tmp"

    result=$(sha1sum "$result_tmp" | awk '{print $1}')

    # this must be moved to the server
    # but for the moment we just clean up
    # the cache after each call
    rm -rf "/wttr.in/cache"

    if grep -Eq "(we are running out of queries|500 Internal Server Error)" "$result_tmp"; then
      echo "$q"
    fi

    if [[ $UPDATE = yes ]]; then
      printf "%s %s %s\\n" "$signature" "$result" "${args[*]} $q" >> test-data/signatures
    elif ! grep -q "$signature $result" test-data/signatures; then
      echo "FAILED: curl -ks ${args[*]} $server$q"
    fi
  done
done <<< "${options}"

rm "$result_tmp"

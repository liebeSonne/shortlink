#!/bin/sh

URL=${1:-"http://localhost:8080"}

SESSION_TOKEN_1="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODA5MDcxNzgsIlVzZXJJRCI6IjhjYjhmOThiLWU2ZmMtNDc1Yy1iZTdjLTEzMmY1NDVmYTAxZiJ9.KcpYiHzbjRB0NJ5b2_NyP0eU7AyVe9uKhNagsPwYS8s"
SESSION_TOKEN_2="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODA5MDc3NTEsIlVzZXJJRCI6ImNiNzAxODk2LTdmZDItNDgyNS1iNWE5LTBjOGNhOWNjNDAyOCJ9.usRwVxaGovsjHpTMvIvOw0-MIiAejqIZvkmex_TVhXE"
SESSION_TOKEN_3="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODA5MDczMjEsIlVzZXJJRCI6IjJiYjVjOTNjLTEzNDQtNGM4Ny1hNWU2LWExYWJiNGNmMjBlYyJ9.W8i7YN3mmOT27dbX0QKB5Net7ECtGpn8cW147Ar8AIk"
SESSION_TOKEN_4="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODA5MDc0MTksIlVzZXJJRCI6ImI4NTkwZmI1LTQ4OWEtNDVkYi1hM2YxLWZlZDY0OWUwNGUxZSJ9.BlHQ8KqdZk2PU8o-fE-1w_yU0jrxgpEXkdjkkWEsZy8"
SESSION_TOKEN_5="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODA5MDcxMTQsIlVzZXJJRCI6ImU5YzdiZGE0LTMyOTMtNDI4Ni05Zjc2LTZmZWE3NmZjYjEyZiJ9.6tlpwv1yAfDTi7AmpjQgMY16hAdRO2dMqYHQxKzYZZ8"

LINK_HOST_1="https://ya.ru"
LINK_HOST_2="https://google.com"

get_random_string() {
    local length=${1:-10}
    tr -dc 'a-zA-Z0-9' < /dev/urandom | head -c "$length"
}

hey -n 100 -c 5 -m POST -d "${LINK_HOST_1}?post=$(get_random_string)" "${URL}"
hey -n 100 -c 5 -m POST -d "${LINK_HOST_1}?post=$(get_random_string)" "${URL}"

for i in $(seq 1 100); do
  hey -n 1 -c 1 -m POST -d "${LINK_HOST_1}?post=$(get_random_string)" "${URL}" &
done

hey -n 100 -c 5 -m POST -d "${LINK_HOST_1}?post=$(get_random_string)" \
  -H "Cookie: session_token=${SESSION_TOKEN_1}" \
  "${URL}"

hey -n 100 -c 5 -m POST -d "${LINK_HOST_1}?post=$(get_random_string)" \
  -H "Cookie: session_token=${SESSION_TOKEN_2}" \
  "${URL}"

for i in $(seq 1 100); do
  hey -n 1 -c 1 -m POST -d "${LINK_HOST_1}?post=$(get_random_string)" \
    -H "Cookie: session_token=${SESSION_TOKEN_3}" \
    "${URL}"
done

hey -n 100 -c 5 -m GET "${URL}/eO6YcxVN"
hey -n 100 -c 5 -m GET "${URL}/BlSbIVpT"
hey -n 100 -c 5 -m GET "${URL}/$(get_random_string)"
hey -n 100 -c 5 -m GET -H "Cookie: session_token=${SESSION_TOKEN_3}" "${URL}/sFDqs3LH"

hey -n 100 -c 5 -m POST \
  -H "Content-Type: application/json" \
  -d '{"url":"'"${LINK_HOST_2}"'?api='"$(get_random_string)"'"}' \
  "${URL}/api/shorten"

for i in $(seq 1 100); do
  hey -n 1 -c 1 -m POST \
    -H "Content-Type: application/json" \
    -d '{"url":"'"${LINK_HOST_2}"'?api='"$(get_random_string)"'"}' \
    "${URL}/api/shorten"
done

hey -n 100 -c 5 -m POST \
  -H "Content-Type: application/json" \
  -d '{"url":"'"${LINK_HOST_2}"'?api='"$(get_random_string)"'"}' \
  "${URL}/api/shorten"

hey -n 100 -c 5 -m POST \
  -H "Content-Type: application/json" \
  -d '{"url":"'"${LINK_HOST_2}"'?api='"$(get_random_string)"'"}' \
  -H "Cookie: session_token=${SESSION_TOKEN_4}" \
  "${URL}/api/shorten"

for i in $(seq 1 100); do
  hey -n 10 -c 3 -m POST \
    -H "Content-Type: application/json" \
    -H "Accept-Encoding: deflate" \
    -d '{"url":"'"${LINK_HOST_1}"'?deflate='"$(get_random_string)"'"}' \
    "${URL}/api/shorten"
done

hey -n 100 -c 5 -m GET "${URL}/ping"

hey -n 100 -c 5 -m POST \
  -H "Content-Type: application/json" \
  -d '[{"correlation_id": "c1", "original_url": "'"${LINK_HOST_1}"'?batch='"$(get_random_string)"'"}, {"correlation_id": "с2", "original_url": "'"${LINK_HOST_2}"'?batch='"$(get_random_string)"'"}]' \
  "${URL}/api/shorten/batch"

for i in $(seq 1 50); do
  hey -n 6 -c 3 -m POST \
    -H "Content-Type: application/json" \
    -d '[{"correlation_id": "c1", "original_url": "'"${LINK_HOST_1}"'?batch='"$(get_random_string)"'"}, {"correlation_id": "с2", "original_url": "'"${LINK_HOST_2}"'?batch='"$(get_random_string)"'"}]' \
    "${URL}/api/shorten/batch"
done

hey -n 100 -c 5 -m GET -H "Content-Type: application/json" "${URL}/api/user/urls"

hey -n 100 -c 5 -m GET \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=${SESSION_TOKEN_2}" \
  "${URL}/api/user/urls"

hey -n 100 -c 5 -m GET \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=${SESSION_TOKEN_5}" \
  "${URL}/api/user/urls"

hey -n 100 -c 5 -m GET \
  -H "Content-Type: application/json" \
  -H "Cookie: session_token=${SESSION_TOKEN_1}" \
  "${URL}/api/user/urls"


#!/usr/bin/env bash
set -euo pipefail

HOST="${SONAR_HOST_URL:-http://localhost:9000}"
ENV_FILE="${SONAR_ENV_FILE:-.env.sonar}"
ADMIN_USER="${SONAR_ADMIN_USER:-admin}"
OLD_PASS="${SONAR_ADMIN_PASSWORD_OLD:-admin}"
NEW_PASS="${SONAR_ADMIN_PASSWORD:-Admin123!}"
TOKEN_NAME="${SONAR_TOKEN_NAME:-local-dev}"

wait_up() {
  echo "Waiting for SonarQube at ${HOST} ..."
  until curl -sf "${HOST}/api/system/status" | grep -q '"status":"UP"'; do
    sleep 3
  done
  echo "SonarQube is UP"
}

ensure_password() {
  # Fresh install still has admin/admin; newer SonarQube requires a password change.
  if curl -sf -u "${ADMIN_USER}:${OLD_PASS}" "${HOST}/api/authentication/validate" | grep -q '"valid":true'; then
    echo "Changing default admin password ..."
    curl -sf -u "${ADMIN_USER}:${OLD_PASS}" -X POST \
      "${HOST}/api/users/change_password" \
      --data-urlencode "login=${ADMIN_USER}" \
      --data-urlencode "previousPassword=${OLD_PASS}" \
      --data-urlencode "password=${NEW_PASS}" >/dev/null || true
  fi
}

generate_token() {
  local pass="${NEW_PASS}"
  if ! curl -sf -u "${ADMIN_USER}:${pass}" "${HOST}/api/authentication/validate" | grep -q '"valid":true'; then
    pass="${OLD_PASS}"
  fi

  # Re-create token with a stable name for local use.
  curl -sf -u "${ADMIN_USER}:${pass}" -X POST \
    "${HOST}/api/user_tokens/revoke" \
    --data-urlencode "name=${TOKEN_NAME}" >/dev/null 2>&1 || true

  local response
  response="$(curl -sf -u "${ADMIN_USER}:${pass}" -X POST \
    "${HOST}/api/user_tokens/generate" \
    --data-urlencode "name=${TOKEN_NAME}")"

  local token
  token="$(printf '%s' "${response}" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')"
  if [[ -z "${token}" ]]; then
    echo "Failed to generate token. Response: ${response}" >&2
    exit 1
  fi

  cat >"${ENV_FILE}" <<EOF
SONAR_HOST_URL=${HOST}
SONAR_TOKEN=${token}
SONAR_ADMIN_USER=${ADMIN_USER}
SONAR_ADMIN_PASSWORD=${pass}
EOF
  chmod 600 "${ENV_FILE}"
  echo "Saved free local token to ${ENV_FILE}"
  echo "UI: ${HOST}  (user: ${ADMIN_USER})"
}

wait_up
ensure_password
generate_token

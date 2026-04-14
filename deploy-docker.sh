#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

log() {
  printf '[INFO] %s\n' "$*"
}

warn() {
  printf '[WARN] %s\n' "$*" >&2
}

die() {
  printf '[ERROR] %s\n' "$*" >&2
  exit 1
}

have_cmd() {
  command -v "$1" >/dev/null 2>&1
}

run_root() {
  if [ "${EUID:-$(id -u)}" -eq 0 ]; then
    "$@"
  elif have_cmd sudo; then
    sudo "$@"
  else
    die "Need root privileges. Re-run as root or install sudo."
  fi
}

ensure_env_file() {
  if [ ! -f .env ]; then
    if [ -f .env.example ]; then
      log "Creating .env from .env.example"
      cp .env.example .env
      warn "Edit .env before deploying, especially passwords and tokens."
    else
      die ".env.example not found"
    fi
  fi
}

install_docker_ubuntu() {
  log "Installing Docker Engine for Debian/Ubuntu"
  run_root apt-get update
  run_root apt-get install -y ca-certificates curl gnupg
  run_root install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/$(. /etc/os-release && echo "$ID")/gpg | run_root gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  run_root chmod a+r /etc/apt/keyrings/docker.gpg
  . /etc/os-release
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$ID $VERSION_CODENAME stable" | run_root tee /etc/apt/sources.list.d/docker.list >/dev/null
  run_root apt-get update
  run_root apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
  run_root systemctl enable --now docker
}

install_docker_rhel() {
  log "Installing Docker Engine for RHEL/CentOS/Rocky/AlmaLinux"
  if have_cmd dnf; then
    run_root dnf -y install dnf-plugins-core
    run_root dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    run_root dnf -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    run_root systemctl enable --now docker
  elif have_cmd yum; then
    run_root yum -y install yum-utils
    run_root yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    run_root yum -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    run_root systemctl enable --now docker
  else
    die "No supported package manager found for Docker installation"
  fi
}

install_docker_fallback() {
  warn "Falling back to the official Docker install script"
  curl -fsSL https://get.docker.com | sh
  run_root systemctl enable --now docker
}

ensure_docker() {
  if have_cmd docker; then
    log "Docker already installed"
    return
  fi

  if [ -r /etc/os-release ]; then
    . /etc/os-release
    case "${ID:-}" in
      ubuntu|debian)
        install_docker_ubuntu
        ;;
      centos|rhel|rocky|almalinux|fedora)
        install_docker_rhel
        ;;
      *)
        install_docker_fallback
        ;;
    esac
  else
    install_docker_fallback
  fi
}

ensure_compose() {
  if docker compose version >/dev/null 2>&1; then
    return
  fi

  if have_cmd docker-compose; then
    return
  fi

  die "Docker Compose is not available"
}

build_with_retry() {
  # Extend API timeouts for slow/unstable servers.
  export DOCKER_CLIENT_TIMEOUT="${DOCKER_CLIENT_TIMEOUT:-600}"
  export COMPOSE_HTTP_TIMEOUT="${COMPOSE_HTTP_TIMEOUT:-600}"
  local build_timeout="${BUILD_TIMEOUT_SECONDS:-1800}"
  local compose_parallel_limit="${COMPOSE_PARALLEL_LIMIT:-2}"

  if [ "$compose_parallel_limit" -lt 2 ]; then
    compose_parallel_limit=2
  fi

  run_build() {
    if have_cmd timeout; then
      timeout --foreground "${build_timeout}s" "$@"
    else
      "$@"
    fi
  }

  if docker compose version >/dev/null 2>&1; then
    log "Building backend (safe mode: non-BuildKit, serialized)"
     if DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 COMPOSE_PARALLEL_LIMIT="$compose_parallel_limit" run_build docker compose build backend && \
       DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 COMPOSE_PARALLEL_LIMIT="$compose_parallel_limit" run_build docker compose build frontend; then
      return
    fi

    warn "Safe mode failed. Retrying with BuildKit serialized mode"
     if COMPOSE_PARALLEL_LIMIT="$compose_parallel_limit" run_build docker compose --progress=plain build backend && \
       COMPOSE_PARALLEL_LIMIT="$compose_parallel_limit" run_build docker compose --progress=plain build frontend; then
      return
    fi

    warn "Retrying full compose build as a final attempt"
    if run_build docker compose --progress=plain build; then
      return
    fi
  else
    log "Building images (docker-compose v1 mode)"
    if DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 COMPOSE_PARALLEL_LIMIT="$compose_parallel_limit" run_build docker-compose build backend && \
       DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 COMPOSE_PARALLEL_LIMIT="$compose_parallel_limit" run_build docker-compose build frontend; then
      return
    fi

    warn "Safe mode failed. Retrying full docker-compose build"
    if run_build docker-compose build; then
      return
    fi
  fi

  die "Image build failed after retries"
}

start_stack() {
  log "Starting GoPan stack"
  build_with_retry

  if docker compose version >/dev/null 2>&1; then
    docker compose up -d
  else
    docker-compose up -d
  fi
}

main() {
  ensure_env_file
  ensure_docker
  ensure_compose
  start_stack
  log "Deployment complete"
}

main "$@"
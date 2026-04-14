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
    COMPOSE_IMPL="v2"
    return
  fi

  if have_cmd docker-compose; then
    COMPOSE_IMPL="v1"
    return
  fi

  die "Docker Compose is not available"
}

compose() {
  if [ "${COMPOSE_IMPL:-}" = "v2" ]; then
    docker compose "$@"
  else
    docker-compose "$@"
  fi
}

log_compose_impl() {
  if [ "${COMPOSE_IMPL:-}" = "v2" ]; then
    log "Using Docker Compose v2 (docker compose)"
  else
    warn "Using Docker Compose v1 (docker-compose). Consider upgrading to v2 for better compatibility."
  fi
}

normalize_build_targets() {
  if [ "$#" -eq 0 ]; then
    BUILD_TARGETS=(backend frontend)
  else
    BUILD_TARGETS=("$@")
  fi
}

build_with_retry() {
  normalize_build_targets "$@"

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

  build_target_serialized() {
    local target
    for target in "${BUILD_TARGETS[@]}"; do
      log "Building ${target}"
      if ! DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 COMPOSE_PARALLEL_LIMIT="$compose_parallel_limit" run_build compose build "$target"; then
        return 1
      fi
    done
    return 0
  }

  if build_target_serialized; then
    return
  fi

  warn "Safe mode failed. Retrying with standard compose build"
  if run_build compose build "${BUILD_TARGETS[@]}"; then
    return
  fi

  die "Image build failed after retries"
}

up_stack() {
  if [ "$#" -eq 0 ]; then
    compose up -d
  else
    compose up -d "$@"
  fi
}

deploy_all() {
  log "Deploying all services"
  build_with_retry backend frontend
  up_stack
  maybe_bootstrap_default_admin
}

rebuild_frontend() {
  log "Rebuilding frontend service"
  build_with_retry frontend
  compose up -d --no-deps --force-recreate frontend
}

rebuild_backend() {
  log "Rebuilding backend service"
  build_with_retry backend
  compose up -d --no-deps --force-recreate backend
}

start_existing_stack() {
  log "Starting existing stack without rebuild"
  compose up -d --no-build
}

stop_stack() {
  log "Stopping stack"
  compose down
}

restart_stack() {
  log "Restarting running services"
  compose restart
}

show_status() {
  compose ps
}

show_logs() {
  local service="${1:-}"
  if [ -n "$service" ]; then
    compose logs --tail=200 "$service"
  else
    compose logs --tail=200
  fi
}

cleanup_dangling_images() {
  log "Cleaning dangling images"
  docker image prune -f
}

current_git_branch() {
  git rev-parse --abbrev-ref HEAD 2>/dev/null || true
}

update_from_git() {
  if ! have_cmd git; then
    die "git is required for update action"
  fi

  local branch="${1:-}"
  if [ -z "$branch" ]; then
    branch="$(current_git_branch)"
  fi
  if [ -z "$branch" ]; then
    die "Unable to detect git branch"
  fi

  log "Updating repository from origin/${branch}"
  git fetch origin "$branch"
  git pull --rebase origin "$branch"
}

insert_default_admin_user() {
  local mysql_container="gopan-mysql"

  if ! docker ps --format '{{.Names}}' | grep -qx "$mysql_container"; then
    warn "MySQL container '$mysql_container' is not running, skip admin bootstrap"
    return
  fi

  local exists
  exists="$(docker exec "$mysql_container" sh -lc 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE" -N -B -e "SELECT COUNT(*) FROM user_basic WHERE name=\"admin\" OR email=\"admin@linux.com\";"' 2>/dev/null || true)"

  if [ -z "$exists" ]; then
    warn "Unable to check existing admin user, skip admin bootstrap"
    return
  fi

  if [ "$exists" != "0" ]; then
    warn "User 'admin' or email 'admin@linux.com' already exists, skip admin bootstrap"
    return
  fi

  if docker exec -i "$mysql_container" sh -lc 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' <<'SQL'
INSERT INTO user_basic
  (identity, name, password, email, status, role, upload_permission, download_permission, share_permission, created_at, updated_at)
VALUES
  (UUID(), 'admin', '$2a$10$N/FJ.Oyey.ak/vxkmdVlXOhpRejoyvaWtd.IDAs6spu735/5tN.Va', 'admin@linux.com', 1, 2, 1, 1, 1, NOW(), NOW());
SQL
  then
    log "Default admin created: name=admin, email=admin@linux.com, password=123456"
    warn "Please change the default admin password immediately after first login"
  else
    warn "Failed to insert default admin user"
  fi
}

maybe_bootstrap_default_admin() {
  if [ ! -t 0 ]; then
    log "Non-interactive shell detected, skip default admin prompt"
    return
  fi

  printf '[INFO] Create default admin user (admin/admin@linux.com, password 123456)? [y/N]: '
  read -r answer
  case "$answer" in
    y|Y|yes|YES)
      insert_default_admin_user
      ;;
    *)
      log "Skip default admin bootstrap"
      ;;
  esac
}

print_menu() {
  cat <<'EOF'

================ GoPan Deploy Menu ================
1) Full deploy (build backend+frontend and up)
2) Rebuild frontend only
3) Rebuild backend only
4) Start stack (no build)
5) Restart stack
6) Stop stack
7) Show status
8) Show logs
9) Update from git and full deploy
10) Create default admin user
11) Clean dangling images
0) Exit
===================================================
EOF
}

interactive_menu() {
  while true; do
    print_menu
    printf '[INFO] Select an option: '
    read -r choice

    case "$choice" in
      1)
        deploy_all
        ;;
      2)
        rebuild_frontend
        ;;
      3)
        rebuild_backend
        ;;
      4)
        start_existing_stack
        ;;
      5)
        restart_stack
        ;;
      6)
        stop_stack
        ;;
      7)
        show_status
        ;;
      8)
        printf '[INFO] Service name (backend/frontend/mysql/redis, empty=all): '
        read -r service
        show_logs "$service"
        ;;
      9)
        printf '[INFO] Branch name (empty=current): '
        read -r branch
        update_from_git "$branch"
        deploy_all
        ;;
      10)
        insert_default_admin_user
        ;;
      11)
        cleanup_dangling_images
        ;;
      0)
        log "Bye"
        return
        ;;
      *)
        warn "Invalid option: $choice"
        ;;
    esac
  done
}

run_action() {
  local action="${1:-deploy}"
  local param="${2:-}"

  case "$action" in
    deploy)
      deploy_all
      ;;
    rebuild-frontend)
      rebuild_frontend
      ;;
    rebuild-backend)
      rebuild_backend
      ;;
    start)
      start_existing_stack
      ;;
    restart)
      restart_stack
      ;;
    stop)
      stop_stack
      ;;
    status)
      show_status
      ;;
    logs)
      show_logs "$param"
      ;;
    update)
      update_from_git "$param"
      deploy_all
      ;;
    create-admin)
      insert_default_admin_user
      ;;
    clean)
      cleanup_dangling_images
      ;;
    menu)
      interactive_menu
      ;;
    *)
      die "Unknown action: $action"
      ;;
  esac
}

main() {
  local action="${1:-deploy}"
  local param="${2:-}"

  ensure_env_file
  ensure_docker
  ensure_compose
  log_compose_impl

  if [ "$action" = "menu" ] || { [ -t 0 ] && [ -z "${1:-}" ]; }; then
    interactive_menu
    return
  fi

  run_action "$action" "$param"
  log "Deployment complete"
}

main "$@"
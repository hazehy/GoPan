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
    die "需要 root 权限，请使用 root 运行或先安装 sudo。"
  fi
}

ensure_env_file() {
  if [ ! -f .env ]; then
    if [ -f .env.example ]; then
      log "未检测到 .env，正在从 .env.example 创建"
      cp .env.example .env
      warn "部署前请先编辑 .env，重点修改密码和令牌等敏感配置。"
    else
      die "未找到 .env.example"
    fi
  fi
}

install_docker_ubuntu() {
  log "正在为 Debian/Ubuntu 安装 Docker Engine"
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
  log "正在为 RHEL/CentOS/Rocky/AlmaLinux 安装 Docker Engine"
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
    die "未找到可用的软件包管理器，无法安装 Docker"
  fi
}

install_docker_fallback() {
  warn "将使用 Docker 官方安装脚本作为兜底方案"
  curl -fsSL https://get.docker.com | sh
  run_root systemctl enable --now docker
}

ensure_docker() {
  if have_cmd docker; then
    log "Docker 已安装"
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

  die "未检测到 Docker Compose"
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
    log "当前使用 Docker Compose v2 (docker compose)"
  else
    warn "当前使用 Docker Compose v1 (docker-compose)，建议升级到 v2 以获得更好兼容性。"
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
      log "正在构建 ${target}"
      if ! DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 COMPOSE_PARALLEL_LIMIT="$compose_parallel_limit" run_build compose build "$target"; then
        return 1
      fi
    done
    return 0
  }

  if build_target_serialized; then
    return
  fi

  warn "安全模式构建失败，正在回退到标准 compose 构建"
  if run_build compose build "${BUILD_TARGETS[@]}"; then
    return
  fi

  die "镜像构建重试后仍失败"
}

up_stack() {
  if [ "$#" -eq 0 ]; then
    compose up -d
  else
    compose up -d "$@"
  fi
}

deploy_all() {
  log "正在执行全量部署"
  build_with_retry backend frontend
  up_stack
  maybe_bootstrap_default_admin
}

rebuild_frontend() {
  log "正在重建前端服务"
  build_with_retry frontend
  compose up -d --no-deps --force-recreate frontend
}

rebuild_backend() {
  log "正在重建后端服务"
  build_with_retry backend
  compose up -d --no-deps --force-recreate backend
}

start_existing_stack() {
  log "正在启动现有服务（不重建）"
  compose up -d --no-build
}

stop_stack() {
  log "正在停止服务栈"
  compose down
}

restart_stack() {
  log "正在重启运行中的服务"
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
  log "正在清理悬空镜像"
  docker image prune -f
}

current_git_branch() {
  git rev-parse --abbrev-ref HEAD 2>/dev/null || true
}

update_from_git() {
  if ! have_cmd git; then
    die "执行更新操作需要先安装 git"
  fi

  local branch="${1:-}"
  if [ -z "$branch" ]; then
    branch="$(current_git_branch)"
  fi
  if [ -z "$branch" ]; then
    die "无法识别当前 Git 分支"
  fi

  log "正在从 origin/${branch} 更新仓库"
  git fetch origin "$branch"
  git pull --rebase origin "$branch"
}

insert_default_admin_user() {
  local mysql_container="gopan-mysql"

  if ! docker ps --format '{{.Names}}' | grep -qx "$mysql_container"; then
    warn "MySQL 容器 '$mysql_container' 未运行，跳过默认管理员初始化"
    return
  fi

  local exists
  exists="$(docker exec "$mysql_container" sh -lc 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE" -N -B -e "SELECT COUNT(*) FROM user_basic WHERE name=\"admin\" OR email=\"admin@linux.com\";"' 2>/dev/null || true)"

  if [ -z "$exists" ]; then
    warn "无法检查管理员是否已存在，跳过默认管理员初始化"
    return
  fi

  if [ "$exists" != "0" ]; then
    warn "admin 用户或 admin@linux.com 邮箱已存在，跳过默认管理员初始化"
    return
  fi

  if docker exec -i "$mysql_container" sh -lc 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' <<'SQL'
INSERT INTO user_basic
  (identity, name, password, email, status, role, upload_permission, download_permission, share_permission, created_at, updated_at)
VALUES
  (UUID(), 'admin', '$2a$10$N/FJ.Oyey.ak/vxkmdVlXOhpRejoyvaWtd.IDAs6spu735/5tN.Va', 'admin@linux.com', 1, 2, 1, 1, 1, NOW(), NOW());
SQL
  then
    log "默认管理员创建成功：用户名 admin，邮箱 admin@linux.com，密码 123456"
    warn "请在首次登录后立即修改默认管理员密码"
  else
    warn "默认管理员创建失败"
  fi
}

maybe_bootstrap_default_admin() {
  if [ ! -t 0 ]; then
    log "检测到非交互环境，跳过默认管理员询问"
    return
  fi

  printf '[INFO] 是否创建默认管理员（admin/admin@linux.com，密码 123456）？[y/N]: '
  read -r answer
  case "$answer" in
    y|Y|yes|YES)
      insert_default_admin_user
      ;;
    *)
      log "已跳过默认管理员初始化"
      ;;
  esac
}

print_menu() {
  cat <<'EOF'

================ GoPan 部署菜单 ================
1) 全量部署（重建后端+前端并启动）
2) 仅重建前端
3) 仅重建后端
4) 启动服务（不重建）
5) 重启服务
6) 停止服务
7) 查看状态
8) 查看日志
9) Git 更新并执行全量部署
10) 创建默认管理员
11) 清理悬空镜像
0) 退出
================================================
EOF
}

interactive_menu() {
  while true; do
    print_menu
    printf '[INFO] 请选择操作编号: '
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
        printf '[INFO] 请输入服务名（backend/frontend/mysql/redis，留空=全部）: '
        read -r service
        show_logs "$service"
        ;;
      9)
        printf '[INFO] 请输入分支名（留空=当前分支）: '
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
        log "已退出"
        return
        ;;
      *)
        warn "无效选项: $choice"
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
      die "未知操作: $action"
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
# Inspired by the Rust programming language installer (https://sh.rustup.rs).
# Some function was coppied from there without changing.

set -u

GOGREP_URL="https://github.com/kgantsov/gogrep/releases/download/v0.1/gogrep"

usage() {
    cat 1>&2 <<EOF
gogrep 0.1
The installer for gogrep
EOF
}

main() {
    need_cmd uname
    need_cmd curl
    need_cmd chmod
    need_cmd mkdir
    need_cmd rm

    get_architecture || return 1

    local _arch="$RETVAL"
    assert_nz "$_arch" "arch"

    local _ext=""
    case "$_arch" in
        *windows*)
            _ext=".exe"
            ;;
    esac

    local _url="$GOGREP_URL-$_arch"

    local _dir="/usr/local/bin"
    local _file="$_dir/gogrep-$_arch"

    echo "Going to download gogrep for $_arch"

    ensure mkdir -p "$_dir"
    ensure curl -sSfL "$_url" -o "$_file"
    ensure mv "$_file" "$_dir/gogrep"

    local _file="$_dir/gogrep"
    ensure chmod 0755 "$_file"

    ignore "$_file" "$@"

    local _retval=$?
    return "$_retval"
}


get_bitness() {
    need_cmd head
    # Architecture detection without dependencies beyond coreutils.
    # ELF files start out "\x7fELF", and the following byte is
    #   0x01 for 32-bit and
    #   0x02 for 64-bit.
    # The printf builtin on some shells like dash only supports octal
    # escape sequences, so we use those.
    local _current_exe_head=$(head -c 5 /proc/self/exe )
    if [ "$_current_exe_head" = "$(printf '\177ELF\001')" ]; then
        echo 32
    elif [ "$_current_exe_head" = "$(printf '\177ELF\002')" ]; then
        echo 64
    else
        err "unknown platform bitness"
    fi
}


get_architecture() {
    local _ostype="$(uname -s)"
    local _cputype="$(uname -m)"

    if [ "$_ostype" = Darwin -a "$_cputype" = i386 ]; then
        # Darwin `uname -s` lies
        if sysctl hw.optional.x86_64 | grep -q ': 1'; then
            local _cputype=x86_64
        fi
    fi

    case "$_ostype" in
        Linux)
            local _ostype=unknown-linux-gnu
            ;;

        Darwin)
            local _ostype=apple-darwin
            ;;

        MINGW* | MSYS* | CYGWIN*)
            local _ostype=pc-windows-gnu
            ;;

        *)
            err "unrecognized OS type: $_ostype"
            ;;

    esac

    case "$_cputype" in

        i386 | i486 | i686 | i786 | x86)
            local _cputype=i386
            ;;

        aarch64)
            local _cputype=aarch64
            ;;

        x86_64 | x86-64 | x64 | amd64)
            local _cputype=x86_64
            ;;

        *)
            err "unknown CPU type: $_cputype"

    esac

    # Detect 64-bit linux with 32-bit userland
    if [ $_ostype = unknown-linux-gnu -a $_cputype = x86_64 ]; then
        if [ "$(get_bitness)" = "32" ]; then
            local _cputype=i386
        fi
    fi

    local _arch="$_cputype-$_ostype"

    RETVAL="$_arch"
}



say() {
    echo "gogrep: $1"
}

err() {
    say "$1" >&2
    exit 1
}

need_cmd() {
    if ! command -v "$1" > /dev/null 2>&1
    then err "need '$1' (command not found)"
    fi
}

need_ok() {
    if [ $? != 0 ]; then err "$1"; fi
}

assert_nz() {
    if [ -z "$1" ]; then err "assert_nz $2"; fi
}

# Run a command that should never fail. If the command fails execution
# will immediately terminate with an error showing the failing
# command.
ensure() {
    "$@"
    need_ok "command failed: $*"
}

# This is just for indicating that commands' results are being
# intentionally ignored. Usually, because it's being executed
# as part of error handling.
ignore() {
    "$@"
}

main "$@" || exit 1

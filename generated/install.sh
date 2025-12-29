#!/bin/sh
# Based in flyctl installer: Copyright 2023 flyctl authors.
# Based on Deno installer: Copyright 2019 the Deno authors. All rights reserved. MIT license.
# TODO(everyone): Keep this script simple and easily auditable.

set -e

main() {
	os=$(uname -s)
	arch=$(uname -m)

	# this is quite ugly, but this way we do not depend on the json parse, nor do we depend
	# on any particular formatting, just the value
	version=$(curl -s https://api.github.com/repos/haoli000/tttns/releases/latest | sed -n 's/.*"tag_name": *"\(v[0-9.]*\)".*/\1/p')

	if [ ! -z $1 ]; then
		version="v$1"
	fi

	release_uri="https://github.com/haoli000/tttns/releases/download/$version/tttns_${os}_${arch}.tar.gz"
	echo "Getting version $version, $release_uri"

	install_path="${CUSTOM_INSTALL:-$HOME}"

	bin_dir="$install_path/bin"
	ts=$(date +%s)
	tmp_dir="$install_path/tttns_tmp$ts"
	exe="$bin_dir/tttns"

	mkdir -p "$bin_dir"
	mkdir -p "$tmp_dir"

	function cleanup {
		rm -rf $tmp_dir
	}
	# be a good citizen and clean up after yourself
	trap cleanup EXIT

	curl -q --fail --location --progress-bar --output "$tmp_dir/tttns.tar.gz" "$release_uri"
	# extract to tmp dir so we don't open existing executable file for writing:
	tar -C "$tmp_dir" -xzf "$tmp_dir/tttns.tar.gz"
	chmod +x "$tmp_dir/tttns"
	# atomically rename into place:
	mv "$tmp_dir/tttns" "$exe"
	rm "$tmp_dir/tttns.tar.gz"

	echo "tttns was installed successfully to $exe"
	if command -v tttns >/dev/null; then
		echo "Run 'tttns help' to get started"
	else
		case $SHELL in
		/bin/zsh) shell_profile=".zshrc" ;;
		*) shell_profile=".bash_profile" ;;
		esac
		echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
		echo "  export install_path=\"$install_path\""
		echo "  export PATH=\"\$install_path/bin:\$PATH\""
		echo "Run '$exe --help' to get started"
	fi
}

main "$1"

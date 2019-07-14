echo 'Installer made with GISC'
if ! [ -x "$(command -v git)" ]; then
	 xcode-select --install>&2
	exit 1
fi

go get github.com/stretchr/testify/assert || exit 1

go build || exit 1

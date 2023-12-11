
clean:
	@rm -rf build

FMT_PATHS = ./*.go ./examples/**/*.go ./fatfs/*.go ./littlefs/*.go

fmt-check:
	@unformatted=$$(gofmt -l $(FMT_PATHS)); [ -z "$$unformatted" ] && exit 0; echo "Unformatted:"; for fn in $$unformatted; do echo "  $$fn"; done; exit 1

smoke-test:
	@mkdir -p build
	tinygo build -size short -o ./build/test.hex -target=itsybitsy-m0 ./examples/simple-fatfs/
	@md5sum ./build/test.hex
	tinygo build -size short -o ./build/test.hex -target=itsybitsy-m0 ./examples/console/fatfs/spi/
	@md5sum ./build/test.hex
	tinygo build -size short -o ./build/test.hex -target=itsybitsy-m4 ./examples/console/fatfs/qspi/
	@md5sum ./build/test.hex
	tinygo build -size short -o ./build/test.hex -target=feather-m4 ./examples/console/fatfs/sdcard/
	@md5sum ./build/test.hex
	tinygo build -size short -o ./build/test.hex -target=itsybitsy-m0 ./examples/console/littlefs/spi/
	@md5sum ./build/test.hex
	tinygo build -size short -o ./build/test.hex -target=itsybitsy-m4 ./examples/console/littlefs/qspi/
	@md5sum ./build/test.hex

test: clean fmt-check smoke-test

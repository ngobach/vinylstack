.PHONY: all clean
all: vinylstack

vinylstack:
	go build -ldflags "-s -w" ./cmd/vinylstack

clean:
	rm -rf vinylstack _dist_
all: vinylstack

.PHONY: vinylstack

vinylstack:
	go build .

clean:
	rm -rf vinylstack _dist_
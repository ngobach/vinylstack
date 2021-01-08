all: vinylstack

.PHONY: vinylstack

vinylstack:
	go build .

clean:
	rm -rf vinylstack _dist_

test:
	go run . -csn 985945

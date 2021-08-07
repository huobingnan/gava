GOC = go build
OUTDIR = ./target
TARGET = mac

build:
	$(GOC) -o $(OUTDIR)/$(TARGET)/gava
run: 
	$(OUTDIR)/$(TARGET)/gava
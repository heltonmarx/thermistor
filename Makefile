GOCC	:= go

GOFILES ?= $(TARGET).go

TARGET	:= thermistor

all:
	@echo "[CC] $(TARGET).go"
	@$(GOCC) build -o $(TARGET) $^

i386:
	@echo "[386-CC] $(TARGET).go"
	@GOOS=linux GOARCH=386 $(GOCC) build -o $(TARGET) $^

arm:
	@echo "[ARM-CC] $(TARGET).go"
	@GOOS=linux GOARCH=arm $(GOCC) build -o $(TARGET) $^

format:
	@find . -name "*.go" -exec go fmt {} \;

clean:
	@echo "clean $(TARGET)"
	@$(GOCC) clean
	@rm -rf $(TARGET)


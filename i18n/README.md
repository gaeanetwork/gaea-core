# Prepare tools(in the directory of this file)
go get -u golang.org/x/text/cmd/gotext

# Coder you error files just like channel.go

## Two points - The first is to implement the errors.go interface
    Error() string
    Code() string

## Two points - The second is that Error() function must use the Printer output
    GetDefaultPrinter().Sprintf(...)

# Generate locale files
go generate

# Copy the out.gotext.json file and translate the `translation` fileds in messages.gotext.json
cp -r locales/zh/out.gotext.json locales/zh/messages.gotext.json

# Start your app
go run main.go





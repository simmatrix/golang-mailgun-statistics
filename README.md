# GoLang Mailgun Statistics (Events API)

> This is my maiden voyage of GoLang, I have a task in my company in which I would need to pull statistics from Mailgun, I have decided to take this opportunity to write it with Go instead. Should you have any feedback or suggestions, please email to simmatrix100@gmail.com and share your thoughts and let's grow together. Thank you. :)

## Configuration

You would need to key in your Mailgun's private and public keys.

```go
var privateKey = "key-xxxxx"
var publicKey = "pubkey-xxxxx"
```

## Usage

Key in your Mailgun domains that you wish to process, you can process multiple domains at one go, thanks to Go routines.

```go
func main() {
  go process("lorem.com")
  go process("ipsum.com")
  var input string
  fmt.Scanln(&input)
}
```

## Thoughts

Through this practice, I have learned how to do the following in Go:

- How to run functions concurrently
- How to use external package (e.g. connecting to Mailgun API)
- How to deal with file system
- How to do essential things (e.g. looping, variable declaration etc.)
- How to handle errors
- How to deal with strings (e.g. replacement, conversion etc.)

## Notes

- Mailgun retains detailed data for two days for free accounts and 30 days for paid accounts.
- For the sample of returned JSON data from Mailgun, you may refer to the [official documentation](https://documentation.mailgun.com/en/latest/api-events.html#event-structure)

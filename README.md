# goatAlerts
Sends a text message of the price of a shoe on GOAT using OCR on a headless browser screenshot!

# Setup
Install Go.

[Set up a Twilio account](https://www.twilio.com/sms) (free trial gives you $15.50 in credits, or about 2000 messages to yourself)  
More reading: [Twilio API guide for Golang](https://www.twilio.com/blog/2017/09/send-text-messages-golang.html)

Clone the repository.  
`git clone https://github.com/max-jardetzky/goatAlerts`

Install dependencies.  
`go get -u github.com/otiai10/gosseract`

Create a file `config.txt` in the directory with 6 lines, according to these directions:  
```
{Name of shoe}  
{Size of shoe}  
{Your Twilio account SID}  
{Your Twilio auth token}  
{Your Twilio phone number}  
{Your phone number}  
```

Current possible shoe names (must be copied exactly):  
  - Yeezy 350 V2 Cinder NRF
  - AJ1 Obsidian
  - Yeezy 700 V3 Alvah
  - SB Dunk Low Travis Scott
  - AJ1 Travis Scott
  - Yeezy 350 V2 Cloud White NRF
  - SB Dunk Low Chunky Dunky
  - Yeezy 350 V2 Black NRF  

If you want to track a shoe that's not there, just hard code it into the map called `urls` in the function `getShoe`.
Sizes must be formatted like `#M` or `#.#M`, where # are digits of a realistic shoe size.
Phone numbers must be formatted as `+12345678900`

Example `config.txt`:
```
Yeezy 350 V2 Cinder NRF
5.5M
AC90asdf7a09e70a9st09e7t09e7t9sd0f
f89s7d89fa9s8798f798ds7fa89se7f9
+13598329855
+14358934356
```

# Usage ideas
Use `cron` or its equivalent in your OS to schedule a daily text. :)

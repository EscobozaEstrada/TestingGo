## _Desktop application binary_

Visual QA is a call quality auditing software for contact centers. 
The application is aimed to assist QA personnel by providing a recording of the agents desktop as they take calls



## Major Contact Center Solutions

- ### Ring Central
- ### Genesys 
- ### 8x8
- ### Zoho 
- ### SalesForce 
- ###  Zendesk 
- ### Twilio 
```
 TwilioPayload{Payload:TilioData{CallID:""}, Action:"connected"}
 TwilioPayload{Payload:TilioData{CallID:""}, Action:"ready"}
 TwilioPayload{Payload:TilioData{CallID:"CA9ee066be506de4c22683de9928fb62d5"}, Action:"ringing"}
 TwilioPayload{Payload:TilioData{CallID:"CA9ee066be506de4c22683de9928fb62d5"}, Action:"answer"}
unexpected end of JSON input
 TwilioPayload{Payload:TilioData{CallID:""}, Action:""}
 TwilioPayload{Payload:TilioData{CallID:"CA9ee066be506de4c22683de9928fb62d5"}, Action:"hangup"}
 TwilioPayload{Payload:TilioData{CallID:""}, Action:"connected"}
 TwilioPayload{Payload:TilioData{CallID:""}, Action:"ready"}
 TwilioPayload{Payload:TilioData{CallID:"CAf802f554888562f0fdad15c561741aeb"}, Action:"ringing"}
 TwilioPayload{Payload:TilioData{CallID:"CAf802f554888562f0fdad15c561741aeb"}, Action:"answer"}
```
- ### Platform28 
- ###  Five9
- ### Nice InContact 
- ### DialPad

## How this plays out

The idea is identify each platform websocket payload like it was done with Twilio to determine when an agent is on a call.
This should be used to create an PhoneStatus interface. 


## POC - Requirements

[x]  Application needs to be able to capture websocket packets based on [Chromium Devtools](https://chromedevtools.github.io/devtools-protocol/tot/Network/#event-webSocketFrameSent)

[x]  Application must be able to digest JSON messages and determine phone status

[ ]  Application must be tested against major contact center solutions (see below)

[ ]  Application must be able to parse SIP messages and determine phone status

[ ]  Application must address multiple browsers spawned from different places of a system

[ ]  Application must be able to determine Edge or Chrome binary path
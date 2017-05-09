# ripcord
A monitoring library to protect your infrastructure during a DOS attack.

## What's the point?

Here's a scenario.  Totally plucked from the air but scary enough that I decided to implement ripcord:

* You run a server that's exposed on the internet.
* You've opted for a price bracket that suits you and enjoy up to 500GB of transfer between you and your users.
* You get DDOS'd, which takes your bandwidth way beyond your allowance.
* It's OK!  You've configured an early-warning system, which emails you if your monthly expenditure exceeds a configured maximum.
* You're in a meeting and you feel an email come through.  You obviously don't open it because that would bring about shame upon you and your kin.
* You leave your meeting to find that your monthly expenditure is currently £78,041 and your cloud provider want blood.

## What it does

Ripcord sits on your web servers and monitors the traffic on any number of configured network interfaces.  If the number of bytes sent or received exceeds a configured maximum, you can perform a task, which might include killing your web server entirely as damage limitation.

### Todo

- [ ] Allow the user to configure command against each network interface
- [ ] Allow for in-process usage
- [x] Allow for out-of-process usage
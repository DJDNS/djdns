djdns
=====

DJDNS is a DNS server written in Go. It uses the DEJE protocol to host and validate a decentralized registry, which helps us move toward a more secure and distributed internet.

There is an earlier version of DJDNS out on the interwebs, written in Python, which is now **deprecated**. Don't use it. All the features will be ported here soon, and the most important feature (using DEJE) is not supported by python-djdns, but *is* supported by this program.

## Why DJDNS?

I'm going to assume you're already basically sold on the world needing something more transparent and secure than ICANN, which is the current "official" DNS registry. So why DJDNS, versus any of the other alternatives?

 * It can unify the altdns world, because it's designed to be really good for reverse-proxying the DNS protocol. You don't have to choose between this or Namecoin or OpenNIC or ICANN. With DJDNS, you get *everything*. By supporting everything, we give even the littlest underdog altdns a large audience who can use it instantly and painlessly.
 * No dogma imposed on anybody. DJDNS is designed to be a political and technical [archipelago](http://slatestarcodex.com/2014/06/07/archipelago-and-atomic-communitarianism/), where anyone can claim a corner of the infinite internet, and declare their own rules (or copyleft-like "unrules") for how that corner should operate.
 * Massive democracy for the pieces that everybody uses. In order to change the `<ROOT>` page, which is the starting point for every DNS resolution, many *many* trusted individuals in many countries must come to majority agreement. The intention is for it to become impractical for any government or private entity to leverage the namespacing of the internet without majority consent.
 * Trivial to include additional metadata, such as Bitcoin addresses, cryptographic pubkeys and certificates, email address, custom message, Twitter handle, sitemap, donation URL, a base64-encoded picture of a giraffe, a recipe for potato salad, per-distro installation instructions, torrent data, and more. If you can think of a way to encode it as JSON, you can include it.
 * Transparent history, like the Bitcoin blockchain.
 * Democratic consent system is *not* susceptible to Sybil attacks.
 * Event chain is *not* susceptible to differing opinions. For every page, there is a canonical content which can be deduced without trusting any of the DEJE peers (even if they all collaborate to lie to you).
 * Underlying technology is used by non-DNS stuff, which means that platform maturity and bugfixing are advanced by a much wider developer audience. It's like how mobile Linux benefits from desktop Linux, and vice versa, and they both share improvements with server Linux. Because everyone is using the same kernel, and everyone gets the benefits from everyone else, they all see a lot of improvements and robustness for their own use.

## Can I run this, right now?

Yes, but it's not really mature enough for production use. Feel free to play around with it, though.

## How to?

You need to set up a Go environment, before anything, and it's not really my job to hold your hand through that. This is not terribly user-friendly yet anyways - the current audience is developers.

```bash
go get github.com/DJDNS/djdns
go install github.com/DJDNS/djdns
djdns # Leave this running and get ready to start some more terminals
```

This will install djdns to your PATH (assuming you have your Go env and PATH set up properly), and start the server. You also need to provide the server with data via the DEJE protocol.

```bash
go install github.com/DJDNS/go-deje/djconvert
djconvert --pretty up $GOPATH/src/github.com/DJDNS/djdns/model/demo.json dns.json
go install github.com/DJDNS/go-deje/demo/router
go install github.com/DJDNS/go-deje/demo/client

# In this terminal
client --file dns.json --topic deje://localhost:8080/root

# In another terminal
router
```

Now we've got a DEJE router running, and a client with DJDNS data available on the router, subscribed to the same topic that DJDNS uses.

Finally, in one more terminal, we test it out.

```bash
dig @localhost -p 9953 ri.hype
```

You should see some activity in the djdns terminal, and the dig command (which does DNS queries) should be able to retrieve some DNS records from the djdns process. Success!

These instructions are currently rough and untested. I'd like to try it in a VM to confirm they work in a clean environment, from scratch.

You can play with this more using the browser client, which ships with go-deje, and is available as long as the router is running. Just point your browser at http://localhost:8080/ and reconnect with the correct topic (deje://localhost:8080/root). You can look at the go-deje documentation for more information, or the series of demo screencasts on DJDNS, for more info on how this works.

## What's the relationship with [CJDNS](http://cjdns.info/)?

They tackle completely different decentralization needs. The name similarity is coincidence.

    cjdns = cjd ns = Caleb James Delisle's Networking Suite.
    djdns = dj dns = DEJE DNS.

They don't even split the same way as acronyms.

That said, it is expected that the two technologies will work really well together, and that Hyperborians will want to use DJDNS as their DNS service once DJDNS is mature enough.

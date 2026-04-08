# Example Configuration
```json
{
    "debug": false,
    "address": "192.0.2.1",
    "udp_port": "53",
    "tls_port": "853",
    "upstream_address": "9.9.9.9",
    "upstream_port": "53",
    "upstream_transport": "udp4",
    "cert":
    {
        "debug": true,
        "active": true,
        "interval": 3600,
        "key": "/path/to/key.pem",
        "cert": "/path/to/cert.pem"
    },
    "matches":
    [
        {
            "qname": ".*example.org.",
            "qtype": "NS",
            "match_outgoing": false,
            "action": "action1"
        },
        {
            "qname": "example.com.",
            "qtype": "AAAA",
            "match_outgoing": true,
            "action": "action2"
        }
    ]
}
```

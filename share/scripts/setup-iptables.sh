#!/usr/bin/env bash

EXT_INTERFACE=eth0

iptables -I INPUT -p tcp -i lo -j ACCEPT
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -A INPUT -p tcp --dport 22024 -j ACCEPT
iptables -A INPUT -p tcp -j REJECT -i "$EXT_INTERFACE" --syn --reject-with tcp-reset

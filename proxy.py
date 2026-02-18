"""
Proxy helper: reads PROXY_URL from environment and applies it to
requests sessions and SMTP connections.
"""

import os
import socket
import smtplib
import ssl
from urllib.parse import urlparse
from dotenv import load_dotenv

load_dotenv()


def get_proxy_url() -> str:
    return os.environ.get("PROXY_URL", "").strip()


def apply_to_session(session) -> None:
    """Apply PROXY_URL to a requests.Session if set."""
    proxy_url = get_proxy_url()
    if proxy_url:
        session.proxies.update({"http": proxy_url, "https": proxy_url})


def smtp_ssl(host: str, port: int) -> smtplib.SMTP_SSL:
    """
    Returns an SMTP_SSL connection.
    If PROXY_URL is set and proxychains is NOT active, routes via the proxy.
    If proxychains is already intercepting (LD_PRELOAD), uses standard SMTP_SSL directly.
    """
    proxy_url = get_proxy_url()
    proxychains_active = "proxychains" in os.environ.get("LD_PRELOAD", "").lower()

    if not proxy_url or proxychains_active:
        # proxychains handles routing transparently, or no proxy needed
        return smtplib.SMTP_SSL(host, port)

    parsed = urlparse(proxy_url)
    scheme = parsed.scheme.lower()

    if scheme in ("socks5", "socks5h", "socks4"):
        import socks
        import socket as _socket
        ptype = socks.SOCKS5 if scheme.startswith("socks5") else socks.SOCKS4
        _orig = _socket.socket
        socks.set_default_proxy(ptype, parsed.hostname, parsed.port,
                                username=parsed.username, password=parsed.password)
        _socket.socket = socks.socksocket
        try:
            smtp = smtplib.SMTP_SSL(host, port)
        finally:
            _socket.socket = _orig
            socks.set_default_proxy()
        return smtp

    elif scheme in ("http", "https"):
        proxy_sock = socket.create_connection((parsed.hostname, parsed.port))
        connect_req = f"CONNECT {host}:{port} HTTP/1.1\r\nHost: {host}:{port}\r\n\r\n"
        proxy_sock.sendall(connect_req.encode())
        resp = b""
        while b"\r\n\r\n" not in resp:
            resp += proxy_sock.recv(4096)
        status_line = resp.split(b"\r\n")[0]
        if b"200" not in status_line:
            raise ConnectionError(f"Proxy CONNECT failed: {status_line}")
        ctx = ssl.create_default_context()
        ssl_sock = ctx.wrap_socket(proxy_sock, server_hostname=host)
        smtp = smtplib.SMTP_SSL.__new__(smtplib.SMTP_SSL)
        smtp._host = host
        smtp.timeout = 30
        smtp.sock = ssl_sock
        smtp.file = ssl_sock.makefile("rb")
        code, msg = smtp.getreply()
        if code != 220:
            raise smtplib.SMTPConnectError(code, msg)
        smtp.ehlo()
        return smtp

    else:
        raise ValueError(f"Unsupported proxy scheme: {scheme}")

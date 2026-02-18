"""
CHSI login module: fetches session cookies required for score queries.
"""

import os
import re
import requests
from dotenv import load_dotenv

load_dotenv()

LOGIN_PAGE_URL = (
    "https://account.chsi.com.cn/passport/login"
    "?entrytype=yzgr&service=https%3A%2F%2Fyz.chsi.com.cn%2Fj_spring_cas_security_check"
)

HEADERS = {
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
    "Accept-Language": "zh-CN,zh-HK;q=0.9,zh;q=0.8,en-US;q=0.7,en;q=0.6",
    "User-Agent": "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",
    "sec-ch-ua": '"Chromium";v="137", "Not/A)Brand";v="24"',
    "sec-ch-ua-mobile": "?1",
    "sec-ch-ua-platform": '"Android"',
}


def _extract_lt_execution(html: str) -> tuple[str, str]:
    lt = re.search(r'name="lt"\s+value="([^"]+)"', html)
    execution = re.search(r'name="execution"\s+value="([^"]+)"', html)
    if not lt or not execution:
        raise RuntimeError("Could not find lt/execution tokens in login page")
    return lt.group(1), execution.group(1)


def login() -> requests.Session:
    """
    Logs in to CHSI using credentials from environment variables:
      CHSI_USERNAME, CHSI_PASSWORD

    Returns a requests.Session with authenticated cookies.
    """
    username = os.environ["CHSI_USERNAME"]
    password = os.environ["CHSI_PASSWORD"]

    session = requests.Session()
    session.headers.update(HEADERS)

    # Step 1: GET login page to obtain lt and execution tokens
    resp = session.get(LOGIN_PAGE_URL)
    resp.raise_for_status()
    lt, execution = _extract_lt_execution(resp.text)

    # Step 2: POST credentials
    payload = {
        "username": username,
        "password": password,
        "lt": lt,
        "execution": execution,
        "_eventId": "submit",
    }
    login_resp = session.post(
        LOGIN_PAGE_URL,
        data=payload,
        headers={"Content-Type": "application/x-www-form-urlencoded",
                 "Origin": "https://account.chsi.com.cn",
                 "Referer": LOGIN_PAGE_URL},
        allow_redirects=True,
    )
    login_resp.raise_for_status()

    # Verify login succeeded by checking we were redirected to yz.chsi.com.cn
    if "account.chsi.com.cn/passport/login" in login_resp.url:
        raise RuntimeError("Login failed: still on login page. Check credentials.")

    return session


if __name__ == "__main__":
    session = login()
    print("Login successful.")
    print("Cookies:", [(c.name, c.value) for c in session.cookies])

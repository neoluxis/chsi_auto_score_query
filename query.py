"""
CHSI score query module: queries master's entrance exam scores.
Requires an authenticated session from login.py.
"""

import os
import re
from requests import Session
from dotenv import load_dotenv

load_dotenv()

QUERY_URL = "https://yz.chsi.com.cn/apply/cjcx/cjcx.do"

QUERY_HEADERS = {
    "Content-Type": "application/x-www-form-urlencoded",
    "Origin": "https://yz.chsi.com.cn",
    "Referer": "https://yz.chsi.com.cn/apply/cjcx/t/{bkdwdm}.dhtml",
    "Upgrade-Insecure-Requests": "1",
    "Sec-Fetch-Dest": "document",
    "Sec-Fetch-Mode": "navigate",
    "Sec-Fetch-Site": "same-origin",
    "Sec-Fetch-User": "?1",
}


def query_score(session: Session) -> str:
    """
    Queries the exam score using credentials from environment variables:
      QUERY_XM      - 姓名 (name)
      QUERY_ZJHM    - 证件号码 (ID card number)
      QUERY_KSBH    - 考生编号 (exam registration number)
      QUERY_BKDWDM  - 报考单位代码 (school code)

    Returns the raw HTML response text.
    """
    xm = os.environ["QUERY_XM"]
    zjhm = os.environ["QUERY_ZJHM"]
    ksbh = os.environ["QUERY_KSBH"]
    bkdwdm = os.environ["QUERY_BKDWDM"]

    payload = {
        "xm": xm,
        "zjhm": zjhm,
        "ksbh": ksbh,
        "bkdwdm": bkdwdm,
        "checkcode": "",
    }

    headers = dict(QUERY_HEADERS)
    headers["Referer"] = headers["Referer"].format(bkdwdm=bkdwdm)

    resp = session.post(QUERY_URL, data=payload, headers=headers, allow_redirects=True)
    resp.raise_for_status()
    return resp.text


def parse_score(html: str) -> dict:
    """
    Parses score information from the query result HTML.
    Returns a dict with available score fields, or raises RuntimeError if not found.
    """
    # Check for "scores not released yet" message
    if "成绩未公布" in html or "暂无成绩" in html:
        raise RuntimeError("成绩尚未公布")

    scores = {}

    # Extract score table rows
    rows = re.findall(r'<td[^>]*>(.*?)</td>', html, re.DOTALL)
    cleaned = [re.sub(r'<[^>]+>', '', cell).strip() for cell in rows]

    # Try to find total score and subject scores
    total_match = re.search(r'总\s*分[：:]\s*(\d+(?:\.\d+)?)', html)
    if total_match:
        scores["总分"] = total_match.group(1)

    return scores if scores else {"raw_cells": cleaned}

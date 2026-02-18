"""
CHSI score query module: queries master's entrance exam scores.
Requires an authenticated session from login.py.
"""

import json
import os
import re
from requests import Session
from dotenv import load_dotenv

load_dotenv()

QUERY_URL = "https://yz.chsi.com.cn/apply/cjcx/cjcx.do"


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

    referer = f"https://yz.chsi.com.cn/apply/cjcx/t/{bkdwdm}.dhtml"

    # Visit the entry page first so the server recognises the navigation flow
    session.get(referer, timeout=15)

    payload = {
        "xm": xm,
        "zjhm": zjhm,
        "ksbh": ksbh,
        "bkdwdm": bkdwdm,
        "checkcode": "",
    }

    headers = {
        "Content-Type": "application/x-www-form-urlencoded",
        "Origin": "https://yz.chsi.com.cn",
        "Referer": referer,
        "Upgrade-Insecure-Requests": "1",
        "Sec-Fetch-Dest": "document",
        "Sec-Fetch-Mode": "navigate",
        "Sec-Fetch-Site": "same-origin",
        "Sec-Fetch-User": "?1",
    }

    resp = session.post(QUERY_URL, data=payload, headers=headers, allow_redirects=True, timeout=15)
    resp.raise_for_status()
    return resp.text


def parse_score(html: str) -> tuple[dict, str]:
    """
    Parses score information from the query result HTML.
    Returns (cj_dict, notice_str). Raises RuntimeError if scores unavailable.
    """
    # Extract the Vue `cj` data object from the inline script
    cj_match = re.search(r'\bcj\s*:\s*(\{.*?\}|null)', html, re.DOTALL)
    if not cj_match:
        raise RuntimeError("无法解析成绩数据")

    raw = cj_match.group(1).strip()
    if raw == "null":
        msg_match = re.search(r'\bmsg\s*:\s*["\']([^"\']*)["\']', html)
        msg = msg_match.group(1) if msg_match else "请检查报考信息或成绩查询尚未开放"
        raise RuntimeError(f"无查询结果：{msg}")

    try:
        cj = json.loads(raw)
    except json.JSONDecodeError:
        raise RuntimeError(f"成绩数据解析失败: {raw[:200]}")

    notice = cj.pop("zsdwsm", "") or ""
    return cj, notice


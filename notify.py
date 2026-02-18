"""
Email notification module: sends score results as an HTML table via SMTP.
"""

import os
import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from dotenv import load_dotenv

load_dotenv()

# Field display order matching the official page
SCORE_FIELDS = [
    ("姓名",    "xm"),
    ("报名号",   "bmh"),
    ("准考证号", "ksbh"),
    ("总分",    "zf"),
    ("第一门",  "km1"),
    ("第二门",  "km2"),
    ("第三门",  "km3"),
    ("第四门",  "km4"),
    ("备注",    "bz"),
]


def build_html(cj: dict, notice: str = "") -> str:
    rows = ""
    for label, key in SCORE_FIELDS:
        val = cj.get(key) or cj.get(label, "—")
        rows += f"""
        <tr>
          <td style="padding:12px 16px;font-weight:600;color:#555;border-bottom:1px solid #eee;white-space:nowrap">{label}</td>
          <td style="padding:12px 16px;color:#333;border-bottom:1px solid #eee">{val}</td>
        </tr>"""

    notice_block = ""
    if notice:
        notice_block = f"""
        <div style="background:#fffbe6;border-left:4px solid #f0c040;padding:10px 14px;margin-bottom:16px;font-size:14px;color:#555">
          招生单位说明：{notice}
        </div>"""

    return f"""<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;background:#f5f5f5;padding:24px">
  <div style="max-width:520px;margin:0 auto;background:#fff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,.08)">
    <div style="background:#1887e0;color:#fff;text-align:center;padding:18px;font-size:18px;font-weight:700">
      2026年部分考生初试成绩查询
    </div>
    <div style="padding:16px">
      {notice_block}
      <table style="width:100%;border-collapse:collapse">
        {rows}
      </table>
    </div>
  </div>
</body>
</html>"""


def send_score_email(cj: dict, notice: str = "") -> None:
    """
    Sends score results via email using environment variables:
      EMAIL_SMTP_SERVER, EMAIL_SMTP_PORT
      EMAIL_SENDER, EMAIL_AUTH_CODE, EMAIL_RECEIVER
    """
    smtp_server  = os.environ["EMAIL_SMTP_SERVER"]
    smtp_port    = int(os.environ["EMAIL_SMTP_PORT"])
    sender       = os.environ["EMAIL_SENDER"]
    auth_code    = os.environ["EMAIL_AUTH_CODE"]
    receiver     = os.environ["EMAIL_RECEIVER"]

    name = cj.get("xm") or cj.get("姓名", "")
    total = cj.get("zf") or cj.get("总分", "")

    msg = MIMEMultipart("alternative")
    msg["Subject"] = f"【研招网成绩】{name} 总分 {total}"
    msg["From"]    = sender
    msg["To"]      = receiver

    msg.attach(MIMEText(build_html(cj, notice), "html", "utf-8"))

    with smtplib.SMTP_SSL(smtp_server, smtp_port) as smtp:
        smtp.login(sender, auth_code)
        smtp.sendmail(sender, receiver, msg.as_string())

    print(f"Email sent to {receiver}")

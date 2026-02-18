"""
Entry point: logs in to CHSI, queries master's entrance exam score, and sends email.
"""

from login import login
from query import query_score, parse_score
from notify import send_score_email


def main():
    print("Logging in to CHSI...")
    session = login()
    print("Login successful.")

    print("Querying score...")
    html = query_score(session)
    cj, notice = parse_score(html)
    print("Score:", {k: cj[k] for k in ("xm", "zf") if k in cj})

    print("Sending email...")
    send_score_email(cj, notice)


if __name__ == "__main__":
    main()

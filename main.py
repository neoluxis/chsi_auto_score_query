"""
Entry point: logs in to CHSI and queries master's entrance exam score.
"""

from login import login
from query import query_score, parse_score


def main():
    print("Logging in to CHSI...")
    session = login()
    print("Login successful.")

    print("Querying score...")
    html = query_score(session)
    result = parse_score(html)
    print("Score result:", result)


if __name__ == "__main__":
    main()

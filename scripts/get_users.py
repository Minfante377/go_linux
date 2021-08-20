import sqlite3
from sys import argv

db = sqlite3.connect(argv[1])
cursor = db.cursor()
cursor.execute("SELECT name FROM sqlite_master WHERE type='table';")
table = cursor.fetchone()[0]
cursor.execute("SELECT user_id FROM {};".format(table))
res = cursor.fetchall()
ids = []
for r in res:
    ids.append(r[0])
print(" ".join(_id for _id in ids))

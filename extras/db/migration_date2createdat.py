#!/bin/env python3

import os
import sqlite3
import pprint
import datetime
import time
import pytz


# Open a connection to the database
database_path = os.environ['DATABASE']
conn = sqlite3.connect(database_path)
c = conn.cursor()

c.execute("SELECT id,date FROM definitions ORDER BY id desc")
results = c.fetchall()

for item in results:
    #pprint.pprint(item)
    try:
        id = item[0]
        original_date = "{} 00:00:00".format(item[1])
        created_at = datetime.datetime.strptime(original_date, '%Y-%m-%d %H:%M:%S')
        created_at = created_at.replace(tzinfo=pytz.UTC)
        created_at_unixtimestamp = int(created_at.timestamp())
        print("Converting id: {} date: {} to {}".format(id, original_date, created_at_unixtimestamp))

        c.execute("UPDATE definitions SET created_at = ? WHERE id = ?", (created_at_unixtimestamp, id))
        conn.commit()


    except ValueError as err:
        print("Error converting {}".format(id))

#!/bin/env python3

import os
import sqlite3
import datetime
import time
import pytz

def convertDateToUnixTimestamp(date_source, timezone="CLT"):
    """
    convertDateToUnixTimestamp date must be in format "yyyy-mm-dd"
    """
    try:
        tz = pytz.timezone('Chile/Continental')
        orignal_date = datetime.datetime.strptime(date_source, "%Y-%m-%d")
        return int(orignal_date.replace(tzinfo=tz).timestamp())
    except Exception as error:
        print(error)
        return int(time.time())

def migrate_users(source, target):
    return True

def migrate_definitions(source, target):
    conn_source = sqlite3.connect(source)
    c_source = conn_source.cursor()
    conn_source.text_factory = lambda x: str(x, 'latin1')

    conn_target = sqlite3.connect(target)
    c_target = conn_target.cursor()

    c_target.execute("DELETE FROM definitions;")

    c_source.execute("SELECT * FROM definitions;")
    results_source = c_source.fetchall()

    for item in results_source:

        try:
            term = item[1]
            meaning = item[2]
            author = item[3]
            locked = None
            active= item[5]
            created_at = convertDateToUnixTimestamp(item[6])
            updated_at = created_at
            hits = int(item[7])
            link = int(item[8])
            locked_by = None
            chat = "#linux"
            deleted_by = None
            deleted_at = None

            stmt = "INSERT INTO definitions (term, meaning, author, chat, hits, link, created_at, updated_at, deleted_at, locked, locked_by, deleted_by) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
            c_target.execute(stmt, (term, meaning, author, chat, hits, link, created_at, updated_at, deleted_at, locked, locked_by, deleted_by))

            conn_target.commit()
            print("Imported id: {}".format(item[0]))
        
        except Exception as err:
            print(err)
            print("Error importing id: {}".format(item[0]))
        
    conn_source.close()
    conn_target.close()


if __name__ == "__main__":

    # Open a connection to the databases
    database_source = os.environ['DATABASE_SOURCE']
    database_target= os.environ['DATABASE_TARGET']

    migrate_users(database_source, database_target)
    migrate_definitions(database_source, database_target)
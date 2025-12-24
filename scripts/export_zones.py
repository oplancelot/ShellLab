import db_config
import mysql.connector
import sqlite3
import os

mysql_config = db_config.get_mysql_config()
SQLITE_DB = os.path.join(os.path.dirname(__file__), '..', 'data', 'shelllab.db')

def export_zones():
    print("Connecting to MySQL...")
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Fetching zones...")
        mysql_cursor.execute("SELECT areatableID, name_loc0 FROM aowow_zones")
        zones = mysql_cursor.fetchall()
        print(f"Found {len(zones)} zones")
        
        mysql_cursor.close()
        mysql_conn.close()
        
        print(f"Connecting to SQLite: {SQLITE_DB}")
        sqlite_conn = sqlite3.connect(SQLITE_DB)
        c = sqlite_conn.cursor()
        
        c.execute("DROP TABLE IF EXISTS zones")
        c.execute("CREATE TABLE zones (id INTEGER PRIMARY KEY, name TEXT)")
        
        print("Inserting zones...")
        for z in zones:
            c.execute("INSERT INTO zones VALUES (?,?)", (z['areatableID'], z['name_loc0']))
            
        # Also need sort IDs (QuestSort.dbc), but for now zones are good start
        
        sqlite_conn.commit()
        sqlite_conn.close()
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_zones()

import db_config
import mysql.connector
import sqlite3
import os

mysql_config = db_config.get_mysql_config()
SQLITE_DB = os.path.join(os.path.dirname(__file__), '..', 'data', 'shelllab.db')

def export_relations():
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Connecting to SQLite...")
        sqlite_conn = sqlite3.connect(SQLITE_DB)
        c = sqlite_conn.cursor()
        
        tables = [
            ('creature_questrelation', 'npc_quest_start'),
            ('creature_involvedrelation', 'npc_quest_end'),
            ('gameobject_questrelation', 'go_quest_start'),
            ('gameobject_involvedrelation', 'go_quest_end')
        ]
        
        for mysql_table, sqlite_table in tables:
            print(f"Exporting {mysql_table} to {sqlite_table}...")
            mysql_cursor.execute(f"SELECT id, quest FROM {mysql_table}")
            rows = mysql_cursor.fetchall()
            
            c.execute(f"DROP TABLE IF EXISTS {sqlite_table}")
            c.execute(f"CREATE TABLE {sqlite_table} (entry INTEGER, quest INTEGER)")
            c.executemany(f"INSERT INTO {sqlite_table} VALUES (?,?)", 
                [(r['id'], r['quest']) for r in rows])
            
            c.execute(f"CREATE INDEX idx_{sqlite_table}_entry ON {sqlite_table}(entry)")
            c.execute(f"CREATE INDEX idx_{sqlite_table}_quest ON {sqlite_table}(quest)")
            print(f"Imported {len(rows)} rows.")
            
        sqlite_conn.commit()
        sqlite_conn.close()
        mysql_cursor.close()
        mysql_conn.close()
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_relations()

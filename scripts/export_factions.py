import db_config
import mysql.connector
import json
import os

mysql_config = db_config.get_mysql_config()
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

def export_factions():
    print("Connecting to MySQL...")
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Fetching factions...")
        # Note: Querying aowow database might need specific handling if db_config defaults to tw_world
        # But get_mysql_config() usually handles env vars. 
        # The original script queried 'aowow_factions', assuming default DB or table prefix.
        # Let's ensure we use the correct database or table qualification.
        
        # If the user configuration points to tw_world, we might need to prefix: aowow.aowow_factions
        # Or switch database. Safer to assume table exists in connected DB or use qualified name.
        # The original script worked, so let's stick to the SQL, but maybe add database context if needed.
        
        # Based on refactoring, let's use db_config.get_aowow_db() to be safe if that's where it lives.
        # Although previous script didn't explicitly switch DB in the query string, it relied on config.
        
        # We will use the 'aowow' database explicitly for this query context.
        aowow_db = db_config.get_aowow_db()
        mysql_cursor.execute(f"USE {aowow_db}")
        
        mysql_cursor.execute("""
            SELECT factionID, name_loc0, description1_loc0, side, team 
            FROM aowow_factions 
            ORDER BY side, name_loc0
        """)
        factions = mysql_cursor.fetchall()
        print(f"Found {len(factions)} factions")
        mysql_cursor.close()
        mysql_conn.close()
        
        os.makedirs(DATA_DIR, exist_ok=True)
        json_path = os.path.join(DATA_DIR, 'factions.json')
        print(f"Exporting to {json_path}...")
        
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(factions, f, ensure_ascii=False)
            
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_factions()

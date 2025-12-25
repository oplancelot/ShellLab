import db_config
import mysql.connector
import json
import os

mysql_config = db_config.get_mysql_config()
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

TABLES = [
    'item_template',
    'creature_template',
    'quest_template',
    'spell_template',
    'gameobject_template'
]

def inspect_schemas():
    print("Connecting to MySQL...")
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        cursor = mysql_conn.cursor(dictionary=True)
        
        schemas = {}
        
        for table in TABLES:
            print(f"Inspecting {table}...")
            # Switch DB for gameobject_template if strictly needed, but usually all in world db for Turtle?
            # Actually gameobject_template might be in a different DB if checks above were using get_tw_world_db?
            # Let's assume standard config first, if table not found, we handle it.
            
            try:
                cursor.execute(f"DESCRIBE {table}")
                columns = cursor.fetchall()
                schemas[table] = columns
            except mysql.connector.Error as err:
                print(f"Error inspecting {table}: {err}")
                # Try switching to tw_world if configured specifically?
                # The users existing export_objects_mysql.py used db_config.get_tw_world_db().
                # Let's try that context if the first fails.
                pass

        cursor.close()
        mysql_conn.close()
        
        # Save to file
        os.makedirs(DATA_DIR, exist_ok=True)
        out_path = os.path.join(DATA_DIR, 'schema_dump.json')
        with open(out_path, 'w', encoding='utf-8') as f:
            json.dump(schemas, f, indent=2, default=str)
            
        print(f"Schema dump saved to {out_path}")
        
    except Exception as e:
        print(f"Fatal Error: {e}")

if __name__ == '__main__':
    inspect_schemas()
